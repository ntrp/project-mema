package nativeutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const MaxBody = 50 << 20

type RequestSpec struct {
	Provider string
	BaseURL  string
	Method   string
	Path     string
	Form     url.Values
	Headers  map[string]string
	Download bool
}

func Do(ctx context.Context, svc providercore.Service, cfg providercore.Config, spec RequestSpec) ([]byte, *http.Response, error) {
	method := spec.Method
	if method == "" {
		method = http.MethodGet
	}
	rawURL := Absolute(cfg, spec.BaseURL, spec.Path)
	var body io.Reader
	if spec.Form != nil {
		if method == http.MethodGet {
			u, err := url.Parse(rawURL)
			if err != nil {
				return nil, nil, err
			}
			q := u.Query()
			for k, values := range spec.Form {
				for _, value := range values {
					q.Add(k, value)
				}
			}
			u.RawQuery = q.Encode()
			rawURL = u.String()
		} else {
			body = strings.NewReader(spec.Form.Encode())
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, nil, err
	}
	if spec.Form != nil && method != http.MethodGet {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for key, value := range spec.Headers {
		req.Header.Set(key, value)
	}
	resp, err := svc.DoProviderRequest(req, spec.Provider, spec.Download)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxBody+1))
	if err != nil {
		return nil, resp, err
	}
	if len(data) > MaxBody {
		return nil, resp, fmt.Errorf("%w: response too large", security.ErrUnsafeArchive)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp, fmt.Errorf("%w: HTTP %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return data, resp, nil
}

func Test(ctx context.Context, svc providercore.Service, cfg providercore.Config, provider, baseURL string) error {
	_, _, err := Do(ctx, svc, cfg, RequestSpec{Provider: provider, BaseURL: baseURL, Path: "/"})
	return err
}

func Absolute(cfg providercore.Config, baseURL, raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	base := strings.TrimRight(providercore.NewConfig(cfg).BaseURL(baseURL), "/") + "/"
	parsedBase, _ := url.Parse(base)
	ref, _ := url.Parse(strings.TrimLeft(raw, "/"))
	if strings.HasPrefix(raw, "/") {
		ref, _ = url.Parse(raw)
	}
	return parsedBase.ResolveReference(ref).String()
}

func Document(data []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(bytes.NewReader(data))
}

func FirstText(sel *goquery.Selection, selectors ...string) string {
	for _, selector := range selectors {
		if value := strings.TrimSpace(sel.Find(selector).First().Text()); value != "" {
			return value
		}
	}
	return strings.TrimSpace(sel.Text())
}

func Attr(sel *goquery.Selection, selector, attr string) string {
	target := sel
	if selector != "" {
		target = sel.Find(selector).First()
	}
	value, _ := target.Attr(attr)
	return strings.TrimSpace(value)
}

func Lang(fallback, raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	switch raw {
	case "1", "en", "eng", "english", "ingles", "английски":
		return "eng"
	case "2", "bg", "bul", "bulgarian", "български":
		return "bul"
	case "fr", "fre", "fra", "french":
		return "fre"
	case "ro", "ron", "romanian":
		return "ron"
	case "lv", "lav", "lva", "latvian":
		return "lav"
	case "es", "spa", "spanish", "español":
		return "spa"
	}
	if fallback != "" {
		return fallback
	}
	return raw
}

func Format(raw string) string {
	ext := strings.TrimPrefix(path.Ext(raw), ".")
	if ext == "" || len(ext) > 5 {
		return "srt"
	}
	return strings.ToLower(ext)
}

func DownloadName(raw string, resp *http.Response) string {
	if resp != nil {
		_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if params["filename"] != "" {
			return params["filename"]
		}
	}
	if parsed, err := url.Parse(raw); err == nil {
		if name := path.Base(parsed.Path); name != "." && name != "/" && name != "" {
			return name
		}
	}
	return "subtitle.srt"
}

func DownloadSubtitle(ctx context.Context, svc providercore.Service, cfg providercore.Config, provider, baseURL string, cand providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(cand.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	data, resp, err := Do(ctx, svc, cfg, RequestSpec{Provider: provider, BaseURL: baseURL, Path: cand.SourceURL, Download: true})
	if err != nil {
		return providercore.Download{}, err
	}
	member, err := providercore.ExtractSubtitle(DownloadName(cand.SourceURL, resp), data, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: Absolute(cfg, baseURL, cand.SourceURL)}, nil
}
