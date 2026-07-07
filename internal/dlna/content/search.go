package content

import (
	"context"
	"errors"
	"regexp"
	"strings"
)

type SearchRequest struct {
	ContainerID    string
	SearchCriteria string
	Filter         string
	StartingIndex  int
	RequestedCount int
	SortCriteria   string
}

type Criterion struct {
	Field string
	Op    string
	Value string
}

var criterionPattern = regexp.MustCompile(`(?i)^\s*([a-z]+:[a-z]+)\s*(contains|derivedfrom|=|>=|<=|>|<)\s*"([^"]*)"\s*$`)
var criteriaAndPattern = regexp.MustCompile(`(?i)\s+and\s+`)

func (t *Tree) Search(ctx context.Context, request SearchRequest) (BrowseResponse, error) {
	criteria, err := ParseCriteria(request.SearchCriteria)
	if err != nil {
		return BrowseResponse{}, err
	}
	objects, err := t.searchScope(ctx, request.ContainerID)
	if err != nil {
		return BrowseResponse{}, err
	}
	matched := make([]Object, 0, len(objects))
	for _, object := range objects {
		if object.Kind == ObjectItem && matchesCriteria(object, criteria) {
			matched = append(matched, object)
		}
	}
	sortObjects(matched, request.SortCriteria)
	total := len(matched)
	matched = pageObjects(matched, request.StartingIndex, request.RequestedCount)
	return BrowseResponse{Objects: matched, NumberReturned: len(matched), TotalMatches: total}, nil
}

func ParseSearchRequest(args map[string]string) (SearchRequest, error) {
	containerID, err := requiredBrowseArg(args, "ContainerID")
	if err != nil {
		return SearchRequest{}, err
	}
	start, err := parseNonNegative(args["StartingIndex"])
	if err != nil {
		return SearchRequest{}, err
	}
	count, err := parseNonNegative(args["RequestedCount"])
	if err != nil {
		return SearchRequest{}, err
	}
	return SearchRequest{
		ContainerID:    containerID,
		SearchCriteria: strings.TrimSpace(args["SearchCriteria"]),
		Filter:         strings.TrimSpace(args["Filter"]),
		StartingIndex:  start,
		RequestedCount: count,
		SortCriteria:   strings.TrimSpace(args["SortCriteria"]),
	}, nil
}

func ParseCriteria(input string) ([]Criterion, error) {
	input = strings.TrimSpace(input)
	if input == "" || input == "*" {
		return nil, nil
	}
	parts := criteriaAndPattern.Split(input, -1)
	criteria := make([]Criterion, 0, len(parts))
	for _, part := range parts {
		match := criterionPattern.FindStringSubmatch(part)
		if len(match) != 4 {
			return nil, errors.New("unsupported search criteria")
		}
		criterion := Criterion{
			Field: strings.ToLower(match[1]),
			Op:    strings.ToLower(match[2]),
			Value: match[3],
		}
		if !supportedCriterion(criterion) {
			return nil, errors.New("unsupported search criteria")
		}
		criteria = append(criteria, criterion)
	}
	return criteria, nil
}

func supportedCriterion(criterion Criterion) bool {
	switch criterion.Field {
	case "dc:title", "upnp:class", "upnp:genre", "dc:creator", "dc:date":
		return true
	default:
		return false
	}
}

func (t *Tree) searchScope(ctx context.Context, containerID string) ([]Object, error) {
	seen := map[string]struct{}{}
	var scoped []Object
	var visit func(string) error
	visit = func(id string) error {
		children, err := t.BrowseChildren(ctx, id)
		if err != nil {
			return err
		}
		for _, child := range children {
			if _, ok := seen[child.ID]; ok {
				continue
			}
			seen[child.ID] = struct{}{}
			scoped = append(scoped, child)
			if child.Kind == ObjectContainer {
				if err := visit(child.ID); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return scoped, visit(containerID)
}

func matchesCriteria(object Object, criteria []Criterion) bool {
	for _, criterion := range criteria {
		if !matchCriterion(object, criterion) {
			return false
		}
	}
	return true
}

func matchCriterion(object Object, criterion Criterion) bool {
	switch criterion.Field {
	case "dc:title":
		return matchString(object.Title, criterion)
	case "upnp:class":
		return matchClass(object.Class, criterion)
	case "upnp:genre":
		return matchStrings(object.Genres, criterion)
	case "dc:creator":
		return matchStrings(object.Artists, criterion)
	case "dc:date":
		return matchDate(stringValue(object.Date), criterion)
	default:
		return false
	}
}

func matchClass(value string, criterion Criterion) bool {
	if criterion.Op == "derivedfrom" {
		return strings.HasPrefix(strings.ToLower(value), strings.ToLower(criterion.Value))
	}
	return matchString(value, criterion)
}

func matchStrings(values []string, criterion Criterion) bool {
	for _, value := range values {
		if matchString(value, criterion) {
			return true
		}
	}
	return false
}

func matchDate(value string, criterion Criterion) bool {
	if criterion.Op == "contains" {
		return strings.Contains(value, criterion.Value)
	}
	return compareByOp(value, criterion.Value, criterion.Op)
}

func matchString(value string, criterion Criterion) bool {
	left := strings.ToLower(value)
	right := strings.ToLower(criterion.Value)
	if criterion.Op == "contains" {
		return strings.Contains(left, right)
	}
	if criterion.Op != "=" {
		return false
	}
	return left == right
}

func compareByOp(left string, right string, op string) bool {
	switch op {
	case "=":
		return left == right
	case ">=":
		return left >= right
	case "<=":
		return left <= right
	case ">":
		return left > right
	case "<":
		return left < right
	default:
		return false
	}
}
