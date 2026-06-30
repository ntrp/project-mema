package storage

import "encoding/json"

type MediaMetadataSnapshot struct {
	CollectionID     *string
	CollectionName   *string
	BackdropPath     *string
	MetadataStatus   *string
	OriginalLanguage *string
	ReleaseDate      *string
	FirstAirDate     *string
	RuntimeMinutes   *int32
	SeasonCount      *int32
	EpisodeCount     *int32
	VoteAverage      *float64
	Genres           []string
	Facts            []MediaFact
	Seasons          []MediaSeason
	Cast             []MediaPerson
}

type MediaFact struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type MediaSeason struct {
	Name         string         `json:"name"`
	EpisodeCount *int32         `json:"episodeCount,omitempty"`
	AirDate      *string        `json:"airDate,omitempty"`
	PosterPath   *string        `json:"posterPath,omitempty"`
	Episodes     []MediaEpisode `json:"episodes,omitempty"`
}

type MediaEpisode struct {
	Name          string  `json:"name"`
	EpisodeNumber int32   `json:"episodeNumber"`
	Overview      *string `json:"overview,omitempty"`
	AirDate       *string `json:"airDate,omitempty"`
	StillPath     *string `json:"stillPath,omitempty"`
}

type MediaPerson struct {
	Name        string  `json:"name"`
	Role        *string `json:"role,omitempty"`
	ProfilePath *string `json:"profilePath,omitempty"`
}

type mediaMetadataPayloads struct {
	genres  []byte
	facts   []byte
	seasons []byte
	cast    []byte
}

func marshalMediaMetadata(input MediaMetadataSnapshot) (mediaMetadataPayloads, error) {
	genres, err := marshalJSONArray(input.Genres)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	facts, err := marshalJSONArray(input.Facts)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	seasons, err := marshalJSONArray(input.Seasons)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	cast, err := marshalJSONArray(input.Cast)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	return mediaMetadataPayloads{genres: genres, facts: facts, seasons: seasons, cast: cast}, nil
}

func scanMediaMetadata(
	target *MediaMetadataSnapshot,
	genres []byte,
	facts []byte,
	seasons []byte,
	cast []byte,
) {
	target.Genres = unmarshalJSONArray[string](genres)
	target.Facts = unmarshalJSONArray[MediaFact](facts)
	target.Seasons = unmarshalJSONArray[MediaSeason](seasons)
	target.Cast = unmarshalJSONArray[MediaPerson](cast)
}

func marshalJSONArray[T any](values []T) ([]byte, error) {
	if values == nil {
		values = []T{}
	}
	return json.Marshal(values)
}

func unmarshalJSONArray[T any](payload []byte) []T {
	values := []T{}
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &values)
	}
	return values
}
