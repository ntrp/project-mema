package content

import (
	"context"
	"os"
	"time"

	"media-manager/internal/storage"

	"github.com/google/uuid"
)

type LibrarySource interface {
	ListMediaItems(context.Context) ([]storage.MediaItem, error)
}

type FileStatFunc func(string) (os.FileInfo, error)

type Tree struct {
	source LibrarySource
	stat   FileStatFunc
}

type ObjectKind string

const (
	ObjectContainer ObjectKind = "container"
	ObjectItem      ObjectKind = "item"
)

type Object struct {
	ID         string
	ParentID   string
	Title      string
	Class      string
	Kind       ObjectKind
	ChildCount int

	MediaType string
	Year      *int32
	CreatedAt time.Time
	UpdatedAt time.Time

	MediaItemID *uuid.UUID
	SeasonID    *uuid.UUID
	EpisodeID   *uuid.UUID
	FileHash    string
	FilePath    string
}

type File struct {
	Path      string
	Hash      string
	Season    int32
	Episode   int32
	HasNumber bool
}

type fileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (f fileInfo) Name() string       { return f.name }
func (f fileInfo) Size() int64        { return f.size }
func (f fileInfo) Mode() os.FileMode  { return 0 }
func (f fileInfo) ModTime() time.Time { return time.Time{} }
func (f fileInfo) IsDir() bool        { return f.isDir }
func (f fileInfo) Sys() any           { return nil }

func NewTree(source LibrarySource) *Tree {
	return &Tree{source: source, stat: os.Stat}
}

func (t *Tree) WithStat(stat FileStatFunc) *Tree {
	t.stat = stat
	return t
}
