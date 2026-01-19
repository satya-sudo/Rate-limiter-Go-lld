package limiter

import "time"

type RateLimiter interface {
	Allow(key string) bool
}

type ManageLimiter interface {
	RateLimiter
	StartCleanUp(idleThreshold time.Duration, cleanUpInterval time.Duration)
	Stop()
}
