package providers

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"media-manager/internal/subtitles/providercore"
)

var (
	registryMu sync.RWMutex
	registry   = map[string]providercore.Adapter{}
)

func Register(key string, adapter providercore.Adapter) {
	normalized := canonicalKey(key)
	if normalized == "" {
		panic("subtitle provider registry: empty key")
	}
	if adapter == nil {
		panic(fmt.Sprintf("subtitle provider registry: nil adapter for %q", normalized))
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := registry[normalized]; exists {
		panic(fmt.Sprintf("subtitle provider registry: duplicate adapter for %q", normalized))
	}
	registry[normalized] = adapter
}

func AdapterFor(key string) (providercore.Adapter, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	adapter, ok := registry[canonicalKey(key)]
	return adapter, ok
}

func RegisteredKeys() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	keys := make([]string, 0, len(registry))
	for key := range registry {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func canonicalKey(key string) string {
	key = strings.ToLower(strings.TrimSpace(key))
	if key == "opensubtitles" {
		return "opensubtitlescom"
	}
	return key
}
