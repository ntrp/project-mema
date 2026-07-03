package downloadclients

import (
	"context"
	"testing"
)

func TestSCNIntegrations003DownloadClientValidationFailuresAreUserFacing(t *testing.T) {
	service := NewService(nil)
	config := Config{Type: "unsupported"}

	added := service.Add(context.Background(), config, AddRequest{})
	if added.Success || added.Message != "Download URL is required" {
		t.Fatalf("add validation result = %#v", added)
	}

	status := service.Status(context.Background(), config, StatusRequest{})
	if status.Success || status.Message != "Download ID is required" {
		t.Fatalf("status validation result = %#v", status)
	}

	cancelled := service.Cancel(context.Background(), config, CancelRequest{})
	if cancelled.Success || cancelled.Message != "Download ID is required" {
		t.Fatalf("cancel validation result = %#v", cancelled)
	}
}

func TestSCNIntegrations003UnsupportedDownloadClientReportsRequestedType(t *testing.T) {
	service := NewService(nil)
	config := Config{Type: "custom-client"}

	added := service.Add(context.Background(), config, AddRequest{URL: "https://indexer.test/release"})
	if added.Success || added.Message != "Unsupported download client type" {
		t.Fatalf("add result = %#v", added)
	}
	if added.Details["type"] != "custom-client" {
		t.Fatalf("add details = %#v", added.Details)
	}

	status := service.Status(context.Background(), config, StatusRequest{DownloadID: "download-1"})
	if status.Success || status.Message != "Unsupported download client type" {
		t.Fatalf("status result = %#v", status)
	}
	if status.Details["type"] != "custom-client" {
		t.Fatalf("status details = %#v", status.Details)
	}

	cancelled := service.Cancel(context.Background(), config, CancelRequest{DownloadID: "download-1"})
	if cancelled.Success || cancelled.Message != "Unsupported download client type" {
		t.Fatalf("cancel result = %#v", cancelled)
	}
	if cancelled.Details["type"] != "custom-client" {
		t.Fatalf("cancel details = %#v", cancelled.Details)
	}
}

func TestSCNIntegrations003FailedTestResultAlwaysHasDetails(t *testing.T) {
	result := NewService(nil).Test(context.Background(), Config{Type: "custom-client"})
	if result.Success || result.Message != "Unsupported download client type" {
		t.Fatalf("test result = %#v", result)
	}
	if result.Details["type"] != "custom-client" {
		t.Fatalf("test details = %#v", result.Details)
	}
	if result.Latency <= 0 {
		t.Fatalf("latency = %s, want recorded duration", result.Latency)
	}
}

func TestSCNIntegrations003ProviderFailureResultsDescribeUserVisibleCause(t *testing.T) {
	status := statusFailedResult(503)
	if status.Success || status.Message != "Unexpected response status" || status.Details["statusCode"] != 503 {
		t.Fatalf("status failure result = %#v", status)
	}

	rejected := formatResultFailure("SABnzbd", "")
	if rejected.Success || rejected.Message != "SABnzbd rejected the test request" {
		t.Fatalf("format failure result = %#v", rejected)
	}
	if rejected.Details["result"] != "empty" {
		t.Fatalf("format failure details = %#v", rejected.Details)
	}

	values := details("ok", true, "", "ignored", 42, "ignored")
	if values["ok"] != true || len(values) != 1 {
		t.Fatalf("details = %#v", values)
	}
}
