package storage

type MediaRenamePreview struct {
	Rows []MediaRenamePreviewRow
}

type MediaRenameApplyResult struct {
	Rows         []MediaRenamePreviewRow
	AppliedCount int32
	SkippedCount int32
	FailedCount  int32
}

type MediaRenamePreviewRow struct {
	CurrentPath  string
	ProposedPath string
	Status       string
	Messages     []string
}
