package storage

import (
	"errors"
	"testing"
)

func TestNormalizeLanguageInputForCreate(t *testing.T) {
	input, err := normalizeLanguageInput(LanguageInput{
		Code:        " jp ",
		DisplayName: "  Scenario   Japanese  ",
		Aliases:     []string{" scenario japanese ", "JP", "", "nihongo", "Nihongo"},
	}, true)
	if err != nil {
		t.Fatalf("normalize language input: %v", err)
	}

	if input.Code != "JP" {
		t.Fatalf("expected code JP, got %q", input.Code)
	}
	if input.DisplayName != "Scenario Japanese" {
		t.Fatalf("expected trimmed display name, got %q", input.DisplayName)
	}
	expectStrings(t, input.Aliases, []string{"JP", "Scenario Japanese", "nihongo"})
}

func TestNormalizeLanguageInputRejectsInvalidValues(t *testing.T) {
	if _, err := normalizeLanguageInput(LanguageInput{Code: "x", DisplayName: "Name"}, true); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid code error, got %v", err)
	}
	if _, err := normalizeLanguageInput(LanguageInput{Code: "TEST", DisplayName: "   "}, true); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected blank display name error, got %v", err)
	}
}

func TestNormalizeLanguageInputForUpdateUsesPathCode(t *testing.T) {
	input, err := normalizeLanguageInput(LanguageInput{
		Code:        "",
		DisplayName: " Path Code Language ",
		Aliases:     []string{" path code language "},
	}, false)
	if err != nil {
		t.Fatalf("normalize language update input: %v", err)
	}

	if input.Code != "" {
		t.Fatalf("expected update input not to require body code, got %q", input.Code)
	}
	expectStrings(t, input.Aliases, []string{"Path Code Language"})
}
