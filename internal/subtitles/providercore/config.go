package providercore

import (
	"net/url"
	"strconv"
	"strings"
)

type ConfigView struct{ config Config }

func NewConfig(config Config) ConfigView { return ConfigView{config: config} }

func (c ConfigView) StringSetting(key string) string {
	value, ok := c.config.Settings[key]
	if !ok || value.StringValue == nil {
		return ""
	}
	return strings.TrimSpace(*value.StringValue)
}

func (c ConfigView) BoolSetting(key string) bool {
	value, ok := c.config.Settings[key]
	return ok && value.BooleanValue != nil && *value.BooleanValue
}

func (c ConfigView) IntSetting(key string) int {
	value, ok := c.config.Settings[key]
	if !ok {
		return 0
	}
	if value.NumberValue != nil {
		return int(*value.NumberValue)
	}
	if value.StringValue != nil {
		parsed, _ := strconv.Atoi(strings.TrimSpace(*value.StringValue))
		return parsed
	}
	return 0
}

func (c ConfigView) StringsSetting(key string) []string {
	value, ok := c.config.Settings[key]
	if !ok {
		return nil
	}
	items := make([]string, 0, len(value.StringValues))
	for _, item := range value.StringValues {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func (c ConfigView) Secret(key string) string {
	if key == "apiKey" && c.config.APIKey != nil {
		return strings.TrimSpace(*c.config.APIKey)
	}
	if key == "password" && c.config.Password != nil {
		return strings.TrimSpace(*c.config.Password)
	}
	return strings.TrimSpace(c.config.SecretSettings[key])
}

func (c ConfigView) RequiredSecret(key string) (string, bool) {
	secret := c.Secret(key)
	return secret, secret != ""
}

func (c ConfigView) CookieString() string { return c.Secret("cookies") }

func (c ConfigView) BaseURL(defaultURL string) string {
	base := strings.TrimSpace(c.config.BaseURL)
	if base == "" {
		base = c.StringSetting("baseUrl")
	}
	if base == "" {
		base = defaultURL
	}
	if parsed, err := url.Parse(base); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		return strings.TrimRight(parsed.String(), "/")
	}
	return strings.TrimRight(base, "/")
}
