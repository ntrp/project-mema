import type { FileNamingSettingsRequest } from './types';

export const defaultFileNamingTemplates: FileNamingSettingsRequest = {
	movieFileFormat: '{movie_title} ({release_year}) {quality_full}',
	movieFolderFormat: '{movie_title} ({release_year})',
	seriesEpisodeFormat:
		'{series_title} - S{season:00}E{episode:00} - {episode_title} {quality_full}',
	dailyEpisodeFormat: '{series_title} - {air_date} - {episode_title} {quality_full}',
	animeEpisodeFormat: '{series_title} - S{season:00}E{episode:00} - {episode_title} {quality_full}',
	seriesFolderFormat: '{series_title}',
	seasonFolderFormat: 'Season {season}',
	specialsFolderFormat: 'Specials'
};

export type FileNamingTemplateParameter = {
	param: string;
	example: string;
	description: string;
};

export const fileNamingTemplateParameterDetails: FileNamingTemplateParameter[] = [
	{
		param: 'movie_title',
		example: 'The Matrix',
		description: 'Movie title from the selected metadata provider'
	},
	{
		param: 'series_title',
		example: 'The Expanse',
		description: 'Series title from the selected metadata provider'
	},
	{
		param: 'release_title',
		example: 'The.Matrix.1999.1080p.BluRay.x264-GROUP',
		description: 'Original release title from the grabbed download'
	},
	{
		param: 'release_year',
		example: '1999',
		description: 'Original release year for the movie or series'
	},
	{ param: 'year', example: '2026', description: 'Year associated with the media item' },
	{ param: 'air_date', example: '2026-07-01', description: 'Episode air date' },
	{ param: 'season:0', example: '1', description: 'Season number padded to at least 1 digit' },
	{ param: 'season:00', example: '01', description: 'Season number padded to at least 2 digits' },
	{ param: 'episode:0', example: '3', description: 'Episode number padded to at least 1 digit' },
	{ param: 'episode:00', example: '03', description: 'Episode number padded to at least 2 digits' },
	{
		param: 'episode:000',
		example: '003',
		description: 'Episode number padded to at least 3 digits'
	},
	{
		param: 'episode_title',
		example: 'Dulcinea',
		description: 'Episode title for series, daily, and anime episodes'
	},
	{ param: 'absolute_episode', example: '124', description: 'Absolute episode number for anime' },
	{ param: 'quality', example: 'Bluray-1080p', description: 'Detected quality name' },
	{
		param: 'quality_full',
		example: 'Bluray-1080p Proper',
		description: 'Quality name with detected modifiers'
	},
	{ param: 'source', example: 'BluRay', description: 'Detected release source' },
	{ param: 'resolution', example: '1080p', description: 'Detected video resolution' },
	{ param: 'video_codec', example: 'x265', description: 'Detected video codec' },
	{ param: 'audio_codec', example: 'DTS', description: 'Detected audio codec' },
	{ param: 'audio_channels', example: '5.1', description: 'Detected audio channel layout' },
	{ param: 'languages', example: 'EN, DE', description: 'Detected audio or subtitle languages' },
	{
		param: 'custom_formats',
		example: 'HDR, IMAX',
		description: 'Matched custom format names included in rename templates'
	},
	{ param: 'edition', example: "Director's Cut", description: 'Detected movie or episode edition' },
	{ param: 'release_group', example: 'GROUP', description: 'Detected release group' },
	{ param: 'release_hash', example: 'A1B2C3D4', description: 'Stable hash for release uniqueness' }
];

export const fileNamingTemplateParameters = fileNamingTemplateParameterDetails.map(
	({ param }) => param
);

export function fileNamingTemplateSuggestions(query: string) {
	const normalized = query.toLowerCase();
	return fileNamingTemplateParameterDetails
		.map((parameter) => ({ parameter, score: fuzzyScore(parameter.param, normalized) }))
		.filter((match) => match.score !== null)
		.sort((left, right) => left.score! - right.score!)
		.map((match) => match.parameter)
		.slice(0, 8);
}

export function fileNamingTemplateExample(template: string) {
	return template.replace(/\{([^{}]+)\}/g, (token, rawKey: string) => {
		const value = exampleValue(rawKey);
		return value ?? token;
	});
}

function exampleValue(rawKey: string) {
	const [key, format] = rawKey.split(':');
	const examples: Record<string, string> = {
		movie_title: 'The Matrix',
		series_title: 'The Expanse',
		episode_title: 'Dulcinea',
		release_title: 'The.Matrix.1999.1080p.BluRay.x264-GROUP',
		release_year: '1999',
		year: '2026',
		air_date: '2026-07-01',
		season: '1',
		episode: '3',
		absolute_episode: '124',
		quality: 'Bluray-1080p',
		quality_full: 'Bluray-1080p Proper',
		source: 'BluRay',
		resolution: '1080p',
		video_codec: 'x265',
		audio_codec: 'DTS',
		audio_channels: '5.1',
		languages: 'EN, DE',
		custom_formats: 'HDR, IMAX',
		edition: "Director's Cut",
		release_group: 'GROUP',
		release_hash: 'A1B2C3D4'
	};
	const value = examples[key];
	if (!value) {
		return undefined;
	}
	if (!format?.startsWith('0')) {
		return value;
	}
	return value.padStart(format.length, '0');
}

function fuzzyScore(value: string, normalizedQuery: string) {
	if (normalizedQuery === '') {
		return 0;
	}

	let score = 0;
	let lastMatchIndex = -1;
	let searchStart = 0;
	const normalizedValue = value.toLowerCase();

	for (const character of normalizedQuery) {
		const matchIndex = normalizedValue.indexOf(character, searchStart);
		if (matchIndex < 0) {
			return null;
		}
		score += matchIndex + (lastMatchIndex >= 0 ? matchIndex - lastMatchIndex - 1 : 0);
		lastMatchIndex = matchIndex;
		searchStart = matchIndex + 1;
	}

	if (normalizedValue.startsWith(normalizedQuery)) {
		score -= 100;
	}
	if (normalizedValue.includes(normalizedQuery)) {
		score -= 50;
	}
	return score;
}
