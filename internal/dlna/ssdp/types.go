package ssdp

import (
	"net"
	"time"
)

const (
	MulticastIPv4 = "239.255.255.250:1900"
	MediaServer   = "urn:schemas-upnp-org:device:MediaServer:1"
	ContentDir    = "urn:schemas-upnp-org:service:ContentDirectory:1"
	Connection    = "urn:schemas-upnp-org:service:ConnectionManager:1"
)

var ServiceTargets = []string{
	"upnp:rootdevice",
	MediaServer,
	ContentDir,
	Connection,
}

type Config struct {
	FriendlyName    string
	HTTPPort        string
	Interfaces      []string
	AnnounceSeconds int32
	UUID            string
	ServerHeader    string
}

type Interface struct {
	Name     string
	Hardware net.Interface
	Addr     net.IP
	Location string
}

type Runtime struct {
	config Config
	ifaces []Interface
	stop   chan struct{}
	done   chan struct{}
}

type Packet struct {
	Target   string
	USN      string
	Location string
	Headers  map[string]string
}

type SearchRequest struct {
	Target string
	MX     time.Duration
}
