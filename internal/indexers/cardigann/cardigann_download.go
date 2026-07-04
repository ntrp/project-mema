package cardigann

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

func (s *Engine) resolveCardigannDownload(
	ctx context.Context,
	def cardigannDefinition,
	config Config,
	templateCtx cardigannContext,
	detailsURL string,
) (string, error) {
	if def.Download == nil || len(def.Download.Selectors) == 0 || detailsURL == "" {
		return detailsURL, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, detailsURL, nil)
	if err != nil {
		return "", err
	}
	req.Header = cardigannDownloadHeaders(def, templateCtx)
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", httpStatusError(resp)
	}
	body, err := readLimitedBody(resp.Body)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	for _, selector := range def.Download.Selectors {
		value, ok, err := selectorFieldValue(doc, selector, templateCtx)
		if err != nil || !ok {
			continue
		}
		return resolveCardigannLink(detailsURL, value), nil
	}
	return "", fmt.Errorf("download selectors did not match")
}

func cardigannDownloadHeaders(def cardigannDefinition, ctx cardigannContext) http.Header {
	headers := renderCardigannHeaders(def.Download.Headers, ctx)
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

func selectorFieldValue(doc *goquery.Document, selector cardigannSelectorField, ctx cardigannContext) (string, bool, error) {
	query, err := renderCardigannTemplate(selector.Selector, ctx)
	if err != nil {
		return "", false, err
	}
	matcher, err := cascadia.Compile(query)
	if err != nil {
		return "", false, err
	}
	match := doc.FindMatcher(matcher).First()
	if match.Length() == 0 {
		return "", false, nil
	}
	value := ""
	if selector.Attribute != "" {
		value, _ = match.Attr(selector.Attribute)
	} else {
		value = match.Text()
	}
	value = applyCardigannFilters(strings.TrimSpace(value), selector.Filters, ctx)
	return value, value != "", nil
}
