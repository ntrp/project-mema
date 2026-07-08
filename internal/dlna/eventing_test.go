package dlna

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEventSubscriptionLifecycleSendsInitialNotify(t *testing.T) {
	notify := make(chan *http.Request, 1)
	events := NewEventManager()
	events.client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		notify <- r
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
	})
	request := httptest.NewRequest("SUBSCRIBE", "/dlna/events/content-directory", nil)
	request.Header.Set("CALLBACK", "<http://renderer.local/events>")
	request.Header.Set("TIMEOUT", "Second-60")
	response := httptest.NewRecorder()

	events.Handle(response, request)

	if response.Code != http.StatusOK || response.Header().Get("SID") == "" || response.Header().Get("TIMEOUT") != "Second-60" {
		t.Fatalf("subscribe = %d sid=%q timeout=%q", response.Code, response.Header().Get("SID"), response.Header().Get("TIMEOUT"))
	}
	var got *http.Request
	select {
	case got = <-notify:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for initial notify")
	}
	if got.Method != "NOTIFY" || got.Header.Get("NTS") != "upnp:propchange" || got.Header.Get("SEQ") != "0" {
		t.Fatalf("notify headers = %#v", got.Header)
	}
	if events.SubscriptionCount() != 1 {
		t.Fatalf("subscriptions = %d", events.SubscriptionCount())
	}

	renew := httptest.NewRequest("SUBSCRIBE", "/dlna/events/content-directory", nil)
	renew.Header.Set("SID", response.Header().Get("SID"))
	renew.Header.Set("TIMEOUT", "Second-120")
	renewed := httptest.NewRecorder()
	events.Handle(renewed, renew)
	if renewed.Code != http.StatusOK || renewed.Header().Get("TIMEOUT") != "Second-120" {
		t.Fatalf("renew = %d timeout=%q", renewed.Code, renewed.Header().Get("TIMEOUT"))
	}

	events.NotifyContentChanged()
	select {
	case got = <-notify:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for changed notify")
	}
	if got.Header.Get("SEQ") != "1" {
		t.Fatalf("renewed notify headers = %#v", got.Header)
	}

	unsubscribe := httptest.NewRequest("UNSUBSCRIBE", "/dlna/events/content-directory", nil)
	unsubscribe.Header.Set("SID", response.Header().Get("SID"))
	done := httptest.NewRecorder()
	events.Handle(done, unsubscribe)
	if done.Code != http.StatusOK {
		t.Fatalf("unsubscribe = %d", done.Code)
	}
	if events.SubscriptionCount() != 0 {
		t.Fatalf("subscriptions after unsubscribe = %d", events.SubscriptionCount())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestEventNotifyContentChangedIncrementsUpdateID(t *testing.T) {
	events := NewEventManager()
	if events.UpdateID() != 0 {
		t.Fatalf("initial update id = %d", events.UpdateID())
	}
	events.NotifyContentChanged()
	if events.UpdateID() != 1 {
		t.Fatalf("updated id = %d", events.UpdateID())
	}
}

func TestEventRouteRejectsMissingCallback(t *testing.T) {
	events := NewEventManager()
	response := httptest.NewRecorder()

	events.Handle(response, httptest.NewRequest("SUBSCRIBE", "/dlna/events/content-directory", nil))

	if response.Code != http.StatusBadRequest || !strings.Contains(response.Body.String(), "missing callback") {
		t.Fatalf("response = %d %s", response.Code, response.Body.String())
	}
}
