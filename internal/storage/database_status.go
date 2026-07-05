package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"
)

type DatabaseStatus struct {
	Type    string
	Version string
}

func (s *SettingsStore) GetDatabaseStatus(ctx context.Context) (DatabaseStatus, error) {
	version, err := storagegen.New(s.pool).GetDatabaseVersion(ctx)
	if err != nil {
		return DatabaseStatus{}, err
	}
	return DatabaseStatus{Type: "PostgreSQL", Version: version}, nil
}
