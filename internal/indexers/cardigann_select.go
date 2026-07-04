package indexers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
	"github.com/tidwall/gjson"
)

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
			return selector.Default, true, nil
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
	return applyCardigannFilters(value, selector.Filters, ctx), true, nil
}

func cardigannJSONRows(body []byte, selector string) []gjson.Result {
	query := normalizeJSONSelector(selector)
	if query == "" || query == "$" {
		result := gjson.ParseBytes(body)
		if result.IsArray() {
			return result.Array()
		}
		return []gjson.Result{result}
	}
	result := gjson.GetBytes(body, query)
	if result.IsArray() {
		return result.Array()
	}
	if result.Exists() {
		return []gjson.Result{result}
	}
	return nil
}

func cardigannJSONValue(row gjson.Result, selector cardigannSelector, ctx cardigannContext, required bool) (string, bool, error) {
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
	result := row.Get(normalizeJSONSelector(query))
	if !result.Exists() {
		if selector.Default != "" {
			return selector.Default, true, nil
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
	return applyCardigannFilters(value, selector.Filters, ctx), true, nil
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
