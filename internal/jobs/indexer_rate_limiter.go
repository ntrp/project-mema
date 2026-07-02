package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/indexers"
)

const (
	indexerQuerySpacing = 1500 * time.Millisecond
	indexer429Backoff   = 5 * time.Minute
)

var errIndexerBackoffActive = errors.New("indexer backoff active")

type indexerRateLimiter struct {
	nextAllowed map[uuid.UUID]time.Time
	blocked     map[uuid.UUID]time.Time
	spacing     time.Duration
	backoff     time.Duration
}

func newIndexerRateLimiter() *indexerRateLimiter {
	return newIndexerRateLimiterWith(indexerQuerySpacing, indexer429Backoff)
}

func newIndexerRateLimiterWith(spacing time.Duration, backoff time.Duration) *indexerRateLimiter {
	return &indexerRateLimiter{
		nextAllowed: map[uuid.UUID]time.Time{},
		blocked:     map[uuid.UUID]time.Time{},
		spacing:     spacing,
		backoff:     backoff,
	}
}

func (l *indexerRateLimiter) wait(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	if blockedUntil, ok := l.blocked[id]; ok && now.Before(blockedUntil) {
		return fmt.Errorf("%w until %s", errIndexerBackoffActive, blockedUntil.Format(time.RFC3339))
	}
	waitUntil := l.nextAllowed[id]
	if waitUntil.After(now) {
		timer := time.NewTimer(time.Until(waitUntil))
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	l.nextAllowed[id] = time.Now().Add(l.spacing)
	return nil
}

func (l *indexerRateLimiter) recordError(id uuid.UUID, err error) {
	status := indexers.StatusCode(err)
	if status == nil || *status != 429 {
		return
	}
	backoff := indexers.RetryAfter(err)
	if backoff <= 0 {
		backoff = l.backoff
	}
	l.blocked[id] = time.Now().Add(backoff)
}
