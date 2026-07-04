package storage

import (
	"errors"
	"math"
	"testing"
)

func TestQualitySizeNumericRoundTrip(t *testing.T) {
	value := 12.34
	converted, err := optionalQualitySizeFloat(optionalQualitySizeNumeric(&value))
	if err != nil {
		t.Fatal(err)
	}
	if converted == nil || math.Abs(*converted-value) > 0.001 {
		t.Fatalf("converted quality size = %#v", converted)
	}

	blank, err := optionalQualitySizeFloat(optionalQualitySizeNumeric(nil))
	if err != nil {
		t.Fatal(err)
	}
	if blank != nil {
		t.Fatalf("expected null quality size to stay nil, got %#v", *blank)
	}
}

func TestQualitySizeSettingsUseSqlcQueries(t *testing.T) {
	ctx, store := testDBStore(t)
	preferred := 15.5
	updated, err := store.SaveQualitySizeSettings(ctx, []QualitySizeSettingInput{{
		QualityID:                "bluray-1080p",
		MinimumSizeMBPerMinute:   10.25,
		PreferredSizeMBPerMinute: &preferred,
	}})
	if err != nil {
		t.Fatal(err)
	}

	setting := qualitySizeByID(t, updated, "bluray-1080p")
	if setting.MinimumSizeMBPerMinute != 10.25 {
		t.Fatalf("minimum size = %v", setting.MinimumSizeMBPerMinute)
	}
	if setting.PreferredSizeMBPerMinute == nil || *setting.PreferredSizeMBPerMinute != preferred {
		t.Fatalf("preferred size = %#v", setting.PreferredSizeMBPerMinute)
	}
	if setting.MaximumSizeMBPerMinute != nil {
		t.Fatalf("expected nil maximum size, got %#v", setting.MaximumSizeMBPerMinute)
	}
}

func TestQualitySizeSettingsRollbackOnGeneratedQueryError(t *testing.T) {
	ctx, store := testDBStore(t)
	if _, err := store.SaveQualitySizeSettings(ctx, []QualitySizeSettingInput{{
		QualityID:              "unknown",
		MinimumSizeMBPerMinute: 9,
	}}); err != nil {
		t.Fatal(err)
	}

	_, err := store.SaveQualitySizeSettings(ctx, []QualitySizeSettingInput{
		{QualityID: "unknown", MinimumSizeMBPerMinute: 12},
		{QualityID: "cam", MinimumSizeMBPerMinute: 1000000000},
	})
	if err == nil || errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected database write error, got %v", err)
	}

	settings, err := store.ListQualitySizeSettings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	setting := qualitySizeByID(t, settings, "unknown")
	if setting.MinimumSizeMBPerMinute != 9 {
		t.Fatalf("expected rollback to keep minimum 9, got %v", setting.MinimumSizeMBPerMinute)
	}
}

func qualitySizeByID(t *testing.T, settings []QualitySizeSetting, id string) QualitySizeSetting {
	t.Helper()
	for _, setting := range settings {
		if setting.ID == id {
			return setting
		}
	}
	t.Fatalf("quality size %q not found", id)
	return QualitySizeSetting{}
}
