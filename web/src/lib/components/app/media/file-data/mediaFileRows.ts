import type { MediaFileUpgradeInfo } from '$lib/components/app/media/files/mediaFileUpgradeState';
import type { MediaItem, MediaItemSubtitle } from '$lib/settings/types';

type MediaFileTrack = NonNullable<NonNullable<MediaItem['files']>[number]['tracks']>[number];
type MediaFileChapter = NonNullable<NonNullable<MediaItem['files']>[number]['chapters']>[number];
type MediaFileSubtitleSatisfaction = NonNullable<
	NonNullable<MediaItem['files']>[number]['subtitleSatisfaction']
>;

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
	subtitleSatisfaction?: MediaFileSubtitleSatisfaction;
	externalSubtitles?: MediaItemSubtitle[];
	upgrade: MediaFileUpgradeInfo;
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
