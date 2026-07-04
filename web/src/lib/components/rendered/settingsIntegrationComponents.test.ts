import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import DownloadClientForm from '$lib/components/settings/download-clients/DownloadClientForm.svelte';
import DownloadClientTable from '$lib/components/settings/download-clients/DownloadClientTable.svelte';
import IndexerHealthStatus from '$lib/components/settings/indexers/IndexerHealthStatus.svelte';
import IndexerTable from '$lib/components/settings/indexers/IndexerTable.svelte';
import {
	downloadClient,
	downloadClientForm,
	indexerRow,
	integrationResult
} from '$lib/components/rendered/settingsIntegrationTestValues';

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
		expect(populated.body).toContain('usenet');
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
		expect(transmission.body).toContain('TORRENT');
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
		expect(sabnzbd.body).toContain('USENET');
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
		expect(populated.body.indexOf('Protocol')).toBeLessThan(populated.body.indexOf('Name'));
		expect(populated.body).toContain('Local Torznab');
		expect(populated.body).toContain('torrent');
		expect(populated.body).toContain('private');
		expect(populated.body).not.toContain('en-US');
		expect(populated.body).toContain('text-emerald-700');
		expect(populated.body).toContain('border-destructive/50');
		expect(populated.body).toContain('5000, 2000');
		expect(populated.body).toContain('Checking');
		expect(populated.body).not.toContain('Base URL');
		expect(populated.body).not.toContain('http://torznab.local');

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
		expect(backingOff.body).toContain('Temp blocked');
		expect(backingOff.body).toContain('text-yellow-700');
		expect(backingOff.body).toContain('Next check');
		expect(backingOff.body).toContain('HTTP 429');

		const permanentlyBlocked = render(IndexerHealthStatus, {
			props: { indexer: indexerRow({ healthStatus: 'disabled' }) }
		});
		expect(permanentlyBlocked.body).toContain('Permanently blocked');
		expect(permanentlyBlocked.body).toContain('text-destructive');
	});
});
