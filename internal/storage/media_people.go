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
	Keywords         []string
	Facts            []MediaFact
	Seasons          []MediaSeason
	Cast             []MediaPerson
	Crew             []MediaPerson
	Recommendations  []MediaRelatedItem
	Similar          []MediaRelatedItem
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
	Monitored    bool           `json:"monitored"`
	Episodes     []MediaEpisode `json:"episodes,omitempty"`
}

type MediaEpisode struct {
	Name          string  `json:"name"`
	EpisodeNumber int32   `json:"episodeNumber"`
	Overview      *string `json:"overview,omitempty"`
	AirDate       *string `json:"airDate,omitempty"`
	StillPath     *string `json:"stillPath,omitempty"`
	Monitored     bool    `json:"monitored"`
}

type MediaPerson struct {
	ExternalProvider *string `json:"externalProvider,omitempty"`
	ExternalID       *string `json:"externalId,omitempty"`
	Name             string  `json:"name"`
	Role             *string `json:"role,omitempty"`
	ProfilePath      *string `json:"profilePath,omitempty"`
}

type MediaRelatedItem struct {
	Title            string  `json:"title"`
	Type             string  `json:"type"`
	Year             *int32  `json:"year,omitempty"`
	ExternalProvider string  `json:"externalProvider"`
	ExternalID       string  `json:"externalId"`
	Overview         *string `json:"overview,omitempty"`
	PosterPath       *string `json:"posterPath,omitempty"`
}

type mediaMetadataPayloads struct {
	genres          []byte
	keywords        []byte
	facts           []byte
	seasons         []byte
	cast            []byte
	crew            []byte
	recommendations []byte
	similar         []byte
}

func marshalMediaMetadata(input MediaMetadataSnapshot) (mediaMetadataPayloads, error) {
	genres, err := marshalJSONArray(input.Genres)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	keywords, err := marshalJSONArray(input.Keywords)
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
	crew, err := marshalJSONArray(input.Crew)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	recommendations, err := marshalJSONArray(input.Recommendations)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	similar, err := marshalJSONArray(input.Similar)
	if err != nil {
		return mediaMetadataPayloads{}, err
	}
	return mediaMetadataPayloads{
		genres:          genres,
		keywords:        keywords,
		facts:           facts,
		seasons:         seasons,
		cast:            cast,
		crew:            crew,
		recommendations: recommendations,
		similar:         similar,
	}, nil
}

func scanMediaMetadata(
	target *MediaMetadataSnapshot,
	genres []byte,
	keywords []byte,
	facts []byte,
	seasons []byte,
	cast []byte,
	crew []byte,
	recommendations []byte,
	similar []byte,
) {
	target.Genres = unmarshalJSONArray[string](genres)
	target.Keywords = unmarshalJSONArray[string](keywords)
	target.Facts = unmarshalJSONArray[MediaFact](facts)
	target.Seasons = unmarshalJSONArray[MediaSeason](seasons)
	target.Cast = unmarshalJSONArray[MediaPerson](cast)
	target.Crew = unmarshalJSONArray[MediaPerson](crew)
	target.Recommendations = unmarshalJSONArray[MediaRelatedItem](recommendations)
	target.Similar = unmarshalJSONArray[MediaRelatedItem](similar)
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
