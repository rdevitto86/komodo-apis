package aws

import (
	"context"
	"errors"
	"komodo-auth-api/internal/config"
	"log"
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
	mem      *localCache
}

const DEFAULT_SESH_TTL int64 = 3600

var ElasticacheClient *ElasticacheConnector
var initOnce sync.Once

// InitElasticacheClient initializes a singleton Elasticache client. In
// prod/staging it will connect to the configured Redis/ElastiCache
// endpoint. In local or when USE_MOCKS=true it creates an in-memory
// fallback that implements GET/SET/DELETE semantics with TTLs.
func InitElasticacheClient() error {
	var initErr error
	initOnce.Do(func() {
		env := config.GetConfigValue("API_ENV")
		useMocks := strings.EqualFold(config.GetConfigValue("USE_MOCKS"), "true")

		// Local/dev/mocks: initialize an in-memory store so callers can
		// use the same API without an external dependency.
		if useMocks || (env != "prod" && env != "staging") {
			log.Println("InitElasticacheClient: initializing in-memory Elasticache client for local/DEV environment")
			ElasticacheClient = &ElasticacheConnector{mem: newLocalCache()}
			return
		}

		endpoint := config.GetConfigValue("ELASTICACHE_ENDPOINT")
		secrets, err := GetSecrets([]string{"ELASTICACHE_PASSWORD"})
		if err != nil {
			initErr = err
			log.Printf("InitElasticacheClient: failed to load ELASTICACHE_PASSWORD secret: %v", err)
			return
		}

		// create redis client
		client := redis.NewClient(&redis.Options{
			Addr:     endpoint,
			Password: secrets["ELASTICACHE_PASSWORD"],
			DB:       0,
		})

		// Ping with timeout to verify connectivity
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			initErr = err
			log.Printf("InitElasticacheClient: ping failed: %v", err)
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
	if ElasticacheClient == nil {
		return "", errors.New("elasticache client not initialized")
	}

	if ElasticacheClient.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		val, err := ElasticacheClient.Client.Get(ctx, key).Result()
		if err == redis.Nil { return "", nil }
		if err != nil { return "", err }
		return val, nil
	}

	// in-memory fallback
	return ElasticacheClient.mem.Get(key), nil
}

// SetCacheItem stores a value with the provided TTL (in seconds). Use ttl<=0
func SetCacheItem(key string, value string, ttl int64) error {
	if ElasticacheClient == nil {
		return errors.New("elasticache client not initialized")
	}

	if ElasticacheClient.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		var dur time.Duration
		if ttl > 0 {
			dur = time.Duration(ttl) * time.Second
		}
		return ElasticacheClient.Client.Set(ctx, key, value, dur).Err()
	}

	// in-memory fallback
	ElasticacheClient.mem.Set(key, value, ttl)
	return nil
}

// DeleteCacheItem removes a key from the store.
func DeleteCacheItem(key string) error {
	if ElasticacheClient == nil {
		return errors.New("elasticache client not initialized")
	}

	if ElasticacheClient.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return ElasticacheClient.Client.Del(ctx, key).Err()
	}

	ElasticacheClient.mem.Delete(key)
	return nil
}

// CloseElasticache closes any underlying clients and stops background goroutines.
func CloseElasticache() error {
	if ElasticacheClient == nil {
		return nil
	}
	if ElasticacheClient.Client != nil {
		return ElasticacheClient.Client.Close()
	}
	if ElasticacheClient.mem != nil {
		ElasticacheClient.mem.Close()
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

// ----------------
// in-memory cache
// ----------------
type memItem struct {
	value     string
	expiresAt time.Time
}

type localCache struct {
	mu      sync.RWMutex
	items   map[string]memItem
	ticker  *time.Ticker
	stopCh  chan struct{}
}

func newLocalCache() *localCache {
	lc := &localCache{
		items:  make(map[string]memItem),
		ticker: time.NewTicker(time.Minute),
		stopCh: make(chan struct{}),
	}
	go lc.janitor()
	return lc
}

func (cache *localCache) janitor() {
	for {
		select {
			case <-cache.ticker.C:
				now := time.Now()
				cache.mu.Lock()
				for key, val := range cache.items {
					if !val.expiresAt.IsZero() && now.After(val.expiresAt) {
						delete(cache.items, key)
					}
				}
				cache.mu.Unlock()
			case <-cache.stopCh:
				cache.ticker.Stop()
				return
		}
	}
}

func (cache *localCache) Get(key string) string {
	cache.mu.RLock()
	itm, ok := cache.items[key]
	cache.mu.RUnlock()
	if !ok { return "" }

	if !itm.expiresAt.IsZero() && time.Now().After(itm.expiresAt) {
		// expired
		cache.mu.Lock()
		delete(cache.items, key)
		cache.mu.Unlock()
		return ""
	}
	return itm.value
}

func (cache *localCache) Set(key string, val string, ttl int64) {
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	cache.mu.Lock()
	cache.items[key] = memItem{value: val, expiresAt: exp}
	cache.mu.Unlock()
}

func (cache *localCache) Delete(key string) {
	cache.mu.Lock()
	delete(cache.items, key)
	cache.mu.Unlock()
}

func (cache *localCache) Close() {
	close(cache.stopCh)
}
