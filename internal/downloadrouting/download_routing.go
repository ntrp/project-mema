package downloadrouting

import (
	"fmt"
	"strings"

	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
)

func TypeSupportsProtocol(clientType string, protocol string) bool {
	return downloadclients.TypeSupportsProtocol(clientType, protocol)
}

func ClientProtocol(client storage.DownloadClient) string {
	if protocol := strings.ToLower(strings.TrimSpace(client.Protocol)); protocol != "" {
		return protocol
	}
	protocol, _ := downloadclients.ProtocolForType(client.Type)
	return protocol
}

func Protocols(clients []storage.DownloadClient) map[string]struct{} {
	protocols := map[string]struct{}{}
	for _, client := range clients {
		if protocol := ClientProtocol(client); protocol != "" {
			protocols[protocol] = struct{}{}
		}
	}
	return protocols
}

func HasProtocol(protocols map[string]struct{}, protocol string) bool {
	_, ok := protocols[strings.ToLower(strings.TrimSpace(protocol))]
	return ok
}

func ClientForProtocol(clients []storage.DownloadClient, protocol string) (storage.DownloadClient, bool) {
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	for _, client := range clients {
		if ClientProtocol(client) == protocol {
			return client, true
		}
	}
	return storage.DownloadClient{}, false
}

func NamedClientForProtocol(
	clients []storage.DownloadClient,
	name string,
	protocol string,
) (storage.DownloadClient, bool) {
	name = strings.TrimSpace(name)
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	for _, client := range clients {
		if client.Name == name && (protocol == "" || ClientProtocol(client) == protocol) {
			return client, true
		}
	}
	return storage.DownloadClient{}, false
}

func ReleaseInputsForClients(
	releases []storage.ReleaseCandidateInput,
	clients []storage.DownloadClient,
) []storage.ReleaseCandidateInput {
	protocols := Protocols(clients)
	filtered := make([]storage.ReleaseCandidateInput, 0, len(releases))
	for _, release := range releases {
		if HasProtocol(protocols, release.IndexerProtocol) {
			filtered = append(filtered, release)
		}
	}
	return filtered
}

func MissingClientMessage(protocol string) string {
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	if protocol == "" {
		return "No enabled download client is configured"
	}
	return fmt.Sprintf("No enabled %s download client is configured", protocol)
}
