package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/indexers/engine"
)

var sizeRegex = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*([KMGT]i?B|B)`)

func FieldString(config engine.Config, names ...string) string {
	fields := fieldMap(config)
	for _, name := range names {
		value, ok := fields[strings.ToLower(name)]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case string:
			return strings.TrimSpace(typed)
		case fmt.Stringer:
			return strings.TrimSpace(typed.String())
		}
	}
	return ""
}

func FieldBool(config engine.Config, name string) bool {
	value, ok := fieldMap(config)[strings.ToLower(name)]
	if !ok {
		return false
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(strings.TrimSpace(typed), "true")
	default:
		return false
	}
}

func FieldFloat(config engine.Config, name string) float64 {
	value, ok := fieldMap(config)[strings.ToLower(name)]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case float64:
		return typed
	case int:
		return float64(typed)
	case string:
		var parsed float64
		if _, err := fmt.Sscanf(typed, "%f", &parsed); err == nil {
			return parsed
		}
	}
	return 0
}

func BaseURL(config engine.Config, fallback string) string {
	if strings.TrimSpace(config.BaseURL) != "" {
		return strings.TrimRight(strings.TrimSpace(config.BaseURL), "/")
	}
	return strings.TrimRight(fallback, "/")
}

func URLWithQuery(baseURL string, path string, values map[string]string) (string, error) {
	parsed, err := engine.ParseBaseURL(strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/"))
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	for key, value := range values {
		if value != "" {
			query.Set(key, value)
		}
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func ResolveURL(baseURL string, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err == nil && parsed.IsAbs() {
		return parsed.String()
	}
	base, err := url.Parse(strings.TrimRight(baseURL, "/") + "/")
	if err != nil {
		return value
	}
	relative, err := url.Parse(value)
	if err != nil {
		return value
	}
	return base.ResolveReference(relative).String()
}

func GetJSON(ctx context.Context, client engine.HTTPDoer, endpoint string, target any) error {
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.HTTPStatusError(resp)
	}
	return decodeJSON(resp.Body, target)
}

func PostJSON(ctx context.Context, client engine.HTTPDoer, endpoint string, payload any, target any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.HTTPStatusError(resp)
	}
	return decodeJSON(resp.Body, target)
}

func TestURL(ctx context.Context, client engine.HTTPDoer, endpoint string, name string) engine.TestResult {
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return engine.FailedResult("Invalid indexer request", "error", err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		return engine.RequestFailedResult(err)
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return engine.StatusFailedResult(resp.StatusCode)
	}
	return engine.SuccessResult(name+" indexer reachable", "endpoint", endpoint)
}

func Magnet(infoHash string) string {
	infoHash = strings.TrimSpace(infoHash)
	if infoHash == "" {
		return ""
	}
	if strings.HasPrefix(infoHash, "magnet:") {
		return infoHash
	}
	return "magnet:?xt=urn:btih:" + url.QueryEscape(infoHash)
}

func Int32Ptr(value int) *int32 {
	converted := int32(value)
	return &converted
}

func UnixTime(seconds int64) *time.Time {
	if seconds <= 0 {
		return nil
	}
	value := time.Unix(seconds, 0).UTC()
	return &value
}

func ParseTime(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	suffix := value
	if len(suffix) > 6 {
		suffix = suffix[len(suffix)-6:]
	}
	if !strings.Contains(value, "Z") && !strings.ContainsAny(suffix, "+-") {
		value += "+01:00"
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05-07:00"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			converted := parsed.UTC()
			return &converted
		}
	}
	return nil
}

func ParseFlexibleTime(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	for _, layout := range []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02",
	} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			converted := parsed.UTC()
			return &converted
		}
	}
	return ParseTime(value)
}

func ParseSizeBytes(value string) int64 {
	match := sizeRegex.FindStringSubmatch(strings.TrimSpace(value))
	if len(match) != 3 {
		return 0
	}
	amount, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0
	}
	unit := strings.ToUpper(match[2])
	multiplier := float64(1)
	switch unit {
	case "KB", "KIB":
		multiplier = 1024
	case "MB", "MIB":
		multiplier = 1024 * 1024
	case "GB", "GIB":
		multiplier = 1024 * 1024 * 1024
	case "TB", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024
	}
	return int64(amount * multiplier)
}

func decodeJSON(body io.Reader, target any) error {
	data, err := engine.ReadLimitedBody(body)
	if err != nil {
		return err
	}
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "[]" {
		return nil
	}
	if err := json.Unmarshal([]byte(trimmed), target); err != nil {
		return fmt.Errorf("parse json response: %w", err)
	}
	return nil
}

func fieldMap(config engine.Config) map[string]any {
	values := map[string]any{}
	if len(config.Fields) == 0 {
		return values
	}
	var fields []struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}
	if err := json.Unmarshal(config.Fields, &fields); err != nil {
		return values
	}
	for _, field := range fields {
		if strings.TrimSpace(field.Name) != "" {
			values[strings.ToLower(strings.TrimSpace(field.Name))] = field.Value
		}
	}
	return values
}
