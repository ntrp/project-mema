package downloadrouting

import (
	"testing"

	"media-manager/internal/storage"
)

func TestClientProtocolRoutingMatchesReleaseProtocol(t *testing.T) {
	clients := []storage.DownloadClient{
		{Name: "SAB", Type: "sabnzbd", Protocol: "usenet"},
		{Name: "Transmission", Type: "transmission", Protocol: "torrent"},
	}

	client, ok := ClientForProtocol(clients, "torrent")
	if !ok || client.Name != "Transmission" {
		t.Fatalf("torrent client = %#v, ok = %v", client, ok)
	}
	client, ok = NamedClientForProtocol(clients, "SAB", "usenet")
	if !ok || client.Name != "SAB" {
		t.Fatalf("named usenet client = %#v, ok = %v", client, ok)
	}
	if TypeSupportsProtocol("sabnzbd", "torrent") {
		t.Fatal("sabnzbd must not support torrent releases")
	}
}

func TestReleaseInputsForClientsFiltersUnsupportedProtocols(t *testing.T) {
	releases := []storage.ReleaseCandidateInput{
		{Title: "Torrent", IndexerProtocol: "torrent"},
		{Title: "Usenet", IndexerProtocol: "usenet"},
	}
	clients := []storage.DownloadClient{{Name: "SAB", Type: "sabnzbd", Protocol: "usenet"}}

	filtered := ReleaseInputsForClients(releases, clients)

	if len(filtered) != 1 || filtered[0].Title != "Usenet" {
		t.Fatalf("filtered releases = %#v", filtered)
	}
}
