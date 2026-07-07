package metadata

type tmdbSearchResponse struct {
	Results []tmdbMedia `json:"results"`
}

type tmdbMedia struct {
	ID           int64   `json:"id"`
	MediaType    string  `json:"media_type"`
	Title        string  `json:"title"`
	Name         string  `json:"name"`
	ReleaseDate  string  `json:"release_date"`
	FirstAirDate string  `json:"first_air_date"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	Popularity   float64 `json:"popularity"`
	VoteAverage  float64 `json:"vote_average"`
	VoteCount    int32   `json:"vote_count"`
	GenreIDs     []int64 `json:"genre_ids"`
	Language     string  `json:"original_language"`
}

type tmdbGenreList struct {
	Genres []tmdbIDName `json:"genres"`
}

type tmdbFacetSearchResponse struct {
	Results []tmdbIDName `json:"results"`
}

type tmdbIDName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type tmdbDetails struct {
	ID               int64              `json:"id"`
	Title            string             `json:"title"`
	Name             string             `json:"name"`
	ReleaseDate      string             `json:"release_date"`
	FirstAirDate     string             `json:"first_air_date"`
	Overview         string             `json:"overview"`
	PosterPath       string             `json:"poster_path"`
	Collection       *tmdbCollection    `json:"belongs_to_collection"`
	BackdropPath     string             `json:"backdrop_path"`
	Status           string             `json:"status"`
	OriginalLanguage string             `json:"original_language"`
	Runtime          int32              `json:"runtime"`
	EpisodeRunTime   []int32            `json:"episode_run_time"`
	NumberOfSeasons  int32              `json:"number_of_seasons"`
	NumberOfEpisodes int32              `json:"number_of_episodes"`
	VoteAverage      float64            `json:"vote_average"`
	Budget           int64              `json:"budget"`
	Revenue          int64              `json:"revenue"`
	Genres           []tmdbName         `json:"genres"`
	Keywords         tmdbKeywords       `json:"keywords"`
	CreatedBy        []tmdbName         `json:"created_by"`
	Networks         []tmdbName         `json:"networks"`
	Production       []tmdbName         `json:"production_companies"`
	Countries        []tmdbCountry      `json:"production_countries"`
	Seasons          []tmdbSeason       `json:"seasons"`
	Credits          tmdbCredits        `json:"credits"`
	ExternalIDs      tmdbExternalIDs    `json:"external_ids"`
	ReleaseDates     tmdbReleaseInfo    `json:"release_dates"`
	ContentRatings   tmdbContentRatings `json:"content_ratings"`
	Videos           tmdbVideos         `json:"videos"`
	Recommendations  tmdbSearchResponse `json:"recommendations"`
	Similar          tmdbSearchResponse `json:"similar"`
}

type tmdbCollection struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	Overview     string      `json:"overview"`
	PosterPath   string      `json:"poster_path"`
	BackdropPath string      `json:"backdrop_path"`
	Parts        []tmdbMedia `json:"parts"`
}

type tmdbName struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProfilePath string `json:"profile_path"`
}

type tmdbKeywords struct {
	Keywords []tmdbName `json:"keywords"`
	Results  []tmdbName `json:"results"`
}

type tmdbCountry struct {
	Name string `json:"name"`
	Code string `json:"iso_3166_1"`
}

type tmdbExternalIDs struct {
	IMDBID      string `json:"imdb_id"`
	WikidataID  string `json:"wikidata_id"`
	FacebookID  string `json:"facebook_id"`
	InstagramID string `json:"instagram_id"`
	TwitterID   string `json:"twitter_id"`
	TVDBID      int64  `json:"tvdb_id"`
}

type tmdbReleaseInfo struct {
	Results []tmdbReleaseCountry `json:"results"`
}

type tmdbReleaseCountry struct {
	Code        string            `json:"iso_3166_1"`
	ReleaseList []tmdbReleaseDate `json:"release_dates"`
}

type tmdbReleaseDate struct {
	Date          string `json:"release_date"`
	Type          int    `json:"type"`
	Certification string `json:"certification"`
}

type tmdbContentRatings struct {
	Results []tmdbContentRating `json:"results"`
}

type tmdbContentRating struct {
	Code   string `json:"iso_3166_1"`
	Rating string `json:"rating"`
}

type tmdbVideos struct {
	Results []tmdbVideo `json:"results"`
}

type tmdbVideo struct {
	Key      string `json:"key"`
	Site     string `json:"site"`
	Type     string `json:"type"`
	Official bool   `json:"official"`
}

type tmdbSeason struct {
	Name         string `json:"name"`
	SeasonNumber int32  `json:"season_number"`
	EpisodeCount int32  `json:"episode_count"`
	AirDate      string `json:"air_date"`
	PosterPath   string `json:"poster_path"`
	Episodes     []tmdbEpisode
}

type tmdbSeasonDetails struct {
	Episodes []tmdbEpisode `json:"episodes"`
}

type tmdbEpisode struct {
	Name          string `json:"name"`
	EpisodeNumber int32  `json:"episode_number"`
	Overview      string `json:"overview"`
	AirDate       string `json:"air_date"`
	StillPath     string `json:"still_path"`
}

type tmdbCredits struct {
	Cast []tmdbCastMember `json:"cast"`
	Crew []tmdbCrewMember `json:"crew"`
}

type tmdbCastMember struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
}

type tmdbCrewMember struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Job         string `json:"job"`
	Department  string `json:"department"`
	ProfilePath string `json:"profile_path"`
}

type tmdbPersonDetails struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	Biography    string              `json:"biography"`
	Birthday     string              `json:"birthday"`
	Deathday     string              `json:"deathday"`
	PlaceOfBirth string              `json:"place_of_birth"`
	ProfilePath  string              `json:"profile_path"`
	AlsoKnownAs  []string            `json:"also_known_as"`
	Credits      tmdbCombinedCredits `json:"combined_credits"`
}

type tmdbPersonSearchResponse struct {
	Results []tmdbPersonSearchResult `json:"results"`
}

type tmdbPersonSearchResult struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	ProfilePath string      `json:"profile_path"`
	Popularity  float64     `json:"popularity"`
	KnownFor    []tmdbMedia `json:"known_for"`
}

type tmdbCombinedCredits struct {
	Cast []tmdbCreditMedia `json:"cast"`
	Crew []tmdbCreditMedia `json:"crew"`
}

type tmdbCreditMedia struct {
	ID           int64   `json:"id"`
	MediaType    string  `json:"media_type"`
	Title        string  `json:"title"`
	Name         string  `json:"name"`
	ReleaseDate  string  `json:"release_date"`
	FirstAirDate string  `json:"first_air_date"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	Character    string  `json:"character"`
	Job          string  `json:"job"`
	Popularity   float64 `json:"popularity"`
}

type tvdbLoginResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Status string `json:"status"`
}

type tvdbSearchResponse struct {
	Data []tvdbSearchResult `json:"data"`
}

type tvdbSearchResult struct {
	ID                 string   `json:"id"`
	ObjectID           string   `json:"objectID"`
	TVDBID             string   `json:"tvdb_id"`
	Slug               string   `json:"slug"`
	Type               string   `json:"type"`
	Name               string   `json:"name"`
	Title              string   `json:"title"`
	Year               string   `json:"year"`
	Overview           string   `json:"overview"`
	OverviewTranslated []string `json:"overview_translated"`
	ImageURL           string   `json:"image_url"`
	Poster             string   `json:"poster"`
	Thumbnail          string   `json:"thumbnail"`
}
