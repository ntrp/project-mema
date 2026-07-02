package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/indexers"
)

func TestIndexerRateLimiterSpacesQueries(t *testing.T) {
	spacing := 20 * time.Millisecond
	limiter := newIndexerRateLimiterWith(spacing, time.Minute)
	id := uuid.New()
	if err := limiter.wait(context.Background(), id); err != nil {
		t.Fatalf("first wait failed: %v", err)
	}
	started := time.Now()
	if err := limiter.wait(context.Background(), id); err != nil {
		t.Fatalf("second wait failed: %v", err)
	}
	if elapsed := time.Since(started); elapsed < spacing {
		t.Fatalf("expected rate limit wait near %s, got %s", spacing, elapsed)
	}
}

func TestIndexerRateLimiterBacksOffAfter429(t *testing.T) {
	limiter := newIndexerRateLimiterWith(time.Millisecond, time.Minute)
	id := uuid.New()
	limiter.recordError(id, indexers.StatusError{StatusCode: 429})
	err := limiter.wait(context.Background(), id)
	if !errors.Is(err, errIndexerBackoffActive) {
		t.Fatalf("expected backoff error, got %v", err)
	}
}

func TestIndexerRateLimiterIgnoresNon429(t *testing.T) {
	limiter := newIndexerRateLimiterWith(time.Millisecond, time.Minute)
	id := uuid.New()
	limiter.recordError(id, indexers.StatusError{StatusCode: 500})
	if err := limiter.wait(context.Background(), id); err != nil {
		t.Fatalf("expected no backoff after non-429: %v", err)
	}
}
