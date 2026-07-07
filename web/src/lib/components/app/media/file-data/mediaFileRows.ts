import type { MediaFileUpgradeInfo } from '$lib/components/app/media/files/mediaFileUpgradeState';
import type { MediaItem, MediaItemSubtitle } from '$lib/settings/types';
import type { MediaFileAudioTargetOption } from '$lib/components/app/media/file-data/mediaFileProfiles';

type MediaFileTrack = NonNullable<NonNullable<MediaItem['files']>[number]['tracks']>[number];
type MediaFileChapter = NonNullable<NonNullable<MediaItem['files']>[number]['chapters']>[number];
type MediaFileSubtitleSatisfaction = NonNullable<
	NonNullable<MediaItem['files']>[number]['subtitleSatisfaction']
>;
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
	externalSubtitles?: MediaItemSubtitle[];
	upgrade: MediaFileUpgradeInfo;
	expectedAudioTargets: MediaFileAudioTargetOption[];
	expectedLanguages: string[];
	expectedRequiredLanguages: string[];
	expectedSubtitleLanguages: string[];
	removeNonEnabledLanguages: boolean;
	removeNonEnabledSubtitleLanguages: boolean;
	score: number;
}

export interface MediaFileGroup {
	key: string;
	title: string;
	rows: MediaFileRow[];
}
