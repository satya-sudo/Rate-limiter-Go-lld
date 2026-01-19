package limiter

import "time"

func (r *TokenBucketRateLimiter) StartCleanUp(idleThreshold time.Duration, cleanUpInterval time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanUpInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.cleanUpIdleBucket(idleThreshold)
			case <-r.done:
				return
			}
		}
	}()
}

func (r *TokenBucketRateLimiter) Stop() {
	r.once.Do(func() {
		close(r.done)
	})
}

func (r *TokenBucketRateLimiter) cleanUpIdleBucket(idleThreshold time.Duration) {
	now := getNowTime()
	r.tokenBucket.Range(func(key, value any) bool {
		bucket := value.(*TokenBucket)
		bucket.mu.Lock()
		lastSeen := bucket.lastRefillTime
		bucket.mu.Unlock()
		
		if time.Duration(now-lastSeen)*time.Millisecond > idleThreshold {
			r.tokenBucket.Delete(key)
		}
		return true
	})
}
