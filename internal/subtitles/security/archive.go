package security

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/nwaples/rardecode/v2"
	"github.com/ulikunitz/xz"
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
	limits = normalizeArchiveLimits(limits)
	lower := strings.ToLower(strings.TrimSpace(name))
	switch {
	case strings.HasSuffix(lower, ".zip") || bytes.HasPrefix(data, []byte("PK")):
		return readZip(data, limits)
	case strings.HasSuffix(lower, ".rar") || bytes.HasPrefix(data, []byte("Rar!")):
		return readRAR(data, limits)
	case strings.HasSuffix(lower, ".gz"):
		return readSingleCompressed(name, data, limits, openGzip)
	case strings.HasSuffix(lower, ".xz"):
		return readSingleCompressed(name, data, limits, openXZ)
	default:
		return nil, ErrUnsupportedArchive
	}
}

func normalizeArchiveLimits(limits ArchiveLimits) ArchiveLimits {
	if limits.MaxMembers <= 0 {
		limits.MaxMembers = 128
	}
	if limits.MaxBytes <= 0 {
		limits.MaxBytes = 50 << 20
	}
	return limits
}

func readZip(data []byte, limits ArchiveLimits) ([]ArchiveMember, error) {
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
		if err := validateArchiveMember(file.Name, file.FileInfo().Mode(), file.FileInfo().IsDir()); err != nil {
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
		content, err := readLimited(rc, limits.MaxBytes-total)
		closeErr := rc.Close()
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, closeErr
		}
		total += int64(len(content))
		members = append(members, ArchiveMember{Name: file.Name, Content: content})
	}
	return members, nil
}

func readRAR(data []byte, limits ArchiveLimits) ([]ArchiveMember, error) {
	rr, err := rardecode.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	members := []ArchiveMember{}
	var total int64
	for {
		header, err := rr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(members)+1 > limits.MaxMembers {
			return nil, fmt.Errorf("%w: too many members", ErrUnsafeArchive)
		}
		if err := validateArchiveMember(header.Name, header.Mode(), header.IsDir); err != nil {
			return nil, err
		}
		if header.IsDir {
			continue
		}
		content, err := readLimited(rr, limits.MaxBytes-total)
		if err != nil {
			return nil, err
		}
		total += int64(len(content))
		members = append(members, ArchiveMember{Name: header.Name, Content: content})
	}
	return members, nil
}

func readSingleCompressed(name string, data []byte, limits ArchiveLimits, opener func([]byte) (io.ReadCloser, error)) ([]ArchiveMember, error) {
	memberName := stripCompressedExtension(filepath.Base(name))
	if err := validateArchiveMember(memberName, 0o644, false); err != nil {
		return nil, err
	}
	rc, err := opener(data)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	content, err := readLimited(rc, limits.MaxBytes)
	if err != nil {
		return nil, err
	}
	return []ArchiveMember{{Name: memberName, Content: content}}, nil
}

func openGzip(data []byte) (io.ReadCloser, error) { return gzip.NewReader(bytes.NewReader(data)) }

func openXZ(data []byte) (io.ReadCloser, error) {
	reader, err := xz.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return io.NopCloser(reader), nil
}

func readLimited(reader io.Reader, maxBytes int64) ([]byte, error) {
	content, err := io.ReadAll(io.LimitReader(reader, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(content)) > maxBytes {
		return nil, fmt.Errorf("%w: uncompressed size limit exceeded", ErrUnsafeArchive)
	}
	return content, nil
}

func stripCompressedExtension(name string) string {
	lower := strings.ToLower(name)
	for _, suffix := range []string{".gz", ".xz"} {
		if strings.HasSuffix(lower, suffix) {
			return name[:len(name)-len(suffix)]
		}
	}
	return name
}

func validateArchiveMember(name string, mode fs.FileMode, isDir bool) error {
	cleaned := strings.ReplaceAll(name, "\\", "/")
	if cleaned == "" || strings.HasPrefix(cleaned, "/") || strings.Contains(cleaned, "../") || strings.HasPrefix(cleaned, "../") || filepath.Clean(cleaned) == "." {
		return fmt.Errorf("%w: path traversal member %q", ErrUnsafeArchive, name)
	}
	lower := strings.ToLower(cleaned)
	for _, suffix := range []string{".zip", ".rar", ".7z", ".gz", ".xz"} {
		if strings.HasSuffix(lower, suffix) {
			return fmt.Errorf("%w: nested archive member %q", ErrUnsafeArchive, name)
		}
	}
	if mode.Type() != 0 && !isDir {
		return fmt.Errorf("%w: non-regular member %q", ErrUnsafeArchive, name)
	}
	return nil
}
