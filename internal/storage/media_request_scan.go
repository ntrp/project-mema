package storage

import "github.com/jackc/pgx/v5"

func scanMediaRequest(row pgx.Row) (MediaRequest, error) {
	var request MediaRequest
	err := row.Scan(
		&request.ID,
		&request.RequestedByUserID,
		&request.RequestedByUsername,
		&request.Type,
		&request.Title,
		&request.Year,
		&request.ExternalProvider,
		&request.ExternalID,
		&request.Overview,
		&request.PosterPath,
		&request.Status,
		&request.MonitorMode,
		&request.MinimumAvailability,
		&request.QualityProfileID,
		&request.LibraryFolderID,
		&request.MediaItemID,
		&request.DecidedAt,
		&request.Tags,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	return request, err
}
