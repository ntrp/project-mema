package satisfaction

import (
	"fmt"
	"path/filepath"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

type WantedRowKind string

const (
	WantedRowMedia               WantedRowKind = "media"
	WantedRowTarget              WantedRowKind = "target"
	WantedRowCustomFormatUpgrade WantedRowKind = "custom_format_upgrade"
)

type WantedRow struct {
	ID                string
	Kind              WantedRowKind
	MediaItemID       string
	MediaTitle        string
	MediaType         string
	SeasonNumber      *int32
	EpisodeNumber     *int32
	FileLabel         string
	FilePath          string
	TargetType        targets.Type
	LanguageID        string
	TargetState       targets.State
	RequiredOperation *targets.Operation
	CurrentScore      *int32
	TargetScore       *int32
}

type WantedRowsInput struct {
	Item                 storage.MediaItem
	HasUsableMedia       bool
	Targets              []WantedTargetInput
	CustomFormatUpgrades []WantedCustomFormatUpgrade
}

type WantedTargetInput struct {
	Target        targets.Target
	FilePath      string
	SeasonNumber  *int32
	EpisodeNumber *int32
}

type WantedCustomFormatUpgrade struct {
	FilePath      string
	CurrentScore  int32
	TargetScore   int32
	SeasonNumber  *int32
	EpisodeNumber *int32
}

func BuildWantedRows(input WantedRowsInput) []WantedRow {
	rows := []WantedRow{}
	if !input.HasUsableMedia {
		rows = append(rows, mediaWantedRow(input.Item))
	}
	for _, target := range input.Targets {
		if !isWantedTargetState(target.Target.State) {
			continue
		}
		rows = append(rows, targetWantedRow(input.Item, target))
	}
	for _, upgrade := range input.CustomFormatUpgrades {
		if upgrade.TargetScore <= upgrade.CurrentScore {
			continue
		}
		rows = append(rows, customFormatWantedRow(input.Item, upgrade))
	}
	return rows
}

func mediaWantedRow(item storage.MediaItem) WantedRow {
	return WantedRow{
		ID:          "media:" + item.ID.String(),
		Kind:        WantedRowMedia,
		MediaItemID: item.ID.String(),
		MediaTitle:  item.Title,
		MediaType:   item.Type,
	}
}

func targetWantedRow(item storage.MediaItem, input WantedTargetInput) WantedRow {
	target := input.Target
	return WantedRow{
		ID:                "target:" + target.ID,
		Kind:              WantedRowTarget,
		MediaItemID:       item.ID.String(),
		MediaTitle:        item.Title,
		MediaType:         item.Type,
		SeasonNumber:      input.SeasonNumber,
		EpisodeNumber:     input.EpisodeNumber,
		FileLabel:         fileLabel(input.FilePath),
		FilePath:          input.FilePath,
		TargetType:        target.Type,
		LanguageID:        target.LanguageID,
		TargetState:       target.State,
		RequiredOperation: target.RequiredOperation,
	}
}

func customFormatWantedRow(item storage.MediaItem, input WantedCustomFormatUpgrade) WantedRow {
	return WantedRow{
		ID:            fmt.Sprintf("custom-format:%s:%s", item.ID.String(), input.FilePath),
		Kind:          WantedRowCustomFormatUpgrade,
		MediaItemID:   item.ID.String(),
		MediaTitle:    item.Title,
		MediaType:     item.Type,
		SeasonNumber:  input.SeasonNumber,
		EpisodeNumber: input.EpisodeNumber,
		FileLabel:     fileLabel(input.FilePath),
		FilePath:      input.FilePath,
		CurrentScore:  &input.CurrentScore,
		TargetScore:   &input.TargetScore,
	}
}

func isWantedTargetState(state targets.State) bool {
	switch state {
	case targets.StateMissing, targets.StatePartial, targets.StatePending, targets.StateBlocked, targets.StateFailed:
		return true
	default:
		return false
	}
}

func fileLabel(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}
