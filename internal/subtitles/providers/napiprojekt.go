package providers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"media-manager/internal/subtitles/providercore"
)

const napiprojektKey = "napiprojekt"

func init() { Register(napiprojektKey, napiprojektAdapter{}) }

type napiprojektAdapter struct{}

func (napiprojektAdapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	endpoint := napiprojektEndpoint(config)
	values := url.Values{"v": {"dreambox"}, "kolejka": {"false"}, "nick": {""}, "pass": {""}, "napios": {"Linux"}, "l": {"PL"}, "f": {"00000000000000000000000000000000"}, "t": {napiprojektSubhash("00000000000000000000000000000000")}}
	_, _, err := providerRequest(ctx, service, http.MethodGet, endpoint+"?"+values.Encode(), napiprojektKey, false, nil)
	return err
}

func (napiprojektAdapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if alpha3Language(request.LanguageID) != "pol" {
		return nil, nil
	}
	hash, err := requestNapiprojektHash(request)
	if err != nil {
		return nil, err
	}
	content, err := napiprojektDownloadByHash(ctx, service, config, hash)
	if err != nil {
		return nil, err
	}
	if len(content) >= 4 && string(content[:4]) == "NPc0" {
		return nil, nil
	}
	return []providercore.Candidate{{ProviderName: config.Name, LanguageID: request.LanguageID, FileID: 0, Format: "srt", ReleaseName: hash, SourceURL: napiprojektEndpoint(config), SourceRef: hash}}, nil
}

func (napiprojektAdapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	hash := strings.TrimSpace(candidate.SourceRef)
	if hash == "" {
		hash = strings.TrimSpace(candidate.ReleaseName)
	}
	if hash == "" {
		return providercore.Download{}, fmt.Errorf("%w: napiprojekt candidate has no hash", providercore.ErrProviderPrerequisiteMissing)
	}
	content, err := napiprojektDownloadByHash(ctx, service, config, hash)
	if err != nil {
		return providercore.Download{}, err
	}
	if len(content) >= 4 && string(content[:4]) == "NPc0" {
		return providercore.Download{}, fmt.Errorf("%w: napiprojekt subtitle no longer available", providercore.ErrProviderBrokenUpstream)
	}
	return providercore.Download{Content: content, URL: napiprojektEndpoint(config)}, nil
}

func requestNapiprojektHash(request providercore.SearchRequest) (string, error) {
	if value := strings.TrimSpace(request.MediaContext.File.Hashes[napiprojektKey]); value != "" {
		return value, nil
	}
	if value := strings.TrimSpace(request.MediaContext.File.Hashes["napi"]); value != "" {
		return value, nil
	}
	path := request.FilePath
	if path == "" {
		path = request.MediaContext.File.Path
	}
	if path == "" {
		return "", fmt.Errorf("%w: napiprojekt requires a file path or napiprojekt hash", providercore.ErrProviderPrerequisiteMissing)
	}
	return napiprojektHash(path)
}

func napiprojektDownloadByHash(ctx context.Context, service providercore.Service, config providercore.Config, hash string) ([]byte, error) {
	values := url.Values{"v": {"dreambox"}, "kolejka": {"false"}, "nick": {""}, "pass": {""}, "napios": {"Linux"}, "l": {"PL"}, "f": {hash}, "t": {napiprojektSubhash(hash)}}
	data, _, err := providerRequest(ctx, service, http.MethodGet, napiprojektEndpoint(config)+"?"+values.Encode(), napiprojektKey, true, nil)
	return data, err
}

func napiprojektEndpoint(config providercore.Config) string {
	base := strings.TrimRight(strings.TrimSpace(config.BaseURL), "/")
	if base == "" {
		base = "https://napiprojekt.pl"
	}
	if strings.HasSuffix(base, "/unit_napisy/dl.php") {
		return base
	}
	return base + "/unit_napisy/dl.php"
}

func napiprojektSubhash(hash string) string {
	idx := []int{0xe, 0x3, 0x6, 0x8, 0x2}
	mul := []int{2, 2, 5, 4, 3}
	add := []int{0, 0xd, 0x10, 0xb, 0x5}
	var out strings.Builder
	for i := range idx {
		if len(hash) <= idx[i] {
			return ""
		}
		t := add[i] + fromHex(hash[idx[i]])
		if len(hash) < t+2 {
			return ""
		}
		value := fromHex(hash[t])*16 + fromHex(hash[t+1])
		out.WriteString(fmt.Sprintf("%x", value*mul[i]%16))
	}
	return out.String()
}

func fromHex(ch byte) int {
	switch {
	case ch >= '0' && ch <= '9':
		return int(ch - '0')
	case ch >= 'a' && ch <= 'f':
		return int(ch-'a') + 10
	case ch >= 'A' && ch <= 'F':
		return int(ch-'A') + 10
	default:
		return 0
	}
}
