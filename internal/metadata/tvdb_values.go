package metadata

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

type tvdbStringNumber string

func (value *tvdbStringNumber) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*value = ""
		return nil
	}
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*value = tvdbStringNumber(strings.TrimSpace(asString))
		return nil
	}
	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err != nil {
		return err
	}
	*value = tvdbStringNumber(asNumber.String())
	return nil
}

func (value tvdbStringNumber) String() string {
	return strings.TrimSpace(string(value))
}

type tvdbStatusValue string

func (value *tvdbStatusValue) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*value = ""
		return nil
	}
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*value = tvdbStatusValue(strings.TrimSpace(asString))
		return nil
	}
	var asObject struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &asObject); err != nil {
		return err
	}
	*value = tvdbStatusValue(strings.TrimSpace(asObject.Name))
	return nil
}

func (value tvdbStatusValue) String() string {
	return strings.TrimSpace(string(value))
}

type tvdbDateValue string

func (value *tvdbDateValue) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*value = ""
		return nil
	}
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*value = tvdbDateValue(strings.TrimSpace(asString))
		return nil
	}
	var asObject struct {
		Date        string `json:"date"`
		ReleaseDate string `json:"releaseDate"`
		FirstAired  string `json:"firstAired"`
	}
	if err := json.Unmarshal(data, &asObject); err != nil {
		return err
	}
	*value = tvdbDateValue(firstNonEmpty(asObject.Date, asObject.ReleaseDate, asObject.FirstAired))
	return nil
}

func (value tvdbDateValue) String() string {
	return strings.TrimSpace(string(value))
}

type tvdbIntPointer struct {
	Value *int64
}

func (value *tvdbIntPointer) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		value.Value = nil
		return nil
	}
	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		parsed, err := asNumber.Int64()
		if err != nil {
			return err
		}
		value.Value = &parsed
		return nil
	}
	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}
	parsed, err := strconv.ParseInt(strings.TrimSpace(asString), 10, 64)
	if err != nil {
		return err
	}
	value.Value = &parsed
	return nil
}
