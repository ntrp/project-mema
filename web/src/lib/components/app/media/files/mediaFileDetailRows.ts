import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import type { MediaFileTrackDeleteRequest } from '$lib/settings/types';

type MediaFileTrack = MediaFileRow['tracks'][number];
type TrackType = MediaFileTrack['type'] | 'chapter';

export type TrackDeleteRequest = Omit<MediaFileTrackDeleteRequest, 'path'>;

export interface MediaFileDetailRow {
	key: string;
	trackNumber: string;
	type: TrackType;
	language: string;
	description: string;
	provenance?: MediaFileTrack['provenance'];
	chapterSummary?: boolean;
	missing?: boolean;
	unwanted?: boolean;
	deleteRequest?: TrackDeleteRequest;
}
