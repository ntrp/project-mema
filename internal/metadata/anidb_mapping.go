package metadata

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type anidbAnime struct {
	ID           int              `xml:"id,attr"`
	Type         string           `xml:"type"`
	EpisodeCount int32            `xml:"episodecount"`
	StartDate    string           `xml:"startdate"`
	EndDate      string           `xml:"enddate"`
	Description  string           `xml:"description"`
	Picture      string           `xml:"picture"`
	Titles       []anidbTitle     `xml:"titles>title"`
	Categories   []anidbNamedFact `xml:"categories>category"`
	Tags         []anidbNamedFact `xml:"tags>tag"`
	Characters   []anidbCharacter `xml:"characters>character"`
	Creators     []anidbCreator   `xml:"creators>name"`
	Resources    []anidbResource  `xml:"resources>resource"`
	RelatedAnime []anidbRelated   `xml:"relatedanime>anime"`
	SimilarAnime []anidbSimilar   `xml:"similaranime>anime"`
	Ratings      anidbRatings     `xml:"ratings"`
	Episodes     []anidbEpisode   `xml:"episodes>episode"`
}

type anidbTitle struct {
	Value    string `xml:",chardata"`
	Language string `xml:"lang,attr"`
	Type     string `xml:"type,attr"`
}

type anidbEpisode struct {
	ID      int                `xml:"id,attr"`
	EpNo    anidbEpisodeNumber `xml:"epno"`
	Length  int32              `xml:"length"`
	AirDate string             `xml:"airdate"`
	Titles  []anidbTitle       `xml:"title"`
}

type anidbEpisodeNumber struct {
	Type   string
	Number int32
}

func (number *anidbEpisodeNumber) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var value struct {
		Type string `xml:"type,attr"`
		Text string `xml:",chardata"`
	}
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	text := strings.TrimSpace(value.Text)
	if len(text) > 1 && strings.Contains("SCTPO", text[:1]) {
		text = text[1:]
	}
	parsed, err := strconv.ParseInt(text, 10, 32)
	if err != nil {
		return err
	}
	number.Type = value.Type
	number.Number = int32(parsed)
	return nil
}

func anidbMainTitle(anime anidbAnime) string {
	for _, title := range anime.Titles {
		if title.Type == "main" && strings.TrimSpace(title.Value) != "" {
			return strings.TrimSpace(title.Value)
		}
	}
	for _, title := range anime.Titles {
		if strings.TrimSpace(title.Value) != "" {
			return strings.TrimSpace(title.Value)
		}
	}
	return ""
}

func anidbSearchResult(anime anidbAnime, posterPath *string) SearchResult {
	id := strconv.Itoa(anime.ID)
	return SearchResult{
		Title:            anidbMainTitle(anime),
		Type:             anidbMediaType(anime.Type),
		Year:             yearFromString(anime.StartDate),
		ExternalProvider: "anidb",
		ExternalID:       id,
		ExternalURL:      optionalString("https://anidb.net/anime/" + id),
		Overview:         optionalString(anidbOverview(anime.Description)),
		PosterPath:       posterPath,
		ReleaseDate:      optionalString(anime.StartDate),
	}
}

