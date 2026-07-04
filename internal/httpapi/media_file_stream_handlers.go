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
)

func (s *Server) StreamMediaItemFile(w http.ResponseWriter, r *http.Request, id ResourceId, params StreamMediaItemFileParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	target, err := s.settings.MediaItemFilePath(r.Context(), uuid.UUID(id), params.Path)
	if err != nil {
		writeSettingsError(w, err, "Could not find media file")
		return
	}
	serveMediaFile(w, r, target)
}

func (s *Server) DownloadMediaItemFilePlaylist(w http.ResponseWriter, r *http.Request, id ResourceId, params DownloadMediaItemFilePlaylistParams) {
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
	w.Header().Set("Content-Type", "audio/x-mpegurl; charset=utf-8")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": playlistFilename(name)}))
	_, _ = fmt.Fprintf(w, "#EXTM3U\n#EXTINF:-1,%s\n%s\n", playlistTitle(name), streamURL(r, params.Path))
}

func serveMediaFile(w http.ResponseWriter, r *http.Request, target string) {
	file, err := os.Open(target)
	if err != nil {
		writeFileOpenError(w, err)
		return
	}
	defer file.Close()
	info, ok := statOpenMediaFile(w, file)
	if !ok {
		return
	}
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", mediaContentType(target))
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

func statMediaFile(w http.ResponseWriter, target string) (os.FileInfo, bool) {
	info, err := os.Stat(target)
	if err != nil {
		writeFileOpenError(w, err)
		return nil, false
	}
	if info.IsDir() {
		writeError(w, http.StatusBadRequest, "invalid_input", "Media file path points to a directory")
		return nil, false
	}
	return info, true
}

func statOpenMediaFile(w http.ResponseWriter, file *os.File) (os.FileInfo, bool) {
	info, err := file.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_file_stat_failed", "Could not inspect media file")
		return nil, false
	}
	if info.IsDir() {
		writeError(w, http.StatusBadRequest, "invalid_input", "Media file path points to a directory")
		return nil, false
	}
	return info, true
}

func writeFileOpenError(w http.ResponseWriter, err error) {
	if os.IsNotExist(err) {
		writeError(w, http.StatusNotFound, "not_found", "Could not find media file")
		return
	}
	writeError(w, http.StatusInternalServerError, "media_file_open_failed", "Could not open media file")
}

func mediaContentType(path string) string {
	if value := mime.TypeByExtension(strings.ToLower(filepath.Ext(path))); value != "" {
		return value
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mkv":
		return "video/x-matroska"
	case ".m4v":
		return "video/x-m4v"
	case ".ts":
		return "video/mp2t"
	default:
		return "application/octet-stream"
	}
}

func streamURL(r *http.Request, filePath string) string {
	query := url.Values{"path": []string{filePath}}
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
	base := strings.TrimSpace(strings.TrimSuffix(name, filepath.Ext(name)))
	if base == "" || base == "." {
		base = "media-stream"
	}
	return base + ".m3u"
}

func playlistTitle(name string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(name)
}
