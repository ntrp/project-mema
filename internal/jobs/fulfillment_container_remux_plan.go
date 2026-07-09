package jobs

import (
	"context"
	"fmt"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

type containerRemuxPlan struct {
	args            FulfillmentActionArgs
	fact            storage.MediaFileFact
	targetContainer string
}

func enqueueContainerRemuxFileJobs(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	enqueue fulfillmentEnqueueFunc,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) (int, error) {
	if enqueue == nil {
		return 0, fmt.Errorf("container remux enqueue function is not configured")
	}
	if args.TargetType != "" && args.TargetType != "video" {
		return 0, nil
	}
	plans := containerRemuxPlans(item, args)
	for _, plan := range plans {
		jobID, err := enqueue(ctx, "container_remux", plan.args)
		if err != nil {
			return 0, err
		}
		initializeQueuedContainerRemuxProgress(ctx, settings, eventBroker, jobID, item, plan.fact)
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "media", "Container remux queued", map[string]any{
			"mediaItemId":     item.ID.String(),
			"title":           item.Title,
			"filePath":        plan.args.FilePath,
			"jobId":           jobID,
			"targetContainer": plan.targetContainer,
		})
	}
	return len(plans), nil
}

func containerRemuxPlans(item storage.MediaItem, args FulfillmentActionArgs) []containerRemuxPlan {
	targetContainer := normalizedContainer(item.FinalContainer)
	if targetContainer == "" {
		return nil
	}
	plans := []containerRemuxPlan{}
	for _, fact := range item.FileFacts {
		if args.FilePath != "" && fact.FilePath != args.FilePath {
			continue
		}
		if normalizedContainer(mediaFactContainer(fact)) == targetContainer {
			continue
		}
		plans = append(plans, containerRemuxPlan{
			args: FulfillmentActionArgs{
				MediaItemID: item.ID.String(),
				FilePath:    fact.FilePath,
				TargetType:  "video",
			},
			fact:            fact,
			targetContainer: targetContainer,
		})
	}
	return plans
}
