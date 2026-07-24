package metadata

import (
	"strconv"
	"strings"
)

type anidbAnime struct {
	ID           int            `xml:"id,attr"`
	Type         string         `xml:"type"`
	EpisodeCount int32          `xml:"episodecount"`
	StartDate    string         `xml:"startdate"`
	EndDate      string         `xml:"enddate"`
	Description  string         `xml:"description"`
	Picture      string         `xml:"picture"`
	Titles       []anidbTitle   `xml:"titles>title"`
	Episodes     []anidbEpisode `xml:"episodes>episode"`
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
	Type   string `xml:"type,attr"`
	Number int32  `xml:",chardata"`
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

func anidbSearchResult(anime anidbAnime) SearchResult {
	id := strconv.Itoa(anime.ID)
	return SearchResult{
		Title:            anidbMainTitle(anime),
		Type:             "serie",
		Year:             yearFromString(anime.StartDate),
		ExternalProvider: "anidb",
		ExternalID:       id,
		ExternalURL:      optionalString("https://anidb.net/anime/" + id),
		Overview:         optionalString(anime.Description),
		PosterPath:       anidbPicture(anime.Picture),
		ReleaseDate:      optionalString(anime.StartDate),
	}
}

func anidbDetails(anime anidbAnime) Details {
	result := anidbSearchResult(anime)
	details := Details{
		Title: result.Title, Type: result.Type, Year: result.Year,
		ExternalProvider: result.ExternalProvider, ExternalID: result.ExternalID,
		ExternalURL: result.ExternalURL, Overview: result.Overview, PosterPath: result.PosterPath,
		ReleaseDate: result.ReleaseDate,
		Facts: []Fact{
			{Label: "AniDB ID", Value: result.ExternalID},
			{Label: "AniDB type", Value: strings.TrimSpace(anime.Type)},
		},
		Aliases:      anidbAliases(anime.Titles),
		EpisodeFacts: anidbEpisodeFacts(anime.Episodes),
	}
	if anime.EpisodeCount > 0 {
		details.EpisodeCount = &anime.EpisodeCount
	}
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

func anidbPicture(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return &value
	}
	value = "https://cdn-eu.anidb.net/images/main/" + strings.TrimLeft(value, "/")
	return &value
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
