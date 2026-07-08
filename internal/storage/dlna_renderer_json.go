package storage

import "encoding/json"

func rendererJSONObject(value json.RawMessage) bool {
	if len(value) == 0 {
		return false
	}
	var decoded map[string]any
	return json.Unmarshal(value, &decoded) == nil && decoded != nil
}
