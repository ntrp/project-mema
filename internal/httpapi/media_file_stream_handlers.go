package httpapi

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"media-manager/internal/delivery"
)

func (s *Server) StreamMediaItemFile(w http.ResponseWriter, r *http.Request, id ResourceId, params StreamMediaItemFileParams) {
	mediaID := uuid.UUID(id)
	if !s.authorizeMediaFileStream(w, r, mediaID, params) {
		return
	}
	target, err := s.settings.MediaItemFilePath(r.Context(), mediaID, params.Path)
	if err != nil {
		writeSettingsError(w, err, "Could not find media file")
		return
	}
	serveMediaFile(w, r, target)
}

func (s *Server) PlayMediaItemFileInVlc(w http.ResponseWriter, r *http.Request, id ResourceId, params PlayMediaItemFileInVlcParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	target, err := s.settings.MediaItemFilePath(r.Context(), uuid.UUID(id), params.Path)
	if err != nil {
		writeSettingsError(w, err, "Could not find media file")
		return
	}
	if _, ok := statMediaFile(w, target); !ok {
		return
	}
	name := filepath.Base(target)
	expires := s.now().Add(streamTokenTTL).Unix()
	token := s.newStreamToken(uuid.UUID(id), params.Path, expires)
	w.Header().Set("Content-Type", "audio/x-mpegurl; charset=utf-8")
	w.Header().Set("Content-Disposition", playlistDisposition(name))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = fmt.Fprintf(w, "#EXTM3U\n#EXTINF:-1,%s\n%s\n", delivery.PlaylistTitle(name), streamURL(r, params.Path, expires, token))
}

func (s *Server) authorizeMediaFileStream(w http.ResponseWriter, r *http.Request, mediaID uuid.UUID, params StreamMediaItemFileParams) bool {
	if s.validStreamToken(mediaID, params.Path, params.StreamExpires, params.StreamToken) {
		return true
	}
	_, ok := s.requireSession(w, r)
	return ok
}

func serveMediaFile(w http.ResponseWriter, r *http.Request, target string) {
	writeFileError(w, delivery.ServeFile(w, r, target))
}

func statMediaFile(w http.ResponseWriter, target string) (os.FileInfo, bool) {
	info, err := delivery.StatFile(target)
	if err != nil {
		writeFileError(w, err)
		return nil, false
	}
	return info, true
}

func writeFileError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	if err == delivery.ErrDirectory {
		writeError(w, http.StatusBadRequest, "invalid_input", "Media file path points to a directory")
		return
	}
	if os.IsNotExist(err) {
		writeError(w, http.StatusNotFound, "not_found", "Could not find media file")
		return
	}
	writeError(w, http.StatusInternalServerError, "media_file_open_failed", "Could not open media file")
}

func mediaContentType(path string) string {
	return delivery.ContentType(path)
}

func streamURL(r *http.Request, filePath string, expires int64, token string) string {
	query := url.Values{
		"path":          []string{filePath},
		"streamExpires": []string{fmt.Sprintf("%d", expires)},
		"streamToken":   []string{token},
	}
	return (&url.URL{
		Scheme:   requestScheme(r),
		Host:     requestHost(r),
		Path:     strings.TrimSuffix(r.URL.Path, "/vlc") + "/stream",
		RawQuery: query.Encode(),
	}).String()
}

func requestScheme(r *http.Request) string {
	if value := forwardedHeader(r, "X-Forwarded-Proto"); value != "" {
		return value
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func requestHost(r *http.Request) string {
	if value := forwardedHeader(r, "X-Forwarded-Host"); value != "" {
		return value
	}
	return r.Host
}

func forwardedHeader(r *http.Request, name string) string {
	value := strings.TrimSpace(r.Header.Get(name))
	if index := strings.Index(value, ","); index >= 0 {
		value = strings.TrimSpace(value[:index])
	}
	return value
}

func playlistFilename(name string) string {
	return delivery.PlaylistFilename(name)
}

func playlistDisposition(name string) string {
	return mime.FormatMediaType("inline", map[string]string{"filename": playlistFilename(name)})
}
