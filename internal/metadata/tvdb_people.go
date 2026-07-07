package metadata

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type tvdbPersonDetailsResponse struct {
	Data tvdbPersonDetails `json:"data"`
}

type tvdbPersonDetails struct {
	ID          tvdbStringNumber      `json:"id"`
	Name        string                `json:"name"`
	Image       string                `json:"image"`
	Birth       string                `json:"birth"`
	BirthPlace  string                `json:"birthPlace"`
	Death       string                `json:"death"`
	Aliases     []tvdbAlias           `json:"aliases"`
	Biographies []tvdbBiography       `json:"biographies"`
	Characters  []tvdbPersonCharacter `json:"characters"`
}

type tvdbAlias struct {
	Name string `json:"name"`
}

type tvdbBiography struct {
	Biography string `json:"biography"`
	Language  string `json:"language"`
}

type tvdbPersonCharacter struct {
	Name       string         `json:"name"`
	MovieID    tvdbIntPointer `json:"movieId"`
	Movie      tvdbRecordInfo `json:"movie"`
	SeriesID   tvdbIntPointer `json:"seriesId"`
	Series     tvdbRecordInfo `json:"series"`
	PeopleType string         `json:"peopleType"`
}

type tvdbRecordInfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Year  string `json:"year"`
}

func (s *Service) personDetailsTVDB(ctx context.Context, config Config, personID string) (PersonDetails, error) {
	token, err := s.tvdbToken(ctx, config)
	if err != nil {
		return PersonDetails{}, err
	}
	config.AccessToken = &token
	personID = strings.TrimSpace(personID)
	if personID == "" {
		return PersonDetails{}, ErrUnsupportedProvider
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "people", personID, "extended")
	if err != nil {
		return PersonDetails{}, err
	}
	values := url.Values{}
	values.Set("meta", "translations")
	endpoint = endpoint + "?" + values.Encode()

	var payload tvdbPersonDetailsResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return PersonDetails{}, err
	}
	details := tvdbPersonDetailsResult(payload.Data, personID)
	if details.Name == "" {
		return PersonDetails{}, errors.New("TVDB person details did not include a name")
	}
	return details, nil
}

func tvdbPersonDetailsResult(item tvdbPersonDetails, personID string) PersonDetails {
	if value := item.ID.String(); value != "" {
		personID = value
	}
	return PersonDetails{
		ID:           personID,
		Name:         strings.TrimSpace(item.Name),
		Biography:    optionalString(tvdbPersonBiography(item.Biographies)),
		Birthday:     optionalString(item.Birth),
		Deathday:     optionalString(item.Death),
		PlaceOfBirth: optionalString(item.BirthPlace),
		ProfilePath:  optionalString(item.Image),
		AlsoKnownAs:  tvdbPersonAliases(item.Aliases),
		Appearances:  tvdbPersonAppearances(item.Characters),
	}
}

func tvdbPersonBiography(items []tvdbBiography) string {
	for _, item := range items {
		if strings.EqualFold(item.Language, "eng") {
			return strings.TrimSpace(item.Biography)
		}
	}
	for _, item := range items {
		if value := strings.TrimSpace(item.Biography); value != "" {
			return value
		}
	}
	return ""
}

func tvdbPersonAliases(items []tvdbAlias) []string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		if value := strings.TrimSpace(item.Name); value != "" {
			values = append(values, value)
		}
	}
	return values
}

func tvdbPersonAppearances(characters []tvdbPersonCharacter) []PersonAppearance {
	appearances := make([]PersonAppearance, 0, len(characters))
	seen := map[string]int{}
	for _, character := range characters {
		appendAppearance(&appearances, seen, tvdbCharacterAppearance(character))
	}
	return appearances
}

func tvdbCharacterAppearance(character tvdbPersonCharacter) PersonAppearance {
	if character.MovieID.Value != nil && *character.MovieID.Value > 0 {
		return tvdbRecordAppearance("movie", *character.MovieID.Value, character.Movie, character.Name)
	}
	if character.SeriesID.Value != nil && *character.SeriesID.Value > 0 {
		return tvdbRecordAppearance("serie", *character.SeriesID.Value, character.Series, character.Name)
	}
	return PersonAppearance{}
}

func tvdbRecordAppearance(mediaType string, id int64, record tvdbRecordInfo, role string) PersonAppearance {
	return PersonAppearance{
		Title:            strings.TrimSpace(record.Name),
		Type:             mediaType,
		Year:             yearFromString(record.Year),
		ExternalProvider: "tvdb",
		ExternalID:       strconv.FormatInt(id, 10),
		PosterPath:       optionalString(record.Image),
		Role:             optionalString(role),
	}
}
