package httpapi

import (
	"net/http"
	"os"

	"media-manager/internal/metadata"
)

func (s *Server) GetAniDBArtwork(w http.ResponseWriter, r *http.Request, filename string) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	path, ok := metadata.AniDBArtworkFilePath(s.cfg.MediaDataDir, filename)
	if !ok {
		http.NotFound(w, r)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeContent(w, r, filename, info.ModTime(), file)
}
