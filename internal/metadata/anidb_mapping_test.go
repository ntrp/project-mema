package metadata

import (
	"encoding/xml"
	"testing"
)

func TestAniDBOVAONAAndSpecialFixturesRemainFacts(t *testing.T) {
	for _, test := range []struct {
		fixture, animeType, episodeType string
		absolute                        bool
	}{
		{fixture: "anidb-ova.xml", animeType: "OVA", episodeType: "1", absolute: true},
		{fixture: "anidb-ona.xml", animeType: "Web", episodeType: "2", absolute: false},
	} {
		t.Run(test.animeType, func(t *testing.T) {
			var anime anidbAnime
			if err := xml.Unmarshal(readAniDBFixture(t, test.fixture), &anime); err != nil {
				t.Fatal(err)
			}
			details := anidbDetails(anime, nil)
			if details.Facts[1].Value != test.animeType || details.EpisodeFacts[0].OriginalType != test.episodeType {
				t.Fatalf("facts = %#v episodeFacts = %#v", details.Facts, details.EpisodeFacts)
			}
			if (details.EpisodeFacts[0].AbsoluteNumber != nil) != test.absolute {
				t.Fatalf("unexpected standard numbering: %#v", details.EpisodeFacts[0])
			}
		})
	}
}

func TestAniDBPrefixedEpisodeNumbers(t *testing.T) {
	for _, value := range []string{"1", "S1", "C1", "T1", "P1", "O1"} {
		var number anidbEpisodeNumber
		if err := xml.Unmarshal([]byte(`<epno type="4">`+value+`</epno>`), &number); err != nil {
			t.Fatalf("%s: %v", value, err)
		}
		if number.Type != "4" || number.Number != 1 {
			t.Fatalf("%s decoded as %#v", value, number)
		}
	}
}

func TestAniDBDetailsExposeTMDBMappingAndCleanOverviewLinks(t *testing.T) {
	var anime anidbAnime
	xmlValue := `<anime id="9541"><title>Example</title><description>Meet http://anidb.net/ch51116 [Mikasa] today.</description><resources><resource type="44"><externalentity><identifier>1429</identifier></externalentity></resource></resources></anime>`
	if err := xml.Unmarshal([]byte(xmlValue), &anime); err != nil {
		t.Fatal(err)
	}
	details := anidbDetails(anime, nil)
	if details.Overview == nil || *details.Overview != "Meet Mikasa today." {
		t.Fatalf("overview = %#v", details.Overview)
	}
	if len(details.ProviderMappings) != 1 || details.ProviderMappings[0].ProviderName != "tmdb" || details.ProviderMappings[0].ExternalID != "1429" {
		t.Fatalf("provider mappings = %#v", details.ProviderMappings)
	}
}

func TestAniDBMovieTypeIsPreserved(t *testing.T) {
	details := anidbDetails(anidbAnime{ID: 7, Type: "Movie", Titles: []anidbTitle{{Type: "main", Value: "Example Movie"}}}, nil)
	if details.Type != "movie" {
		t.Fatalf("type = %q, want movie", details.Type)
	}
}

func TestAniDBDetailsPopulateAvailableResourceFields(t *testing.T) {
	var anime anidbAnime
	xmlValue := `<anime id="1"><type>TV Series</type><episodecount>2</episodecount><startdate>2024-01-01</startdate><enddate>2024-03-01</enddate><titles><title type="main">Example</title></titles><relatedanime><anime id="2" type="sequel">Example 2</anime></relatedanime><similaranime><anime id="3" approval="15" total="20">Similar Example</anime></similaranime><ratings><permanent count="100">8.25</permanent><temporary count="20">7.5</temporary></ratings><categories><category><name>Action</name></category></categories><tags><tag weight="400"><name>military</name></tag><tag weight="300"><name>minor theme</name></tag><tag weight="600" globalspoiler="true"><name>spoiler</name></tag></tags><characters><character id="10"><name>Hero</name><seiyuu id="20" picture="actor.jpg">Actor</seiyuu></character></characters><creators><name id="30" type="Direction">Director</name></creators><episodes><episode id="100"><epno type="1">1</epno><length>24</length><airdate>2024-01-01</airdate><title xml:lang="en">Pilot</title></episode><episode id="101"><epno type="2">1</epno><airdate>2024-01-02</airdate><title xml:lang="en">Special</title></episode></episodes></anime>`
	if err := xml.Unmarshal([]byte(xmlValue), &anime); err != nil {
		t.Fatal(err)
	}
	details := anidbDetails(anime, nil)
	if len(details.Genres) != 1 || details.Genres[0] != "Action" || len(details.Keywords) != 1 || details.Keywords[0] != "military" {
		t.Fatalf("genres=%#v keywords=%#v", details.Genres, details.Keywords)
	}
	if len(details.Seasons) != 2 || details.SeasonCount == nil || *details.SeasonCount != 1 || details.RuntimeMinutes == nil || *details.RuntimeMinutes != 24 {
		t.Fatalf("season/runtime details = %#v", details)
	}
	if len(details.Cast) != 1 || details.Cast[0].Name != "Actor" || details.Cast[0].Role == nil || *details.Cast[0].Role != "Hero" {
		t.Fatalf("cast = %#v", details.Cast)
	}
	if len(details.Crew) != 1 || details.Crew[0].Name != "Director" || details.Status == nil || *details.Status != "Ended" {
		t.Fatalf("crew/status = %#v %#v", details.Crew, details.Status)
	}
	if details.Type != "serie" || details.VoteAverage == nil || *details.VoteAverage != 8.25 {
		t.Fatalf("type/rating = %q %#v", details.Type, details.VoteAverage)
	}
	if len(details.Recommendations) != 1 || details.Recommendations[0].Title != "Example 2" || details.Recommendations[0].ExternalID != "2" {
		t.Fatalf("related = %#v", details.Recommendations)
	}
	if len(details.Similar) != 1 || details.Similar[0].Title != "Similar Example" || details.Similar[0].ExternalID != "3" {
		t.Fatalf("similar = %#v", details.Similar)
	}
}
