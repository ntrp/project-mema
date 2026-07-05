package storage

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestIndexerProxiesUseGeneratedQueries(t *testing.T) {
	ctx, store := testDBStore(t)
	name := "proxy-" + uuid.NewString()

	created, err := store.CreateIndexerProxy(ctx, IndexerProxyInput{
		Name:                  " " + name + " ",
		Implementation:        "flarerSolverr",
		Link:                  " http://proxy.example.test ",
		Enabled:               true,
		OnHealthIssue:         true,
		SupportsOnHealthIssue: true,
		IncludeHealthWarnings: true,
		TestCommand:           "test",
		Fields:                json.RawMessage(`[{"name":"url","value":"http://proxy.example.test"}]`),
	})
	if err != nil {
		t.Fatalf("create indexer proxy: %v", err)
	}
	if created.Name != name || created.Link != "http://proxy.example.test" {
		t.Fatalf("created proxy should trim name and link, got %#v", created)
	}
	if len(created.Fields) == 0 {
		t.Fatal("created proxy should retain fields")
	}

	updated, err := store.UpdateIndexerProxy(ctx, created.ID, IndexerProxyInput{
		Name:                  name + "-updated",
		Implementation:        "custom",
		Link:                  "http://updated-proxy.example.test",
		Enabled:               false,
		OnHealthIssue:         false,
		SupportsOnHealthIssue: true,
		IncludeHealthWarnings: false,
		TestCommand:           "ping",
		Fields:                json.RawMessage(`[{"name":"url","value":"http://updated-proxy.example.test"}]`),
	})
	if err != nil {
		t.Fatalf("update indexer proxy: %v", err)
	}
	if updated.Name != name+"-updated" || updated.Enabled {
		t.Fatalf("updated proxy = %#v", updated)
	}

	found, err := store.GetIndexerProxy(ctx, created.ID)
	if err != nil {
		t.Fatalf("get indexer proxy: %v", err)
	}
	if found.ID != created.ID || found.Name != updated.Name {
		t.Fatalf("found proxy = %#v, want updated %#v", found, updated)
	}

	proxies, err := store.ListIndexerProxies(ctx)
	if err != nil {
		t.Fatalf("list indexer proxies: %v", err)
	}
	if !indexerProxyListHas(proxies, created.ID) {
		t.Fatalf("created proxy missing from list: %#v", proxies)
	}

	if err := store.DeleteIndexerProxy(ctx, created.ID); err != nil {
		t.Fatalf("delete indexer proxy: %v", err)
	}
	if _, err := store.GetIndexerProxy(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected deleted proxy to be missing, got %v", err)
	}
	if err := store.DeleteIndexerProxy(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected second delete to report not found, got %v", err)
	}
}

func indexerProxyListHas(proxies []IndexerProxy, id uuid.UUID) bool {
	for _, proxy := range proxies {
		if proxy.ID == id {
			return true
		}
	}
	return false
}
