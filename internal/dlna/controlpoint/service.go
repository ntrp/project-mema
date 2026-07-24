package controlpoint

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

const rendererSearch = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: " + MediaRendererDevice + "\r\n\r\n"

type expectedSourceKey struct{}

type Interface struct {
	Name     string
	Hardware net.Interface
	Address  net.IP
}

type Service struct {
	registry *Registry
	client   *http.Client
	networks []*net.IPNet
	ifaces   []Interface
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewService(ifaces []Interface, networks []*net.IPNet, client *http.Client) *Service {
	if client == nil {
		client = &http.Client{Timeout: 4 * time.Second}
	} else {
		clone := *client
		client = &clone
	}
	previousRedirect := client.CheckRedirect
	client.CheckRedirect = func(request *http.Request, via []*http.Request) error {
		if source, ok := request.Context().Value(expectedSourceKey{}).(net.IP); ok {
			if err := validateEndpointSource(request.URL.String(), source, networks); err != nil {
				return err
			}
		} else if err := ValidateEndpoint(request.URL.String(), networks); err != nil {
			return err
		}
		if previousRedirect != nil {
			return previousRedirect(request, via)
		}
		return nil
	}
	return &Service{
		registry: NewRegistry(), client: client,
		networks: append([]*net.IPNet{}, networks...), ifaces: append([]Interface{}, ifaces...),
	}
}
func (s *Service) Registry() *Registry { return s.registry }
func (s *Service) Client() Client      { return NewSecureClient(s.client, s.networks) }
func (s *Service) Start(ctx context.Context) {
	if s == nil || s.cancel != nil {
		return
	}
	ctx, s.cancel = context.WithCancel(ctx)
	for _, iface := range s.ifaces {
		s.wg.Add(2)
		go s.searchLoop(ctx, iface)
		go s.notifyLoop(ctx, iface)
	}
}
func (s *Service) Stop() {
	if s == nil || s.cancel == nil {
		return
	}
	s.cancel()
	s.wg.Wait()
	s.cancel = nil
}
func (s *Service) searchLoop(ctx context.Context, iface Interface) {
	defer s.wg.Done()
	target, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: iface.Address, Port: 0})
	if err != nil {
		return
	}
	defer conn.Close()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	buffer := make([]byte, 4096)
	for {
		_, _ = conn.WriteToUDP([]byte(rendererSearch), target)
		deadline := time.Now().Add(3 * time.Second)
		for time.Now().Before(deadline) {
			_ = conn.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
			n, source, readErr := conn.ReadFromUDP(buffer)
			if readErr == nil {
				s.process(ctx, append([]byte{}, buffer[:n]...), source.IP, iface.Address)
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
func (s *Service) notifyLoop(ctx context.Context, iface Interface) {
	defer s.wg.Done()
	group, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return
	}
	conn, err := net.ListenMulticastUDP("udp4", &iface.Hardware, group)
	if err != nil {
		return
	}
	defer conn.Close()
	buffer := make([]byte, 4096)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, source, readErr := conn.ReadFromUDP(buffer)
		if readErr == nil && sourceOnInterface(source.IP, iface) {
			s.process(ctx, append([]byte{}, buffer[:n]...), source.IP, iface.Address)
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
func (s *Service) process(ctx context.Context, payload []byte, source, local net.IP) {
	advertisement, ok := parseAdvertisement(payload)
	if !ok || !ipAllowed(source, s.networks) {
		return
	}
	if advertisement.Byebye {
		s.registry.RemoveFrom(advertisement.UDN, source)
		return
	}
	if err := validateEndpointSource(advertisement.Location, source, s.networks); err != nil {
		return
	}
	if existing, found := s.registry.Get(advertisement.UDN, time.Time{}); found &&
		existing.SourceIP.Equal(source) && existing.DescriptionURL == advertisement.Location {
		existing.LocalIP = append(net.IP(nil), local...)
		s.registry.Upsert(existing, source, advertisement.MaxAge)
		return
	}
	go s.fetch(ctx, advertisement, source, local)
}
func (s *Service) fetch(ctx context.Context, ad advertisement, source, local net.IP) {
	payload, err := s.fetchURL(ctx, ad.Location, source)
	if err != nil {
		return
	}
	renderer, err := ParseRendererDescription(ad.Location, payload)
	if err != nil || renderer.UDN != ad.UDN {
		return
	}
	renderer.DescriptionURL = ad.Location
	renderer.LocalIP = append(net.IP(nil), local...)
	if err := validateRenderer(renderer, source, s.networks); err != nil {
		return
	}
	renderer.AVTransportActions, renderer.AVTransportActionsKnown = s.fetchActions(ctx, renderer.AVTransportSCPDURL, source)
	renderer.RenderingActions, renderer.RenderingActionsKnown = s.fetchActions(ctx, renderer.RenderingSCPDURL, source)
	s.registry.Upsert(renderer, source, ad.MaxAge)
}

func (s *Service) fetchActions(ctx context.Context, endpoint string, source net.IP) ([]string, bool) {
	if endpoint == "" {
		return nil, false
	}
	document, err := s.fetchURL(ctx, endpoint, source)
	if err != nil {
		return nil, false
	}
	actions, err := ParseServiceActions(document)
	return actions, err == nil
}

func (s *Service) fetchURL(ctx context.Context, endpoint string, expected ...net.IP) ([]byte, error) {
	if len(expected) > 0 {
		if err := validateEndpointSource(endpoint, expected[0], s.networks); err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, expectedSourceKey{}, append(net.IP(nil), expected[0]...))
	} else if err := ValidateEndpoint(endpoint, s.networks); err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("renderer description returned status %d", response.StatusCode)
	}
	return io.ReadAll(io.LimitReader(response.Body, 256<<10))
}

func sourceOnInterface(source net.IP, iface Interface) bool {
	addresses, err := iface.Hardware.Addrs()
	if err != nil {
		return false
	}
	for _, address := range addresses {
		network, ok := address.(*net.IPNet)
		if ok && network.Contains(iface.Address) && network.Contains(source) {
			return true
		}
	}
	return false
}

func ipAllowed(ip net.IP, networks []*net.IPNet) bool {
	for _, network := range networks {
		if network != nil && network.Contains(ip) {
			return true
		}
	}
	return false
}
