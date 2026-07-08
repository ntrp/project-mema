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
	episodeKey,
	missingRow,
	seasonNumberFromName
} from '$lib/components/app/media/files/mediaFileMissing';
import {
	audioInfo,
	matchToken,
	qualityInfo
} from '$lib/components/app/media/files/mediaFileParsing';
import { mediaFileUpgradeInfo } from '$lib/components/app/media/files/mediaFileUpgradeState';
import type {
	MediaFileGroup,
	MediaFileRow
} from '$lib/components/app/media/file-data/mediaFileRows';
import type { MediaItem, MediaItemSubtitle } from '$lib/settings/types';

export type {
	MediaFileGroup,
	MediaFileRow
} from '$lib/components/app/media/file-data/mediaFileRows';
export {
	episodeKey,
	missingRow,
	seasonNumberFromName
} from '$lib/components/app/media/files/mediaFileMissing';

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
		sizeBytes: info?.sizeBytes,
		languages: mediaFileLanguageInfo(name),
		quality,
		formats,
		tracks: info?.tracks ?? [],
		chapters: info?.chapters ?? [],
		otherFiles: info?.otherFiles ?? [],
		subtitleSatisfaction: info?.subtitleSatisfaction,
		requirements: info?.requirements,
		missingTracks: info?.missingTracks ?? [],
		rollup: info?.rollup,
		externalSubtitles: externalSubtitlesForPath(item.externalSubtitles ?? [], path),
		upgrade,
		score: 0
	};
}
function externalSubtitlesForPath(subtitles: MediaItemSubtitle[], path: string) {
	return subtitles.filter((subtitle) => sameSubtitleMediaBase(subtitle.filePath, path));
}

function sameSubtitleMediaBase(subtitlePath: string, mediaPath: string) {
	const subtitleBase = baseWithoutExtension(subtitlePath);
	const mediaBase = baseWithoutExtension(mediaPath);
	return subtitleBase.toLowerCase().startsWith(`${mediaBase.toLowerCase()}.`);
}

function baseWithoutExtension(path: string) {
	const name = fileName(path);
	const extensionIndex = name.lastIndexOf('.');
	return extensionIndex > 0 ? name.slice(0, extensionIndex) : name;
}
