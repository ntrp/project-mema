package cardigann

import (
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/tidwall/gjson"
	"golang.org/x/text/unicode/norm"
)

func applyCardigannFilters(value string, filters []cardigannFilter, ctx cardigannContext) string {
	result := value
	for _, filter := range filters {
		result = applyCardigannFilter(result, filter, ctx)
	}
	return strings.TrimSpace(result)
}

func applyCardigannFilter(value string, filter cardigannFilter, ctx cardigannContext) string {
	switch strings.ToLower(filter.Name) {
	case "trim":
		if cutset := filterArg(filter.Args); cutset != "" {
			return strings.Trim(value, cutset)
		}
		return strings.TrimSpace(value)
	case "tolower":
		return strings.ToLower(value)
	case "toupper":
		return strings.ToUpper(value)
	case "replace":
		args := filterArgs(filter.Args)
		if len(args) < 2 {
			return value
		}
		return strings.ReplaceAll(value, args[0], renderFilterArg(args[1], ctx))
	case "re_replace":
		args := filterArgs(filter.Args)
		if len(args) < 2 {
			return value
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			return value
		}
		return re.ReplaceAllString(value, renderFilterArg(args[1], ctx))
	case "regexp":
		args := filterArgs(filter.Args)
		if len(args) == 0 {
			return value
		}
		re, err := regexp.Compile(args[0])
		if err != nil {
			return value
		}
		match := re.FindStringSubmatch(value)
		if len(match) > 1 {
			return match[1]
		}
		if len(match) == 1 {
			return match[0]
		}
		return ""
	case "split":
		args := filterArgs(filter.Args)
		if len(args) < 2 {
			return value
		}
		index, err := strconv.Atoi(args[1])
		if err != nil {
			return value
		}
		parts := strings.Split(value, args[0])
		if index < 0 {
			index += len(parts)
		}
		if index < 0 || index >= len(parts) {
			return ""
		}
		return parts[index]
	case "prepend":
		return renderFilterArg(filterArg(filter.Args), ctx) + value
	case "append":
		return value + renderFilterArg(filterArg(filter.Args), ctx)
	case "urldecode":
		decoded, err := url.QueryUnescape(value)
		if err != nil {
			return value
		}
		return decoded
	case "urlencode":
		return url.QueryEscape(value)
	case "htmldecode":
		return html.UnescapeString(value)
	case "htmlencode":
		return html.EscapeString(value)
	case "querystring":
		parsed, err := url.Parse(value)
		if err != nil {
			return value
		}
		return parsed.Query().Get(filterArg(filter.Args))
	case "dateparse":
		if parsed, ok := parseCardigannDate(value, filterArg(filter.Args)); ok {
			return parsed.Format(time.RFC3339)
		}
		return value
	case "timeago", "reltime":
		if parsed, ok := parseFuzzyTime(value, time.Now()); ok {
			return parsed.Format(time.RFC3339)
		}
		return value
	case "fuzzytime":
		if parsed, ok := parseFuzzyTime(value, time.Now()); ok {
			return parsed.Format(time.RFC3339)
		}
		return value
	case "validfilename":
		return validCardigannFilename(value)
	case "diacritics":
		if filterArg(filter.Args) == "replace" {
			return replaceDiacritics(value)
		}
		return value
	case "jsonjoinarray":
		args := filterArgs(filter.Args)
		if len(args) < 2 {
			return value
		}
		items := gjson.Get(value, normalizeJSONSelector(args[0])).Array()
		values := make([]string, 0, len(items))
		for _, item := range items {
			values = append(values, item.String())
		}
		return strings.Join(values, args[1])
	default:
		return value
	}
}

func validCardigannFilename(value string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]+`)
	return strings.TrimSpace(re.ReplaceAllString(value, "_"))
}

func replaceDiacritics(value string) string {
	decomposed := norm.NFD.String(value)
	out := make([]rune, 0, len(decomposed))
	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		out = append(out, r)
	}
	return norm.NFC.String(string(out))
}

func renderFilterArg(value string, ctx cardigannContext) string {
	rendered, err := renderCardigannTemplate(value, ctx)
	if err != nil {
		return value
	}
	return rendered
}

func filterArg(value any) string {
	args := filterArgs(value)
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func filterArgs(value any) []string {
	switch typed := value.(type) {
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			values = append(values, scalarString(item))
		}
		return values
	case []string:
		return typed
	case string:
		return []string{typed}
	case int:
		return []string{strconv.Itoa(typed)}
	case int64:
		return []string{strconv.FormatInt(typed, 10)}
	case bool:
		return []string{strconv.FormatBool(typed)}
	default:
		return nil
	}
}

func scalarString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return ""
	}
}
