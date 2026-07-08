package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

var (
	ErrFixedScheduleNotFound = errors.New("fixed schedule not found")
	ErrFixedScheduleActive   = errors.New("fixed schedule already has an active execution")
)

func (c *Client) Start(ctx context.Context) error {
	return c.river.Start(ctx)
}

func (c *Client) Stop(ctx context.Context) error {
	return c.river.Stop(ctx)
}

func (c *Client) AbortJob(ctx context.Context, id int64) error {
	_, err := c.river.JobCancel(ctx, id)
	return err
}

func (c *Client) EnqueueFixedSchedule(ctx context.Context, scheduleID string) (int64, error) {
	definition, ok := fixedJobDefinitionByID(scheduleID)
	if !ok {
		return 0, ErrFixedScheduleNotFound
	}
	if c.settings != nil {
		schedules, err := c.settings.ListSystemJobSchedules(ctx)
		if err != nil {
			return 0, err
		}
		found := false
		for _, schedule := range schedules {
			if schedule.ID != definition.ID {
				continue
			}
			found = true
			if schedule.ActiveRiverJobID != nil {
				return 0, ErrFixedScheduleActive
			}
			break
		}
		if !found {
			return 0, ErrFixedScheduleNotFound
		}
	}
	metadata, _ := json.Marshal(map[string]any{
		"app:manual_schedule_run": true,
		"app:system_schedule_id":  definition.ID,
	})
	result, err := c.river.Insert(ctx, definition.args(), &river.InsertOpts{
		Queue:    definition.Queue,
		Metadata: metadata,
		UniqueOpts: river.UniqueOpts{
			ByQueue: true,
			ByState: []rivertype.JobState{
				rivertype.JobStateAvailable,
				rivertype.JobStatePending,
				rivertype.JobStateRunning,
				rivertype.JobStateScheduled,
				rivertype.JobStateRetryable,
			},
		},
	})
	if err != nil {
		return 0, err
	}
	recordJobUpdated(ctx, c.settings, c.events, result.Job, "")
	return result.Job.ID, nil
}

func (c *Client) EnqueueReleaseSearch(ctx context.Context, mediaItemID uuid.UUID, query string) (int64, error) {
	result, err := c.river.Insert(ctx, ReleaseSearchArgs{MediaItemID: mediaItemID.String(), Query: strings.TrimSpace(query)}, &river.InsertOpts{
		Queue: queueMediaSearch,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	recordJobUpdated(ctx, c.settings, c.events, result.Job, "")
	return result.Job.ID, nil
}

func (c *Client) EnqueueAutoSearchDownload(ctx context.Context, mediaItemID uuid.UUID) (int64, error) {
	result, err := c.river.Insert(ctx, AutoSearchDownloadArgs{MediaItemID: mediaItemID.String()}, &river.InsertOpts{
		Queue: queueMediaSearch,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	recordJobUpdated(ctx, c.settings, c.events, result.Job, "")
	return result.Job.ID, nil
}

func (c *Client) EnqueueSubtitleSearch(ctx context.Context, args SubtitleSearchArgs) (int64, error) {
	result, err := c.river.Insert(ctx, args, &river.InsertOpts{
		Queue: queueMediaSearch,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	recordJobUpdated(ctx, c.settings, c.events, result.Job, "")
	return result.Job.ID, nil
}

func (c *Client) EnqueueGrabRelease(ctx context.Context, args GrabReleaseArgs) (int64, error) {
	result, err := c.river.Insert(ctx, args, &river.InsertOpts{
		Queue: queueDownloads,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	recordJobUpdated(ctx, c.settings, c.events, result.Job, "")
	return result.Job.ID, nil
}

func (c *Client) EnqueueMediaComponentExtraction(ctx context.Context, artifactID uuid.UUID) (int64, error) {
	result, err := c.river.Insert(ctx, MediaComponentExtractionArgs{ArtifactID: artifactID.String()}, &river.InsertOpts{
		Queue: queueMediaAssembly,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	recordJobUpdated(ctx, c.settings, c.events, result.Job, "")
	return result.Job.ID, nil
}

func (c *Client) EnqueueMediaComponentMux(ctx context.Context, runID uuid.UUID) (int64, error) {
	result, err := c.river.Insert(ctx, MediaComponentMuxArgs{RunID: runID.String()}, &river.InsertOpts{
		Queue: queueMediaAssembly,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	recordJobUpdated(ctx, c.settings, c.events, result.Job, "")
	return result.Job.ID, nil
}
