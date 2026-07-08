package dlna

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type EventManager struct {
	mu     sync.Mutex
	update atomic.Uint32
	subs   map[string]eventSub
	client *http.Client
}

type eventSub struct {
	Callback string
	Expires  time.Time
	Seq      uint32
}

func NewEventManager() *EventManager {
	return &EventManager{subs: map[string]eventSub{}, client: &http.Client{Timeout: 2 * time.Second}}
}

func (e *EventManager) UpdateID() int {
	return int(e.update.Load())
}

func (e *EventManager) SubscriptionCount() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.purgeExpiredLocked(time.Now())
	return len(e.subs)
}

func (e *EventManager) NotifyContentChanged() {
	e.update.Add(1)
	e.notifyAll()
}

func (e *EventManager) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "SUBSCRIBE":
		e.subscribe(w, r)
	case "UNSUBSCRIBE":
		e.unsubscribe(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (e *EventManager) subscribe(w http.ResponseWriter, r *http.Request) {
	if sid := strings.TrimSpace(r.Header.Get("SID")); sid != "" {
		e.renew(w, r, sid)
		return
	}
	callback := strings.Trim(r.Header.Get("CALLBACK"), "<>")
	if callback == "" {
		http.Error(w, "missing callback", http.StatusBadRequest)
		return
	}
	sid := "uuid:" + uuid.NewString()
	timeout := requestedTimeout(r.Header.Get("TIMEOUT"))
	e.mu.Lock()
	e.subs[sid] = eventSub{Callback: callback, Expires: time.Now().Add(timeout)}
	e.mu.Unlock()
	w.Header().Set("SID", sid)
	w.Header().Set("TIMEOUT", "Second-"+fmt.Sprint(int(timeout.Seconds())))
	w.WriteHeader(http.StatusOK)
	e.notifySID(sid)
}

func (e *EventManager) renew(w http.ResponseWriter, r *http.Request, sid string) {
	timeout := requestedTimeout(r.Header.Get("TIMEOUT"))
	e.mu.Lock()
	sub, ok := e.subs[sid]
	if ok {
		sub.Expires = time.Now().Add(timeout)
		e.subs[sid] = sub
	}
	e.mu.Unlock()
	if !ok {
		http.Error(w, "unknown subscription", http.StatusPreconditionFailed)
		return
	}
	w.Header().Set("SID", sid)
	w.Header().Set("TIMEOUT", "Second-"+fmt.Sprint(int(timeout.Seconds())))
	w.WriteHeader(http.StatusOK)
}

func (e *EventManager) unsubscribe(w http.ResponseWriter, r *http.Request) {
	sid := r.Header.Get("SID")
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.subs[sid]; !ok {
		http.Error(w, "unknown subscription", http.StatusPreconditionFailed)
		return
	}
	delete(e.subs, sid)
	w.WriteHeader(http.StatusOK)
}

func (e *EventManager) notifyAll() {
	e.mu.Lock()
	e.purgeExpiredLocked(time.Now())
	sids := make([]string, 0, len(e.subs))
	for sid := range e.subs {
		sids = append(sids, sid)
	}
	e.mu.Unlock()
	for _, sid := range sids {
		e.notifySID(sid)
	}
}

func (e *EventManager) notifySID(sid string) {
	e.mu.Lock()
	sub, ok := e.subs[sid]
	if !ok || time.Now().After(sub.Expires) {
		delete(e.subs, sid)
		e.mu.Unlock()
		return
	}
	seq := sub.Seq
	sub.Seq++
	e.subs[sid] = sub
	e.mu.Unlock()
	payload := []byte(`<?xml version="1.0"?><e:propertyset xmlns:e="urn:schemas-upnp-org:event-1-0"><e:property><SystemUpdateID>` + fmt.Sprint(e.UpdateID()) + `</SystemUpdateID></e:property></e:propertyset>`)
	request, err := http.NewRequest("NOTIFY", sub.Callback, bytes.NewReader(payload))
	if err != nil {
		return
	}
	request.Header.Set("NT", "upnp:event")
	request.Header.Set("NTS", "upnp:propchange")
	request.Header.Set("SID", sid)
	request.Header.Set("SEQ", fmt.Sprint(seq))
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	_, _ = e.client.Do(request)
}

func (e *EventManager) purgeExpiredLocked(now time.Time) {
	for sid, sub := range e.subs {
		if now.After(sub.Expires) {
			delete(e.subs, sid)
		}
	}
}

func requestedTimeout(value string) time.Duration {
	value = strings.TrimPrefix(strings.TrimSpace(value), "Second-")
	if value == "" || strings.EqualFold(value, "infinite") {
		return 30 * time.Minute
	}
	var seconds int
	if _, err := fmt.Sscanf(value, "%d", &seconds); err != nil || seconds <= 0 {
		return 30 * time.Minute
	}
	if seconds > 3600 {
		seconds = 3600
	}
	return time.Duration(seconds) * time.Second
}
