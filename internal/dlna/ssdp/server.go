package ssdp

import (
	"context"
	"net"
	"time"
)

type udpConn interface {
	ReadFromUDP([]byte) (int, *net.UDPAddr, error)
	WriteToUDP([]byte, *net.UDPAddr) (int, error)
	Close() error
}

func Start(ctx context.Context, config Config) (*Runtime, error) {
	ifaces, err := DiscoverInterfaces(config)
	if err != nil {
		return nil, err
	}
	runtime := &Runtime{config: config, ifaces: ifaces, stop: make(chan struct{}), done: make(chan struct{})}
	if len(ifaces) == 0 {
		close(runtime.done)
		return runtime, nil
	}
	go runtime.run(ctx)
	return runtime, nil
}

func (r *Runtime) Stop(ctx context.Context) error {
	select {
	case <-r.stop:
	default:
		close(r.stop)
	}
	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Runtime) Interfaces() []Interface {
	return append([]Interface{}, r.ifaces...)
}

func (r *Runtime) run(ctx context.Context) {
	defer close(r.done)
	done := make(chan struct{})
	for _, iface := range r.ifaces {
		go r.runInterface(ctx, iface, done)
	}
	ticker := time.NewTicker(time.Duration(r.config.AnnounceSeconds) * time.Second)
	defer ticker.Stop()
	r.sendAlive()
	for {
		select {
		case <-ctx.Done():
			close(done)
			r.sendByebye()
			return
		case <-r.stop:
			close(done)
			r.sendByebye()
			return
		case <-ticker.C:
			r.sendAlive()
		}
	}
}

func (r *Runtime) runInterface(ctx context.Context, iface Interface, done <-chan struct{}) {
	addr, _ := net.ResolveUDPAddr("udp4", MulticastIPv4)
	conn, err := net.ListenMulticastUDP("udp4", &iface.Hardware, addr)
	if err != nil {
		return
	}
	defer conn.Close()
	buffer := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		default:
		}
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, remote, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		request, ok := ParseSearch(buffer[:n])
		if !ok {
			continue
		}
		packet, ok := SearchResponse(r.config, iface, request.Target)
		if !ok {
			continue
		}
		_, _ = conn.WriteToUDP([]byte(PacketText(packet, "HTTP/1.1 200 OK")), remote)
	}
}

func (r *Runtime) sendAlive() {
	r.sendNotify("NOTIFY * HTTP/1.1", AlivePackets)
}

func (r *Runtime) sendByebye() {
	r.sendNotify("NOTIFY * HTTP/1.1", ByebyePackets)
}

func (r *Runtime) sendNotify(statusLine string, packets func(Config, Interface) []Packet) {
	target, err := net.ResolveUDPAddr("udp4", MulticastIPv4)
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp4", nil, target)
	if err != nil {
		return
	}
	defer conn.Close()
	for _, iface := range r.ifaces {
		for _, packet := range packets(r.config, iface) {
			_, _ = conn.Write([]byte(PacketText(packet, statusLine)))
		}
	}
}
