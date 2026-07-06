package httpapi

import (
	"testing"

	"media-manager/internal/storage"
)

func TestCredentialResponsesIncludeRawSecrets(t *testing.T) {
	secret := "raw-secret-value"
	download := downloadClientResponse(storage.DownloadClient{
		Password: &secret,
		APIKey:   &secret,
	})
	if !download.PasswordSet || !download.ApiKeySet {
		t.Fatalf("download client secret flags = %#v", download)
	}
	if download.Password == nil || *download.Password != secret || download.ApiKey == nil || *download.ApiKey != secret {
		t.Fatalf("download client secrets = %#v", download)
	}

	indexer := indexerResponse(storage.Indexer{APIKey: &secret}, newCatalogLanguageMapper(nil))
	if !indexer.ApiKeySet {
		t.Fatalf("indexer secret flags = %#v", indexer)
	}
	if indexer.ApiKey == nil || *indexer.ApiKey != secret {
		t.Fatalf("indexer secret = %#v", indexer)
	}

	provider := metadataProviderResponse(storage.MetadataProvider{
		APIKey:      &secret,
		PIN:         &secret,
		AccessToken: &secret,
	})
	if !provider.ApiKeySet || !provider.PinSet || !provider.AccessTokenSet {
		t.Fatalf("metadata provider secret flags = %#v", provider)
	}
	if provider.ApiKey == nil || *provider.ApiKey != secret ||
		provider.Pin == nil || *provider.Pin != secret ||
		provider.AccessToken == nil || *provider.AccessToken != secret {
		t.Fatalf("metadata provider secrets = %#v", provider)
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