func anidbDetails(anime anidbAnime, posterPath *string) Details {
	result := anidbSearchResult(anime, posterPath)
	details := Details{
		Title: result.Title, Type: result.Type, Year: result.Year,
		ExternalProvider: result.ExternalProvider, ExternalID: result.ExternalID,
		ExternalURL: result.ExternalURL, Overview: result.Overview, PosterPath: result.PosterPath,
		Status: anidbStatus(anime.EndDate), OriginalLanguage: optionalString("ja"),
		ReleaseDate: result.ReleaseDate, FirstAirDate: result.ReleaseDate,
		RuntimeMinutes: anidbRuntime(anime.Episodes),
		Facts: []Fact{
			{Label: "AniDB ID", Value: result.ExternalID},
			{Label: "AniDB type", Value: strings.TrimSpace(anime.Type)},
		},
		Aliases:          anidbAliases(anime.Titles),
		ProviderMappings: anidbProviderMappings(anime.Resources),
		Genres:           anidbNames(anime.Categories),
		Keywords:         anidbTags(anime.Tags),
		Seasons:          anidbSeasons(anime.Episodes),
		Cast:             anidbCast(anime.Characters),
		Crew:             anidbCrew(anime.Creators),
		EpisodeFacts:     anidbEpisodeFacts(anime.Episodes),
		Recommendations:  anidbRelatedResults(anime.RelatedAnime, result.Type),
		Similar:          anidbSimilarResults(anime.SimilarAnime, result.Type),
	}
	if rating := anidbVoteAverage(anime.Ratings); rating != nil {
		details.VoteAverage = rating
	}
	if anime.EpisodeCount > 0 {
		details.EpisodeCount = &anime.EpisodeCount
	}
	details.SeasonCount = anidbSeasonCount(details.Seasons)
	if strings.TrimSpace(anime.EndDate) != "" {
		details.Facts = append(details.Facts, Fact{Label: "AniDB end date", Value: strings.TrimSpace(anime.EndDate)})
	}
	for _, episode := range details.EpisodeFacts {
		details.Facts = append(details.Facts, Fact{Label: "AniDB episode " + episode.ExternalID, Value: anidbEpisodeFactValue(episode)})
	}
	return details
}

func anidbAliases(titles []anidbTitle) []Alias {
	aliases := make([]Alias, 0, len(titles))
	for _, title := range titles {
		value := strings.TrimSpace(title.Value)
		if value == "" {
			continue
		}
		language := strings.TrimSpace(title.Language)
		aliases = append(aliases, Alias{Value: value, Language: optionalString(language), Kind: anidbAliasKind(title)})
	}
	return aliases
}

func anidbAliasKind(title anidbTitle) string {
	if title.Type == "syn" || title.Type == "short" {
		return "synonym"
	}
	switch strings.ToLower(title.Language) {
	case "ja":
		return "native"
	case "en":
		return "english"
	case "x-jat":
		return "romaji"
	default:
		return "synonym"
	}
}

func anidbEpisodeFacts(episodes []anidbEpisode) []EpisodeFact {
	facts := make([]EpisodeFact, 0, len(episodes))
	for _, episode := range episodes {
		if episode.ID <= 0 || episode.EpNo.Number <= 0 || strings.TrimSpace(episode.EpNo.Type) == "" {
			continue
		}
		fact := EpisodeFact{
			ExternalID: strconv.Itoa(episode.ID), OriginalType: episode.EpNo.Type,
			OriginalNumber: episode.EpNo.Number, Title: anidbEpisodeTitle(episode.Titles),
			AirDate: optionalString(episode.AirDate),
		}
		if episode.EpNo.Type == "1" {
			absolute := episode.EpNo.Number
			fact.AbsoluteNumber = &absolute
		}
		facts = append(facts, fact)
	}
	return facts
}

func anidbEpisodeTitle(titles []anidbTitle) string {
	for _, language := range []string{"en", "x-jat", "ja"} {
		for _, title := range titles {
			if title.Language == language && strings.TrimSpace(title.Value) != "" {
				return strings.TrimSpace(title.Value)
			}
		}
	}
	if len(titles) > 0 {
		return strings.TrimSpace(titles[0].Value)
	}
	return ""
}

func anidbEpisodeFactValue(fact EpisodeFact) string {
	parts := []string{"type=" + fact.OriginalType, "number=" + strconv.Itoa(int(fact.OriginalNumber))}
	if fact.Title != "" {
		parts = append(parts, "title="+fact.Title)
	}
	if fact.AirDate != nil {
		parts = append(parts, "airdate="+*fact.AirDate)
	}
	return strings.Join(parts, "; ")
}
