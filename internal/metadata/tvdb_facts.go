package metadata

import (
	"strconv"
	"strings"
)

func tvdbPreferredTranslation(values []tvdbTranslation) tvdbTranslation {
	for _, item := range values {
		if strings.EqualFold(item.Language, "eng") && strings.TrimSpace(item.Overview) != "" {
			return item
		}
	}
	for _, item := range values {
		if item.IsPrimary && strings.TrimSpace(item.Overview) != "" {
			return item
		}
	}
	for _, item := range values {
		if strings.TrimSpace(item.Overview) != "" {
			return item
		}
	}
	return tvdbTranslation{}
}

func tvdbCertificationFacts(ratings []tvdbContentRating) []Fact {
	for _, rating := range ratings {
		if strings.EqualFold(rating.Country, "usa") || strings.EqualFold(rating.Country, "us") {
			if value := strings.TrimSpace(rating.Name); value != "" {
				return []Fact{{Label: "Certification", Value: value}}
			}
		}
	}
	for _, rating := range ratings {
		if value := strings.TrimSpace(rating.Name); value != "" {
			return []Fact{{Label: "Certification", Value: value}}
		}
	}
	return nil
}

func tvdbReleaseFacts(releases []tvdbMovieRelease) []Fact {
	facts := []Fact{}
	for _, release := range releases {
		if strings.TrimSpace(release.Date) == "" {
			continue
		}
		if strings.EqualFold(release.Country, "usa") || strings.EqualFold(release.Country, "us") {
			facts = append(facts, Fact{Label: "Release Date", Value: release.Date})
			break
		}
	}
	return facts
}

func tvdbMoneyFacts(item tvdbDetails) []Fact {
	facts := []Fact{}
	if value := tvdbMoneyValue(firstNonEmpty(item.BoxOffice, item.BoxOfficeUS)); value != "" {
		facts = append(facts, Fact{Label: "Revenue", Value: value})
	}
	if value := tvdbMoneyValue(item.Budget); value != "" {
		facts = append(facts, Fact{Label: "Budget", Value: value})
	}
	return facts
}

func tvdbProductionCountries(countries []tvdbProduction) []string {
	values := make([]string, 0, len(countries))
	for _, country := range countries {
		if value := tvdbCountryDisplay(country.Country, country.Name); value != "" {
			values = append(values, value)
		}
	}
	return values
}

func tvdbMoneyValue(value string) string {
	value = strings.TrimSpace(strings.TrimPrefix(value, "$"))
	value = strings.ReplaceAll(value, ",", "")
	if value == "" {
		return ""
	}
	whole, cents, ok := strings.Cut(value, ".")
	if !ok {
		cents = "00"
	}
	dollars, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || dollars <= 0 {
		return strings.TrimSpace(value)
	}
	cents = (cents + "00")[:2]
	if _, err := strconv.Atoi(cents); err != nil {
		return strings.TrimSpace(value)
	}
	return strings.TrimSuffix(formatMoney(dollars), ".00") + "." + cents
}

func tvdbCountryDisplay(code string, name string) string {
	code = tvdbCountryCode(code)
	name = firstNonEmpty(name, tvdbCountryName(code), strings.TrimSpace(code))
	return countryFlag(code, strings.TrimSpace(name))
}

func tvdbCountryCode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "usa":
		return "US"
	case "gbr", "uk":
		return "GB"
	default:
		if len(strings.TrimSpace(value)) == 2 {
			return strings.ToUpper(strings.TrimSpace(value))
		}
		return ""
	}
}

func tvdbCountryName(code string) string {
	switch strings.ToUpper(strings.TrimSpace(code)) {
	case "US":
		return "United States"
	case "GB":
		return "United Kingdom"
	default:
		return ""
	}
}

func tvdbCompanyNames(item tvdbDetails) []string {
	values := []string{}
	values = append(values, tvdbNames(item.Studios)...)
	values = append(values, tvdbNames(item.Companies.Studio)...)
	values = append(values, tvdbNames(item.Companies.Production)...)
	return uniqueStrings(values)
}

func tvdbKeywords(item tvdbDetails) []string {
	values := []string{}
	for _, tag := range item.TagOptions {
		values = append(values, firstNonEmpty(tag.Name, tag.TagName))
	}
	for _, inspiration := range item.Inspirations {
		values = append(values, firstNonEmpty(inspiration.TypeName, inspiration.Type))
	}
	return uniqueStrings(values)
}

func tvdbRemoteIDFacts(ids []tvdbRemoteID) []Fact {
	facts := []Fact{}
	for _, id := range ids {
		value := strings.TrimSpace(id.ID)
		source := tvdbRemoteSourceName(id.SourceName)
		if value == "" || source == "" {
			continue
		}
		facts = append(facts, Fact{Label: source + " ID", Value: value})
	}
	return facts
}

func tvdbRemoteSourceName(source string) string {
	switch strings.ToLower(strings.TrimSpace(source)) {
	case "imdb":
		return "IMDb"
	case "wikidata":
		return "Wikidata"
	case "tmdb":
		return "TMDB"
	default:
		return strings.TrimSpace(source)
	}
}

func tvdbStringFact(label string, values []string) []Fact {
	values = uniqueStrings(values)
	if len(values) == 0 {
		return nil
	}
	return []Fact{{Label: label, Value: strings.Join(values, "\n")}}
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	unique := []string{}
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, trimmed)
	}
	return unique
}
