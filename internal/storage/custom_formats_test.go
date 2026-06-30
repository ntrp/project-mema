package storage

import "testing"

func TestMarshalCustomFormatSpecs(t *testing.T) {
	includeSpecs, excludeSpecs, err := marshalCustomFormatSpecs(CustomFormatInput{
		IncludeSpecs: []CustomFormatSpec{{
			ID:       "release-title",
			Name:     "Release Title",
			Type:     "releaseTitle",
			Value:    "WEB",
			Required: true,
		}},
		ExcludeSpecs: []CustomFormatSpec{{
			ID:       "not-cam",
			Name:     "Not CAM",
			Type:     "source",
			Value:    "CAM",
			Required: true,
		}},
	})
	if err != nil {
		t.Fatalf("marshal custom format specs: %v", err)
	}
	if string(includeSpecs) != `[{"id":"release-title","name":"Release Title","type":"releaseTitle","value":"WEB","required":true}]` {
		t.Fatalf("unexpected include specs JSON: %s", includeSpecs)
	}
	if string(excludeSpecs) != `[{"id":"not-cam","name":"Not CAM","type":"source","value":"CAM","required":true}]` {
		t.Fatalf("unexpected exclude specs JSON: %s", excludeSpecs)
	}
}
