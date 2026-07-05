package httpapi

func libraryFolderKindAllowsMediaKind(folderKind string, mediaKind string) bool {
	switch folderKind {
	case "movie":
		return mediaKind == "movie" || mediaKind == "anime_movie"
	case "series":
		return mediaKind == "series" || mediaKind == "anime_series"
	default:
		return false
	}
}
