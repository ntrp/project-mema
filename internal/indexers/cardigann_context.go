package indexers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"text/template"
)

type cardigannContext struct {
	Config     map[string]any
	Query      map[string]any
	Result     map[string]any
	Categories []string
	Keywords   string
	True       bool
	False      bool
}

func newCardigannContext(def cardigannDefinition, config Config, query string, mediaType string) cardigannContext {
	values := fieldValueMap(config.Fields)
	siteLink := strings.TrimRight(config.BaseURL, "/") + "/"
	values["sitelink"] = siteLink
	values["apikey"] = stringValue(config.APIKey)
	for _, setting := range def.Settings {
		if _, ok := values[setting.Name]; !ok {
			values[setting.Name] = setting.Default
		}
	}
	queryMap := map[string]any{
		"Type":        searchType(mediaType),
		"Q":           query,
		"Categories":  config.Categories,
		"Limit":       "",
		"Offset":      "",
		"Extended":    "",
		"APIKey":      "",
		"Genre":       "",
		"Movie":       "",
		"Year":        "",
		"IMDBID":      "",
		"IMDBIDShort": "",
		"TMDBID":      "",
		"Series":      "",
		"Ep":          "",
		"Season":      "",
	}
	keywords := strings.TrimSpace(query)
	return cardigannContext{
		Config:     values,
		Query:      queryMap,
		Result:     map[string]any{},
		Categories: trackerCategories(def, config.Categories),
		Keywords:   keywords,
		True:       true,
		False:      false,
	}
}

func fieldValueMap(raw json.RawMessage) map[string]any {
	values := map[string]any{}
	if len(raw) == 0 {
		return values
	}
	var fields []struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}
	if err := json.Unmarshal(raw, &fields); err != nil {
		return values
	}
	for _, field := range fields {
		if strings.TrimSpace(field.Name) != "" {
			values[field.Name] = field.Value
		}
	}
	return values
}

func searchType(mediaType string) string {
	switch mediaType {
	case "movie":
		return "movie"
	case "series":
		return "tvsearch"
	default:
		return "search"
	}
}

func renderCardigannTemplate(input string, ctx cardigannContext) (string, error) {
	if !strings.Contains(input, "{{") {
		return input, nil
	}
	tpl, err := template.New("cardigann").Funcs(cardigannTemplateFuncs()).Option("missingkey=zero").Parse(repairTemplate(input))
	if err != nil {
		return "", err
	}
	var out strings.Builder
	if err := tpl.Execute(&out, ctx); err != nil {
		return "", err
	}
	return out.String(), nil
}

func repairTemplate(input string) string {
	return strings.ReplaceAll(input, ".False))", ".False)")
}

func cardigannTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"join": func(values any, separator string) string {
			return strings.Join(stringList(values), separator)
		},
		"re_replace": func(value any, pattern string, replacement string) string {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Sprint(value)
			}
			return re.ReplaceAllString(fmt.Sprint(value), replacement)
		},
		"urlquery": func(value any) string {
			return url.QueryEscape(fmt.Sprint(value))
		},
	}
}

func stringList(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []int32:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			values = append(values, fmt.Sprint(item))
		}
		return values
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			values = append(values, fmt.Sprint(item))
		}
		return values
	default:
		if value == nil {
			return nil
		}
		return []string{fmt.Sprint(value)}
	}
}
