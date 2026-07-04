package downloadclients

import "strings"

const (
	ProtocolTorrent = "torrent"
	ProtocolUsenet  = "usenet"
)

func ProtocolForType(clientType string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(clientType)) {
	case "transmission":
		return ProtocolTorrent, true
	case "sabnzbd":
		return ProtocolUsenet, true
	default:
		return "", false
	}
}

func TypeSupportsProtocol(clientType string, protocol string) bool {
	expected, ok := ProtocolForType(clientType)
	return ok && expected == strings.ToLower(strings.TrimSpace(protocol))
}
