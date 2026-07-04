package jobs

type GrabReleaseArgs struct {
	ActivityID  string `json:"activity_id" river:"unique"`
	MediaItemID string `json:"media_item_id"`
	Title       string `json:"title"`
	DownloadURL string `json:"download_url"`
	IndexerName string `json:"indexer_name"`
	Protocol    string `json:"protocol,omitempty"`
}

func (GrabReleaseArgs) Kind() string {
	return "media.grab_release"
}
