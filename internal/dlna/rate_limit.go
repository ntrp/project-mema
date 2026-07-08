package dlna

import (
	"sync"
	"time"
)

const (
	defaultDLNARateLimit  = 240
	defaultDLNARateWindow = time.Minute
)

type dlnaRateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	buckets map[string]dlnaRateBucket
}

type dlnaRateBucket struct {
	Count    int
	ResetAt  time.Time
	LastSeen time.Time
}

func newDLNARateLimiter(limit int, window time.Duration) *dlnaRateLimiter {
	return &dlnaRateLimiter{limit: limit, window: window, buckets: map[string]dlnaRateBucket{}}
}

func (l *dlnaRateLimiter) Allow(key string, now time.Time) bool {
	if l == nil || l.limit <= 0 || l.window <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prune(now)
	bucket := l.buckets[key]
	if bucket.ResetAt.IsZero() || !now.Before(bucket.ResetAt) {
		bucket = dlnaRateBucket{ResetAt: now.Add(l.window)}
	}
	if bucket.Count >= l.limit {
		l.buckets[key] = bucket
		return false
	}
	bucket.Count++
	bucket.LastSeen = now
	l.buckets[key] = bucket
	return true
}

func (l *dlnaRateLimiter) prune(now time.Time) {
	for key, bucket := range l.buckets {
		if now.Sub(bucket.LastSeen) > 2*l.window {
			delete(l.buckets, key)
		}
	}
}
