package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const animeToshoKey = "animetosho"

func init() { Register(animeToshoKey, animeToshoAdapter{}) }

type animeToshoAdapter struct{}

type animeToshoEntry struct {
	ID        int64  `json:"id"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Title     string `json:"title"`
}

type animeToshoTorrent struct {
	Files []animeToshoFile `json:"files"`
}

type animeToshoFile struct {
	Filename    string                 `json:"filename"`
	Attachments []animeToshoAttachment `json:"attachments"`
}

type animeToshoAttachment struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	Info struct {
		Lang string `json:"lang"`
		Name string `json:"name"`
	} `json:"info"`
}

func (animeToshoAdapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	endpoint := animeToshoFeedURL(config, url.Values{"show": {"torrent"}, "id": {"0"}})
	_, _, err := providerRequest(ctx, service, http.MethodGet, endpoint, animeToshoKey, false, nil)
	return err
}

func (animeToshoAdapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	episodeID := animeToshoEpisodeID(request)
	if episodeID == "" {
		return nil, fmt.Errorf("%w: animetosho requires an AniDB episode id", providercore.ErrProviderPrerequisiteMissing)
	}
	entries, err := animeToshoEntries(ctx, service, config, episodeID)
	if err != nil {
		return nil, err
	}
	want := alpha3Language(request.LanguageID)
	candidates := []providercore.Candidate{}
	for _, entry := range entries {
		torrent, err := animeToshoTorrentFiles(ctx, service, config, entry.ID)
		if err != nil {
			return nil, err
		}
		for _, file := range torrent.Files {
			for _, attachment := range file.Attachments {
				if attachment.Type != "subtitle" {
					continue
				}
				lang := alpha3Language(attachment.Info.Lang)
				if lang == "" {
					lang = "eng"
				}
				if want != "" && lang != want && !(want == "pob" && lang == "por" && strings.Contains(strings.ToLower(attachment.Info.Name), "brazil")) {
					continue
				}
				name := file.Filename
				if strings.TrimSpace(name) == "" {
					name = entry.Title
				}
				candidates = append(candidates, providercore.Candidate{ProviderName: config.Name, LanguageID: request.LanguageID, FileID: attachment.ID, Format: "srt", ReleaseName: name, SourceURL: animeToshoDownloadURL(attachment.ID), SourceRef: strconv.FormatInt(attachment.ID, 10)})
			}
		}
	}
	return candidates, nil
}

func (animeToshoAdapter) Download(ctx context.Context, service providercore.Service, _ providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	link := strings.TrimSpace(candidate.SourceURL)
	if link == "" && candidate.FileID > 0 {
		link = animeToshoDownloadURL(candidate.FileID)
	}
	if link == "" {
		return providercore.Download{}, fmt.Errorf("%w: animetosho candidate has no attachment id", providercore.ErrProviderPrerequisiteMissing)
	}
	data, _, err := providerRequest(ctx, service, http.MethodGet, link, animeToshoKey, true, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	archiveName := strings.TrimSuffix(link, ".xz") + ".srt.xz"
	member, err := providercore.ExtractSubtitle(archiveName, data, security.ArchiveLimits{MaxBytes: providerReadLimit})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: link}, nil
}

func animeToshoEntries(ctx context.Context, service providercore.Service, config providercore.Config, episodeID string) ([]animeToshoEntry, error) {
	data, _, err := providerRequest(ctx, service, http.MethodGet, animeToshoFeedURL(config, url.Values{"eid": {episodeID}}), animeToshoKey, false, nil)
	if err != nil {
		return nil, err
	}
	var entries []animeToshoEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	complete := entries[:0]
	for _, entry := range entries {
		if entry.Status == "complete" {
			complete = append(complete, entry)
		}
	}
	sort.SliceStable(complete, func(i, j int) bool { return complete[i].Timestamp > complete[j].Timestamp })
	return complete, nil
}

func animeToshoTorrentFiles(ctx context.Context, service providercore.Service, config providercore.Config, id int64) (animeToshoTorrent, error) {
	data, _, err := providerRequest(ctx, service, http.MethodGet, animeToshoFeedURL(config, url.Values{"show": {"torrent"}, "id": {strconv.FormatInt(id, 10)}}), animeToshoKey, false, nil)
	if err != nil {
		return animeToshoTorrent{}, err
	}
	var torrent animeToshoTorrent
	if err := json.Unmarshal(data, &torrent); err != nil {
		return animeToshoTorrent{}, err
	}
	return torrent, nil
}

func animeToshoEpisodeID(request providercore.SearchRequest) string {
	for _, ids := range []map[string]string{request.MediaContext.EpisodeExternalIDs, request.MediaContext.ExternalIDs} {
		for _, key := range []string{"anidb_episode_id", "anidbEpisodeID", "anidb_eid", "eid"} {
			if value := strings.TrimSpace(ids[key]); value != "" {
				return value
			}
		}
	}
	return ""
}

func animeToshoFeedURL(config providercore.Config, values url.Values) string {
	base := strings.TrimRight(strings.TrimSpace(config.BaseURL), "/")
	if base == "" {
		base = "https://feed.animetosho.org"
	}
	endpoint := base + "/json"
	if len(values) > 0 {
		endpoint += "?" + values.Encode()
	}
	return endpoint
}

func animeToshoDownloadURL(id int64) string {
	hexID := fmt.Sprintf("%08x", id)
	return fmt.Sprintf("https://animetosho.org/storage/attach/%s/%d.xz", hexID, id)
}
