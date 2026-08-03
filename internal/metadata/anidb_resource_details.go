package metadata

import (
	"sort"
	"strconv"
	"strings"
)

const anidbTwoStarTagWeight = 400

type anidbNamedFact struct {
	Name          string `xml:"name"`
	Weight        int    `xml:"weight,attr"`
	LocalSpoiler  bool   `xml:"localspoiler,attr"`
	GlobalSpoiler bool   `xml:"globalspoiler,attr"`
}

type anidbCharacter struct {
	ID      int           `xml:"id,attr"`
	Name    string        `xml:"name"`
	Picture string        `xml:"picture"`
	Seiyuu  []anidbSeiyuu `xml:"seiyuu"`
}

type anidbSeiyuu struct {
	ID      int    `xml:"id,attr"`
	Picture string `xml:"picture,attr"`
	Name    string `xml:",chardata"`
}

type anidbCreator struct {
	ID   int    `xml:"id,attr"`
	Type string `xml:"type,attr"`
	Name string `xml:",chardata"`
}

func anidbNames(values []anidbNamedFact) []string {
	return anidbNamesWithMinimumWeight(values, 0)
}

func anidbTags(values []anidbNamedFact) []string {
	return anidbNamesWithMinimumWeight(values, anidbTwoStarTagWeight)
}

func anidbNamesWithMinimumWeight(values []anidbNamedFact, minimumWeight int) []string {
	items := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		name := strings.TrimSpace(value.Name)
		if name == "" || value.Weight < minimumWeight || value.LocalSpoiler || value.GlobalSpoiler || seen[name] {
			continue
		}
		seen[name] = true
		items = append(items, name)
	}
	sort.Strings(items)
	return items
}

func anidbSeasons(episodes []anidbEpisode) []Season {
	bySeason := map[int32][]Episode{}
	for _, episode := range episodes {
		seasonNumber := int32(1)
		if episode.EpNo.Type == "2" {
			seasonNumber = 0
		} else if episode.EpNo.Type != "1" {
			continue
		}
		bySeason[seasonNumber] = append(bySeason[seasonNumber], Episode{
			Name: anidbEpisodeTitle(episode.Titles), EpisodeNumber: episode.EpNo.Number,
			AirDate: optionalString(episode.AirDate),
		})
	}
	seasons := make([]Season, 0, len(bySeason))
	for _, number := range []int32{0, 1} {
		episodes := bySeason[number]
		if len(episodes) == 0 {
			continue
		}
		count := int32(len(episodes))
		name := "Season 1"
		if number == 0 {
			name = "Specials"
		}
		seasons = append(seasons, Season{
			Name: name, SeasonNumber: number, EpisodeCount: &count,
			AirDate: episodes[0].AirDate, Episodes: episodes,
		})
	}
	return seasons
}

func anidbRuntime(episodes []anidbEpisode) *int32 {
	for _, episode := range episodes {
		if episode.EpNo.Type == "1" && episode.Length > 0 {
			return &episode.Length
		}
	}
	return nil
}

func anidbStatus(endDate string) *string {
	if strings.TrimSpace(endDate) == "" {
		return optionalString("Continuing")
	}
	return optionalString("Ended")
}

func anidbSeasonCount(seasons []Season) *int32 {
	var count int32
	for _, season := range seasons {
		if season.SeasonNumber > 0 {
			count++
		}
	}
	if count == 0 {
		return nil
	}
	return &count
}

func anidbCast(characters []anidbCharacter) []Person {
	items := []Person{}
	for _, character := range characters {
		role := strings.TrimSpace(character.Name)
		for _, actor := range character.Seiyuu {
			name := strings.TrimSpace(actor.Name)
			if name == "" {
				continue
			}
			id, provider := strconv.Itoa(actor.ID), "anidb"
			items = append(items, Person{
				ExternalProvider: &provider, ExternalID: &id, Name: name,
				Role: optionalString(role),
			})
		}
	}
	return items
}

func anidbCrew(creators []anidbCreator) []Person {
	items := make([]Person, 0, len(creators))
	for _, creator := range creators {
		name := strings.TrimSpace(creator.Name)
		if name == "" {
			continue
		}
		id, provider := strconv.Itoa(creator.ID), "anidb"
		items = append(items, Person{
			ExternalProvider: &provider, ExternalID: &id, Name: name,
			Role: optionalString(strings.TrimSpace(creator.Type)),
		})
	}
	return items
}
