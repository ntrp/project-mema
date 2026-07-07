package httpapi

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/delivery"
	mediatools "media-manager/internal/tools"
)

type mediaFileTrackDeleteInputError struct {
	message string
}

func (e mediaFileTrackDeleteInputError) Error() string {
	return e.message
}

func isMediaFileTrackDeleteInputError(err error) bool {
	var inputErr mediaFileTrackDeleteInputError
	return errors.As(err, &inputErr)
}

func deleteMediaFileTrack(ctx context.Context, path string, request MediaFileTrackDeleteRequest) error {
	output, cleanup, err := trackDeleteOutput(path)
	if err != nil {
		return err
	}
	defer cleanup()
	commandArgs, err := mediaFileTrackDeleteArgs(path, output, request)
	if err != nil {
		return err
	}
	defer commandArgs.cleanup()
	if _, err := mediatools.RunOutput(ctx, mediatools.CommandSpec{
		Name:           "ffmpeg",
		Args:           commandArgs.args,
		Timeout:        30 * time.Minute,
		MaxOutputBytes: 64 * 1024,
		MaxStderrBytes: 64 * 1024,
	}); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err == nil {
		_ = os.Chmod(output, info.Mode().Perm())
	}
	return os.Rename(output, path)
}

type mediaFileTrackDeleteCommand struct {
	args    []string
	cleanup func()
}

func mediaFileTrackDeleteArgs(
	input string,
	output string,
	request MediaFileTrackDeleteRequest,
) (mediaFileTrackDeleteCommand, error) {
	if err := mediatools.SafePathArg(input); err != nil {
		return mediaFileTrackDeleteCommand{}, err
	}
	if err := mediatools.SafePathArg(output); err != nil {
		return mediaFileTrackDeleteCommand{}, err
	}
	base := []string{"-hide_banner", "-loglevel", "error", "-y", "-i", input}
	switch request.TargetType {
	case MediaFileTrackDeleteRequestTargetTypeAudio, MediaFileTrackDeleteRequestTargetTypeSubtitle:
		if request.TrackIndex == nil {
			return mediaFileTrackDeleteCommand{}, inputError("track index is required")
		}
		return mediaFileTrackDeleteCommand{args: append(base,
			"-map", "0",
			"-map", "-0:"+strconv.FormatInt(int64(*request.TrackIndex), 10),
			"-c", "copy",
			output,
		), cleanup: func() {}}, nil
	case MediaFileTrackDeleteRequestTargetTypeChapters:
		return mediaFileTrackDeleteCommand{args: append(base,
			"-map", "0",
			"-map_chapters", "-1",
			"-c", "copy",
			output,
		), cleanup: func() {}}, nil
	case MediaFileTrackDeleteRequestTargetTypeChapter:
		if request.ChapterIndex == nil {
			return mediaFileTrackDeleteCommand{}, inputError("chapter index is required")
		}
		metadataPath, err := chapterMetadataWithout(input, *request.ChapterIndex)
		if err != nil {
			return mediaFileTrackDeleteCommand{}, err
		}
		return mediaFileTrackDeleteCommand{args: append(base,
			"-i", metadataPath,
			"-map", "0",
			"-map_metadata", "0",
			"-map_chapters", "1",
			"-c", "copy",
			output,
		), cleanup: func() { _ = os.Remove(metadataPath) }}, nil
	default:
		return mediaFileTrackDeleteCommand{}, inputError("unsupported track delete target")
	}
}

func trackDeleteOutput(path string) (string, func(), error) {
	if err := mediatools.SafePathArg(path); err != nil {
		return "", func() {}, err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".mema-track-delete-*"+filepath.Ext(path))
	if err != nil {
		return "", func() {}, err
	}
	output := file.Name()
	if err := file.Close(); err != nil {
		_ = os.Remove(output)
		return "", func() {}, err
	}
	return output, func() { _ = os.Remove(output) }, nil
}

func chapterMetadataWithout(path string, chapterIndex int32) (string, error) {
	probe := delivery.Probe(path)
	probeChapters := mediaFileChaptersFromDelivery(probe.Chapters)
	chapters := make([]MediaFileChapter, 0, len(probeChapters))
	found := false
	for _, chapter := range probeChapters {
		if chapter.Index == chapterIndex {
			found = true
			continue
		}
		chapters = append(chapters, chapter)
	}
	if !found {
		return "", inputError("chapter was not found")
	}
	metadata, err := ffmetadataChapters(chapters)
	if err != nil {
		return "", err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".mema-chapters-*.ffmetadata")
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := file.WriteString(metadata); err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

func ffmetadataChapters(chapters []MediaFileChapter) (string, error) {
	var builder strings.Builder
	builder.WriteString(";FFMETADATA1\n")
	for _, chapter := range chapters {
		start, err := chapterTimeMillis(chapter.StartTime)
		if err != nil {
			return "", err
		}
		end, err := chapterTimeMillis(chapter.EndTime)
		if err != nil {
			return "", err
		}
		if end <= start {
			return "", inputError("chapter end time must be after start time")
		}
		builder.WriteString("[CHAPTER]\nTIMEBASE=1/1000\n")
		builder.WriteString("START=" + strconv.FormatInt(start, 10) + "\n")
		builder.WriteString("END=" + strconv.FormatInt(end, 10) + "\n")
		if chapter.Title != nil && strings.TrimSpace(*chapter.Title) != "" {
			builder.WriteString("title=" + ffmetadataEscape(*chapter.Title) + "\n")
		}
	}
	return builder.String(), nil
}

func chapterTimeMillis(value *string) (int64, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return 0, inputError("chapter time is missing")
	}
	seconds, err := strconv.ParseFloat(strings.TrimSpace(*value), 64)
	if err != nil {
		return 0, inputError(fmt.Sprintf("invalid chapter time %q", *value))
	}
	return int64(seconds * 1000), nil
}

func ffmetadataEscape(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "\n", "\\\n", "=", "\\=", ";", "\\;", "#", "\\#")
	return replacer.Replace(value)
}

func inputError(message string) error {
	return mediaFileTrackDeleteInputError{message: message}
}
