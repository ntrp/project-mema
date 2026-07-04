import type {
	DownloadClient,
	DownloadClientForm as DownloadClientFormValue,
	Indexer,
	IntegrationTestResponse
} from '$lib/settings/types';

export function downloadClient(overrides: Partial<DownloadClient> = {}): DownloadClient {
	return {
		id: 'client-1',
		name: 'Local SABnzbd',
		type: 'sabnzbd',
		protocol: 'usenet',
		baseUrl: 'http://sabnzbd.local',
		enabled: true,
		priority: 50,
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

export function downloadClientForm(
	overrides: Partial<DownloadClientFormValue> = {}
): DownloadClientFormValue {
	return {
		name: 'Local client',
		type: 'sabnzbd',
		protocol: 'usenet',
		baseUrl: 'http://client.local',
		apiKey: 'secret',
		username: 'admin',
		password: 'password',
		category: 'movies',
		enabled: true,
		priority: 50,
		...overrides
	};
}

export function indexerRow(overrides: Partial<Indexer> = {}): Indexer {
	return {
		id: 'indexer-1',
		name: 'Local Torznab',
		definitionId: 'generic-torznab',
		baseUrl: 'http://torznab.local',
		categories: [5000, 2000],
		protocol: 'torrent',
		privacy: 'private',
		language: 'en-US',
		supportsRss: true,
		supportsSearch: true,
		supportsRedirect: true,
		supportsPagination: true,
		capabilities: {
			categories: [],
			supportsRawSearch: true,
			searchParams: ['q'],
			tvSearchParams: ['q'],
			movieSearchParams: ['q']
		},
		enabled: true,
		priority: 10,
		healthStatus: 'healthy',
		failureCount: 0,
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

export function integrationResult(
	overrides: Partial<IntegrationTestResponse> = {}
): IntegrationTestResponse {
	return {
		success: true,
		message: 'Connection ok',
		latencyMs: 42,
		checkedAt: '2026-07-03T00:00:00Z',
		details: {},
		...overrides
	};
}
