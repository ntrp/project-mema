package indexers

import (
	"context"
	"net/http"
)

type capsDocument struct {
	Server struct {
		Title   string `xml:"title,attr"`
		Version string `xml:"version,attr"`
	} `xml:"server"`
	Limits struct {
		Max     int `xml:"max,attr"`
		Default int `xml:"default,attr"`
	} `xml:"limits"`
	Categories []capsCategory `xml:"categories>category"`
}

type capsCategory struct {
	ID       int            `xml:"id,attr"`
	Name     string         `xml:"name,attr"`
	Children []capsCategory `xml:"subcat"`
}

func (s *Service) testCaps(ctx context.Context, config Config) TestResult {
	endpoint, err := endpointWithQuery(config.BaseURL, map[string]string{
		"t":      "caps",
		"apikey": stringValue(config.APIKey),
	})
	if err != nil {
		return failedResult("Invalid indexer URL", "error", err.Error())
	}

	req, err := get(ctx, endpoint)
	if err != nil {
		return failedResult("Invalid indexer request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return requestFailedResult(err)
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return statusFailedResult(resp.StatusCode)
	}

	var payload capsDocument
	if err := decodeLimitedXML(resp.Body, &payload); err != nil {
		return failedResult("Indexer caps response could not be parsed", "error", err.Error())
	}
	if len(payload.Categories) == 0 && payload.Server.Title == "" {
		return failedResult("Indexer caps response did not include capabilities")
	}

	return successResult(
		"Indexer capabilities OK",
		"title", payload.Server.Title,
		"version", payload.Server.Version,
		"categoryCount", countCategories(payload.Categories),
		"maxResults", payload.Limits.Max,
		"defaultResults", payload.Limits.Default,
	)
}

func countCategories(categories []capsCategory) int {
	count := len(categories)
	for _, category := range categories {
		count += countCategories(category.Children)
	}
	return count
}
