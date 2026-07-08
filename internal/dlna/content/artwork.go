package content

import (
	"net/url"
	"path"
	"strings"
)

func ApplyArtworkURLs(baseURL string, objects []Object) []Object {
	if strings.TrimSpace(baseURL) == "" {
		return objects
	}
	updated := make([]Object, len(objects))
	copy(updated, objects)
	for index := range updated {
		if updated[index].Kind == ObjectItem && updated[index].FilePath != "" {
			artwork := ArtworkURL(baseURL, updated[index].ID)
			artwork += "?kind=thumbnail"
			updated[index].Artwork = &artwork
			continue
		}
		if updated[index].Artwork != nil {
			artwork := ArtworkURL(baseURL, updated[index].ID)
			updated[index].Artwork = &artwork
		}
	}
	return updated
}

func ArtworkURL(baseURL string, objectID string) string {
	base, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "/dlna/artwork/" + url.PathEscape(objectID)
	}
	base.Path = path.Join(base.Path, "/dlna/artwork", objectID)
	return base.String()
}
