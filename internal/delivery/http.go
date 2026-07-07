package delivery

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ServeFile(w http.ResponseWriter, r *http.Request, target string) error {
	file, err := os.Open(target)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		return ErrDirectory
	}
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", ContentType(target))
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	return nil
}

func StatFile(target string) (os.FileInfo, error) {
	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, ErrDirectory
	}
	return info, nil
}

func ContentType(path string) string {
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

func PlaylistFilename(name string) string {
	base := strings.TrimSpace(strings.TrimSuffix(name, filepath.Ext(name)))
	if base == "" || base == "." {
		base = "media-stream"
	}
	return base + ".m3u"
}

func PlaylistTitle(name string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(name)
}
