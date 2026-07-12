package providers

import (
	"context"
	"crypto/md5" // #nosec G501 -- NapiProjekt's public API requires this legacy content hash.
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"media-manager/internal/subtitles/providercore"
)

const providerReadLimit = 10 << 20

func providerRequest(ctx context.Context, service providercore.Service, method, rawURL, provider string, download bool, body io.Reader) ([]byte, *http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, nil, err
	}
	response, err := service.DoProviderRequest(request, provider, download)
	if err != nil {
		return nil, response, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, providerReadLimit+1))
	if err != nil {
		return nil, response, err
	}
	if len(data) > providerReadLimit {
		return nil, response, fmt.Errorf("provider response size limit exceeded")
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, response, fmt.Errorf("provider returned HTTP %d", response.StatusCode)
	}
	return data, response, nil
}

func napiprojektHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New() // #nosec G401 -- provider-specific lookup hash, not a security control.
	if _, err := io.Copy(hash, io.LimitReader(file, 10<<20)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
