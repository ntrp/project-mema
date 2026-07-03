package indexers

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestSCNIntegrations001StatusErrorExposesRetryAndFailureClassification(t *testing.T) {
	retryAfter := 3 * time.Second
	err := StatusError{StatusCode: http.StatusTooManyRequests, RetryAfter: retryAfter}

	code := StatusCode(err)
	if code == nil || *code != http.StatusTooManyRequests {
		t.Fatalf("status code = %#v, want 429", code)
	}
	if RetryAfter(err) != retryAfter {
		t.Fatalf("retry after = %s, want %s", RetryAfter(err), retryAfter)
	}
	if IsPermanentFailure(code) {
		t.Fatal("429 should be retryable")
	}

	notFound := int32(http.StatusNotFound)
	if !IsPermanentFailure(&notFound) {
		t.Fatal("404 should be a permanent failure")
	}
	if IsPermanentFailure(nil) {
		t.Fatal("nil status should not be a permanent failure")
	}
}

func TestSCNIntegrations001StatusCodeFromDetailsNormalizesNumericTypes(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  int32
	}{
		{name: "int", value: int(500), want: 500},
		{name: "int32", value: int32(502), want: 502},
		{name: "int64", value: int64(503), want: 503},
		{name: "float64", value: float64(504), want: 504},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code := StatusCodeFromDetails(map[string]interface{}{"statusCode": tc.value})
			if code == nil || *code != tc.want {
				t.Fatalf("status code = %#v, want %d", code, tc.want)
			}
		})
	}

	if code := StatusCodeFromDetails(map[string]interface{}{"statusCode": "500"}); code != nil {
		t.Fatalf("string status code = %#v, want nil", code)
	}
	if code := StatusCode(errors.New("plain error")); code != nil {
		t.Fatalf("plain error status code = %#v, want nil", code)
	}
}
