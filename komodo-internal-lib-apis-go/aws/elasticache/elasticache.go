package elasticache

import (
	"context"
	"errors"
	sm "komodo-internal-lib-apis-go/aws/secrets-manager"
	"komodo-internal-lib-apis-go/config"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type ElasticacheConnector struct {
	Endpoint string
	Password string
	Client   *redis.Client
}

const DEFAULT_SESH_TTL int64 = 3600

var (
	ElasticacheClient *ElasticacheConnector
	initOnce          sync.Once
)

// Initializes a singleton Elasticache client. In
// prod/staging it will connect to the configured Redis/ElastiCache
// endpoint. 
// fallback that implements GET/SET/DELETE semantics with TTLs.
func InitElasticacheClient() error {
	var initErr error

	// only initialize once
	initOnce.Do(func() {
		logger.Info("initializing AWS Elasticache client")
		return

		endpoint := config.GetConfigValue("ELASTICACHE_ENDPOINT")
		secrets, err := sm.GetSecrets([]string{"ELASTICACHE_PASSWORD"})
		if err != nil {
			initErr = err
			logger.Error("failed to load ELASTICACHE_PASSWORD secret", err)
			return
		}

		// create redis client
		client := redis.NewClient(&redis.Options{
			Addr:     endpoint,
			Password: secrets["ELASTICACHE_PASSWORD"],
			DB:       0,
		})

		// Ping with timeout to verify connectivity
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
	
		if err := client.Ping(ctx).Err(); err != nil {
			initErr = err
			logger.Error("failed to ping Elasticache", err)
			return
		}

		ElasticacheClient = &ElasticacheConnector{
			Endpoint: endpoint,
			Password: secrets["ELASTICACHE_PASSWORD"],
			Client:   client,
		}
	})
	return initErr
}

// GetCacheItem returns the string value stored at key. If the key does not exist, it returns an error.
func GetCacheItem(key string) (string, error) {
	if ElasticacheClient != nil && ElasticacheClient.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
		defer cancel()
		val, err := ElasticacheClient.Client.Get(ctx, key).Result()

		if err == redis.Nil { return "", nil }
		if err != nil { return "", err }
		return val, nil
	}
	logger.Error("elasticache client not available")
	return "", errors.New("elasticache client not available")
}

// SetCacheItem stores a value with the provided TTL (in seconds). Use ttl<=0
func SetCacheItem(key string, value string, ttl int64) error {
	if ElasticacheClient != nil && ElasticacheClient.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
		defer cancel()

		var dur time.Duration
		if ttl > 0 {
			dur = time.Duration(ttl) * time.Second
		}
		return ElasticacheClient.Client.Set(ctx, key, value, dur).Err()
	}
	logger.Error("elasticache client not available")
	return errors.New("elasticache client not available")
}

// DeleteCacheItem removes a key from the store.
func DeleteCacheItem(key string) error {
	if ElasticacheClient != nil && ElasticacheClient.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
		defer cancel()
		return ElasticacheClient.Client.Del(ctx, key).Err()
	}
	logger.Error("elasticache client not available")
	return errors.New("elasticache client not available")
}

// CloseElasticache closes any underlying clients and stops background goroutines.
func CloseElasticache() error {
	if ElasticacheClient == nil { return nil }
	if ElasticacheClient.Client != nil {
		return ElasticacheClient.Client.Close()
	}
	return nil
}

// token bucket Lua script (atomic): returns {allowed, wait_ms}
var tokenBucketScript = redis.NewScript(`
local now = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local burst = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local ttl = tonumber(ARGV[5])

local data = redis.call('HMGET', KEYS[1], 'tokens', 'ts')
local tokens = tonumber(data[1])
local ts = tonumber(data[2])
if tokens == nil then
  tokens = burst
  ts = now
end
local elapsed = (now - ts) / 1000.0
if elapsed < 0 then elapsed = 0 end
local new_tokens = tokens + elapsed * rate
if new_tokens > burst then new_tokens = burst end
local allowed = 0
local wait_ms = 0
if new_tokens >= requested then
  new_tokens = new_tokens - requested
  allowed = 1
else
  local deficit = requested - new_tokens
  if rate > 0 then
	wait_ms = math.ceil((deficit / rate) * 1000)
  else
	wait_ms = 0
  end
end
redis.call('HMSET', KEYS[1], 'tokens', tostring(new_tokens), 'ts', tostring(now))
redis.call('EXPIRE', KEYS[1], ttl)
return {allowed, tostring(wait_ms)}
`)

// AllowDistributed attempts to consume a single token from a distributed
// token bucket stored in Elasticache/Redis. It returns (allowed, retryAfter, err).
// If Elasticache is not configured/available it returns an error.
func AllowDistributed(ctx context.Context, key string) (bool, time.Duration, error) {
	if ElasticacheClient == nil {
		return false, 0, errors.New("elasticache client not initialized")
	}
	if ElasticacheClient.Client == nil {
		return false, 0, errors.New("elasticache client not available (local fallback in use)")
	}

	// Read rate config from env with defaults (match middleware defaults)
	rate := float64(10)
	burst := float64(20)
	if v := strings.TrimSpace(config.GetConfigValue("RATE_LIMIT_RPS")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 { rate = f }
	}
	if v := strings.TrimSpace(config.GetConfigValue("RATE_LIMIT_BURST")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 { burst = f }
	}

	ttlSec := 300
	if v := strings.TrimSpace(config.GetConfigValue("RATE_LIMIT_BUCKET_TTL_SEC")); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 { ttlSec = i }
	}

	now := time.Now().UnixMilli()
	// Execute script
	res, err := tokenBucketScript.Run(ctx, ElasticacheClient.Client, []string{key}, now, rate, burst, 1, ttlSec).Result()
	if err != nil {
		return false, 0, err
	}

	// script returns [allowed, wait_ms]
	arr, ok := res.([]interface{})
	if !ok || len(arr) < 2 {
		return false, 0, errors.New("unexpected script result")
	}

	// allowed may be number or string
	var allowed bool
	switch v := arr[0].(type) {
		case int64:
			allowed = v == 1
		case string:
			allowed = v == "1"
		default:
			allowed = false
	}

	var waitMs int64
	switch v := arr[1].(type) {
		case int64:
			waitMs = v
		case string:
			if parsed, err := strconv.ParseInt(v, 10, 64); err == nil { waitMs = parsed }
	}

	return allowed, time.Duration(waitMs) * time.Millisecond, nil
}
