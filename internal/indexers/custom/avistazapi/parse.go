package avistazapi

import (
	"encoding/json"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

func parseReleases(config engine.Config, options Options, body []byte) ([]engine.Release, error) {
	var decoded apiResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}
	releases := make([]engine.Release, 0, len(decoded.Data))
	for _, item := range decoded.Data {
		release := releaseFromAPI(config, options, item)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func releaseFromAPI(config engine.Config, options Options, item apiRelease) engine.Release {
	title := strings.TrimSpace(item.FileName)
	if options.PreferRelease && strings.TrimSpace(item.ReleaseTitle) != "" {
		title = strings.TrimSpace(item.ReleaseTitle)
	}
	seeders := nullableInt32Ptr(item.Seed)
	peers := nullablePeers(item.Seed, item.Leech)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           title,
		DownloadURL:     engine.FirstNonEmpty(item.Download, common.Magnet(item.InfoHash)),
		InfoURL:         strings.TrimSpace(item.URL),
		GUID:            engine.FirstNonEmpty(item.URL, item.Download, item.InfoHash),
		SizeBytes:       nullableInt64Value(item.FileSize),
		Seeders:         seeders,
		Peers:           peers,
		PublishedAt:     common.ParseFlexibleTime(item.CreatedAtISO),
	}
}

func nullableInt32Ptr(value nullableInt32) *int32 {
	if !value.Valid {
		return nil
	}
	result := value.Value
	return &result
}

func nullablePeers(seed nullableInt32, leech nullableInt32) *int32 {
	if !seed.Valid && !leech.Valid {
		return nil
	}
	var total int32
	if seed.Valid {
		total += seed.Value
	}
	if leech.Valid {
		total += leech.Value
	}
	return &total
}

func nullableInt64Value(value nullableInt64) int64 {
	if !value.Valid {
		return 0
	}
	return value.Value
}
