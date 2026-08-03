import { mount, tick, unmount, type Component } from 'svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	listIndexerCatalog: vi.fn(),
	listIndexerAppProfiles: vi.fn(),
	listIndexerProxies: vi.fn()
}));

vi.mock('$lib/settings/api', () => mocks);

import RenderWithTooltip from '$lib/components/rendered/RenderWithTooltip.svelte';
import IndexersSettingsSection from '../IndexersSettingsSection.svelte';
import { emptyIndexerForm } from '$lib/settings/forms';

let component: ReturnType<typeof mount> | undefined;
const target = document.createElement('div');

afterEach(async () => {
	if (component) await unmount(component);
	target.remove();
	component = undefined;
	vi.clearAllMocks();
});

describe('IndexersSettingsSection catalog', () => {
	it('shows loaded indexers and filter choices in the add modal', async () => {
		mocks.listIndexerCatalog.mockResolvedValue({
			entries: [catalogEntry()],
			protocols: ['torrent'],
			languages: ['English'],
			privacy: ['public'],
			categories: [{ id: 2000, name: 'Movies', children: [] }]
		});
		mocks.listIndexerAppProfiles.mockResolvedValue([]);
		mocks.listIndexerProxies.mockResolvedValue([]);
		document.body.append(target);
		component = mount(RenderWithTooltip, {
			target,
			props: {
				component: IndexersSettingsSection as unknown as Component<Record<string, unknown>>,
				componentProps: sectionProps()
			}
		});
		await tick();

		button('Add indexer').click();

		await vi.waitFor(() => {
			expect(document.body.textContent).toContain('Test Indexer');
		});
		button('All protocols', document.body).click();
		await tick();
		expect(document.body.textContent).toContain('torrent');
	});
});

function button(label: string, scope: Pick<typeof document.body, 'querySelectorAll'> = target) {
	const element = [...scope.querySelectorAll('button')].find(
		(candidate) => candidate.textContent?.trim() === label
	);
	if (!element) throw new Error(`missing button: ${label}`);
	return element;
}

function sectionProps() {
	return {
		indexers: [],
		indexerSearch: {
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
		},
		tags: [],
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
		onTestConfig: vi.fn(),
		onClearIndexerSearchCache: vi.fn(),
		onSaveIndexerSearchSettings: vi.fn()
	};
}

function catalogEntry() {
	return {
		definitionId: 'test-indexer',
		name: 'Test Indexer',
		implementation: 'Test',
		implementationName: 'Test',
		description: 'Test catalog entry',
		language: 'English',
		protocol: 'torrent',
		privacy: 'public',
		supportsRss: true,
		supportsSearch: true,
		supportsRedirect: false,
		supportsPagination: true,
		capabilities: {
			categories: [{ id: 2000, name: 'Movies', children: [] }],
			searchParams: ['q'],
			movieSearchParams: ['q'],
			tvSearchParams: ['q'],
			musicSearchParams: ['q'],
			bookSearchParams: ['q']
		},
		fields: []
	};
}
