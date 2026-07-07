package dlna

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/delivery"
	mediatools "media-manager/internal/tools"
)

var fallbackPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
	0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
	0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x60, 0x00, 0x00, 0x00,
	0x02, 0x00, 0x01, 0xe2, 0x21, 0xbc, 0x33, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
}

func (m *Manager) artwork(w http.ResponseWriter, r *http.Request) {
	id := artworkIDFromPath(r.URL.Path)
	object, err := m.contentTree().BrowseMetadata(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if r.URL.Query().Get("kind") == "thumbnail" && object.FilePath != "" {
		m.thumbnail(w, r, object.FilePath)
		return
	}
	if object.Artwork != nil && serveKnownArtwork(w, r, *object.Artwork) {
		return
	}
	serveFallbackIcon(w, r)
}

func serveKnownArtwork(w http.ResponseWriter, r *http.Request, artwork string) bool {
	artwork = strings.TrimSpace(artwork)
	if strings.HasPrefix(artwork, "http://") || strings.HasPrefix(artwork, "https://") {
		http.Redirect(w, r, artwork, http.StatusFound)
		return true
	}
	if filepath.IsAbs(artwork) {
		return delivery.ServeFile(w, r, artwork) == nil
	}
	return false
}

func serveFallbackIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	if r.Method != http.MethodHead {
		_, _ = w.Write(fallbackPNG)
	}
}

func (m *Manager) thumbnail(w http.ResponseWriter, r *http.Request, target string) {
	w.Header().Set("Content-Type", "image/jpeg")
	if r.Method == http.MethodHead {
		return
	}
	info, err := delivery.StatFile(target)
	if err != nil {
		writeFileError(w, err)
		return
	}
	cachePath := thumbnailCachePath(m.thumbDir, target, info.ModTime(), info.Size())
	if _, err := os.Stat(cachePath); err == nil {
		writeFileError(w, delivery.ServeFile(w, r, cachePath))
		return
	}
	if !acquireTranscodeSlot(w, r) {
		return
	}
	defer func() { <-dlnaTranscodeSlots }()
	if err := generateThumbnail(r, target, cachePath); err != nil {
		http.Error(w, "could not generate thumbnail", http.StatusInternalServerError)
		return
	}
	writeFileError(w, delivery.ServeFile(w, r, cachePath))
}

func generateThumbnail(r *http.Request, target string, cachePath string) error {
	if err := mediatools.SafePathArg(target); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return err
	}
	args := []string{"-hide_banner", "-loglevel", "error", "-y", "-ss", "1", "-i", target, "-frames:v", "1", "-vf", "scale=320:-2", cachePath}
	_, err := mediatools.RunOutput(r.Context(), mediatools.CommandSpec{
		Name:           "ffmpeg",
		Args:           args,
		Timeout:        30 * time.Second,
		MaxStderrBytes: 64 * 1024,
	})
	return err
}

func artworkIDFromPath(path string) string {
	path = strings.TrimPrefix(path, "/dlna")
	path = strings.TrimPrefix(path, "/artwork/")
	id, err := url.PathUnescape(strings.Trim(path, "/"))
	if err != nil {
		return ""
	}
	return id
}

func thumbnailCachePath(dir string, target string, modTime time.Time, size int64) string {
	key := target + "\x00" + modTime.UTC().Format(time.RFC3339Nano) + "\x00" + strconv.FormatInt(size, 10)
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(dir, hex.EncodeToString(sum[:16])+".jpg")
}
