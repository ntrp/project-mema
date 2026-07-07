package metadata

type tvdbDetailsResponse struct {
	Data tvdbDetails `json:"data"`
}

type tvdbDetails struct {
	ID                  tvdbStringNumber    `json:"id"`
	Slug                string              `json:"slug"`
	Name                string              `json:"name"`
	Title               string              `json:"title"`
	Year                tvdbStringNumber    `json:"year"`
	Overview            string              `json:"overview"`
	Image               string              `json:"image"`
	FirstAired          string              `json:"firstAired"`
	FirstRelease        tvdbDateValue       `json:"first_release"`
	Runtime             int32               `json:"runtime"`
	AverageRuntime      int32               `json:"averageRuntime"`
	Status              tvdbStatusValue     `json:"status"`
	OriginalCountry     string              `json:"originalCountry"`
	OriginalLanguage    string              `json:"originalLanguage"`
	Score               float64             `json:"score"`
	BoxOffice           string              `json:"boxOffice"`
	BoxOfficeUS         string              `json:"boxOfficeUS"`
	Budget              string              `json:"budget"`
	Artworks            []tvdbArtwork       `json:"artworks"`
	ContentRatings      []tvdbContentRating `json:"contentRatings"`
	Genres              []tvdbNamedEntity   `json:"genres"`
	Inspirations        []tvdbInspiration   `json:"inspirations"`
	Companies           tvdbCompanies       `json:"companies"`
	ProductionCountries []tvdbProduction    `json:"production_countries"`
	Studios             []tvdbNamedEntity   `json:"studios"`
	RemoteIDs           []tvdbRemoteID      `json:"remoteIds"`
	Trailers            []tvdbTrailer       `json:"trailers"`
	Translations        tvdbTranslations    `json:"translations"`
	TagOptions          []tvdbTagOption     `json:"tagOptions"`
	Seasons             []tvdbSeason        `json:"seasons"`
	Episodes            []tvdbEpisode       `json:"episodes"`
	Characters          []tvdbCharacter     `json:"characters"`
	Releases            []tvdbMovieRelease  `json:"releases"`
	SpokenLanguages     []string            `json:"spoken_languages"`
	SubtitleLanguages   []string            `json:"subtitleLanguages"`
}

type tvdbNamedEntity struct {
	Name string `json:"name"`
}

type tvdbArtwork struct {
	Image        string  `json:"image"`
	Thumbnail    string  `json:"thumbnail"`
	Type         int32   `json:"type"`
	Width        int32   `json:"width"`
	Height       int32   `json:"height"`
	Score        float64 `json:"score"`
	IncludesText bool    `json:"includesText"`
	Language     string  `json:"language"`
}

type tvdbContentRating struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type tvdbCompanies struct {
	Studio     []tvdbNamedEntity `json:"studio"`
	Network    []tvdbNamedEntity `json:"network"`
	Production []tvdbNamedEntity `json:"production"`
}

type tvdbProduction struct {
	Country string `json:"country"`
	Name    string `json:"name"`
}

type tvdbInspiration struct {
	Type     string `json:"type"`
	TypeName string `json:"type_name"`
}

type tvdbTagOption struct {
	Name    string `json:"name"`
	TagName string `json:"tagName"`
}

type tvdbRemoteID struct {
	ID         string `json:"id"`
	SourceName string `json:"sourceName"`
}

type tvdbTrailer struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Language string `json:"language"`
}

type tvdbTranslations struct {
	NameTranslations     []tvdbTranslation `json:"nameTranslations"`
	OverviewTranslations []tvdbTranslation `json:"overviewTranslations"`
}

type tvdbTranslationResponse struct {
	Data tvdbTranslation `json:"data"`
}

type tvdbTranslation struct {
	Language  string `json:"language"`
	Name      string `json:"name"`
	Overview  string `json:"overview"`
	Tagline   string `json:"tagline"`
	IsPrimary bool   `json:"isPrimary"`
}

type tvdbSeason struct {
	Name         string `json:"name"`
	Number       int32  `json:"number"`
	EpisodeCount int32  `json:"episodeCount"`
	Image        string `json:"image"`
}

type tvdbEpisode struct {
	Name          string `json:"name"`
	Number        int32  `json:"number"`
	EpisodeNumber int32  `json:"episodeNumber"`
	Overview      string `json:"overview"`
	Aired         string `json:"aired"`
	Image         string `json:"image"`
}

type tvdbCharacter struct {
	PeopleID   tvdbStringNumber `json:"peopleId"`
	PersonID   tvdbStringNumber `json:"personId"`
	PersonName string           `json:"personName"`
	Name       string           `json:"name"`
	Role       string           `json:"role"`
	Image      string           `json:"image"`
	PeopleType string           `json:"peopleType"`
	TypeName   string           `json:"typeName"`
	Type       int32            `json:"type"`
	IsFeatured bool             `json:"isFeatured"`
	Sort       int32            `json:"sort"`
}

type tvdbMovieRelease struct {
	Date    string `json:"date"`
	Country string `json:"country"`
}
