package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) ListMediaComponentArtifactsForSource(
	ctx context.Context,
	sourceID uuid.UUID,
) ([]MediaComponentArtifact, error) {
	return listMediaComponentArtifactsForSource(ctx, s.pool, sourceID)
}

func (s *SettingsStore) GetMediaComponentArtifact(
	ctx context.Context,
	artifactID uuid.UUID,
) (MediaComponentArtifact, error) {
	row, err := storagegen.New(s.pool).GetMediaComponentArtifact(ctx, artifactID)
	return mediaComponentArtifactRow(row, err)
}

func (s *SettingsStore) CreateMediaComponentArtifact(
	ctx context.Context,
	mediaItemID uuid.UUID,
	sourceID uuid.UUID,
	input MediaComponentArtifactInput,
) (MediaComponentArtifact, error) {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return MediaComponentArtifact{}, err
	}
	source, err := s.GetMediaComponentSource(ctx, mediaItemID, sourceID)
	if err != nil {
		return MediaComponentArtifact{}, err
	}
	target, id, err := componentArtifactTarget(item, source, input)
	if err != nil {
		return MediaComponentArtifact{}, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return MediaComponentArtifact{}, err
	}
	row, err := storagegen.New(s.pool).CreateMediaComponentArtifact(ctx, storagegen.CreateMediaComponentArtifactParams{
		ID:          id,
		MediaItemID: mediaItemID,
		SourceID:    sourceID,
		StreamID:    input.StreamID,
		StreamType:  normalizeComponentArtifactStreamType(input.StreamType),
		Language:    textValue(input.Language),
		OutputPath:  target,
		JobID:       textValue(input.JobID),
	})
	return mediaComponentArtifactRow(row, err)
}

func (s *SettingsStore) AssignMediaComponentArtifactJob(
	ctx context.Context,
	artifactID uuid.UUID,
	jobID string,
) (MediaComponentArtifact, error) {
	row, err := storagegen.New(s.pool).AssignMediaComponentArtifactJob(ctx, storagegen.AssignMediaComponentArtifactJobParams{
		ID:    artifactID,
		JobID: textValue(&jobID),
	})
	return mediaComponentArtifactRow(row, err)
}

func (s *SettingsStore) StartMediaComponentArtifact(
	ctx context.Context,
	artifactID uuid.UUID,
) (MediaComponentArtifact, error) {
	row, err := storagegen.New(s.pool).StartMediaComponentArtifact(ctx, artifactID)
	return mediaComponentArtifactRow(row, err)
}

func (s *SettingsStore) CompleteMediaComponentArtifact(
	ctx context.Context,
	artifactID uuid.UUID,
	toolSummary string,
) (MediaComponentArtifact, error) {
	artifact, err := s.GetMediaComponentArtifact(ctx, artifactID)
	if err != nil {
		return MediaComponentArtifact{}, err
	}
	info, err := os.Stat(artifact.OutputPath)
	if err != nil || info.IsDir() {
		return MediaComponentArtifact{}, ErrInvalidInput
	}
	size := info.Size()
	row, err := storagegen.New(s.pool).CompleteMediaComponentArtifact(ctx, storagegen.CompleteMediaComponentArtifactParams{
		ID:          artifactID,
		ToolSummary: strings.TrimSpace(toolSummary),
		SizeBytes:   int8Value(&size),
	})
	return mediaComponentArtifactRow(row, err)
}

func (s *SettingsStore) FailMediaComponentArtifact(
	ctx context.Context,
	artifactID uuid.UUID,
	toolSummary string,
	errMessage string,
) (MediaComponentArtifact, error) {
	errMessage = strings.TrimSpace(errMessage)
	row, err := storagegen.New(s.pool).FailMediaComponentArtifact(ctx, storagegen.FailMediaComponentArtifactParams{
		ID:           artifactID,
		ToolSummary:  strings.TrimSpace(toolSummary),
		ErrorMessage: textValue(&errMessage),
	})
	return mediaComponentArtifactRow(row, err)
}

