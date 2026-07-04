package cardigann

import (
	"strconv"
	"strings"
)

func trackerCategories(def cardigannDefinition, selected []int32) []string {
	if len(selected) == 0 {
		return []string{}
	}
	want := map[int32]bool{}
	for _, id := range selected {
		want[id] = true
	}
	values := []string{}
	for _, mapping := range def.Caps.CategoryMappings {
		if want[standardCategoryID(mapping.Cat)] {
			values = append(values, mapping.ID)
		}
	}
	if len(values) > 0 {
		return values
	}
	for _, id := range selected {
		values = append(values, strconv.FormatInt(int64(id), 10))
	}
	return values
}

func standardCategoryID(value string) int32 {
	top := strings.Split(value, "/")[0]
	switch top {
	case "Console":
		return 1000
	case "Movies":
		return 2000
	case "Audio":
		return 3000
	case "PC":
		return 4000
	case "TV":
		return 5000
	case "XXX":
		return 6000
	case "Books":
		return 7000
	case "Other":
		return 8000
	default:
		return 0
	}
}
