package storage

import (
	"errors"
	"testing"
)

func TestLogFileSettingsUseGeneratedQueries(t *testing.T) {
	ctx, store := testDBStore(t)

	initial, err := store.GetLogFileSettings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if initial.Directory != DefaultLogFileDirectory || initial.RetentionDays != DefaultLogRetentionDays {
		t.Fatalf("initial log file settings = %#v", initial)
	}

	updated, err := store.UpdateLogFileSettings(ctx, LogFileSettingsInput{
		Enabled:       true,
		Directory:     " .data/custom-logs ",
		RetentionDays: 14,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !updated.Enabled || updated.Directory != ".data/custom-logs" || updated.RetentionDays != 14 {
		t.Fatalf("updated log file settings = %#v", updated)
	}

	if _, err := store.UpdateLogFileSettings(ctx, LogFileSettingsInput{
		Enabled:       true,
		Directory:     ".data/logs",
		RetentionDays: 366,
	}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("invalid log file settings error = %v, want %v", err, ErrInvalidInput)
	}
}
