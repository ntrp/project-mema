package speedappapi

import (
	"encoding/json"
	"strconv"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

type torrent struct {
	ID                   int     `json:"id"`
	URL                  string  `json:"url"`
	Name                 string  `json:"name"`
	ShortDescription     string  `json:"short_description"`
	Size                 int64   `json:"size"`
	CreatedAt            string  `json:"created_at"`
	TimesCompleted       int     `json:"times_completed"`
	Leechers             int     `json:"leechers"`
	Seeders              int     `json:"seeders"`
	Poster               string  `json:"poster"`
	IMDBID               string  `json:"imdb_id"`
	DownloadVolumeFactor float64 `json:"download_volume_factor"`
	UploadVolumeFactor   float64 `json:"upload_volume_factor"`
	Category             struct {
		ID int `json:"id"`
	} `json:"category"`
}

func parseReleases(config engine.Config, options Options, body []byte) ([]engine.Release, error) {
	var decoded []torrent
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}
	releases := make([]engine.Release, 0, len(decoded))
	for _, item := range decoded {
		release := item.toRelease(config, options)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func (t torrent) toRelease(config engine.Config, options Options) engine.Release {
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           cleanTitle(t.Name),
		DownloadURL:     downloadURL(config, options, t.ID),
		InfoURL:         t.URL,
		GUID:            engine.FirstNonEmpty(t.URL, strconv.Itoa(t.ID)),
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(t.Seeders),
		Peers:           common.Int32Ptr(t.Seeders + t.Leechers),
		PublishedAt:     common.ParseFlexibleTime(t.CreatedAt),
	}
}

func downloadURL(config engine.Config, options Options, id int) string {
	if id == 0 {
		return ""
	}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, options.DefaultBaseURL), "/api/torrent/"+strconv.Itoa(id)+"/download", nil)
	if err != nil {
		return ""
	}
	return endpoint
}
