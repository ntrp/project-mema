package httpapi

import (
	"reflect"
	"testing"

	"media-manager/internal/decisions"
)

func TestCustomFormatParsingResponseIncludesOnlyRenameTemplateFormats(t *testing.T) {
	response := customFormatParsingResponse(decisions.ParsedRelease{}, []decisions.CustomFormatMatch{
		{ID: "included", Name: "Included", IncludeInRenameTemplate: true},
		{ID: "excluded", Name: "Excluded"},
	}, nil)

	want := []string{"Included"}
	if !reflect.DeepEqual(response.Details.CustomFormatNames, want) {
		t.Fatalf("CustomFormatNames = %v, want %v", response.Details.CustomFormatNames, want)
	}
}
