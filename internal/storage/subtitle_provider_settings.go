package storage

import (
	"encoding/json"
	"sort"
)

type SubtitleProviderSettingValue struct {
	StringValue  *string  `json:"stringValue,omitempty"`
	NumberValue  *float64 `json:"numberValue,omitempty"`
	BooleanValue *bool    `json:"booleanValue,omitempty"`
	StringValues []string `json:"stringValues,omitempty"`
}

type SubtitleProviderSettings map[string]SubtitleProviderSettingValue

type SubtitleProviderSecretSettings map[string]string

func subtitleProviderSettingsJSON(settings SubtitleProviderSettings) []byte {
	if len(settings) == 0 {
		return []byte("{}")
	}
	data, err := json.Marshal(settings)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func subtitleProviderSecretsJSON(secrets SubtitleProviderSecretSettings) []byte {
	if len(secrets) == 0 {
		return []byte("{}")
	}
	data, err := json.Marshal(secrets)
	if err != nil {
		return []byte("{}")
	}
	return data
}

func subtitleProviderSettingsFromJSON(data []byte) SubtitleProviderSettings {
	var settings SubtitleProviderSettings
	if len(data) == 0 || json.Unmarshal(data, &settings) != nil || settings == nil {
		return SubtitleProviderSettings{}
	}
	return settings
}

func subtitleProviderSecretsFromJSON(data []byte) SubtitleProviderSecretSettings {
	var secrets SubtitleProviderSecretSettings
	if len(data) == 0 || json.Unmarshal(data, &secrets) != nil || secrets == nil {
		return SubtitleProviderSecretSettings{}
	}
	return secrets
}

func subtitleProviderSecretKeys(secrets SubtitleProviderSecretSettings) []string {
	keys := make([]string, 0, len(secrets))
	for key, value := range secrets {
		if value != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}

func subtitleProviderSettingString(settings SubtitleProviderSettings, key string) *string {
	value, ok := settings[key]
	if !ok || value.StringValue == nil {
		return nil
	}
	return value.StringValue
}

func subtitleProviderSettingValue(value string) SubtitleProviderSettingValue {
	return SubtitleProviderSettingValue{StringValue: &value}
}

func normalizedSubtitleProviderInput(input SubtitleProviderInput) SubtitleProviderInput {
	if input.Settings == nil {
		input.Settings = SubtitleProviderSettings{}
	}
	if input.SecretSettings == nil {
		input.SecretSettings = SubtitleProviderSecretSettings{}
	}
	if input.BaseURL == "" {
		if value := subtitleProviderSettingString(input.Settings, "baseUrl"); value != nil {
			input.BaseURL = *value
		}
	}
	if input.Username == nil {
		input.Username = subtitleProviderSettingString(input.Settings, "username")
	}
	if input.Password == nil {
		if value, ok := input.SecretSettings["password"]; ok && value != "" {
			input.Password = &value
		}
	}
	if input.APIKey == nil {
		if value, ok := input.SecretSettings["apiKey"]; ok && value != "" {
			input.APIKey = &value
		}
	}
	mirrorSubtitleProviderFields(input.Settings, input.SecretSettings, input)
	return input
}

func preserveSubtitleProviderUpdateSecrets(input SubtitleProviderInput, current SubtitleProvider) SubtitleProviderInput {
	hasNewAPIKey := input.SecretSettings != nil && input.SecretSettings["apiKey"] != ""
	hasNewPassword := input.SecretSettings != nil && input.SecretSettings["password"] != ""
	if input.SecretSettings == nil {
		input.SecretSettings = cloneSubtitleProviderSecrets(current.SecretSettings)
	} else {
		for key, value := range current.SecretSettings {
			if _, ok := input.SecretSettings[key]; !ok {
				input.SecretSettings[key] = value
			}
		}
	}
	if input.APIKey == nil && !hasNewAPIKey {
		input.APIKey = current.APIKey
	}
	if input.Password == nil && !hasNewPassword {
		input.Password = current.Password
	}
	for _, field := range input.ClearSecretFields {
		delete(input.SecretSettings, field)
		if field == "apiKey" {
			input.APIKey = nil
		}
		if field == "password" {
			input.Password = nil
		}
	}
	return input
}

func cloneSubtitleProviderSecrets(secrets SubtitleProviderSecretSettings) SubtitleProviderSecretSettings {
	clone := SubtitleProviderSecretSettings{}
	for key, value := range secrets {
		clone[key] = value
	}
	return clone
}

func normalizedSubtitleProvider(provider SubtitleProvider) SubtitleProvider {
	if provider.Settings == nil {
		provider.Settings = SubtitleProviderSettings{}
	}
	if provider.SecretSettings == nil {
		provider.SecretSettings = SubtitleProviderSecretSettings{}
	}
	mirrorSubtitleProviderFields(provider.Settings, provider.SecretSettings, SubtitleProviderInput{
		BaseURL:  provider.BaseURL,
		Username: provider.Username,
		Password: provider.Password,
		APIKey:   provider.APIKey,
	})
	provider.SecretFieldsSet = subtitleProviderSecretKeys(provider.SecretSettings)
	return provider
}

func mirrorSubtitleProviderFields(
	settings SubtitleProviderSettings,
	secrets SubtitleProviderSecretSettings,
	input SubtitleProviderInput,
) {
	if input.BaseURL != "" {
		settings["baseUrl"] = subtitleProviderSettingValue(input.BaseURL)
	}
	if input.Username != nil && *input.Username != "" {
		settings["username"] = subtitleProviderSettingValue(*input.Username)
	}
	if input.Password != nil && *input.Password != "" {
		secrets["password"] = *input.Password
	}
	if input.APIKey != nil && *input.APIKey != "" {
		secrets["apiKey"] = *input.APIKey
	}
}
