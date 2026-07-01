package decisions

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestMatchCustomFormatsCarriesRenameTemplateFlag(t *testing.T) {
	formatID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	matches := MatchCustomFormats(ParsedRelease{Source: "WEB-DL"}, []storage.CustomFormat{{
		ID:                      formatID,
		Name:                    "WEB",
		IncludeInRenameTemplate: true,
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID:       "source",
			Name:     "Source",
			Type:     "source",
			Value:    "WEB",
			Required: true,
		}},
	}})

	if len(matches) != 1 {
		t.Fatalf("expected one match, got %d", len(matches))
	}
	if !matches[0].IncludeInRenameTemplate {
		t.Fatal("expected rename template flag to be carried into match")
	}
}
