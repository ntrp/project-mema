package content

import (
	"context"
	"errors"
	"testing"

	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestBrowsePaginatesWithoutChangingTotalMatches(t *testing.T) {
	ctx := context.Background()
	paths := []string{"/media/A.mkv", "/media/B.mkv", "/media/C.mkv"}
	tree := NewTree(fakeSource{items: []storage.MediaItem{
		{ID: uuid.New(), Type: "movie", Title: "B Movie", FilePaths: []string{paths[1]}},
		{ID: uuid.New(), Type: "movie", Title: "A Movie", FilePaths: []string{paths[0]}},
		{ID: uuid.New(), Type: "movie", Title: "C Movie", FilePaths: []string{paths[2]}},
	}}).WithStat(fakeStat(paths...))

	response, err := tree.Browse(ctx, BrowseRequest{
		ObjectID:       EncodeID(RootContainerRef("movies")),
		BrowseFlag:     BrowseDirectChildren,
		StartingIndex:  1,
		RequestedCount: 1,
		SortCriteria:   "+dc:title",
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.TotalMatches != 3 || response.NumberReturned != 1 {
		t.Fatalf("counts = returned %d total %d", response.NumberReturned, response.TotalMatches)
	}
	if len(response.Objects) != 1 || response.Objects[0].Title != "B Movie" {
		t.Fatalf("paged objects = %#v", response.Objects)
	}
}

func TestBrowseSortsByMultipleCriteria(t *testing.T) {
	paths := []string{"/media/A.mkv", "/media/B.mkv", "/media/C.mkv"}
	tree := NewTree(fakeSource{items: []storage.MediaItem{
		{ID: uuid.New(), Type: "movie", Title: "Same", FilePaths: []string{paths[0]}},
		{ID: uuid.New(), Type: "movie", Title: "Alpha", FilePaths: []string{paths[1]}},
		{ID: uuid.New(), Type: "movie", Title: "Same", FilePaths: []string{paths[2]}},
	}}).WithStat(fakeStat(paths...))

	response, err := tree.Browse(context.Background(), BrowseRequest{
		ObjectID:     EncodeID(RootContainerRef("movies")),
		BrowseFlag:   BrowseDirectChildren,
		SortCriteria: "+dc:title,-dc:date",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Objects) != 3 || response.Objects[0].Title != "Alpha" {
		t.Fatalf("objects = %#v", response.Objects)
	}
}

func TestBrowseMetadataReturnsExactlyOneObject(t *testing.T) {
	ctx := context.Background()
	path := "/media/Scenario.mkv"
	item := storage.MediaItem{ID: uuid.New(), Type: "movie", Title: "Scenario", FilePaths: []string{path}}
	tree := NewTree(fakeSource{items: []storage.MediaItem{item}}).WithStat(fakeStat(path))

	response, err := tree.Browse(ctx, BrowseRequest{
		ObjectID:   EncodeID(MediaItemRef(item.ID)),
		BrowseFlag: BrowseMetadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	if response.TotalMatches != 1 || response.NumberReturned != 1 || len(response.Objects) != 1 {
		t.Fatalf("metadata response = %#v", response)
	}
	if response.Objects[0].Title != "Scenario" {
		t.Fatalf("metadata object = %#v", response.Objects[0])
	}
}

func TestBrowseInvalidObjectReturnsNotFound(t *testing.T) {
	tree := NewTree(fakeSource{})

	_, err := tree.Browse(context.Background(), BrowseRequest{
		ObjectID:   "not-an-object-id",
		BrowseFlag: BrowseDirectChildren,
	})
	if !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("err = %v", err)
	}
}
