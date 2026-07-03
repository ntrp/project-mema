import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import DownloadClientForm from '$lib/components/settings/download-clients/DownloadClientForm.svelte';
import DownloadClientTable from '$lib/components/settings/download-clients/DownloadClientTable.svelte';
import IndexerHealthStatus from '$lib/components/settings/indexers/IndexerHealthStatus.svelte';
import IndexerTable from '$lib/components/settings/indexers/IndexerTable.svelte';
import type {
	DownloadClient,
	DownloadClientForm as DownloadClientFormValue,
	Indexer,
	IntegrationTestResponse
} from '$lib/settings/types';

describe('rendered download client settings components (SCN-INTEGRATIONS-003)', () => {
	it('renders download client rows and empty state', () => {
		const populated = render(DownloadClientTable, {
			props: {
				clients: [downloadClient()],
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});
		expect(populated.body).toContain('Local SABnzbd');
		expect(populated.body).toContain('sabnzbd');
		expect(populated.body).toContain('http://sabnzbd.local');
		expect(populated.body).toContain('Delete Local SABnzbd');

		const empty = render(DownloadClientTable, {
			props: { clients: [], onEdit: vi.fn(), onDelete: vi.fn() }
		});
		expect(empty.body).toContain('No download clients configured');
	});

	it('renders provider-specific download client form fields and test status', () => {
		const transmission = render(DownloadClientForm, {
			props: {
				form: downloadClientForm({ type: 'transmission' }),
				saving: false,
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onTest: vi.fn(),
				showTypeSelect: false
			}
		});
		expect(transmission.body).toContain('Username');
		expect(transmission.body).toContain('Password');
		expect(transmission.body).not.toContain('API key');

		const sabnzbd = render(DownloadClientForm, {
			props: {
				form: downloadClientForm({ id: 'client-1', type: 'sabnzbd' }),
				saving: true,
				testing: true,
				testResult: integrationResult(),
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onTest: vi.fn()
			}
		});
		expect(sabnzbd.body).toContain('Edit download client');
		expect(sabnzbd.body).toContain('API key');
		expect(sabnzbd.body).toContain('Testing');

		const tested = render(DownloadClientForm, {
			props: {
				form: downloadClientForm({ id: 'client-1', type: 'sabnzbd' }),
				saving: false,
				testing: false,
				testResult: integrationResult(),
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onTest: vi.fn()
			}
		});
		expect(tested.body).toContain('Test OK');
		expect(tested.body).toContain('Connection ok - 42 ms');
	});
});

describe('rendered indexer settings components (SCN-INTEGRATIONS-001)', () => {
	it('renders indexer health details, test state, and empty state', () => {
		const indexer = indexerRow();
		const populated = render(IndexerTable, {
			props: {
				indexers: [indexer],
				testingId: indexer.id,
				testResults: { [indexer.id]: integrationResult() },
				onEdit: vi.fn(),
				onDelete: vi.fn(),
				onTest: vi.fn()
			}
		});
		expect(populated.body).toContain('Local Torznab');
		expect(populated.body).toContain('5000, 2000');
		expect(populated.body).toContain('Checking');

		const checked = render(IndexerTable, {
			props: {
				indexers: [indexer],
				testResults: { [indexer.id]: integrationResult() },
				onEdit: vi.fn(),
				onDelete: vi.fn(),
				onTest: vi.fn()
			}
		});
		expect(checked.body).toContain('Connection ok - 42 ms');

		const empty = render(IndexerTable, {
			props: { indexers: [], testResults: {}, onEdit: vi.fn(), onDelete: vi.fn(), onTest: vi.fn() }
		});
		expect(empty.body).toContain('No indexers configured');
	});

	it('renders disabled and backing-off indexer health labels', () => {
		const disabled = render(IndexerHealthStatus, {
			props: { indexer: indexerRow({ enabled: false }) }
		});
		expect(disabled.body).toContain('Disabled');
		expect(disabled.body).toContain('No query yet');

		const backingOff = render(IndexerHealthStatus, {
			props: {
				indexer: indexerRow({
					healthStatus: 'temporary_disabled',
					nextCheckAt: '2026-07-03T05:00:00Z',
					lastStatusCode: 429
				})
			}
		});
		expect(backingOff.body).toContain('Backing off');
		expect(backingOff.body).toContain('Next check');
		expect(backingOff.body).toContain('HTTP 429');
	});
});

function downloadClient(overrides: Partial<DownloadClient> = {}): DownloadClient {
	return {
		id: 'client-1',
		name: 'Local SABnzbd',
		type: 'sabnzbd',
		baseUrl: 'http://sabnzbd.local',
		enabled: true,
		priority: 50,
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

function downloadClientForm(
	overrides: Partial<DownloadClientFormValue> = {}
): DownloadClientFormValue {
	return {
		name: 'Local client',
		type: 'sabnzbd',
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

function indexerRow(overrides: Partial<Indexer> = {}): Indexer {
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

function integrationResult(
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