func listMediaComponentArtifactsForSource(
	ctx context.Context,
	q storagegen.DBTX,
	sourceID uuid.UUID,
) ([]MediaComponentArtifact, error) {
	rows, err := storagegen.New(q).ListMediaComponentArtifactsForSource(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	artifacts := make([]MediaComponentArtifact, 0, len(rows))
	for _, row := range rows {
		artifacts = append(artifacts, mediaComponentArtifactFromRow(row))
	}
	return artifacts, nil
}

func componentArtifactTarget(
	item MediaItem,
	source MediaComponentSource,
	input MediaComponentArtifactInput,
) (string, uuid.UUID, error) {
	if source.RetentionState != "retained" || !strings.EqualFold(filepath.Ext(source.RetainedPath), ".mkv") {
		return "", uuid.Nil, ErrInvalidInput
	}
	streamType := normalizeComponentArtifactStreamType(input.StreamType)
	if streamType == "" || input.StreamID < 0 || !componentArtifactStreamAllowed(source, input) {
		return "", uuid.Nil, ErrInvalidInput
	}
	if _, err := mediaComponentSourceTarget(item, source.RetainedPath); err != nil {
		return "", uuid.Nil, err
	}
	id := uuid.New()
	extension := ".mka"
	if streamType == "subtitle" {
		extension = ".mks"
	}
	target, err := mediaComponentSourceTarget(
		item,
		filepath.Join(
			".mema",
			"component-artifacts",
			source.ID.String(),
			id.String(),
			"stream-"+strconv.Itoa(int(input.StreamID))+extension,
		),
	)
	return target, id, err
}

func normalizeComponentArtifactStreamType(value string) string {
	switch strings.TrimSpace(value) {
	case "audio", "subtitle":
		return strings.TrimSpace(value)
	default:
		return ""
	}
}

func componentArtifactStreamAllowed(source MediaComponentSource, input MediaComponentArtifactInput) bool {
	streamType := normalizeComponentArtifactStreamType(input.StreamType)
	inventory := strings.TrimSpace(source.StreamInventory)
	if inventory == "" || (!strings.HasPrefix(inventory, "{") && !strings.HasPrefix(inventory, "[")) {
		return true
	}
	streams, ok := componentArtifactInventoryStreams(inventory)
	if !ok {
		return false
	}
	for _, stream := range streams {
		if stream.ID != input.StreamID || normalizeComponentArtifactStreamType(stream.Type) != streamType {
			continue
		}
		if input.Language == nil || strings.EqualFold(strings.TrimSpace(stream.Language), strings.TrimSpace(*input.Language)) {
			return true
		}
	}
	return false
}

func componentArtifactInventoryStreams(inventory string) ([]componentArtifactInventoryStream, bool) {
	var list []componentArtifactInventoryStream
	if err := json.Unmarshal([]byte(inventory), &list); err == nil {
		return normalizeComponentArtifactInventory(list), true
	}
	var payload struct {
		Streams []componentArtifactInventoryStream `json:"streams"`
	}
	if err := json.Unmarshal([]byte(inventory), &payload); err != nil {
		return nil, false
	}
	return normalizeComponentArtifactInventory(payload.Streams), true
}

type componentArtifactInventoryStream struct {
	ID        int32  `json:"id"`
	Index     *int32 `json:"index"`
	Type      string `json:"type"`
	CodecType string `json:"codec_type"`
	Language  string `json:"language"`
}

func normalizeComponentArtifactInventory(values []componentArtifactInventoryStream) []componentArtifactInventoryStream {
	streams := make([]componentArtifactInventoryStream, 0, len(values))
	for _, value := range values {
		if value.Index != nil {
			value.ID = *value.Index
		}
		if value.Type == "" {
			value.Type = value.CodecType
		}
		streams = append(streams, value)
	}
	return streams
}

func mediaComponentArtifactRow(row storagegen.AppMediaComponentArtifact, err error) (MediaComponentArtifact, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaComponentArtifact{}, ErrNotFound
	}
	if err != nil {
		return MediaComponentArtifact{}, err
	}
	return mediaComponentArtifactFromRow(row), nil
}

func mediaComponentArtifactFromRow(row storagegen.AppMediaComponentArtifact) MediaComponentArtifact {
	return MediaComponentArtifact{
		ID:           row.ID,
		MediaItemID:  row.MediaItemID,
		SourceID:     row.SourceID,
		StreamID:     row.StreamID,
		StreamType:   row.StreamType,
		Language:     textPtr(row.Language),
		OutputPath:   row.OutputPath,
		Status:       row.Status,
		ToolName:     row.ToolName,
		ToolSummary:  row.ToolSummary,
		ErrorMessage: textPtr(row.ErrorMessage),
		JobID:        textPtr(row.JobID),
		SizeBytes:    int8Ptr(row.SizeBytes),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		CompletedAt:  row.CompletedAt,
	}
}
