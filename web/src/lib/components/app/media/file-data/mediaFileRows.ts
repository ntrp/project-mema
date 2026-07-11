import type { MediaFileUpgradeInfo } from '$lib/components/app/media/files/mediaFileUpgradeState';
import type { MediaItem, MediaItemSubtitle } from '$lib/settings/types';
type MediaFileTrack = NonNullable<NonNullable<MediaItem['files']>[number]['tracks']>[number];
type MediaFileChapter = NonNullable<NonNullable<MediaItem['files']>[number]['chapters']>[number];
type MediaFileSubtitleSatisfaction = NonNullable<
	NonNullable<MediaItem['files']>[number]['subtitleSatisfaction']
>;
type MediaFileRollup = NonNullable<NonNullable<MediaItem['files']>[number]['rollup']>;
type MediaFileRequirements = NonNullable<NonNullable<MediaItem['files']>[number]['requirements']>;
type MediaFileMissingTrack = NonNullable<
	NonNullable<MediaItem['files']>[number]['missingTracks']
>[number];
type MediaFileOtherFile = NonNullable<
	NonNullable<MediaItem['files']>[number]['otherFiles']
>[number];

export interface MediaFileRow {
	key: string;
	path?: string;
	relativePath: string;
	exists: boolean;
	seasonNumber?: number;
	episodeNumber?: number;
	episodeTitle?: string;
	videoCodec: string;
	audioInfo: string;
	size: string;
	sizeBytes?: number;
	languages: string;
	quality: string;
	formats: string[];
	tracks: MediaFileTrack[];
	chapters: MediaFileChapter[];
	otherFiles: MediaFileOtherFile[];
	subtitleSatisfaction?: MediaFileSubtitleSatisfaction;
	requirements?: MediaFileRequirements;
	missingTracks: MediaFileMissingTrack[];
	rollup?: MediaFileRollup;
	externalSubtitles?: MediaItemSubtitle[];
	upgrade: MediaFileUpgradeInfo;
	score: number;
}

export interface MediaFileGroup {
	key: string;
	title: string;
	rows: MediaFileRow[];
}
