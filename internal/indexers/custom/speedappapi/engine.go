package speedappapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

type Options struct {
	Name           string
	DefaultBaseURL string
}

type Engine struct {
	client  engine.HTTPDoer
	options Options
}

func New(options Options, clients ...engine.HTTPDoer) *Engine {
	var client engine.HTTPDoer
	if len(clients) > 0 {
		client = clients[0]
	}
	return &Engine{client: client, options: options}
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	token, err := apiToken(ctx, e.client, config, e.options)
	if err != nil {
		return engine.FailedResult("Invalid "+e.options.Name+" request", "error", err.Error())
	}
	endpoint, err := searchURL(config, e.options, "test")
	if err != nil {
		return engine.FailedResult("Invalid "+e.options.Name+" request", "error", err.Error())
	}
	req, err := authenticatedRequest(ctx, endpoint, token)
	if err != nil {
		return engine.FailedResult("Invalid "+e.options.Name+" request", "error", err.Error())
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return engine.RequestFailedResult(err)
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.StatusFailedResult(resp.StatusCode)
	}
	return engine.SuccessResult(e.options.Name+" indexer reachable", "endpoint", endpoint)
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	token, err := apiToken(ctx, e.client, config, e.options)
	if err != nil {
		return nil, err
	}
	endpoint, err := searchURL(config, e.options, query)
	if err != nil {
		return nil, err
	}
	req, err := authenticatedRequest(ctx, endpoint, token)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, engine.HTTPStatusError(resp)
	}
	body, err := engine.ReadLimitedBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseReleases(config, e.options, body)
}

func apiToken(ctx context.Context, client engine.HTTPDoer, config engine.Config, options Options) (string, error) {
	if token := common.FieldString(config, "apiKey", "apikey", "token"); token != "" {
		return token, nil
	}
	email := common.FieldString(config, "email", "username")
	password := common.FieldString(config, "password")
	if email == "" || password == "" {
		return "", fmt.Errorf("%s requires API key or email/password", options.Name)
	}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, options.DefaultBaseURL), "/api/login", nil)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(map[string]string{"username": email, "password": password})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", engine.HTTPStatusError(resp)
	}
	var decoded struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if strings.TrimSpace(decoded.Token) == "" {
		return "", fmt.Errorf("%s authentication response did not include a token", options.Name)
	}
	return decoded.Token, nil
}

func searchURL(config engine.Config, options Options, query string) (string, error) {
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, options.DefaultBaseURL) + "/api/torrent")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("itemsPerPage", "100")
	values.Set("sort", "torrent.createdAt")
	values.Set("direction", "desc")
	if strings.TrimSpace(query) != "" {
		values.Set("search", strings.TrimSpace(query))
	}
	for _, category := range config.Categories {
		values.Add("categories[]", strconv.FormatInt(int64(category), 10))
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func authenticatedRequest(ctx context.Context, endpoint string, token string) (*http.Request, error) {
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

func cleanTitle(value string) string {
	value = regexp.MustCompile(`(?i)\[REQUEST(ED)?\]`).ReplaceAllString(value, "")
	return strings.Trim(value, " .")
}
