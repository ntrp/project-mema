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

export const fileNamingTemplateParameters = [
	'movie_title',
	'series_title',
	'episode_title',
	'release_title',
	'release_year',
	'year',
	'air_date',
	'season',
	'episode',
	'absolute_episode',
	'quality',
	'quality_full',
	'source',
	'resolution',
	'video_codec',
	'audio_codec',
	'audio_channels',
	'languages',
	'custom_formats',
	'edition',
	'release_group',
	'release_hash'
];
