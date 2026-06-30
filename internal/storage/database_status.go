package storage

import "context"

type DatabaseStatus struct {
	Type    string
	Version string
}

func (s *SettingsStore) GetDatabaseStatus(ctx context.Context) (DatabaseStatus, error) {
	status := DatabaseStatus{Type: "PostgreSQL"}
	err := s.pool.QueryRow(ctx, `select current_setting('server_version')`).Scan(&status.Version)
	return status, err
}
