package catalog

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed subtitle_provider_catalog.generated.json
var catalogFS embed.FS

const artifactName = "subtitle_provider_catalog.generated.json"

func All() ([]Entry, error) {
	data, err := catalogFS.ReadFile(artifactName)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func MustAll() []Entry {
	entries, err := All()
	if err != nil {
		panic(err)
	}
	return entries
}

func Lookup(key string) (Entry, bool) {
	entries, err := All()
	if err != nil {
		return Entry{}, false
	}
	for _, entry := range entries {
		if entry.Key == key {
			return entry, true
		}
	}
	return Entry{}, false
}

func Require(key string) (Entry, error) {
	entry, ok := Lookup(key)
	if !ok {
		return Entry{}, fmt.Errorf("subtitle provider catalog key %q not found", key)
	}
	return entry, nil
}

func PickerKeys() ([]string, error) {
	entries, err := All()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(entries))
	for _, entry := range entries {
		keys = append(keys, entry.Key)
	}
	return keys, nil
}
