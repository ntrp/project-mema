package gazelleapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

func searchURL(config engine.Config, options Options, query string) (string, error) {
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, options.DefaultBaseURL) + "/ajax.php")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("action", "browse")
	values.Set("order_by", "time")
	values.Set("order_way", "desc")
	if search := strings.TrimSpace(query); search != "" {
		values.Set("searchstr", strings.ReplaceAll(search, ".", " "))
	}
	for _, category := range config.Categories {
		if category > 0 && category < 1000 {
			values.Set(fmt.Sprintf("filter_cat[%d]", category), "1")
		}
	}
	if options.FreeleechParam != "" && (common.FieldBool(config, "freeleechOnly") ||
		common.FieldBool(config, "freeleech") ||
		common.FieldBool(config, "freeloadOnly")) {
		values.Set(options.FreeleechParam, freeleechValue(options))
	}
	if options.ExcludeScene && common.FieldBool(config, "excludeScene") {
		values.Set("scene", "0")
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func cookie(ctx context.Context, client engine.HTTPDoer, config engine.Config, options Options) (string, error) {
	if configured := common.FieldString(config, "cookie"); configured != "" {
		return configured, nil
	}
	username := common.FieldString(config, "username")
	password := common.FieldString(config, "password")
	if username == "" || password == "" {
		return "", fmt.Errorf("%s requires cookie or username/password", options.Name)
	}
	endpoint, err := loginURL(config, options)
	if err != nil {
		return "", err
	}
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	form.Set("keeplogged", "1")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", endpoint)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", engine.HTTPStatusError(resp)
	}
	cookies := make([]string, 0, len(resp.Cookies()))
	for _, value := range resp.Cookies() {
		cookies = append(cookies, value.Name+"="+value.Value)
	}
	if len(cookies) == 0 {
		return "", fmt.Errorf("%s authentication did not return cookies", options.Name)
	}
	return strings.Join(cookies, "; "), nil
}

func authorizedRequest(ctx context.Context, client engine.HTTPDoer, config engine.Config, options Options, endpoint string) (*http.Request, error) {
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if options.AuthHeader != "" {
		key := apiKey(config)
		if key == "" {
			return nil, fmt.Errorf("%s requires an API key", options.Name)
		}
		req.Header.Set(options.AuthHeader, options.AuthTokenPrefix+key)
		return req, nil
	}
	cookie, err := cookie(ctx, client, config, options)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", cookie)
	return req, nil
}

func freeleechValue(options Options) string {
	if options.FreeleechValue != "" {
		return options.FreeleechValue
	}
	return "1"
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

func loginURL(config engine.Config, options Options) (string, error) {
	return common.URLWithQuery(common.BaseURL(config, options.DefaultBaseURL), "/login.php", nil)
}
