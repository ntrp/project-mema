package content

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"
)

type BrowseFlag string

const (
	BrowseMetadata       BrowseFlag = "BrowseMetadata"
	BrowseDirectChildren BrowseFlag = "BrowseDirectChildren"
)

type BrowseRequest struct {
	ObjectID       string
	BrowseFlag     BrowseFlag
	Filter         string
	StartingIndex  int
	RequestedCount int
	SortCriteria   string
}

type BrowseResponse struct {
	Objects        []Object
	NumberReturned int
	TotalMatches   int
	UpdateID       int
}

func (t *Tree) Browse(ctx context.Context, request BrowseRequest) (BrowseResponse, error) {
	objects, err := t.browseObjects(ctx, request)
	if err != nil {
		return BrowseResponse{}, err
	}
	sortObjects(objects, request.SortCriteria)
	total := len(objects)
	objects = pageObjects(objects, request.StartingIndex, request.RequestedCount)
	return BrowseResponse{
		Objects:        objects,
		NumberReturned: len(objects),
		TotalMatches:   total,
		UpdateID:       0,
	}, nil
}

func (t *Tree) browseObjects(ctx context.Context, request BrowseRequest) ([]Object, error) {
	switch request.BrowseFlag {
	case BrowseMetadata:
		object, err := t.BrowseMetadata(ctx, request.ObjectID)
		if err != nil {
			return nil, err
		}
		return []Object{object}, nil
	case BrowseDirectChildren:
		return t.BrowseChildren(ctx, request.ObjectID)
	default:
		return nil, errors.New("unsupported browse flag")
	}
}

func ParseBrowseRequest(args map[string]string) (BrowseRequest, error) {
	objectID, err := requiredBrowseArg(args, "ObjectID")
	if err != nil {
		return BrowseRequest{}, err
	}
	flag, err := requiredBrowseArg(args, "BrowseFlag")
	if err != nil {
		return BrowseRequest{}, err
	}
	start, err := parseNonNegative(args["StartingIndex"])
	if err != nil {
		return BrowseRequest{}, err
	}
	count, err := parseNonNegative(args["RequestedCount"])
	if err != nil {
		return BrowseRequest{}, err
	}
	return BrowseRequest{
		ObjectID:       objectID,
		BrowseFlag:     BrowseFlag(flag),
		Filter:         strings.TrimSpace(args["Filter"]),
		StartingIndex:  start,
		RequestedCount: count,
		SortCriteria:   strings.TrimSpace(args["SortCriteria"]),
	}, nil
}

func requiredBrowseArg(args map[string]string, name string) (string, error) {
	value := strings.TrimSpace(args[name])
	if value == "" {
		return "", errors.New("missing argument: " + name)
	}
	return value, nil
}

func parseNonNegative(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, errors.New("invalid non-negative integer")
	}
	return parsed, nil
}

func pageObjects(objects []Object, start int, count int) []Object {
	if start >= len(objects) {
		return []Object{}
	}
	end := len(objects)
	if count > 0 && start+count < end {
		end = start + count
	}
	return objects[start:end]
}

func sortObjects(objects []Object, criteria string) {
	orders := parseSortCriteria(criteria)
	if len(orders) == 0 {
		return
	}
	sort.SliceStable(objects, func(i, j int) bool {
		for _, order := range orders {
			compare := compareObjects(objects[i], objects[j], order.Field)
			if compare == 0 {
				continue
			}
			if order.Descending {
				return compare > 0
			}
			return compare < 0
		}
		return false
	})
}

type sortOrder struct {
	Field      string
	Descending bool
}

func parseSortCriteria(criteria string) []sortOrder {
	orders := []sortOrder{}
	for _, part := range strings.Split(criteria, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		orders = append(orders, sortOrder{
			Field:      strings.TrimLeft(part, "+-"),
			Descending: strings.HasPrefix(part, "-"),
		})
	}
	return orders
}

func compareObjects(left Object, right Object, field string) int {
	switch field {
	case "dc:date":
		return strings.Compare(stringValue(left.Date), stringValue(right.Date))
	case "upnp:class":
		return strings.Compare(left.Class, right.Class)
	case "dc:title", "":
		return strings.Compare(left.Title, right.Title)
	default:
		return strings.Compare(left.Title, right.Title)
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
