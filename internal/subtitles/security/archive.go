package security

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedArchive = errors.New("unsupported archive format")
	ErrUnsafeArchive      = errors.New("unsafe archive")
)

type ArchiveLimits struct {
	MaxMembers int
	MaxBytes   int64
}

type ArchiveMember struct {
	Name    string
	Content []byte
}

func ReadArchive(name string, data []byte, limits ArchiveLimits) ([]ArchiveMember, error) {
	lower := strings.ToLower(strings.TrimSpace(name))
	if strings.HasSuffix(lower, ".rar") || strings.HasSuffix(lower, ".7z") || bytes.HasPrefix(data, []byte("Rar!")) || bytes.HasPrefix(data, []byte("7z\xbc\xaf\x27\x1c")) {
		return nil, ErrUnsupportedArchive
	}
	if !strings.HasSuffix(lower, ".zip") && !bytes.HasPrefix(data, []byte("PK")) {
		return nil, ErrUnsupportedArchive
	}
	if limits.MaxMembers <= 0 {
		limits.MaxMembers = 128
	}
	if limits.MaxBytes <= 0 {
		limits.MaxBytes = 50 << 20
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	if len(zr.File) > limits.MaxMembers {
		return nil, fmt.Errorf("%w: too many members", ErrUnsafeArchive)
	}
	members := make([]ArchiveMember, 0, len(zr.File))
	var total int64
	for _, file := range zr.File {
		if err := validateArchiveMember(file); err != nil {
			return nil, err
		}
		if file.FileInfo().IsDir() {
			continue
		}
		if total+int64(file.UncompressedSize64) > limits.MaxBytes {
			return nil, fmt.Errorf("%w: uncompressed size limit exceeded", ErrUnsafeArchive)
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		content, readErr := io.ReadAll(io.LimitReader(rc, limits.MaxBytes-total+1))
		closeErr := rc.Close()
		if readErr != nil {
			return nil, readErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		total += int64(len(content))
		if total > limits.MaxBytes {
			return nil, fmt.Errorf("%w: uncompressed size limit exceeded", ErrUnsafeArchive)
		}
		members = append(members, ArchiveMember{Name: file.Name, Content: content})
	}
	return members, nil
}

func validateArchiveMember(file *zip.File) error {
	name := strings.ReplaceAll(file.Name, "\\", "/")
	if name == "" || strings.HasPrefix(name, "/") || strings.Contains(name, "../") || strings.HasPrefix(name, "../") || filepath.Clean(name) == "." {
		return fmt.Errorf("%w: path traversal member %q", ErrUnsafeArchive, file.Name)
	}
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, ".zip") || strings.HasSuffix(lower, ".rar") || strings.HasSuffix(lower, ".7z") {
		return fmt.Errorf("%w: nested archive member %q", ErrUnsafeArchive, file.Name)
	}
	mode := file.FileInfo().Mode()
	if mode.Type() != 0 && !file.FileInfo().IsDir() {
		return fmt.Errorf("%w: non-regular member %q", ErrUnsafeArchive, file.Name)
	}
	return nil
}
