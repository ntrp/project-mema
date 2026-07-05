package storage

type MediaRenamePreview struct {
	Rows []MediaRenamePreviewRow
}

type MediaRenamePreviewRow struct {
	CurrentPath  string
	ProposedPath string
	Status       string
	Messages     []string
}
