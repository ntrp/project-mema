package jobs

import (
	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func fulfillmentActionDetails(operation targets.OperationType, args FulfillmentActionArgs) map[string]any {
	details := map[string]any{
		"operation":          operation,
		"mediaItemId":        args.MediaItemID,
		"filePath":           args.FilePath,
		"targetType":         args.TargetType,
		"languageId":         args.LanguageID,
		"trackId":            args.TrackID,
		"otherFileId":        args.OtherFileID,
		"externalSubtitleId": args.ExternalSubtitleID,
	}
	if args.Manual {
		details["manual"] = true
	}
	return details
}

func fulfillmentTrackCount(item storage.MediaItem) int {
	count := 0
	for _, fact := range item.FileFacts {
		count += len(fact.Tracks)
	}
	return count
}

func fulfillmentSkipReason(operation targets.OperationType, args FulfillmentActionArgs, target targets.Target) string {
	if args.TargetType != "" && string(target.Type) != args.TargetType {
		return "target type does not match request"
	}
	if args.LanguageID != "" && !satisfaction.LanguageMatches(target.LanguageID, args.LanguageID) {
		return "language does not match request"
	}
	if target.RequiredOperation == nil {
		return "target has no required operation"
	}
	if target.RequiredOperation.Type != operation {
		return "target requires different operation"
	}
	return "target not eligible"
}
