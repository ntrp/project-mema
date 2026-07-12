package security

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"media-manager/internal/subtitles/catalog"
)

var (
	ErrOutboundURLBlocked = errors.New("subtitle provider outbound URL blocked")
	ErrDownloadHostClosed = errors.New("subtitle provider download hosts are not configured")
)

func ValidateProviderURL(providerKey string, rawURL string, download bool) error {
	entry, ok := catalog.Lookup(providerKey)
	if !ok {
		return fmt.Errorf("%w: unknown provider %q", ErrOutboundURLBlocked, providerKey)
	}
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("%w: invalid URL", ErrOutboundURLBlocked)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return fmt.Errorf("%w: unsupported scheme %q", ErrOutboundURLBlocked, parsed.Scheme)
	}
	host := strings.ToLower(parsed.Hostname())
	if !entry.OutboundPolicy.AllowLocalHosts && isLocalHost(host) {
		return fmt.Errorf("%w: local host %q is not allowed", ErrOutboundURLBlocked, host)
	}
	allowed := entry.OutboundPolicy.AllowedBaseHosts
	if download {
		allowed = entry.OutboundPolicy.AllowedDownloadHosts
		if len(allowed) == 0 {
			return ErrDownloadHostClosed
		}
	}
	if len(allowed) == 0 {
		return nil
	}
	if !hostAllowed(host, allowed) {
		return fmt.Errorf("%w: host %q is outside provider allowlist", ErrOutboundURLBlocked, host)
	}
	return nil
}

func ValidateRedirect(providerKey string, fromURL string, toURL string, download bool) error {
	if err := ValidateProviderURL(providerKey, fromURL, download); err != nil {
		return err
	}
	return ValidateProviderURL(providerKey, toURL, download)
}

func hostAllowed(host string, allowed []string) bool {
	for _, item := range allowed {
		pattern := strings.ToLower(strings.TrimSpace(item))
		if pattern == "" {
			continue
		}
		if strings.HasPrefix(pattern, "*.") {
			suffix := strings.TrimPrefix(pattern, "*")
			if strings.HasSuffix(host, suffix) && host != strings.TrimPrefix(suffix, ".") {
				return true
			}
			continue
		}
		if host == pattern {
			return true
		}
	}
	return false
}

func isLocalHost(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified()
}
