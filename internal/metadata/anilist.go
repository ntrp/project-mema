package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type anilistResponse struct {
	Data   *anilistData `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type anilistData struct {
	Page  *anilistPage  `json:"Page"`
	Media *anilistMedia `json:"Media"`
}

type anilistPage struct {
	Media []anilistMedia `json:"media"`
}

type anilistMedia struct {
	Raw      json.RawMessage `json:"-"`
	ID       int             `json:"id"`
	IDMal    int             `json:"idMal"`
	Type     string          `json:"type"`
	Format   string          `json:"format"`
	Synonyms []string        `json:"synonyms"`
	Title    struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
		Native  string `json:"native"`
	} `json:"title"`
	Description     string                         `json:"description"`
	CountryOfOrigin string                         `json:"countryOfOrigin"`
	StartDate       struct{ Year, Month, Day int } `json:"startDate"`
	CoverImage      struct {
		Large      string `json:"large"`
		ExtraLarge string `json:"extraLarge"`
	} `json:"coverImage"`
	ExternalLinks []struct {
		Site string `json:"site"`
		URL  string `json:"url"`
	} `json:"externalLinks"`
	SiteURL  string `json:"siteUrl"`
	Duration int    `json:"duration"`
}

func (m *anilistMedia) UnmarshalJSON(payload []byte) error {
	type media anilistMedia
	if err := json.Unmarshal(payload, (*media)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0], payload...)
	return nil
}

type anilistQuery struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

const anilistFields = `id idMal type format title{romaji english native} synonyms description countryOfOrigin startDate{year month day} coverImage{large extraLarge} siteUrl duration externalLinks{site url}`

func (s *Service) anilist(ctx context.Context, config Config, query string, vars map[string]any, out *anilistResponse) error {
	endpoint := strings.TrimRight(config.BaseURL, "/")
	if endpoint == "" {
		endpoint = "https://graphql.anilist.co"
	}
	if err := s.doJSON(ctx, config, http.MethodPost, endpoint, anilistQuery{Query: query, Variables: vars}, out); err != nil {
		var httpErr ProviderHTTPError
		if errors.As(err, &httpErr) || IsRateLimited(err) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return ProviderPayloadError{Message: err.Error()}
	}
	if len(out.Errors) > 0 {
		message := strings.TrimSpace(out.Errors[0].Message)
		if message == "" {
			message = "unknown error"
		}
		return ProviderGraphQLError{Message: message}
	}
	return nil
}

func anilistType(format string) string {
	switch format {
	case "MOVIE":
		return "movie"
	case "TV", "TV_SHORT", "ONA", "OVA", "SPECIAL":
		return "serie"
	default:
		return ""
	}
}

func anilistDate(m anilistMedia) string {
	if m.StartDate.Year <= 0 || m.StartDate.Month <= 0 || m.StartDate.Day <= 0 {
		return ""
	}
	date := time.Date(m.StartDate.Year, time.Month(m.StartDate.Month), m.StartDate.Day, 0, 0, 0, 0, time.UTC)
	if date.Year() != m.StartDate.Year || int(date.Month()) != m.StartDate.Month || date.Day() != m.StartDate.Day {
		return ""
	}
	return date.Format("2006-01-02")
}

func anilistTitle(m anilistMedia) string {
	return firstNonEmpty(m.Title.English, m.Title.Romaji, m.Title.Native)
}

func (s *Service) searchAniList(ctx context.Context, config Config, req SearchRequest) ([]SearchResult, error) {
	if req.MediaType != "movie" && req.MediaType != "serie" {
		return nil, nil
	}
	query := `query($search:String,$type:MediaType){Page(perPage:10){media(search:$search,type:$type){` + anilistFields + `}}}`
	var payload anilistResponse
	if err := s.anilist(ctx, config, query, map[string]any{"search": req.Query, "type": "ANIME"}, &payload); err != nil {
		return nil, err
	}
	if payload.Data == nil || payload.Data.Page == nil {
		return nil, ProviderPayloadError{Message: "missing search page"}
	}
	results := []SearchResult{}
	for _, media := range payload.Data.Page.Media {
		mediaType := anilistType(media.Format)
		if media.ID <= 0 || media.Type != "ANIME" || mediaType != req.MediaType || anilistTitle(media) == "" {
			continue
		}
		result := SearchResult{
			Title:            anilistTitle(media),
			Aliases:          anilistAliases(media),
			Type:             mediaType,
			ExternalProvider: "anilist",
			ExternalID:       strconv.Itoa(media.ID),
			ExternalURL:      optionalString(media.SiteURL),
			Overview:         optionalString(media.Description),
			PosterPath:       optionalString(firstNonEmpty(media.CoverImage.ExtraLarge, media.CoverImage.Large)),
			OriginalLanguage: nativeLanguage(media.CountryOfOrigin),
		}
		if media.StartDate.Year > 0 {
			year := int32(media.StartDate.Year)
			result.Year = &year
		}
		if date := anilistDate(media); date != "" {
			result.ReleaseDate = &date
		}
		results = append(results, result)
	}
	return results, nil
}

func (s *Service) detailsAniList(ctx context.Context, config Config, req DetailsRequest) (Details, error) {
	id, err := strconv.Atoi(strings.TrimSpace(req.ExternalID))
	if err != nil || id <= 0 {
		return Details{}, ProviderPayloadError{Message: "invalid AniList media id"}
	}
	query := `query($id:Int){Media(id:$id){` + anilistFields + `}}`
	var payload anilistResponse
	if err := s.anilist(ctx, config, query, map[string]any{"id": id}, &payload); err != nil {
		return Details{}, err
	}
	if payload.Data == nil || payload.Data.Media == nil || payload.Data.Media.ID <= 0 {
		return Details{}, ProviderPayloadError{Message: "missing media"}
	}
	media := *payload.Data.Media
	if media.Type != "ANIME" {
		return Details{}, ProviderPayloadError{Message: "media is not anime"}
	}
	mediaType := anilistType(media.Format)
	if mediaType == "" || (req.MediaType != "" && req.MediaType != mediaType) || anilistTitle(media) == "" {
		return Details{}, ProviderPayloadError{Message: "unsupported or untitled anime media"}
	}
	details := Details{
		Title:            anilistTitle(media),
		Aliases:          anilistAliases(media),
		ProviderMappings: anilistMappings(media, mediaType),
		Type:             mediaType,
		ExternalProvider: "anilist",
		ExternalID:       strconv.Itoa(media.ID),
		ExternalURL:      optionalString(media.SiteURL),
		Overview:         optionalString(media.Description),
		PosterPath:       optionalString(firstNonEmpty(media.CoverImage.ExtraLarge, media.CoverImage.Large)),
		OriginalLanguage: nativeLanguage(media.CountryOfOrigin),
	}
	if media.StartDate.Year > 0 {
		year := int32(media.StartDate.Year)
		details.Year = &year
	}
	if date := anilistDate(media); date != "" {
		details.ReleaseDate = &date
		if mediaType == "serie" {
			details.FirstAirDate = &date
		}
	}
	if media.Duration > 0 {
		duration := int32(media.Duration)
		details.RuntimeMinutes = &duration
	}
	return details, nil
}
