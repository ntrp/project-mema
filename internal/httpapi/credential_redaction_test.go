package httpapi

import (
	"encoding/json"
	"strings"
	"testing"

	"media-manager/internal/storage"
)

func TestCredentialResponsesMaskRawSecrets(t *testing.T) {
	secret := "raw-secret-value"
	download := downloadClientResponse(storage.DownloadClient{
		Password: &secret,
		APIKey:   &secret,
	})
	assertNoRawSecret(t, download, secret)
	if !download.PasswordSet || !download.ApiKeySet || download.Password != nil || download.ApiKey != nil {
		t.Fatalf("download client secret flags = %#v", download)
	}

	indexer := indexerResponse(storage.Indexer{APIKey: &secret}, newCatalogLanguageMapper(nil))
	assertNoRawSecret(t, indexer, secret)
	if !indexer.ApiKeySet || indexer.ApiKey != nil {
		t.Fatalf("indexer secret flags = %#v", indexer)
	}

	provider := metadataProviderResponse(storage.MetadataProvider{
		APIKey:      &secret,
		PIN:         &secret,
		AccessToken: &secret,
	})
	assertNoRawSecret(t, provider, secret)
	if !provider.ApiKeySet || !provider.PinSet || !provider.AccessTokenSet {
		t.Fatalf("metadata provider secret flags = %#v", provider)
	}
	if provider.ApiKey != nil || provider.Pin != nil || provider.AccessToken != nil {
		t.Fatalf("metadata provider leaked secret fields = %#v", provider)
	}
}

func TestCredentialUpdatesPreserveOmittedSecrets(t *testing.T) {
	secret := "raw-secret-value"
	download := preserveDownloadClientSecrets(
		storage.DownloadClientInput{},
		DownloadClientRequest{},
		storage.DownloadClient{Password: &secret, APIKey: &secret},
	)
	if download.Password == nil || *download.Password != secret || download.APIKey == nil || *download.APIKey != secret {
		t.Fatalf("download secrets = %#v", download)
	}

	indexer := preserveIndexerSecrets(storage.IndexerInput{}, IndexerRequest{}, storage.Indexer{APIKey: &secret})
	if indexer.APIKey == nil || *indexer.APIKey != secret {
		t.Fatalf("indexer secret = %#v", indexer)
	}

	provider := preserveMetadataProviderSecrets(
		storage.MetadataProviderInput{},
		MetadataProviderRequest{},
		storage.MetadataProvider{APIKey: &secret, PIN: &secret, AccessToken: &secret},
	)
	if provider.APIKey == nil || provider.PIN == nil || provider.AccessToken == nil {
		t.Fatalf("metadata provider secrets = %#v", provider)
	}
}

func assertNoRawSecret(t *testing.T, value any, secret string) {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), secret) {
		t.Fatalf("response leaked raw secret: %s", payload)
	}
}
