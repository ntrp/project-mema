package jobs

import (
	"fmt"

	"media-manager/internal/storage"
)

func fulfillmentApplyTrackScope(
	item storage.MediaItem,
	args FulfillmentActionArgs,
) (storage.MediaItem, FulfillmentActionArgs, *storage.MediaFileTrackFact, error) {
	if args.TrackID == "" {
		return item, args, nil, nil
	}
	for _, fact := range item.FileFacts {
		if args.FilePath != "" && fact.FilePath != args.FilePath {
			continue
		}
		for _, track := range fact.Tracks {
			if track.ID.String() != args.TrackID {
				continue
			}
			scopedFact := fact
			scopedTrack := track
			scopedFact.Tracks = []storage.MediaFileTrackFact{scopedTrack}
			item.FileFacts = []storage.MediaFileFact{scopedFact}
			args.FilePath = scopedTrack.FilePath
			if args.TargetType == "" {
				args.TargetType = scopedTrack.TrackType
			}
			if args.LanguageID == "" && scopedTrack.LanguageID != nil {
				args.LanguageID = *scopedTrack.LanguageID
			}
			return item, args, &scopedTrack, nil
		}
	}
	return item, args, nil, fmt.Errorf("fulfillment track not found: %s", args.TrackID)
}
