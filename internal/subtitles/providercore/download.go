package providercore

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ulikunitz/xz"

	"media-manager/internal/subtitles/security"
)

var subtitleExtensions = map[string]bool{
	".srt": true, ".ass": true, ".ssa": true, ".vtt": true, ".sub": true,
}

func ExtractSubtitle(name string, data []byte, limits security.ArchiveLimits) (security.ArchiveMember, error) {
	members, err := security.ReadArchive(name, data, limits)
	if err == nil {
		return BestSubtitleMember(members)
	}
	if !errors.Is(err, security.ErrUnsupportedArchive) {
		return security.ArchiveMember{}, err
	}
	lower := strings.ToLower(strings.TrimSpace(name))
	if strings.HasSuffix(lower, ".gz") {
		return decompressSingle(name, data, limits, openGzipReader)
	}
	if strings.HasSuffix(lower, ".xz") {
		return decompressSingle(name, data, limits, openXZReader)
	}
	return security.ArchiveMember{Name: filepath.Base(name), Content: data}, nil
}

func BestSubtitleMember(members []security.ArchiveMember) (security.ArchiveMember, error) {
	candidates := make([]security.ArchiveMember, 0, len(members))
	for _, member := range members {
		if subtitleExtensions[strings.ToLower(filepath.Ext(member.Name))] {
			candidates = append(candidates, member)
		}
	}
	if len(candidates) == 0 {
		return security.ArchiveMember{}, fmt.Errorf("%w: no subtitle member", ErrProviderPrerequisiteMissing)
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return subtitleScore(candidates[i]) > subtitleScore(candidates[j])
	})
	return candidates[0], nil
}

func subtitleScore(member security.ArchiveMember) int {
	name := strings.ToLower(filepath.Base(member.Name))
	score := len(member.Content)
	if strings.Contains(name, "forced") || strings.Contains(name, "sdh") || strings.Contains(name, "commentary") {
		score -= 1 << 20
	}
	if strings.HasSuffix(name, ".srt") {
		score += 1000
	}
	return score
}

func decompressSingle(name string, data []byte, limits security.ArchiveLimits, opener func([]byte) (io.ReadCloser, error)) (security.ArchiveMember, error) {
	if limits.MaxBytes <= 0 {
		limits.MaxBytes = 50 << 20
	}
	rc, err := opener(data)
	if err != nil {
		return security.ArchiveMember{}, err
	}
	defer rc.Close()
	content, err := io.ReadAll(io.LimitReader(rc, limits.MaxBytes+1))
	if err != nil {
		return security.ArchiveMember{}, err
	}
	if int64(len(content)) > limits.MaxBytes {
		return security.ArchiveMember{}, fmt.Errorf("%w: decompressed size limit exceeded", security.ErrUnsafeArchive)
	}
	return security.ArchiveMember{Name: stripDownloadSuffix(filepath.Base(name)), Content: content}, nil
}

func openGzipReader(data []byte) (io.ReadCloser, error) { return gzip.NewReader(bytes.NewReader(data)) }

func openXZReader(data []byte) (io.ReadCloser, error) {
	reader, err := xz.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return io.NopCloser(reader), nil
}

func stripDownloadSuffix(name string) string {
	lower := strings.ToLower(name)
	for _, suffix := range []string{".gz", ".xz"} {
		if strings.HasSuffix(lower, suffix) {
			return name[:len(name)-len(suffix)]
		}
	}
	return name
}
