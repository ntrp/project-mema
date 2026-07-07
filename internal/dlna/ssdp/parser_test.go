package ssdp

import (
	"testing"
	"time"
)

func TestParseSearchAcceptsImperfectHeaders(t *testing.T) {
	request, ok := ParseSearch([]byte("M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 9\r\nST: ssdp:all"))
	if !ok {
		t.Fatal("expected search request")
	}
	if request.Target != "ssdp:all" {
		t.Fatalf("target = %q", request.Target)
	}
	if request.MX != 5*time.Second {
		t.Fatalf("MX = %s, want clamped 5s", request.MX)
	}
}

func TestParseSearchRejectsNonDiscover(t *testing.T) {
	if _, ok := ParseSearch([]byte("NOTIFY * HTTP/1.1\r\nST: ssdp:all\r\n")); ok {
		t.Fatal("expected non-search packet to be ignored")
	}
	if _, ok := ParseSearch([]byte("M-SEARCH * HTTP/1.1\r\nMAN: other\r\nST: ssdp:all\r\n")); ok {
		t.Fatal("expected non-discover search to be ignored")
	}
}
