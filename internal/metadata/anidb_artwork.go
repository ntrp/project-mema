package metadata

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	anidbArtworkURLPrefix = "/api/media/artwork/anidb/"
	anidbArtworkMaxBytes  = 10 << 20
)

func (s *Service) ConfigureAniDBArtworkCache(mediaDataDir string) {
	s.mu.Lock()
	s.anidbArtworkDir = filepath.Join(mediaDataDir, ".mema", "anidb", "artwork")
	s.mu.Unlock()
}

func AniDBArtworkFilePath(mediaDataDir, name string) (string, bool) {
	if !validAniDBArtworkName(name) {
		return "", false
	}
	return filepath.Join(mediaDataDir, ".mema", "anidb", "artwork", name), true
}

func (s *Service) cacheAniDBArtwork(ctx context.Context, picture, userAgent string) *string {
	name, ok := anidbArtworkName(picture)
	if !ok {
		return nil
	}

	s.mu.Lock()
	dir := s.anidbArtworkDir
	baseURL := s.anidbImageBaseURL
	if dir == "" {
		s.mu.Unlock()
		return nil
	}
	path := filepath.Join(dir, name)
	if fileExists(path) {
		s.mu.Unlock()
		return optionalString(anidbArtworkURLPrefix + name)
	}
	if flight := s.anidbArtworkFlights[name]; flight != nil {
		s.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil
		case <-flight:
			if fileExists(path) {
				return optionalString(anidbArtworkURLPrefix + name)
			}
			return nil
		}
	}
	flight := make(chan struct{})
	s.anidbArtworkFlights[name] = flight
	s.mu.Unlock()

	err := s.downloadAniDBArtwork(ctx, strings.TrimRight(baseURL, "/")+"/"+name, path, userAgent)
	s.mu.Lock()
	delete(s.anidbArtworkFlights, name)
	close(flight)
	s.mu.Unlock()
	if err != nil {
		return nil
	}
	return optionalString(anidbArtworkURLPrefix + name)
}

func (s *Service) downloadAniDBArtwork(ctx context.Context, sourceURL, destination, userAgent string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "image/jpeg")
	req.Header.Set("User-Agent", userAgent)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ProviderHTTPError{StatusCode: resp.StatusCode}
	}
	contents, err := io.ReadAll(io.LimitReader(resp.Body, anidbArtworkMaxBytes+1))
	if err != nil {
		return err
	}
	if len(contents) > anidbArtworkMaxBytes {
		return fmt.Errorf("AniDB artwork exceeds size limit")
	}
	imageConfig, err := jpeg.DecodeConfig(bytes.NewReader(contents))
	if err != nil {
		return fmt.Errorf("decode AniDB artwork: %w", err)
	}
	if imageConfig.Width <= 0 || imageConfig.Height <= 0 {
		return fmt.Errorf("AniDB artwork has invalid dimensions")
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(destination), ".anidb-artwork-*")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	written, err := temp.Write(contents)
	if err == nil && written != len(contents) {
		err = io.ErrShortWrite
	}
	if err == nil {
		err = temp.Chmod(0o644)
	}
	if closeErr := temp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tempName, destination)
}

func anidbArtworkName(value string) (string, bool) {
	name := filepath.Base(strings.TrimSpace(value))
	return name, validAniDBArtworkName(name)
}

func validAniDBArtworkName(name string) bool {
	if !strings.HasSuffix(name, ".jpg") || len(name) <= len(".jpg") {
		return false
	}
	for _, char := range strings.TrimSuffix(name, ".jpg") {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular() && info.Size() > 0
}
