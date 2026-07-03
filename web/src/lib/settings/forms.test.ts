import { describe, expect, it } from 'vitest';
import {
	normalizeDownloadClientForm,
	normalizeIndexerForm,
	normalizeLibraryFolderForm,
	normalizeMetadataProviderForm,
	normalizePathMappingForm
} from './forms';

describe('settings form normalization', () => {
	it('SCN-SETTINGS-006 trims required values and omits blank optionals', () => {
		expect(
			normalizeDownloadClientForm({
				name: ' SAB ',
				type: 'sabnzbd',
				baseUrl: ' http://sab.local ',
				username: ' ',
				password: ' secret ',
				apiKey: '',
				category: ' movies ',
				enabled: true,
				priority: 5
			})
		).toEqual({
			name: 'SAB',
			type: 'sabnzbd',
			baseUrl: 'http://sab.local',
			username: undefined,
			password: 'secret',
			apiKey: undefined,
			category: 'movies',
			enabled: true,
			priority: 5
		});

		expect(
			normalizeIndexerForm({
				name: ' Torznab ',
				type: 'torznab',
				baseUrl: ' http://indexer.local ',
				apiKey: ' key ',
				categoriesText: '2000, bad, 2040',
				enabled: false,
				priority: 20
			})
		).toMatchObject({
			name: 'Torznab',
			baseUrl: 'http://indexer.local',
			apiKey: 'key',
			categories: [2000, 2040]
		});

		expect(
			normalizeMetadataProviderForm({
				name: ' TMDB ',
				type: 'tmdb',
				baseUrl: ' http://metadata.local ',
				apiKey: '',
				pin: ' 1234 ',
				accessToken: ' ',
				enabled: true,
				priority: 1
			})
		).toMatchObject({ name: 'TMDB', pin: '1234', accessToken: undefined });

		expect(normalizeLibraryFolderForm({ path: ' /media/movies ' })).toEqual({
			path: '/media/movies'
		});
		expect(normalizePathMappingForm({ clientPath: ' /downloads ', appPath: ' /data ' })).toEqual({
			clientPath: '/downloads',
			appPath: '/data'
		});
	});
});
