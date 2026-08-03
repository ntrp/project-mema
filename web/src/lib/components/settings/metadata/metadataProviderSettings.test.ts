import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MetadataProviderCard from './MetadataProviderCard.svelte';
import MetadataProviderSettings from './MetadataProviderSettings.svelte';
import { metadataProviderDefinitions } from './metadataProviderDefinitions';
import type {
	IntegrationTestResults,
	MetadataProvider,
	MetadataProviderForm
} from '$lib/settings/types';

const now = '2026-07-03T00:00:00Z';

function provider(overrides: Partial<MetadataProvider> = {}): MetadataProvider {
	return {
		id: 'provider-tmdb',
		name: 'TMDB Local',
		type: 'tmdb',
		baseUrl: 'http://metadata.test/tmdb',
		apiKey: 'api-key',
		accessToken: 'access-token',
		enabled: true,
		priority: 100,
		createdAt: now,
		updatedAt: now,
		...overrides,
		fallbackEnabled: overrides.fallbackEnabled ?? false,
		apiKeySet: overrides.apiKeySet ?? false,
		pinSet: overrides.pinSet ?? false,
		accessTokenSet: overrides.accessTokenSet ?? false
	};
}

describe('metadata provider settings (SCN-SETTINGS-014)', () => {
	it('renders configured and default provider cards with test state', () => {
		const testResults: IntegrationTestResults = {
			'provider-tmdb': {
				success: true,
				message: 'Mock metadata ready',
				latencyMs: 42,
				checkedAt: now,
				details: {}
			}
		};

		const { body } = render(MetadataProviderSettings, {
			props: {
				metadataProviders: [provider()],
				onSave: vi.fn((_form: MetadataProviderForm) => undefined),
				onTest: vi.fn(),
				testingId: 'provider-tvdb',
				savingId: 'provider-tmdb',
				testResults
			}
		});

		expect(body).toContain('TMDB');
		expect(body).toContain('TVDB');
		expect(body).toContain('AniList');
		expect(body).toContain('AniDB');
		expect(body).toContain('No credentials required. AniList uses its public GraphQL API.');
		expect(body).toContain('https://graphql.anilist.co');
		expect(body).toContain('Registered client name');
		expect(body).toContain('Client version');
		expect(body).toContain('http://api.anidb.net:9001');
		expect(body).toContain('http://metadata.test/tmdb');
		expect(body).toContain('https://api4.thetvdb.com/v4');
		expect(body).toContain('Enabled');
		expect(body).toContain('Fill missing metadata');
		expect(body).toContain('Test OK');
		expect(body).toContain('Mock metadata ready - 42 ms');
		expect(body).toContain('Saving');
		expect(body).toContain('PIN');
	});

	it('renders AniList without credential fields', () => {
		const definition = metadataProviderDefinitions.find(({ type }) => type === 'anilist');
		if (!definition) throw new Error('missing AniList definition');

		const { body } = render(MetadataProviderCard, {
			props: {
				definition,
				onSave: vi.fn(),
				onTest: vi.fn()
			}
		});

		expect(body).toContain('AniList');
		expect(body).toContain('https://graphql.anilist.co');
		expect(body).toContain('No credentials required');
		expect(body).toContain('Save');
		expect(body).not.toContain('Saving');
		expect(body).not.toContain('API key');
		expect(body).not.toContain('Access token');
	});

	it('renders disabled provider failure status and disabled test action', () => {
		const { body } = render(MetadataProviderCard, {
			props: {
				definition: {
					type: 'tvdb',
					name: 'TVDB',
					baseUrl: 'https://api4.thetvdb.com/v4',
					priority: 110,
					fields: 'tvdb'
				},
				provider: provider({
					id: 'provider-tvdb',
					name: 'TVDB Local',
					type: 'tvdb',
					baseUrl: 'http://metadata.test/tvdb',
					apiKey: undefined,
					pin: '1234',
					enabled: false,
					priority: 110
				}),
				onSave: vi.fn(),
				onTest: vi.fn(),
				testResult: {
					success: false,
					message: 'Unauthorized',
					latencyMs: 9,
					checkedAt: now,
					details: {}
				}
			}
		});

		expect(body).toContain('Disabled');
		expect(body).toContain('Test failed');
		expect(body).toContain('Unauthorized - 9 ms');
		expect(body).toContain('http://metadata.test/tvdb');
		expect(body).toContain('value="1234"');
	});
});
