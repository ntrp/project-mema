package events

import (
	"context"
	"testing"
	"time"
)

func TestScenarioSCNSystem002SubscriberReceivesPublishedEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	broker := NewBroker()
	events := broker.Subscribe(ctx)

	broker.Publish("settings.updated", map[string]string{"id": "one"})

	select {
	case event := <-events:
		if event.Type != "settings.updated" {
			t.Fatalf("event type = %q", event.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber did not receive event")
	}

	cancel()
	select {
	case _, ok := <-events:
		if ok {
			t.Fatal("subscriber channel should close after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("subscriber channel did not close")
	}
}

func TestScenarioSCNSystem002PublishIgnoresNilOrEmptyEvents(t *testing.T) {
	var broker *Broker
	broker.Publish("ignored", nil)

	active := NewBroker()
	events := active.Subscribe(context.Background())
	active.Publish("", "ignored")
	select {
	case event := <-events:
		t.Fatalf("unexpected event = %#v", event)
	default:
	}
}
