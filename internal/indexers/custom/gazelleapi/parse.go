package gazelleapi

import (
	"encoding/json"
	"html"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

func parseReleases(config engine.Config, options Options, body []byte) ([]engine.Release, error) {
	var decoded response
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}
	if !strings.EqualFold(strings.TrimSpace(decoded.Status), "success") {
		return []engine.Release{}, nil
	}
	releases := []engine.Release{}
	for _, group := range decoded.Response.Results {
		if len(group.Torrents) == 0 && group.TorrentID != 0 {
			group.Torrents = []torrent{group.asTorrent()}
		}
		for _, item := range group.Torrents {
			if skipFreeleechToken(config, options, item.CanUseToken) || skipFreeload(config, item.IsFreeload) {
				continue
			}
			release := releaseFromTorrent(config, options, group, item)
			if release.Title != "" && release.DownloadURL != "" {
				releases = append(releases, release)
			}
		}
	}
	return releases, nil
}

func (g releaseGroup) asTorrent() torrent {
	return torrent{
		TorrentID:   g.TorrentID,
		FileCount:   g.FileCount,
		Size:        g.Size,
		Snatches:    g.Snatches,
		Seeders:     g.Seeders,
		Leechers:    g.Leechers,
		Category:    g.Category,
		Time:        g.GroupTime,
		IsFreeLeech: g.IsFreeLeech,
		IsFreeleech: g.IsFreeleech,
		IsFreeload:  g.IsFreeload,
		IsNeutral:   g.IsNeutral,
		IsPersonal:  g.IsPersonal,
		CanUseToken: g.CanUseToken,
	}
}

func releaseFromTorrent(config engine.Config, options Options, group releaseGroup, item torrent) engine.Release {
	seeders := item.Seeders.Value
	leechers := item.Leechers.Value
	title := titleFrom(group, item, options)
	infoURL := infoURL(config, options, group.GroupID.Value, item.TorrentID)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           title,
		DownloadURL:     downloadURL(config, options, item),
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       item.Size.Value,
		Seeders:         common.Int32Ptr(int(seeders)),
		Peers:           common.Int32Ptr(int(seeders + leechers)),
		PublishedAt:     gazelleTime(item.Time, group.GroupTime),
	}
}

func titleFrom(group releaseGroup, item torrent, options Options) string {
	if item.FileName != "" {
		return strings.TrimSpace(html.UnescapeString(item.FileName))
	}
	if options.PreferMusicTitle {
		return musicTitle(group, item)
	}
	groupName := strings.TrimSpace(html.UnescapeString(group.GroupName))
	if options.PreferGroupTitle || len(group.Torrents) == 0 {
		return groupName
	}
	parts := []string{}
	if strings.TrimSpace(group.Artist) != "" {
		parts = append(parts, strings.TrimSpace(group.Artist))
	}
	if groupName != "" {
		parts = append(parts, groupName)
	}
	title := strings.Join(parts, " - ")
	if strings.TrimSpace(group.GroupYear.Value) != "" {
		title += " (" + strings.TrimSpace(group.GroupYear.Value) + ")"
	}
	format := strings.TrimSpace(strings.Join([]string{item.Format, item.Encoding}, " "))
	if format != "" {
		title += " [" + format + "]"
	}
	if strings.TrimSpace(item.Media) != "" {
		title += " [" + strings.TrimSpace(item.Media) + "]"
	}
	if item.HasCue {
		title += " [Cue]"
	}
	return strings.TrimSpace(title)
}

func musicTitle(group releaseGroup, item torrent) string {
	title := strings.TrimSpace(html.UnescapeString(group.Artist))
	groupName := strings.TrimSpace(html.UnescapeString(group.GroupName))
	if title != "" && groupName != "" {
		title += " - "
	}
	title += groupName
	if strings.TrimSpace(group.GroupYear.Value) != "" {
		title += " (" + strings.TrimSpace(group.GroupYear.Value) + ")"
	}
	if strings.TrimSpace(group.ReleaseType) != "" && !strings.EqualFold(group.ReleaseType, "Unknown") {
		title += " [" + strings.TrimSpace(group.ReleaseType) + "]"
	}
	if strings.TrimSpace(item.RemasterTitle) != "" {
		remaster := strings.TrimSpace(strings.TrimSpace(item.RemasterTitle) + " " + strings.TrimSpace(item.RemasterYear))
		title += " [" + remaster + "]"
	}
	flags := compact([]string{
		strings.TrimSpace(strings.TrimSpace(item.Format) + " " + strings.TrimSpace(item.Encoding)),
		strings.TrimSpace(item.Media),
	})
	if item.HasLog {
		flags = append(flags, "Log ("+strconv.Itoa(item.LogScore)+"%)")
	}
	if item.HasCue {
		flags = append(flags, "Cue")
	}
	if len(flags) > 0 {
		title += " [" + strings.Join(flags, " / ") + "]"
	}
	return strings.TrimSpace(title)
}

func compact(values []string) []string {
	out := []string{}
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, strings.TrimSpace(value))
		}
	}
	return out
}

func infoURL(config engine.Config, options Options, groupID string, torrentID int) string {
	values := map[string]string{"id": strings.TrimSpace(groupID)}
	if torrentID > 0 {
		values["torrentid"] = strconv.Itoa(torrentID)
	}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, options.DefaultBaseURL), "/torrents.php", values)
	if err != nil {
		return ""
	}
	return endpoint
}

func downloadURL(config engine.Config, options Options, item torrent) string {
	if item.TorrentID == 0 {
		return ""
	}
	path := options.DownloadPath
	if path == "" {
		path = "/torrents.php"
	}
	values := map[string]string{
		"action": "download",
		"id":     strconv.Itoa(item.TorrentID),
	}
	if useFreeleechToken(config, options, item) {
		values["usetoken"] = "1"
	}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, options.DefaultBaseURL), path, values)
	if err != nil {
		return ""
	}
	return endpoint
}

func skipFreeleechToken(config engine.Config, options Options, canUseToken bool) bool {
	return options.SupportsFreeleechTokens && freeleechTokenMode(config) == 2 && !canUseToken
}

func skipFreeload(config engine.Config, isFreeload bool) bool {
	return common.FieldBool(config, "freeloadOnly") && !isFreeload
}

func useFreeleechToken(config engine.Config, options Options, item torrent) bool {
	if !options.SupportsFreeleechTokens || !item.CanUseToken || isFree(item) {
		return false
	}
	mode := freeleechTokenMode(config)
	return mode == 1 || mode == 2
}

func freeleechTokenMode(config engine.Config) int {
	return int(common.FieldFloat(config, "useFreeleechToken"))
}

func isFree(item torrent) bool {
	return item.IsFreeLeech || item.IsFreeleech || item.IsNeutral || item.IsFreeload || item.IsPersonal
}

func gazelleTime(values ...string) *time.Time {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
			return common.UnixTime(seconds)
		}
		if parsed := common.ParseFlexibleTime(value); parsed != nil {
			return parsed
		}
	}
	return nil
}
