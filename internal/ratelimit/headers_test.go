package ratelimit

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestScenarioSCNSystem003DelayFromHeadersUsesFirstValidDelay(t *testing.T) {
	headers := http.Header{}
	headers.Set("Retry-After", "1.5")
	headers.Set("X-RateLimit-Reset-After", "10")

	if got := DelayFromHeaders(headers); got != 1500*time.Millisecond {
		t.Fatalf("delay = %s", got)
	}
}

func TestScenarioSCNSystem003DelayFromHeadersParsesDatesAndResetEpochs(t *testing.T) {
	future := time.Now().Add(2 * time.Second).UTC().Format(http.TimeFormat)
	headers := http.Header{"Retry-After": []string{future}}
	if got := DelayFromHeaders(headers); got <= 0 || got > 3*time.Second {
		t.Fatalf("http date delay = %s", got)
	}

	resetHeaders := http.Header{}
	resetHeaders.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(5*time.Second).Unix(), 10))
	if got := DelayFromHeaders(resetHeaders); got <= 0 || got > 6*time.Second {
		t.Fatalf("epoch reset delay = %s", got)
	}
}

func TestScenarioSCNSystem003DelayFromHeadersIgnoresInvalidValues(t *testing.T) {
	headers := http.Header{"Retry-After": []string{"later"}}
	if got := DelayFromHeaders(headers); got != 0 {
		t.Fatalf("delay = %s", got)
	}
}
