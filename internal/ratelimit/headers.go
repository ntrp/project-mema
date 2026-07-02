package ratelimit

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

func DelayFromHeaders(headers http.Header) time.Duration {
	for _, name := range []string{"Retry-After", "X-RateLimit-Reset-After", "RateLimit-Reset"} {
		if delay := parseDelay(headers.Get(name)); delay > 0 {
			return delay
		}
	}
	if delay := parseResetAt(headers.Get("X-RateLimit-Reset")); delay > 0 {
		return delay
	}
	return 0
}

func parseDelay(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.ParseFloat(value, 64); err == nil {
		return time.Duration(seconds * float64(time.Second))
	}
	if when, err := http.ParseTime(value); err == nil {
		return time.Until(when)
	}
	return 0
}

func parseResetAt(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Until(time.Unix(seconds, 0))
	}
	if when, err := http.ParseTime(value); err == nil {
		return time.Until(when)
	}
	return 0
}
