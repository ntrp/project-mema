package avistazapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

func token(ctx context.Context, client engine.HTTPDoer, config engine.Config, options Options) (string, error) {
	if value := common.FieldString(config, "token", "apiKey", "apikey"); value != "" {
		return value, nil
	}
	username := common.FieldString(config, "username")
	password := common.FieldString(config, "password")
	pid := common.FieldString(config, "pid")
	if username == "" || password == "" || pid == "" {
		return "", fmt.Errorf("%s requires token or username/password/pid", options.Name)
	}
	endpoint, err := authURL(config, options)
	if err != nil {
		return "", err
	}
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	form.Set("pid", strings.TrimSpace(pid))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", engine.HTTPStatusError(resp)
	}
	var decoded authResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if strings.TrimSpace(decoded.Token) == "" {
		return "", fmt.Errorf("%s authentication response did not include a token", options.Name)
	}
	return decoded.Token, nil
}

func authURL(config engine.Config, options Options) (string, error) {
	return common.URLWithQuery(common.BaseURL(config, options.DefaultBaseURL), "/api/v1/jackett/auth", nil)
}

func searchURL(config engine.Config, options Options, query string) (string, error) {
	base, err := engine.ParseBaseURL(common.BaseURL(config, options.DefaultBaseURL) + "/api/v1/jackett/torrents")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("limit", "50")
	if options.AnimeCategories {
		addAnimeCategories(values, config.Categories)
	} else {
		values.Set("in", "1")
		values.Set("type", trackerType(config.Categories))
	}
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		values.Add("discount[]", "1")
	}
	if strings.TrimSpace(query) != "" {
		values.Set("search", strings.TrimSpace(query))
	}
	base.RawQuery = values.Encode()
	return base.String(), nil
}

func authenticatedJSONRequest(ctx context.Context, endpoint string, token string) (*http.Request, error) {
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req, nil
}

func trackerType(categories []int32) string {
	foundMovie := false
	foundTV := false
	for _, category := range categories {
		switch category {
		case 2000, 2030, 2040, 2045, 2048, 2050, 2060:
			foundMovie = true
		case 5000, 5030, 5040, 5045, 5048, 5050, 5060, 5070:
			foundTV = true
		}
	}
	if foundMovie && !foundTV {
		return "1"
	}
	if foundTV && !foundMovie {
		return "2"
	}
	return "0"
}

func addAnimeCategories(values url.Values, categories []int32) {
	for _, category := range categories {
		switch category {
		case 2000:
			values.Add("format[]", "MOVIE")
		case 5000, 5070:
			values.Add("format[]", "TV")
			values.Add("format[]", "TV_SHORT")
			values.Add("format[]", "SPECIAL")
			values.Add("format[]", "OVA")
			values.Add("format[]", "ONA")
		case 7030:
			values.Add("format[]", "MANGA")
		}
	}
}
