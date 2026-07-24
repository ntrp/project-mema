import type { MetadataProviderType } from '$lib/settings/types';

export interface MetadataProviderDefinition {
	type: MetadataProviderType;
	name: string;
	baseUrl: string;
	priority: number;
	fields: 'tmdb' | 'tvdb' | 'anilist' | 'anidb';
}

export const metadataProviderDefinitions: MetadataProviderDefinition[] = [
	{
		type: 'tmdb',
		name: 'TMDB',
		baseUrl: 'https://api.themoviedb.org/3',
		priority: 100,
		fields: 'tmdb'
	},
	{
		type: 'tvdb',
		name: 'TVDB',
		baseUrl: 'https://api4.thetvdb.com/v4',
		priority: 110,
		fields: 'tvdb'
	},
	{
		type: 'anilist',
		name: 'AniList',
		baseUrl: 'https://graphql.anilist.co',
		priority: 120,
		fields: 'anilist'
	},
	{
		type: 'anidb',
		name: 'AniDB',
		baseUrl: 'http://api.anidb.net:9001',
		priority: 120,
		fields: 'anidb'
	}
];
