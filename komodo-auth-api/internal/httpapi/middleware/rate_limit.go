package middleware

import (
	"komodo-auth-api/internal/config"
	"komodo-auth-api/internal/logger"
	"komodo-auth-api/internal/thirdparty/aws"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type bucket struct {
	mu      sync.Mutex
	tokens  float64
	last    time.Time
	created time.Time
}

var (
	rlOnce     sync.Once
	rps        float64
	burst      float64
	buckets    sync.Map // map[string]*bucket
	evictOnce  sync.Once
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		env := strings.ToLower(config.GetConfigValue("API_ENV"))
		key := clientKey(req)

		if !strings.EqualFold(config.GetConfigValue("USE_MOCKS"), "true") && (env == "prod" || env == "staging") {
			// Try distributed token consume via Elasticache/Redis. If Redis
			// is not available, fall back to the local in-process bucket.
			ctx := req.Context()
			allowed, wait, err := aws.AllowDistributed(ctx, key)

			if err != nil {
				// If the Elasticache client isn't initialized/available, fall back
				// to local bucket behavior. Otherwise, decide fail-open vs fail-closed.
				logger.Error("Elasticache rate limiter error: "+err.Error(), req)

				if strings.Contains(err.Error(), "not available") || strings.Contains(err.Error(), "not initialized") {
					// local process bucket fallback
					bucket := getBucket(key)
					if !bucket.allow() {
						if wait := bucket.retryAfter(); wait > 0 {
							wtr.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds()+0.5)))
						}
						logger.Error("rate limit exceeded for client: " + key, req)
						http.Error(wtr, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
						return
					}
				} else {
					// Elasticache is present but errored. Honor fail-open config.
					if shouldFailOpen() {
						logger.Error("rate limiter failing open for client: " + key, req)
					} else {
						http.Error(wtr, "Service unavailable", http.StatusServiceUnavailable)
						return
					}
				}
			} else {
				if !allowed {
					if wait > 0 {
						// convert to seconds for Retry-After header
						wtr.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds()+0.5)))
					}
					logger.Error("rate limit exceeded for client: " + key, req)
					http.Error(wtr, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
					return
				}
				// allowed via distributed token bucket
			}
		} else {
			// local process bucket
			bucket := getBucket(key)

			if !bucket.allow() {
				// Rate limit exceeded
				if wait := bucket.retryAfter(); wait > 0 {
					wtr.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds()+0.5)))
				}
				logger.Error("rate limit exceeded for client: " + key, req)
				http.Error(wtr, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(wtr, req)
	})
}

// allow checks and updates the bucket token count
func (bkt *bucket) allow() bool {
	rps, burst := rateConfig()
	now := time.Now()
	bkt.mu.Lock()

	// Refill tokens based on elapsed time
	if !bkt.last.IsZero() {
		elapsed := now.Sub(bkt.last).Seconds()
		if elapsed > 0 {
			bkt.tokens += elapsed * rps
			if bkt.tokens > burst {
				bkt.tokens = burst
			}
		}
	} else {
		bkt.tokens = burst
	}

	allowed := false
	if bkt.tokens >= 1 {
		bkt.tokens -= 1
		allowed = true
	}

	bkt.last = now
	bkt.mu.Unlock()

	return allowed
}

// retryAfter estimates how long until the next token is available
func (bkt *bucket) retryAfter() time.Duration {
	rps, _ := rateConfig()
	if rps <= 0 {
		return time.Second
	}

	bkt.mu.Lock()
	defer bkt.mu.Unlock()
	deficit := 1 - bkt.tokens
	if deficit <= 0 {
		return 0
	}
	secs := deficit / rps
	return time.Duration(secs * float64(time.Second))
}

// rateConfig reads and caches rate limit settings from env vars
func rateConfig() (float64, float64) {
	rlOnce.Do(func() {
		// Helper to parse float env var with default
		parseFloatEnv := func(key string, dflt float64) float64 {
			if val := strings.TrimSpace(config.GetConfigValue(key)); val != "" {
				if f, err := strconv.ParseFloat(val, 64); err == nil {
					return f
				}
			}
			return dflt
		}

		rps = parseFloatEnv("RATE_LIMIT_RPS", 10)    // default 10 req/sec
		burst = parseFloatEnv("RATE_LIMIT_BURST", 20) // default burst 20
		// stricter validation: treat non-positive rps as invalid and reset
		if rps <= 0 {
			rps = 10
		}
		if burst < 1 {
			burst = 20
		}
	})
	return rps, burst
}

// helper to compute client key (IP)
func clientKey(req *http.Request) string {
	// prefer first X-Forwarded-For entry when present
	if xf := req.Header.Get("X-Forwarded-For"); xf != "" {
		parts := strings.Split(xf, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}
	// fallback to remote addr host
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return req.RemoteAddr
}

// getBucket retrieves or creates a rate limit bucket for the given key
func getBucket(key string) *bucket {
	// ensure the background evictor is running
	evictOnce.Do(startBucketEvictor)

	if v, ok := buckets.Load(key); ok {
		return v.(*bucket)
	}
	b := &bucket{tokens: 0, last: time.Time{}, created: time.Now()}
	actual, _ := buckets.LoadOrStore(key, b)
	return actual.(*bucket)
}

// startBucketEvictor removes idle buckets after configured TTL
func startBucketEvictor() {
	ttlSec := 300
	if val := strings.TrimSpace(config.GetConfigValue("RATE_LIMIT_BUCKET_TTL_SEC")); val != "" {
		if i, err := strconv.Atoi(val); err == nil && i > 0 {
			ttlSec = i
		}
	}

	ttl := time.Duration(ttlSec) * time.Second
	ticker := time.NewTicker(time.Minute)

	go func() {
		for range ticker.C {
			now := time.Now()
			buckets.Range(func(key, val any) bool {
				bucket := val.(*bucket)

				bucket.mu.Lock()
				lastActive := bucket.last
				if lastActive.IsZero() {
					lastActive = bucket.created
				}
				bucket.mu.Unlock()

				if now.Sub(lastActive) > ttl {
					buckets.Delete(key)
				}
				return true
			})
		}
	}()
}

// shouldFailOpen decides fail-open vs fail-closed when the distributed store is unavailable
func shouldFailOpen() bool {
	v := strings.ToLower(strings.TrimSpace(config.GetConfigValue("RATE_LIMIT_FAIL_OPEN")))
	if v == "" {
		return true // default to fail-open to reduce customer impact
	}
	return v == "true" || v == "1" || v == "yes"
}
