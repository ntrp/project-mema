package decisions

import (
	"regexp"
	"strings"

	"media-manager/internal/storage"
)

type ReleaseDecision struct {
	Release storage.ReleaseCandidateInput
	Quality string
}

type Engine struct {
	qualities []qualityRule
}

type qualityRule struct {
	id        string
	name      string
	sortOrder int32
	tokens    []string
}

var nonAlphaNumeric = regexp.MustCompile(`[^a-z0-9]+`)

func NewEngine() Engine {
	definitions := storage.QualitySizeDefinitions()
	qualities := make([]qualityRule, 0, len(definitions))
	for _, definition := range definitions {
		qualities = append(qualities, qualityRule{
			id:        definition.ID,
			name:      definition.Name,
			sortOrder: definition.SortOrder,
			tokens: uniqueTokens(
				normalizedToken(definition.ID),
				normalizedToken(definition.Name),
			),
		})
	}
	return Engine{qualities: qualities}
}

func (e Engine) ChooseRelease(candidates []storage.ReleaseCandidateInput) (ReleaseDecision, bool) {
	if len(candidates) == 0 {
		return ReleaseDecision{}, false
	}

	best := candidates[0]
	bestQuality := e.detectQuality(best.Title)
	for _, candidate := range candidates[1:] {
		candidateQuality := e.detectQuality(candidate.Title)
		if betterRelease(candidate, candidateQuality, best, bestQuality) {
			best = candidate
			bestQuality = candidateQuality
		}
	}
	return ReleaseDecision{Release: best, Quality: bestQuality.name}, true
}

func betterRelease(left storage.ReleaseCandidateInput, leftQuality qualityRule, right storage.ReleaseCandidateInput, rightQuality qualityRule) bool {
	if leftQuality.sortOrder != rightQuality.sortOrder {
		return leftQuality.sortOrder > rightQuality.sortOrder
	}
	leftSeeders := int32(-1)
	rightSeeders := int32(-1)
	if left.Seeders != nil {
		leftSeeders = *left.Seeders
	}
	if right.Seeders != nil {
		rightSeeders = *right.Seeders
	}
	if leftSeeders != rightSeeders {
		return leftSeeders > rightSeeders
	}
	if left.SizeBytes != right.SizeBytes {
		return left.SizeBytes > right.SizeBytes
	}
	return strings.ToLower(left.Title) < strings.ToLower(right.Title)
}

func (e Engine) detectQuality(title string) qualityRule {
	normalizedTitle := normalizedToken(title)
	best := qualityRule{}
	for _, quality := range e.qualities {
		for _, token := range quality.tokens {
			if token == "" || !strings.Contains(normalizedTitle, token) {
				continue
			}
			if quality.sortOrder > best.sortOrder {
				best = quality
			}
		}
	}
	return best
}

func normalizedToken(value string) string {
	return nonAlphaNumeric.ReplaceAllString(strings.ToLower(value), "")
}

func uniqueTokens(values ...string) []string {
	seen := map[string]struct{}{}
	tokens := []string{}
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		tokens = append(tokens, value)
	}
	return tokens
}
