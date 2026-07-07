package metadata

import "strings"

const (
	tvdbMoviePosterType   int32 = 14
	tvdbMovieBackdropType int32 = 15
)

func tvdbPoster(artworks []tvdbArtwork) string {
	return tvdbArtworkImage(artworks, func(item tvdbArtwork) bool {
		if item.Type == tvdbMoviePosterType {
			return true
		}
		return item.Height > item.Width && item.Width > 0
	})
}

func tvdbBackdrop(artworks []tvdbArtwork) string {
	return tvdbArtworkImage(artworks, func(item tvdbArtwork) bool {
		if item.Type == tvdbMovieBackdropType {
			return true
		}
		return item.Width > item.Height && item.Height > 0
	})
}

func tvdbArtworkImage(artworks []tvdbArtwork, accepts func(tvdbArtwork) bool) string {
	best := tvdbArtwork{}
	for _, item := range artworks {
		if strings.TrimSpace(item.Image) == "" || !accepts(item) {
			continue
		}
		if best.Image == "" || item.Score > best.Score {
			best = item
		}
	}
	return strings.TrimSpace(best.Image)
}

func tvdbTrailerURL(trailers []tvdbTrailer) string {
	for _, trailer := range trailers {
		if value := strings.TrimSpace(trailer.URL); value != "" {
			return value
		}
	}
	return ""
}
