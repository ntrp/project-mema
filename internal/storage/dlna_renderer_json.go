package storage

import (
	"encoding/json"
	"regexp"
)

func rendererJSONObject(value json.RawMessage) bool {
	if len(value) == 0 {
		return false
	}
	var decoded map[string]any
	return json.Unmarshal(value, &decoded) == nil && decoded != nil
}

func validateRendererMatchRegexes(value json.RawMessage) error {
	var decoded any
	if err := json.Unmarshal(value, &decoded); err != nil {
		return ErrInvalidInput
	}
	if !rendererRegexesValid(decoded) {
		return ErrInvalidInput
	}
	return nil
}

func rendererRegexesValid(value any) bool {
	switch typed := value.(type) {
	case map[string]any:
		if !rendererRegexRuleValid(typed) {
			return false
		}
		for _, child := range typed {
			if !rendererRegexesValid(child) {
				return false
			}
		}
	case []any:
		for _, child := range typed {
			if !rendererRegexesValid(child) {
				return false
			}
		}
	}
	return true
}

func rendererRegexRuleValid(rule map[string]any) bool {
	kind, _ := rule["kind"].(string)
	ruleType, _ := rule["type"].(string)
	if kind != "regex" && ruleType != "regex" {
		return true
	}
	pattern, _ := rule["value"].(string)
	if pattern == "" {
		pattern, _ = rule["contains"].(string)
	}
	if pattern == "" {
		return false
	}
	_, err := regexp.Compile(pattern)
	return err == nil
}
