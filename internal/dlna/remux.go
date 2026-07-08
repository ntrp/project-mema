package dlna

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"media-manager/internal/delivery"
	mediatools "media-manager/internal/tools"
)

var errNoDLNATranscodeSlot = errors.New("no DLNA transcode slot available")

func (m *Manager) cachedMatroskaRemux(r *http.Request, target string, probe delivery.ProbeResult) (string, error) {
	cachePath, err := m.matroskaRemuxCachePath(target)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}
	select {
	case dlnaTranscodeSlots <- struct{}{}:
	case <-r.Context().Done():
		return "", r.Context().Err()
	default:
		return "", errNoDLNATranscodeSlot
	}
	defer func() { <-dlnaTranscodeSlots }()
	if err := generateMatroskaRemux(r, target, cachePath, matroskaAudioTranscodeDecision(probe)); err != nil {
		return "", err
	}
	return cachePath, nil
}

func (m *Manager) existingMatroskaRemux(target string) (string, bool) {
	cachePath, err := m.matroskaRemuxCachePath(target)
	if err != nil {
		return "", false
	}
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, true
	}
	return cachePath, false
}

func (m *Manager) matroskaRemuxCachePath(target string) (string, error) {
	info, err := delivery.StatFile(target)
	if err != nil {
		return "", err
	}
	return remuxCachePath(m.remuxDir, target, info.ModTime(), info.Size()), nil
}

func generateMatroskaRemux(r *http.Request, target string, cachePath string, decision delivery.Decision) error {
	if err := mediatools.SafePathArg(target); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(cachePath), filepath.Base(cachePath)+".*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmpPath) }()
	_, err = mediatools.RunOutput(r.Context(), mediatools.CommandSpec{
		Name:           "ffmpeg",
		Args:           matroskaRemuxArgs(target, tmpPath, decision),
		MaxStderrBytes: 64 * 1024,
	})
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, cachePath)
}

func matroskaRemuxArgs(target string, output string, decision delivery.Decision) []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-y",
		"-fflags", "+genpts",
		"-i", target,
		"-map", "0:v:0",
		"-map", "0:a:0?",
		"-sn", "-dn",
		"-c:v", decision.Plan.VideoCodec,
	}
	if decision.Plan.VideoCodec != "copy" {
		args = append(args, "-preset", "veryfast", "-pix_fmt", "yuv420p", "-profile:v", "high")
	}
	args = append(args, "-c:a", decision.Plan.AudioCodec)
	if decision.Plan.AudioCodec != "copy" {
		args = append(args, "-ac", "2")
	}
	return append(args, "-f", "matroska", output)
}

func matroskaRemuxStreamArgs(target string, decision delivery.Decision) []string {
	args := matroskaRemuxArgs(target, "pipe:1", decision)
	return args
}

func remuxCachePath(dir string, target string, modTime time.Time, size int64) string {
	key := target + "\x00" + modTime.UTC().Format(time.RFC3339Nano) + "\x00" + strconv.FormatInt(size, 10)
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(dir, hex.EncodeToString(sum[:16])+".mkv")
}
