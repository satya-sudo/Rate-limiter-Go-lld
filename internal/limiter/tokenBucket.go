package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens         float64
	lastRefillTime int64
	mu             sync.Mutex
}

type TokenBucketRateLimiter struct {
	capacity    float64
	refillRate  float64
	tokenBucket sync.Map
	done        chan struct{}
	once        sync.Once
}

func NewTokenBucket(capacity float64) *TokenBucket {
	return &TokenBucket{
		tokens:         capacity,
		lastRefillTime: getNowTime(),
	}
}

func NewTokenBucketRateLimiter(capacity float64, refillRate float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		capacity:    capacity,
		refillRate:  refillRate,
		tokenBucket: sync.Map{},
		done:        make(chan struct{}),
		once:        sync.Once{},
	}
}

func (r *TokenBucketRateLimiter) Allow(key string) bool {
	bucket := r.getOrCreateBucket(key)
	bucket.mu.Lock() // lock once for all updates
	defer bucket.mu.Unlock()
	now := getNowTime()
	bucket.refill(now, r.capacity, r.refillRate)
	return bucket.tryConsume()

}
func (r *TokenBucketRateLimiter) getOrCreateBucket(key string) *TokenBucket {
	bucket, ok := r.tokenBucket.Load(key)
	if !ok {
		bucket = NewTokenBucket(r.capacity)
		r.tokenBucket.Store(key, bucket)
	}
	return bucket.(*TokenBucket)
}

func (b *TokenBucket) refill(now int64, capacity float64, refillRate float64) {
	timeLapse := float64(now-b.lastRefillTime) / 1000.0
	newTokenCount := b.tokens + (timeLapse * refillRate)
	b.tokens = min(capacity, newTokenCount)
	b.lastRefillTime = now
}

func (b *TokenBucket) tryConsume() bool {
	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}
	return false
}

func getNowTime() int64 {
	return time.Now().UnixMilli()
}
