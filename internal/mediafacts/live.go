package mediafacts

import (
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/delivery"
	"media-manager/internal/storage"
)

func TrackID(mediaItemID uuid.UUID, filePath string, trackType string, streamIndex int32) uuid.UUID {
	key := strings.Join([]string{mediaItemID.String(), filePath, trackType, strconv.FormatInt(int64(streamIndex), 10)}, "\x00")
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(key))
}

func OtherFileID(mediaItemID uuid.UUID, mediaPath string, otherPath string, fileType string) uuid.UUID {
	key := strings.Join([]string{mediaItemID.String(), mediaPath, otherPath, fileType}, "\x00")
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(key))
}

func WithLiveFileFacts(item storage.MediaItem, filePath string) storage.MediaItem {
	paths := item.FilePaths
	if filePath != "" {
		paths = []string{filePath}
	}
	for _, path := range paths {
		if fileFactHasTracks(item.FileFacts, path) {
			continue
		}
		fact, ok := liveFact(item.ID, path)
		if !ok {
			continue
		}
		item.FileFacts = replaceFileFact(item.FileFacts, fact)
	}
	return item
}

func fileFactHasTracks(facts []storage.MediaFileFact, path string) bool {
	for _, fact := range facts {
		if fact.FilePath == path && len(fact.Tracks) > 0 {
			return true
		}
	}
	return false
}

func liveFact(mediaItemID uuid.UUID, path string) (storage.MediaFileFact, bool) {
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		return storage.MediaFileFact{}, false
	}
	probe := delivery.Probe(path)
	if len(probe.Tracks) == 0 {
		return storage.MediaFileFact{}, false
	}
	size := stat.Size()
	fact := storage.MediaFileFact{
		ID:          uuid.New(),
		MediaItemID: mediaItemID,
		FilePath:    path,
		SizeBytes:   &size,
	}
	for _, track := range probe.Tracks {
		fact.Tracks = append(fact.Tracks, liveTrack(mediaItemID, fact.ID, path, track))
	}
	return fact, true
}

func liveTrack(
	mediaItemID uuid.UUID,
	factID uuid.UUID,
	path string,
	track delivery.Track,
) storage.MediaFileTrackFact {
	return storage.MediaFileTrackFact{
		ID:              TrackID(mediaItemID, path, string(track.Type), int32Value(track.Index)),
		MediaFileFactID: factID,
		MediaItemID:     mediaItemID,
		FilePath:        path,
		StreamIndex:     int32Value(track.Index),
		TrackType:       string(track.Type),
		LanguageID:      track.Language,
		Codec:           track.Codec,
		Channels:        track.ChannelLayout,
		DurationMs:      durationMs(track.Duration),
		BitrateKbps:     bitrateKbps(track.BitRate),
		Width:           track.Width,
		Height:          track.Height,
		PixelFormat:     track.PixelFormat,
		Title:           track.Title,
	}
}

func durationMs(value *float64) *int64 {
	if value == nil || *value <= 0 {
		return nil
	}
	ms := int64(*value * 1000)
	return &ms
}

func int32Value(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}

func bitrateKbps(value *string) *int32 {
	if value == nil {
		return nil
	}
	bits, err := strconv.ParseInt(*value, 10, 32)
	if err != nil || bits <= 0 {
		return nil
	}
	kbps := int32(bits / 1000)
	return &kbps
}

func replaceFileFact(facts []storage.MediaFileFact, fact storage.MediaFileFact) []storage.MediaFileFact {
	for index := range facts {
		if facts[index].FilePath == fact.FilePath {
			facts[index] = fact
			return facts
		}
	}
	return append(facts, fact)
}
