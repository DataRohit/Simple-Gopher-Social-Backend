package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]time.Time
	mu       sync.Mutex
	limit    time.Duration
}

func NewRateLimiter(limit time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]time.Time),
		limit:    limit,
	}
}

func (rl *RateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if lastVisit, found := rl.visitors[ip]; found {
		if time.Since(lastVisit) < rl.limit {
			return false, rl.limit - time.Since(lastVisit)
		}
	}

	rl.visitors[ip] = time.Now()
	return true, 0
}
