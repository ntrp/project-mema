package cardigann

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type cardigannSearchRequest struct {
	URL        string
	Method     string
	Body       url.Values
	Headers    http.Header
	SearchPath cardigannSearchPath
}

func cardigannSearchRequests(def cardigannDefinition, config Config, ctx cardigannContext) ([]cardigannSearchRequest, error) {
	paths := def.Search.Paths
	if len(paths) == 0 && def.Search.Path != "" {
		paths = []cardigannSearchPath{{cardigannRequest: cardigannRequest{Path: def.Search.Path}}}
	}
	requests := make([]cardigannSearchRequest, 0, len(paths))
	for _, path := range paths {
		req, err := buildCardigannSearchRequest(def, config, ctx, path)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func buildCardigannSearchRequest(
	def cardigannDefinition,
	config Config,
	ctx cardigannContext,
	path cardigannSearchPath,
) (cardigannSearchRequest, error) {
	renderedPath, err := renderCardigannTemplate(path.Path, ctx)
	if err != nil {
		return cardigannSearchRequest{}, err
	}
	renderedPath = strings.ReplaceAll(renderedPath, " ", "%20")
	endpoint, err := resolveCardigannURL(config.BaseURL, renderedPath)
	if err != nil {
		return cardigannSearchRequest{}, err
	}
	method := strings.ToUpper(strings.TrimSpace(path.Method))
	if method == "" {
		method = http.MethodGet
	}
	inputs := mergedCardigannInputs(def.Search.Inputs, path.Inputs, path.InheritInputs)
	body := url.Values{}
	rawQuery := []string{}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return cardigannSearchRequest{}, err
	}
	query := parsed.Query()
	for key, value := range inputs {
		rendered, err := renderCardigannTemplate(value, ctx)
		if err != nil {
			return cardigannSearchRequest{}, err
		}
		if rendered == "" && !def.Search.AllowEmptyInputs {
			continue
		}
		if key == "$raw" {
			rawQuery = append(rawQuery, strings.TrimLeft(rendered, "&?"))
			continue
		}
		if method == http.MethodPost {
			body.Set(key, rendered)
		} else {
			query.Set(key, rendered)
		}
	}
	parsed.RawQuery = query.Encode()
	for _, raw := range rawQuery {
		if raw == "" {
			continue
		}
		if parsed.RawQuery != "" {
			parsed.RawQuery += "&"
		}
		parsed.RawQuery += raw
	}
	return cardigannSearchRequest{
		URL:        parsed.String(),
		Method:     method,
		Body:       body,
		Headers:    cardigannRequestHeaders(def, ctx),
		SearchPath: path,
	}, nil
}

func mergedCardigannInputs(base map[string]string, override map[string]string, inherit *bool) map[string]string {
	values := map[string]string{}
	if inherit == nil || *inherit {
		for key, value := range base {
			values[key] = value
		}
	}
	for key, value := range override {
		values[key] = value
	}
	return values
}

func renderCardigannHeaders(headers map[string][]string, ctx cardigannContext) http.Header {
	rendered := http.Header{}
	for key, values := range headers {
		for _, value := range values {
			if out, err := renderCardigannTemplate(value, ctx); err == nil && strings.TrimSpace(out) != "" {
				rendered.Add(key, out)
			}
		}
	}
	return rendered
}

func cardigannRequestHeaders(def cardigannDefinition, ctx cardigannContext) http.Header {
	headers := renderCardigannHeaders(def.Search.Headers, ctx)
	if def.Login == nil || !strings.EqualFold(def.Login.Method, "cookie") {
		return headers
	}
	cookie := ""
	if def.Login.Inputs != nil {
		cookie = def.Login.Inputs["cookie"]
	}
	if cookie == "" {
		cookie = "{{ .Config.cookie }}"
	}
	if rendered, err := renderCardigannTemplate(cookie, ctx); err == nil && strings.TrimSpace(rendered) != "" {
		headers.Set("Cookie", rendered)
	}
	return headers
}

func resolveCardigannURL(baseURL string, path string) (string, error) {
	base, err := parseBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return base.String(), nil
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", err
	}
	if parsed.IsAbs() {
		return parsed.String(), nil
	}
	return base.ResolveReference(parsed).String(), nil
}

func (s *Engine) executeCardigannRequest(ctx context.Context, request cardigannSearchRequest) ([]byte, error) {
	var body *strings.Reader
	if request.Method == http.MethodPost {
		body = strings.NewReader(request.Body.Encode())
	} else {
		body = strings.NewReader("")
	}
	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, body)
	if err != nil {
		return nil, err
	}
	req.Header = request.Headers.Clone()
	if request.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, httpStatusError(resp)
	}
	data, err := readLimitedBody(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("empty indexer response")
	}
	return data, nil
}
