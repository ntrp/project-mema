package imports

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
)

type completedDownloadSelection struct {
	SelectedSources    []string
	RejectedCandidates []rejectedDownloadCandidate
}

type rejectedDownloadCandidate struct {
	SourcePath string
	Reason     string
}

type downloadCandidate struct {
	SourcePath string
	SizeBytes  int64
	ExtRank    int
	Depth      int
}

func selectCompletedDownloadCandidates(files []downloadclients.StatusFile, mappings []storage.PathMapping) (completedDownloadSelection, error) {
	candidates := []downloadCandidate{}
	rejected := []rejectedDownloadCandidate{}
	for _, file := range files {
		if !file.Complete {
			rejected = append(rejected, rejectedDownloadCandidate{SourcePath: file.Path, Reason: "incomplete"})
			continue
		}
		mapped := mapPath(file.Path, mappings)
		info, err := os.Stat(mapped)
		if err != nil {
			return completedDownloadSelection{}, fmt.Errorf("download file is not visible to the app: %s", mapped)
		}
		if info.IsDir() {
			found, err := videoCandidatesInDir(mapped)
			if err != nil {
				return completedDownloadSelection{}, err
			}
			if len(found) == 0 {
				rejected = append(rejected, rejectedDownloadCandidate{SourcePath: mapped, Reason: "no_video_files"})
			}
			candidates = append(candidates, found...)
			continue
		}
		if !isVideoFile(mapped) {
			rejected = append(rejected, rejectedDownloadCandidate{SourcePath: mapped, Reason: "not_video_file"})
			continue
		}
		candidates = append(candidates, newDownloadCandidate(mapped, file.SizeBytes, info))
	}
	if len(candidates) == 0 {
		return completedDownloadSelection{RejectedCandidates: rejected}, nil
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return betterCandidate(candidates[i], candidates[j])
	})
	for _, candidate := range candidates[1:] {
		rejected = append(rejected, rejectedDownloadCandidate{
			SourcePath: candidate.SourcePath,
			Reason:     "lower_scoring_candidate",
		})
	}
	return completedDownloadSelection{
		SelectedSources:    []string{candidates[0].SourcePath},
		RejectedCandidates: rejected,
	}, nil
}

func videoCandidatesInDir(root string) ([]downloadCandidate, error) {
	candidates := []downloadCandidate{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !isVideoFile(path) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		candidates = append(candidates, newDownloadCandidate(path, 0, info))
		return nil
	})
	return candidates, err
}

func newDownloadCandidate(path string, reportedSize int64, info os.FileInfo) downloadCandidate {
	size := info.Size()
	if reportedSize > 0 {
		size = reportedSize
	}
	return downloadCandidate{
		SourcePath: path,
		SizeBytes:  size,
		ExtRank:    videoExtensionRank(path),
		Depth:      strings.Count(filepath.Clean(path), string(os.PathSeparator)),
	}
}

func betterCandidate(left downloadCandidate, right downloadCandidate) bool {
	if left.SizeBytes != right.SizeBytes {
		return left.SizeBytes > right.SizeBytes
	}
	if left.ExtRank != right.ExtRank {
		return left.ExtRank > right.ExtRank
	}
	if left.Depth != right.Depth {
		return left.Depth < right.Depth
	}
	return left.SourcePath < right.SourcePath
}

func videoExtensionRank(path string) int {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mkv":
		return 90
	case ".mp4", ".m4v":
		return 80
	case ".mov":
		return 70
	case ".ts":
		return 60
	case ".avi", ".wmv":
		return 50
	case ".mpeg", ".mpg":
		return 40
	case ".webm":
		return 30
	default:
		return 0
	}
}
