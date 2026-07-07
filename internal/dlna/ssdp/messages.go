package ssdp

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

func AlivePackets(config Config, iface Interface) []Packet {
	return notifyPackets(config, iface, "ssdp:alive")
}

func ByebyePackets(config Config, iface Interface) []Packet {
	return notifyPackets(config, iface, "ssdp:byebye")
}

func SearchResponse(config Config, iface Interface, target string) (Packet, bool) {
	if !SupportsTarget(config.UUID, target) {
		return Packet{}, false
	}
	return Packet{
		Target:   normalizedTarget(target),
		USN:      usn(config.UUID, normalizedTarget(target)),
		Location: iface.Location,
		Headers: map[string]string{
			"CACHE-CONTROL": "max-age=" + announceSeconds(config),
			"EXT":           "",
			"LOCATION":      iface.Location,
			"SERVER":        serverHeader(config),
			"ST":            normalizedTarget(target),
			"USN":           usn(config.UUID, normalizedTarget(target)),
		},
	}, true
}

func PacketText(packet Packet, statusLine string) string {
	var builder strings.Builder
	builder.WriteString(statusLine)
	builder.WriteString("\r\n")
	keys := make([]string, 0, len(packet.Headers))
	for key := range packet.Headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		builder.WriteString(key)
		if packet.Headers[key] == "" {
			builder.WriteString(":\r\n")
			continue
		}
		builder.WriteString(": ")
		builder.WriteString(packet.Headers[key])
		builder.WriteString("\r\n")
	}
	builder.WriteString("\r\n")
	return builder.String()
}

func SupportsTarget(uuid string, target string) bool {
	target = normalizedTarget(target)
	if target == "ssdp:all" || target == "upnp:rootdevice" || target == uuidURN(uuid) {
		return true
	}
	for _, value := range ServiceTargets[1:] {
		if target == value {
			return true
		}
	}
	return false
}

func notifyPackets(config Config, iface Interface, nts string) []Packet {
	packets := make([]Packet, 0, len(ServiceTargets)+1)
	for _, target := range append([]string{uuidURN(config.UUID)}, ServiceTargets...) {
		headers := map[string]string{
			"CACHE-CONTROL": "max-age=" + announceSeconds(config),
			"HOST":          MulticastIPv4,
			"LOCATION":      iface.Location,
			"NT":            target,
			"NTS":           nts,
			"SERVER":        serverHeader(config),
			"USN":           usn(config.UUID, target),
		}
		packets = append(packets, Packet{Target: target, USN: headers["USN"], Location: iface.Location, Headers: headers})
	}
	return packets
}

func DiscoverInterfaces(config Config) ([]Interface, error) {
	system, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	allowed := map[string]bool{}
	for _, name := range config.Interfaces {
		allowed[strings.TrimSpace(name)] = true
	}
	results := []Interface{}
	for _, item := range system {
		if len(allowed) > 0 && !allowed[item.Name] {
			continue
		}
		if item.Flags&net.FlagUp == 0 || item.Flags&net.FlagMulticast == 0 || item.Flags&net.FlagLoopback != 0 {
			continue
		}
		ip := firstUsableIPv4(item)
		if ip == nil {
			continue
		}
		results = append(results, Interface{
			Name:     item.Name,
			Hardware: item,
			Addr:     ip,
			Location: "http://" + net.JoinHostPort(ip.String(), config.HTTPPort) + "/dlna/rootDesc.xml",
		})
	}
	return results, nil
}

func firstUsableIPv4(item net.Interface) net.IP {
	addrs, err := item.Addrs()
	if err != nil {
		return nil
	}
	for _, addr := range addrs {
		var ip net.IP
		switch value := addr.(type) {
		case *net.IPNet:
			ip = value.IP
		case *net.IPAddr:
			ip = value.IP
		}
		ip = ip.To4()
		if ip != nil && !ip.IsLoopback() && !ip.IsLinkLocalUnicast() {
			return ip
		}
	}
	return nil
}

func normalizedTarget(target string) string {
	target = strings.TrimSpace(target)
	if strings.HasPrefix(target, "uuid:") {
		return target
	}
	return target
}

func usn(uuid string, target string) string {
	root := uuidURN(uuid)
	if target == root {
		return root
	}
	return root + "::" + target
}

func uuidURN(uuid string) string {
	if strings.HasPrefix(uuid, "uuid:") {
		return uuid
	}
	return "uuid:" + uuid
}

func serverHeader(config Config) string {
	if strings.TrimSpace(config.ServerHeader) != "" {
		return config.ServerHeader
	}
	return "Mema/0.0 UPnP/1.0 MemaDLNA/0.1"
}

func announceSeconds(config Config) string {
	if config.AnnounceSeconds <= 0 {
		return "1800"
	}
	return fmt.Sprintf("%d", config.AnnounceSeconds)
}
