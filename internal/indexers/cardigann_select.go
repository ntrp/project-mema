package indexers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"github.com/tidwall/gjson"
)

type htmlRow = goquery.Selection

type jsonRow struct {
	current gjson.Result
	parent  *gjson.Result
}

func cardigannHTMLRows(body []byte, selector string) ([]*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	matcher, err := cascadia.Compile(selector)
	if err != nil {
		return nil, err
	}
	rows := []*goquery.Selection{}
	doc.FindMatcher(matcher).Each(func(_ int, row *goquery.Selection) {
		rows = append(rows, row)
	})
	return rows, nil
}

func cardigannHTMLValue(row *goquery.Selection, selector cardigannSelector, ctx cardigannContext, required bool) (string, bool, error) {
	if selector.Text != "" {
		rendered, err := renderCardigannTemplate(selector.Text, ctx)
		if err != nil {
			return "", false, err
		}
		return applyCardigannFilters(rendered, selector.Filters, ctx), true, nil
	}
	query, err := renderCardigannTemplate(selector.Selector, ctx)
	if err != nil {
		return "", false, err
	}
	target := row
	if strings.TrimSpace(query) != "" {
		matcher, err := cascadia.Compile(query)
		if err != nil {
			return "", false, err
		}
		if goquery.NodeName(row) != "" && row.IsMatcher(matcher) {
			target = row
		} else {
			target = row.FindMatcher(matcher).First()
		}
	}
	if target.Length() == 0 {
		if selector.Default != "" {
			rendered, err := renderCardigannTemplate(selector.Default, ctx)
			if err != nil {
				return "", false, err
			}
			return rendered, true, nil
		}
		if required {
			return "", false, fmt.Errorf("selector %q did not match", query)
		}
		return "", false, nil
	}
	if selector.Remove != "" {
		target = target.Clone()
		target.Find(selector.Remove).Remove()
	}
	value := ""
	if selector.Attribute != "" {
		value, _ = target.Attr(selector.Attribute)
	} else {
		value = target.Text()
	}
	value = applyCardigannCase(value, selector.Case)
	return applyCardigannFilters(value, selector.Filters, ctx), true, nil
}

func cardigannJSONRows(body []byte, selector string, rows cardigannRows) []jsonRow {
	query := normalizeJSONSelector(selector)
	parentRows := []gjson.Result{}
	if query == "" || query == "$" {
		result := gjson.ParseBytes(body)
		if result.IsArray() {
			parentRows = result.Array()
		} else {
			parentRows = []gjson.Result{result}
		}
	} else {
		result := gjson.GetBytes(body, query)
		if result.IsArray() {
			parentRows = result.Array()
		} else if result.Exists() {
			parentRows = []gjson.Result{result}
		}
	}
	releaseRows := []jsonRow{}
	for _, parent := range parentRows {
		parentCopy := parent
		selected := parent
		if rows.Attribute != "" {
			selected = parent.Get(normalizeJSONSelector(rows.Attribute))
			if !selected.Exists() {
				if rows.MissingAttributeNoResults {
					continue
				}
				releaseRows = append(releaseRows, jsonRow{current: parent, parent: &parentCopy})
				continue
			}
		}
		if rows.Multiple {
			for _, item := range selected.Array() {
				releaseRows = append(releaseRows, jsonRow{current: item, parent: &parentCopy})
			}
			continue
		}
		releaseRows = append(releaseRows, jsonRow{current: selected, parent: &parentCopy})
	}
	return releaseRows
}

func cardigannJSONValue(row jsonRow, selector cardigannSelector, ctx cardigannContext, required bool) (string, bool, error) {
	if selector.Text != "" {
		rendered, err := renderCardigannTemplate(selector.Text, ctx)
		if err != nil {
			return "", false, err
		}
		return applyCardigannFilters(rendered, selector.Filters, ctx), true, nil
	}
	query, err := renderCardigannTemplate(selector.Selector, ctx)
	if err != nil {
		return "", false, err
	}
	source := row.current
	normalized := normalizeJSONSelector(query)
	if strings.HasPrefix(strings.TrimSpace(query), "..") && row.parent != nil {
		source = *row.parent
		normalized = normalizeJSONSelector(strings.TrimPrefix(strings.TrimSpace(query), ".."))
	}
	result := source.Get(normalized)
	if !result.Exists() {
		if selector.Default != "" {
			rendered, err := renderCardigannTemplate(selector.Default, ctx)
			if err != nil {
				return "", false, err
			}
			return rendered, true, nil
		}
		if required {
			return "", false, fmt.Errorf("selector %q did not match", query)
		}
		return "", false, nil
	}
	value := result.String()
	if selector.Attribute != "" {
		value = result.Get(selector.Attribute).String()
	}
	value = applyCardigannCase(value, selector.Case)
	return applyCardigannFilters(value, selector.Filters, ctx), true, nil
}

func applyCardigannCase(value string, cases map[string]string) string {
	if len(cases) == 0 {
		return value
	}
	if mapped, ok := cases[value]; ok {
		return mapped
	}
	if mapped, ok := cases["*"]; ok {
		return mapped
	}
	return value
}

func normalizeJSONSelector(selector string) string {
	selector = strings.TrimSpace(selector)
	selector, _, _ = strings.Cut(selector, ":has(")
	selector = strings.TrimPrefix(selector, "$.")
	if selector == "$" {
		return ""
	}
	selector = strings.TrimPrefix(selector, "$")
	selector = strings.TrimPrefix(selector, ".")
	return selector
}
