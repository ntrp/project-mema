import { describe, expect, it } from 'vitest';
import {
	emptySubtitleProviderForm,
	normalizeSubtitleProviderForm,
	subtitleProviderFormFromProvider
} from './providerIntegrationForms';
import type { SubtitleProvider, SubtitleProviderCatalogEntry } from './types';

const entry: SubtitleProviderCatalogEntry = {
	key: 'opensubtitlescom',
	displayName: 'OpenSubtitles.com',
	runtimeStatus: 'supported',
	runtimeMessage: 'Runtime supported.',
	mediaTypes: ['movie'],
	dependencies: {},
	outboundPolicy: {},
	fields: [
		{
			key: 'baseUrl',
			label: 'Base URL',
			type: 'text',
			persisted: true,
			options: ['https://api.opensubtitles.com']
		}
	]
};

describe('subtitle provider catalog forms', () => {
	it('seeds catalog defaults and keeps unsupported entries disabled', () => {
		const form = emptySubtitleProviderForm('opensubtitlescom', entry);
		expect(form.name).toBe('OpenSubtitles.com');
		expect(form.baseUrl).toBe('https://api.opensubtitles.com');
		expect(form.settings?.baseUrl?.stringValue).toBe('https://api.opensubtitles.com');
		expect(form.enabled).toBe(true);

		const disabled = emptySubtitleProviderForm('whisperai', {
			...entry,
			key: 'whisperai',
			runtimeStatus: 'unsupported'
		});
		expect(disabled.enabled).toBe(false);
	});

	it('does not echo saved provider secrets into editable fields', () => {
		const provider = {
			id: 'subtitle-1',
			name: 'OpenSubtitles',
			type: 'opensubtitlescom',
			catalogKey: 'opensubtitlescom',
			baseUrl: 'https://api.opensubtitles.com',
			settings: {},
			enabled: true,
			priority: 100,
			apiKeySet: true,
			passwordSet: true,
			secretFieldsSet: ['apiKey', 'password'],
			runtimeStatus: 'supported',
			runtimeMessage: 'Runtime supported.',
			mockSubtitles: [],
			createdAt: '2026-07-03T00:00:00Z',
			updatedAt: '2026-07-03T00:00:00Z'
		} satisfies SubtitleProvider;
		const form = subtitleProviderFormFromProvider(provider);
		expect(form.apiKey).toBe('');
		expect(form.password).toBe('');
		expect(form.secretFieldsSet).toEqual(['apiKey', 'password']);
	});

	it('normalizes secret settings and explicit clears', () => {
		const request = normalizeSubtitleProviderForm({
			...emptySubtitleProviderForm('opensubtitlescom', entry),
			secretSettings: { apiKey: ' new-key ' },
			clearSecretFields: ['password'],
			settings: { useHash: { booleanValue: true } }
		});
		expect(request.secretSettings).toEqual({ apiKey: 'new-key' });
		expect(request.clearSecretFields).toEqual(['password']);
		expect(request.settings?.useHash?.booleanValue).toBe(true);
	});
});
