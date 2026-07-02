package decisions

import (
	"regexp"
	"strings"
	"sync"

	"media-manager/internal/storage"
)

type CustomFormatMatch struct {
	ID                      string
	Name                    string
	IncludeInRenameTemplate bool
	MatchedSpecs            []CustomFormatSpecMatch
}

type CustomFormatSpecMatch struct {
	ID       string
	Name     string
	Type     string
	Value    string
	Required bool
}

func MatchCustomFormats(parsed ParsedRelease, formats []storage.CustomFormat) []CustomFormatMatch {
	matches := []CustomFormatMatch{}
	for _, format := range formats {
		match, ok := matchCustomFormat(parsed, format)
		if ok {
			matches = append(matches, match)
		}
	}
	return matches
}

func matchCustomFormat(parsed ParsedRelease, format storage.CustomFormat) (CustomFormatMatch, bool) {
	if len(format.IncludeSpecs) == 0 {
		return CustomFormatMatch{}, false
	}
	requiredCount := 0
	requiredMatches := 0
	optionalCount := 0
	optionalMatches := 0
	strongOptionalMatches := 0
	specMatches := []CustomFormatSpecMatch{}
	for _, spec := range format.ExcludeSpecs {
		if specMatchesRelease(parsed, spec) {
			return CustomFormatMatch{}, false
		}
	}
	for _, spec := range format.IncludeSpecs {
		matched := specMatchesRelease(parsed, spec)
		if spec.Required {
			requiredCount++
			if matched {
				requiredMatches++
				specMatches = append(specMatches, specMatch(spec))
			}
			continue
		}
		optionalCount++
		if matched {
			optionalMatches++
			if isStrongOptionalSpec(spec) {
				strongOptionalMatches++
			}
			specMatches = append(specMatches, specMatch(spec))
		}
	}
	if requiredMatches != requiredCount {
		return CustomFormatMatch{}, false
	}
	if requiredCount == 0 && optionalCount > 0 && optionalMatches == 0 {
		return CustomFormatMatch{}, false
	}
	if requiredCount == 0 && optionalCount > 1 && strongOptionalMatches == 0 {
		return CustomFormatMatch{}, false
	}
	return CustomFormatMatch{
		ID:                      format.ID.String(),
		Name:                    format.Name,
		IncludeInRenameTemplate: format.IncludeInRenameTemplate,
		MatchedSpecs:            specMatches,
	}, true
}

func isStrongOptionalSpec(spec storage.CustomFormatSpec) bool {
	switch spec.Type {
	case "source", "resolution", "quality":
		return false
	default:
		return true
	}
}

func specMatch(spec storage.CustomFormatSpec) CustomFormatSpecMatch {
	return CustomFormatSpecMatch{
		ID:       spec.ID,
		Name:     spec.Name,
		Type:     spec.Type,
		Value:    spec.Value,
		Required: spec.Required,
	}
}

func specMatchesRelease(parsed ParsedRelease, spec storage.CustomFormatSpec) bool {
	switch spec.Type {
	case "quality":
		return valueMatchesAny(spec.Value, parsed.Quality, parsed.QualityID)
	case "source":
		return valueMatchesAny(spec.Value, parsed.Source)
	case "resolution":
		return valueMatchesAny(spec.Value, parsed.Resolution)
	case "videoCodec":
		return valueMatchesAny(spec.Value, parsed.VideoCodec)
	case "audioCodec":
		return valueMatchesAny(spec.Value, parsed.AudioCodec, parsed.AudioChannels)
	case "releaseGroup":
		return valueMatchesAny(spec.Value, parsed.ReleaseGroup)
	case "releaseType":
		return valueMatchesAny(spec.Value, parsed.ReleaseType)
	case "edition":
		return valueMatchesAny(spec.Value, parsed.Edition)
	case "language":
		return valueMatchesAny(spec.Value, parsed.Languages...)
	case "releaseTitle", "indexerFlag":
		return valueMatchesAny(spec.Value, parsed.ReleaseTitle)
	default:
		return valueMatchesAny(spec.Value, parsed.ReleaseTitle)
	}
}

func valueMatchesAny(pattern string, candidates ...string) bool {
	for _, candidate := range candidates {
		if valueMatches(pattern, candidate) {
			return true
		}
	}
	return false
}

func valueMatches(pattern string, candidate string) bool {
	if strings.TrimSpace(pattern) == "" || strings.TrimSpace(candidate) == "" {
		return false
	}
	compiled := cachedPattern(pattern)
	if compiled.regex != nil {
		return compiled.regex.MatchString(candidate)
	}
	if compiled.literal {
		return containsAnyNormalized(candidate, pattern)
	}
	return false
}

type cachedRegexPattern struct {
	regex   *regexp.Regexp
	literal bool
}

var regexPatternCache sync.Map

func cachedPattern(pattern string) cachedRegexPattern {
	if cached, ok := regexPatternCache.Load(pattern); ok {
		return cached.(cachedRegexPattern)
	}
	compiled := cachedRegexPattern{}
	if regex, err := regexp.Compile("(?i)" + pattern); err == nil {
		compiled.regex = regex
	} else {
		compiled.literal = literalPattern.MatchString(pattern)
	}
	actual, _ := regexPatternCache.LoadOrStore(pattern, compiled)
	return actual.(cachedRegexPattern)
}

var literalPattern = regexp.MustCompile(`^[A-Za-z0-9 ._+-]+$`)
