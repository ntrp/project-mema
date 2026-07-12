package providers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"media-manager/internal/subtitles/providercore"
)

const bsplayerKey = "bsplayer"
const bsplayerAPI = "https://s1.api.bsplayer-subtitles.com/v1.php"

func init() { Register(bsplayerKey, bsplayerAdapter{}) }

type bsplayerAdapter struct{}

func (bsplayerAdapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	_, err := bsplayerSearch(ctx, service, config, providercore.SearchRequest{})
	if err != nil && strings.Contains(err.Error(), "hash/size") {
		return nil
	}
	return err
}

func (bsplayerAdapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	return bsplayerSearch(ctx, service, config, request)
}

func (bsplayerAdapter) Download(ctx context.Context, service providercore.Service, _ providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	link := strings.TrimSpace(candidate.SourceURL)
	if link == "" {
		return providercore.Download{}, fmt.Errorf("%w: bsplayer candidate has no download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	request.Header.Set("User-Agent", "Mozilla/4.0 (compatible; Synapse)")
	response, err := service.DoProviderRequest(request, bsplayerKey, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return providercore.Download{}, fmt.Errorf("%w: provider returned HTTP %d", providercore.ErrProviderBrokenUpstream, response.StatusCode)
	}
	reader, err := gzip.NewReader(io.LimitReader(response.Body, providerReadLimit+1))
	if err != nil {
		return providercore.Download{}, err
	}
	defer reader.Close()
	content, err := io.ReadAll(io.LimitReader(reader, providerReadLimit+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(content) > providerReadLimit {
		return providercore.Download{}, fmt.Errorf("provider response size limit exceeded")
	}
	return providercore.Download{Content: content, URL: link}, nil
}

func bsplayerSearch(ctx context.Context, service providercore.Service, config providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	hash := firstHash(sr.MediaContext.File.Hashes, "bsplayer", "moviehash", "opensubtitles")
	if hash == "" || sr.MediaContext.File.SizeBytes <= 0 {
		return nil, fmt.Errorf("%w: bsplayer search requires file hash/size", providercore.ErrProviderPrerequisiteMissing)
	}
	body := bsplayerEnvelope(hash, sr.MediaContext.File.SizeBytes, sr.LanguageID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, providercore.NewConfig(config).BaseURL(bsplayerAPI), bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "\"searchSubtitles\"")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; Synapse)")
	resp, err := service.DoProviderRequest(req, bsplayerKey, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return parseBSPlayer(data, sr.LanguageID), nil
}

func bsplayerEnvelope(hash string, size int64, lang string) string {
	return fmt.Sprintf(`<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="http://api.bsplayer-subtitles.com/v1.php"><SOAP-ENV:Body><ns1:searchSubtitles><hash>%s</hash><size>%d</size><language>%s</language></ns1:searchSubtitles></SOAP-ENV:Body></SOAP-ENV:Envelope>`, xmlEscape(hash), size, xmlEscape(lang))
}

type bsNode struct {
	XMLName xml.Name
	Text    string   `xml:",chardata"`
	Nodes   []bsNode `xml:",any"`
}

func parseBSPlayer(data []byte, fallbackLang string) []providercore.Candidate {
	var root bsNode
	if xml.Unmarshal(data, &root) != nil {
		return nil
	}
	items := findBSItems(root)
	out := []providercore.Candidate{}
	for _, item := range items {
		values := flattenBS(item)
		link := firstValue(values, "downloadlink", "downloadurl", "subdownloadlink")
		if link == "" {
			continue
		}
		lang := firstValue(values, "language", "sublang", "isolanguage")
		if lang == "" {
			lang = fallbackLang
		}
		out = append(out, providercore.Candidate{ProviderName: bsplayerKey, LanguageID: lang, Format: "srt", ReleaseName: firstValue(values, "subname", "filename", "movie"), SourceURL: link})
	}
	return out
}

func findBSItems(n bsNode) []bsNode {
	name := strings.ToLower(n.XMLName.Local)
	if name == "subtitle" || name == "item" || name == "sub" {
		return []bsNode{n}
	}
	var out []bsNode
	for _, child := range n.Nodes {
		out = append(out, findBSItems(child)...)
	}
	return out
}

func flattenBS(n bsNode) map[string]string {
	out := map[string]string{}
	var walk func(bsNode)
	walk = func(node bsNode) {
		if text := strings.TrimSpace(node.Text); text != "" {
			out[strings.ToLower(node.XMLName.Local)] = text
		}
		for _, child := range node.Nodes {
			walk(child)
		}
	}
	walk(n)
	return out
}

func firstHash(hashes map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(hashes[key]); value != "" {
			return value
		}
	}
	return ""
}

func firstValue(values map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(values[key]); value != "" {
			return value
		}
	}
	return ""
}

func xmlEscape(value string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(value))
	return b.String()
}
