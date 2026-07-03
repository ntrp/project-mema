import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import DownloadClientsSettingsSection from '$lib/components/settings/download-clients/DownloadClientsSettingsSection.svelte';
import IndexersSettingsSection from '$lib/components/settings/indexers/IndexersSettingsSection.svelte';
import QualitySizeSettings from '$lib/components/settings/quality/QualitySizeSettings.svelte';
import SystemGeneralSettings from '$lib/components/settings/system/SystemGeneralSettings.svelte';
import { emptyDownloadClientForm, emptyIndexerForm } from '$lib/settings/forms';
import type { IndexerSearchResponse } from '$lib/settings/types';

describe('rendered integration settings sections (SCN-SETTINGS-018, SCN-SETTINGS-004)', () => {
	it('renders download client section actions and empty state', () => {
		const { body } = render(DownloadClientsSettingsSection, {
			props: {
				clients: [],
				form: emptyDownloadClientForm(),
				saving: false,
				onSave: vi.fn(),
				onTestConfig: vi.fn(),
				onCancel: vi.fn(),
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});

		expect(body).toContain('Download clients');
		expect(body).toContain('Add download client');
		expect(body).toContain('No download clients configured');
	});

	it('renders indexer section search settings and empty state', () => {
		const { body } = render(IndexersSettingsSection, {
			props: {
				indexers: [],
				indexerSearch: emptyIndexerSearch(),
				form: emptyIndexerForm(),
				saving: false,
				clearingIndexerSearchCache: false,
				savingIndexerSearchSettings: false,
				testResults: {},
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onEdit: vi.fn(),
				onDelete: vi.fn(),
				onTest: vi.fn(),
				onClearIndexerSearchCache: vi.fn(),
				onSaveIndexerSearchSettings: vi.fn()
			}
		});

		expect(body).toContain('Indexers');
		expect(body).toContain('Add indexer');
		expect(body).toContain('No indexers configured');
		expect(body).toContain('Indexer search settings');
		expect(body).toContain('Reset cache');
	});
});

describe('rendered general system settings sections (SCN-SYSTEM-006, SCN-SETTINGS-008)', () => {
	it('renders general system settings forms before async data loads', () => {
		const { body } = render(SystemGeneralSettings);

		expect(body).toContain('General');
		expect(body).toContain('Event retention days');
		expect(body).toContain('Write logs to files');
		expect(body).toContain('Log directory');
		expect(body).toContain('Save settings');
	});

	it('renders quality size loading table and save controls', () => {
		const { body } = render(QualitySizeSettings);

		expect(body).toContain('Quality sizes');
		expect(body).toContain('Release scoring');
		expect(body).toContain('Reload');
		expect(body).toContain('Save sizes');
		expect(body).toContain('Loading quality sizes');
	});
});

function emptyIndexerSearch(): IndexerSearchResponse {
	return {
		settings: {
			cacheDurationMinutes: 60,
			historyRetentionDays: 14,
			automaticBlocklistExpiryDays: 7
		},
		stats: { totalEntries: 0, activeEntries: 0, expiredEntries: 0, indexerCount: 0 },
		cacheEntries: [],
		historyEntries: [],
		historyTotalEntries: 0,
		historyStats: { totalEntries: 0, cacheHits: 0, cacheMisses: 0, failures: 0 }
	};
}
