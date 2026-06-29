package metadata

type tmdbSearchResponse struct {
	Results []tmdbMedia `json:"results"`
}

type tmdbMedia struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Name         string `json:"name"`
	ReleaseDate  string `json:"release_date"`
	FirstAirDate string `json:"first_air_date"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
}

type tmdbDetails struct {
	ID               int64        `json:"id"`
	Title            string       `json:"title"`
	Name             string       `json:"name"`
	ReleaseDate      string       `json:"release_date"`
	FirstAirDate     string       `json:"first_air_date"`
	Overview         string       `json:"overview"`
	PosterPath       string       `json:"poster_path"`
	BackdropPath     string       `json:"backdrop_path"`
	Status           string       `json:"status"`
	OriginalLanguage string       `json:"original_language"`
	Runtime          int32        `json:"runtime"`
	EpisodeRunTime   []int32      `json:"episode_run_time"`
	NumberOfSeasons  int32        `json:"number_of_seasons"`
	NumberOfEpisodes int32        `json:"number_of_episodes"`
	VoteAverage      float64      `json:"vote_average"`
	Genres           []tmdbName   `json:"genres"`
	CreatedBy        []tmdbName   `json:"created_by"`
	Networks         []tmdbName   `json:"networks"`
	Seasons          []tmdbSeason `json:"seasons"`
	Credits          tmdbCredits  `json:"credits"`
}

type tmdbName struct {
	Name string `json:"name"`
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
}

type tmdbCastMember struct {
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
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
