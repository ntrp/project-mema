package avistazapi

import "encoding/json"

type Options struct {
	Name            string
	DefaultBaseURL  string
	AnimeCategories bool
	PreferRelease   bool
}

type apiResponse struct {
	Data []apiRelease `json:"data"`
}

type authResponse struct {
	Token string `json:"token"`
}

type apiRelease struct {
	URL              string             `json:"url"`
	Download         string             `json:"download"`
	Category         map[string]string  `json:"category"`
	CreatedAtISO     string             `json:"created_at_iso"`
	FileName         string             `json:"file_name"`
	ReleaseTitle     string             `json:"release_title"`
	InfoHash         string             `json:"info_hash"`
	Leech            nullableInt32      `json:"leech"`
	Completed        nullableInt32      `json:"completed"`
	Seed             nullableInt32      `json:"seed"`
	FileSize         nullableInt64      `json:"file_size"`
	FileCount        nullableInt32      `json:"file_count"`
	DownloadMultiply nullableFloat      `json:"download_multiply"`
	UploadMultiply   nullableFloat      `json:"upload_multiply"`
	VideoQuality     string             `json:"video_quality"`
	Type             string             `json:"type"`
	Format           string             `json:"format"`
	MovieTV          *apiExternalIDInfo `json:"movie_tv"`
}

type apiExternalIDInfo struct {
	TMDB string `json:"tmdb"`
	TVDB string `json:"tvdb"`
	IMDB string `json:"imdb"`
}

type nullableInt32 struct {
	Value int32
	Valid bool
}

func (n *nullableInt32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		n.Valid = false
		return nil
	}
	var value int32
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	n.Value = value
	n.Valid = true
	return nil
}

type nullableInt64 struct {
	Value int64
	Valid bool
}

func (n *nullableInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		n.Valid = false
		return nil
	}
	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	n.Value = value
	n.Valid = true
	return nil
}

type nullableFloat struct {
	Value float64
	Valid bool
}

func (n *nullableFloat) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		n.Valid = false
		return nil
	}
	var value float64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	n.Value = value
	n.Valid = true
	return nil
}
