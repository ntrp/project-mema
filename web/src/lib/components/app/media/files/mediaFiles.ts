import { mediaFileLanguageInfo } from '$lib/components/app/media/file-data/mediaFileLanguages';
import { matchedFormats } from '$lib/components/app/media/file-data/mediaFileFormats';
import {
	fileProfileSettings,
	type MediaFileProfileOption
} from '$lib/components/app/media/file-data/mediaFileProfiles';
import { mediaFileInfo, mediaFileSize } from '$lib/components/app/media/file-data/mediaFileSize';
import {
	episodeParts,
	fileName,
	relativePath
} from '$lib/components/app/media/files/mediaFilePath';
import {
	audioInfo,
	matchToken,
	qualityInfo
} from '$lib/components/app/media/files/mediaFileParsing';
import {
	mediaFileUpgradeInfo,
	type MediaFileUpgradeInfo
} from '$lib/components/app/media/files/mediaFileUpgradeState';
import type { MediaItem } from '$lib/settings/types';
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
	upgrade: MediaFileUpgradeInfo;
	expectedLanguages: string[];
	removeNonEnabledLanguages: boolean;
	score: number;
}
export interface MediaFileGroup {
	key: string;
	title: string;
	rows: MediaFileRow[];
}
export function mediaFileGroups(
	item: MediaItem,
	qualityProfiles: MediaFileProfileOption[] = []
): MediaFileGroup[] {
	return item.type === 'serie'
		? seriesGroups(item, qualityProfiles)
		: movieGroups(item, qualityProfiles);
}
function movieGroups(item: MediaItem, qualityProfiles: MediaFileProfileOption[]): MediaFileGroup[] {
	const rows = item.filePaths.length
		? item.filePaths.map((path) => fileRow(item, path, qualityProfiles))
		: [missingRow('movie-missing', item.title)];
	return [{ key: 'movie', title: 'Movie file', rows }];
}

function seriesGroups(
	item: MediaItem,
	qualityProfiles: MediaFileProfileOption[]
): MediaFileGroup[] {
	const rows = item.filePaths.map((path) => fileRow(item, path, qualityProfiles));
	const byEpisode = new Map(
		rows.map((row) => [episodeKey(row.seasonNumber, row.episodeNumber), row])
	);
	const groups = (item.seasons ?? [])
		.map((season, index) => {
			const seasonNumber = seasonNumberFromName(season.name) ?? index + 1;
			const rows = (season.episodes ?? []).map(
				(episode) =>
					byEpisode.get(episodeKey(seasonNumber, episode.episodeNumber)) ??
					missingRow(
						`s${seasonNumber}e${episode.episodeNumber}`,
						episode.name,
						seasonNumber,
						episode.episodeNumber
					)
			);
			return { key: `season-${seasonNumber}`, title: season.name, rows };
		})
		.filter((group) => group.rows.length > 0);

	for (const row of rows) {
		const key = `season-${row.seasonNumber ?? 0}`;
		if (groups.some((group) => group.rows.includes(row))) continue;
		let group = groups.find((item) => item.key === key);
		if (!group) {
			group = {
				key,
				title: row.seasonNumber ? `Season ${row.seasonNumber}` : 'Unmatched files',
				rows: []
			};
			groups.push(group);
		}
		group.rows.push(row);
	}
	return groups.length
		? groups
		: [
				{
					key: 'series-missing',
					title: 'Season 1',
					rows: [missingRow('series-missing', item.title, 1, 1)]
				}
			];
}
export function fileRow(
	item: MediaItem,
	path: string,
	qualityProfiles: MediaFileProfileOption[] = []
): MediaFileRow {
	const name = fileName(path);
	const formats = matchedFormats(name);
	const info = mediaFileInfo(item, path);
	const sizeBytes = info?.sizeBytes;
	const profile = fileProfileSettings(item, qualityProfiles);
	const exists = info?.status !== 'missing';
	const quality = qualityInfo(name);
	const upgrade = mediaFileUpgradeInfo(exists, quality, formats, profile.profile);
	return {
		key: path,
		path,
		relativePath: relativePath(item.mediaFolderPath, path),
		exists,
		...episodeParts(path),
		videoCodec: matchToken(name, ['x265', 'h265', 'hevc', 'x264', 'h264', 'avc']),
		audioInfo: audioInfo(name),
		size: mediaFileSize(item, path),
		sizeBytes,
		languages: mediaFileLanguageInfo(name),
		quality,
		formats,
		tracks: info?.tracks ?? [],
		chapters: info?.chapters ?? [],
		subtitleSatisfaction: info?.subtitleSatisfaction,
		upgrade,
		expectedLanguages: profile.expectedLanguages,
		removeNonEnabledLanguages: profile.removeNonEnabledLanguages,
		score: 0
	};
}
export function missingRow(
	key: string,
	title: string,
	seasonNumber?: number,
	episodeNumber?: number
): MediaFileRow {
	return {
		key,
		relativePath: '-',
		exists: false,
		seasonNumber,
		episodeNumber,
		episodeTitle: title,
		videoCodec: '-',
		audioInfo: '-',
		size: '-',
		sizeBytes: undefined,
		languages: '-',
		quality: '-',
		formats: [],
		tracks: [],
		chapters: [],
		subtitleSatisfaction: {
			state: 'missing',
			wantedLanguages: [],
			matchedLanguages: [],
			missingLanguages: []
		},
		upgrade: { state: 'missing', label: 'Missing', reasons: ['File is missing'] },
		expectedLanguages: [],
		removeNonEnabledLanguages: false,
		score: 0
	};
}
export function episodeKey(season?: number, episode?: number) {
	return `${season ?? 0}:${episode ?? 0}`;
}

export function seasonNumberFromName(name: string) {
	if (name.trim().toLowerCase() === 'specials') {
		return 0;
	}
	const match = /(\d+)/.exec(name);
	return match ? Number(match[1]) : undefined;
}
