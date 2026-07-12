package providers

import (
	"context"
	"net/http"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type stubAdapter struct{}

func (stubAdapter) Test(context.Context, providercore.Service, providercore.Config) error { return nil }
func (stubAdapter) Search(context.Context, providercore.Service, providercore.Config, providercore.SearchRequest) ([]providercore.Candidate, error) {
	return nil, nil
}
func (stubAdapter) Download(context.Context, providercore.Service, providercore.Config, providercore.Candidate) (providercore.Download, error) {
	return providercore.Download{}, nil
}

type stubService struct{}

func (stubService) DoProviderRequest(*http.Request, string, bool) (*http.Response, error) {
	return nil, nil
}

func TestRegisterAndLookupAdapter(t *testing.T) {
	key := "registrytest"
	Register(key, stubAdapter{})
	adapter, ok := AdapterFor(" RegistryTest ")
	if !ok {
		t.Fatal("registered adapter not found")
	}
	if err := adapter.Test(context.Background(), stubService{}, providercore.Config{Type: key}); err != nil {
		t.Fatalf("adapter test failed: %v", err)
	}
}

func TestRegisterRejectsDuplicateKeys(t *testing.T) {
	key := "duplicatetest"
	Register(key, stubAdapter{})
	defer func() {
		if recover() == nil {
			t.Fatal("expected duplicate registration panic")
		}
	}()
	Register(" DuplicateTest ", stubAdapter{})
}

func TestOpenSubtitlesAliasIsCanonicalized(t *testing.T) {
	key := "opensubtitlescom-alias-test"
	Register(key, stubAdapter{})
	keys := RegisteredKeys()
	for _, registered := range keys {
		if registered == key {
			return
		}
	}
	t.Fatalf("registered keys %v do not include %q", keys, key)
}
