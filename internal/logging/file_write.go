package logging

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func (m *Manager) writeFileRecord(record slog.Record) {
	if m.fileSink == nil {
		return
	}
	_ = m.fileSink.write(record)
}

func (s *fileSink) write(record slog.Record) error {
	now := time.Now()
	day := now.Format(time.DateOnly)
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.settings.Enabled {
		return nil
	}
	if err := s.rotate(day); err != nil {
		return err
	}
	if err := s.cleanup(day); err != nil {
		return err
	}
	payload := Entry{
		ID:         "",
		Time:       record.Time.UTC(),
		Level:      levelFromSlog(record.Level),
		Message:    record.Message,
		Attributes: attrsFromRecord(record),
	}
	if payload.Time.IsZero() {
		payload.Time = now.UTC()
	}
	line, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = s.currentFile.Write(append(line, '\n'))
	return err
}

func (s *fileSink) rotate(day string) error {
	if s.currentFile != nil && s.currentDate == day {
		return nil
	}
	if s.currentFile != nil {
		_ = s.currentFile.Close()
		s.currentFile = nil
	}
	if err := os.MkdirAll(s.settings.Directory, 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(s.settings.Directory, logFilePrefix+day+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	s.currentFile = file
	s.currentDate = day
	return nil
}

func (s *fileSink) cleanup(day string) error {
	if s.lastCleanupDay == day {
		return nil
	}
	s.lastCleanupDay = day
	files, err := ListFiles(s.settings.Directory)
	if err != nil {
		return err
	}
	cutoff := time.Now().AddDate(0, 0, -s.settings.RetentionDays+1)
	for _, file := range files {
		if file.ModifiedAt.Before(cutoff) {
			_ = os.Remove(file.Path)
		}
	}
	return nil
}
