import type { MediaItem } from '$lib/settings/types';

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
	languages: string;
	quality: string;
	formats: string[];
	score: number;
}

export interface MediaFileGroup {
	key: string;
	title: string;
	rows: MediaFileRow[];
}

const formatPatterns = [
	'Remux',
	'BluRay',
	'WEB-DL',
	'WEBDL',
	'WEBRip',
	'HDTV',
	'AMZN',
	'DSNP',
	'NF',
	'ATVP',
	'DD+',
	'TrueHD',
	'Atmos',
	'DTS',
	'HDR',
	'DV',
	'Proper',
	'Repack'
];

export function mediaFileGroups(item: MediaItem): MediaFileGroup[] {
	return item.type === 'series' ? seriesGroups(item) : movieGroups(item);
}

function movieGroups(item: MediaItem): MediaFileGroup[] {
	const rows = item.filePaths.length
		? item.filePaths.map((path) => fileRow(item, path))
		: [missingRow('movie-missing', item.title)];
	return [{ key: 'movie', title: 'Movie file', rows }];
}

function seriesGroups(item: MediaItem): MediaFileGroup[] {
	const rows = item.filePaths.map((path) => fileRow(item, path));
	const byEpisode = new Map(
		rows.map((row) => [episodeKey(row.seasonNumber, row.episodeNumber), row])
	);
	const groups = (item.seasons ?? [])
		.filter((season) => season.name.toLowerCase() !== 'specials')
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

function fileRow(item: MediaItem, path: string): MediaFileRow {
	const name = fileName(path);
	const formats = matchedFormats(name);
	return {
		key: path,
		path,
		relativePath: relativePath(item.mediaFolderPath, path),
		exists: true,
		...episodeParts(path),
		videoCodec: matchToken(name, ['x265', 'h265', 'hevc', 'x264', 'h264', 'avc']),
		audioInfo: audioInfo(name),
		size: '-',
		languages: languageInfo(name),
		quality: qualityInfo(name),
		formats,
		score: 0
	};
}

function missingRow(
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
		languages: '-',
		quality: '-',
		formats: [],
		score: 0
	};
}

function fileName(path: string) {
	return path.replaceAll('\\', '/').split('/').filter(Boolean).pop() ?? path;
}

function relativePath(root: string | undefined, path: string) {
	if (!root) return fileName(path);
	const normalizedRoot = root.replaceAll('\\', '/').replace(/\/+$/, '');
	const normalizedPath = path.replaceAll('\\', '/');
	return normalizedPath.startsWith(`${normalizedRoot}/`)
		? normalizedPath.slice(normalizedRoot.length + 1)
		: fileName(path);
}

function episodeParts(path: string) {
	const match = /s(\d{1,2})e(\d{1,3})/i.exec(path);
	if (!match) return {};
	return { seasonNumber: Number(match[1]), episodeNumber: Number(match[2]) };
}

function episodeKey(season?: number, episode?: number) {
	return `${season ?? 0}:${episode ?? 0}`;
}

function seasonNumberFromName(name: string) {
	const match = /(\d+)/.exec(name);
	return match ? Number(match[1]) : undefined;
}

function qualityInfo(value: string) {
	return matchToken(value, ['2160p', '1080p', '720p', '576p', '480p']);
}

function languageInfo(value: string) {
	if (/\bmulti\b/i.test(value)) return 'Multi';
	if (/\bdual\b/i.test(value)) return 'Dual';
	if (/\benglish\b/i.test(value)) return 'English';
	return '-';
}

function audioInfo(value: string) {
	const tokens = ['TrueHD', 'Atmos', 'DTS-HD', 'DTS', 'DDP', 'DD+', 'EAC3', 'AC3', 'AAC']
		.map((token) => matchToken(value, [token]))
		.filter((token) => token !== '-');
	return tokens.join(' ') || '-';
}

function matchedFormats(value: string) {
	return formatPatterns.filter((format) => new RegExp(format.replace('+', '\\+'), 'i').test(value));
}

function matchToken(value: string, tokens: string[]) {
	return tokens.find((token) => new RegExp(token, 'i').test(value)) ?? '-';
}
