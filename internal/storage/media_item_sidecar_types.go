package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaItemSidecar struct {
	ID            uuid.UUID
	MediaItemID   uuid.UUID
	MediaFilePath string
	FilePath      string
	SidecarType   MediaSidecarType
	Subtype       *string
	LanguageID    *string
	Format        *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type MediaItemSidecarInput struct {
	MediaItemID   uuid.UUID
	MediaFilePath string
	FilePath      string
	SidecarType   MediaSidecarType
	Subtype       string
	LanguageID    string
	Format        string
}
