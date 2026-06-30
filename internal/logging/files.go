package logging

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const logFilePrefix = "app-"

type FileSettings struct {
	Enabled       bool
	Directory     string
	RetentionDays int
}

type FileInfo struct {
	Name       string
	SizeBytes  int64
	ModifiedAt time.Time
	Path       string
}

type fileSink struct {
	mu             sync.Mutex
	settings       FileSettings
	currentDate    string
	currentFile    *os.File
	lastCleanupDay string
}

func newFileSink() *fileSink {
	return &fileSink{settings: FileSettings{RetentionDays: 7}}
}

func (m *Manager) ConfigureFile(settings FileSettings) error {
	settings = normalizeFileSettings(settings)
	return m.fileSink.configure(settings)
}

func (m *Manager) LogFileSettings() FileSettings {
	return m.fileSink.currentSettings()
}

func ListFiles(directory string) ([]FileInfo, error) {
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return []FileInfo{}, nil
	}
	if err != nil {
		return nil, err
	}
	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !isLogFileName(entry.Name()) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, FileInfo{
			Name:       entry.Name(),
			SizeBytes:  info.Size(),
			ModifiedAt: info.ModTime().UTC(),
			Path:       filepath.Join(directory, entry.Name()),
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name > files[j].Name
	})
	return files, nil
}

func FilePath(directory, name string) (string, bool) {
	if !isLogFileName(name) || filepath.Base(name) != name {
		return "", false
	}
	return filepath.Join(directory, name), true
}

func isLogFileName(name string) bool {
	return strings.HasPrefix(name, logFilePrefix) && strings.HasSuffix(name, ".log")
}

func normalizeFileSettings(settings FileSettings) FileSettings {
	settings.Directory = strings.TrimSpace(settings.Directory)
	if settings.RetentionDays < 1 {
		settings.RetentionDays = 7
	}
	return settings
}

func (s *fileSink) configure(settings FileSettings) error {
	if settings.Enabled {
		if settings.Directory == "" {
			return errors.New("log directory is required")
		}
		if err := os.MkdirAll(settings.Directory, 0o755); err != nil {
			return fmt.Errorf("create log directory: %w", err)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.currentFile != nil && s.settings.Directory != settings.Directory {
		_ = s.currentFile.Close()
		s.currentFile = nil
		s.currentDate = ""
	}
	if !settings.Enabled && s.currentFile != nil {
		_ = s.currentFile.Close()
		s.currentFile = nil
		s.currentDate = ""
	}
	s.settings = settings
	return nil
}

func (s *fileSink) currentSettings() FileSettings {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.settings
}
