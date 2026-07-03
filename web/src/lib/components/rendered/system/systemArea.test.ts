import { describe, expect, it, vi } from 'vitest';

import SystemArea from '$lib/components/app/system/SystemArea.svelte';
import { renderWithTooltip } from '$lib/components/rendered/renderHelpers';
import { emptyIndexerSearch, emptyMetadataCache } from '$lib/settings/api';
import type { SystemSection } from '$lib/settings/types';

describe('rendered system area sections (SCN-SYSTEM-001)', () => {
	it('renders the system status surface while status data loads', () => {
		const { body } = renderSystemArea('status');

		expect(body).toContain('System');
		expect(body).toContain('Status');
		expect(body).toContain('Loading system status');
		expect(body).toContain('Refresh');
	});

	it('renders live jobs filters and empty table state', () => {
		const { body } = renderSystemArea('jobs');

		expect(body).toContain('Jobs');
		expect(body).toContain('Live');
		expect(body).toContain('Status');
		expect(body).toContain('Queue');
		expect(body).toContain('Kind');
		expect(body).toContain('Search');
		expect(body).toContain('No jobs match the filters.');
	});

	it('renders live events controls and initial loading state', () => {
		const { body } = renderSystemArea('events');

		expect(body).toContain('Events');
		expect(body).toContain('Live');
		expect(body).toContain('Info');
		expect(body).toContain('Severity');
		expect(body).toContain('Category');
		expect(body).toContain('Loading events');
		expect(body).toContain('Clear all events');
	});

	it('renders log stream controls, reconnecting state, and retained file table', () => {
		const { body } = renderSystemArea('logs');

		expect(body).toContain('Logs');
		expect(body).toContain('Log stream reconnecting');
		expect(body).toContain('Clear logs');
		expect(body).toContain('Follow logs');
		expect(body).toContain('Verbosity');
		expect(body).toContain('Waiting for log entries');
		expect(body).toContain('Log files');
		expect(body).toContain('No log files retained.');
	});
});

function renderSystemArea(activeSection: SystemSection) {
	return renderWithTooltip(SystemArea, {
		activeSection,
		indexerSearch: emptyIndexerSearch(),
		metadataCache: emptyMetadataCache(),
		loadingIndexerSearch: false,
		loadingMetadataCache: false,
		clearingIndexerSearchCache: false,
		clearingMetadataCache: false,
		onClearIndexerSearchCache: vi.fn(),
		onClearIndexerSearchCachePattern: vi.fn(),
		onDeleteIndexerSearchCacheEntry: vi.fn(),
		onClearIndexerSearchHistory: vi.fn(),
		onLoadMoreIndexerSearchCache: vi.fn(),
		onLoadMoreIndexerSearchHistory: vi.fn(),
		onClearMetadataCache: vi.fn(),
		onClearMetadataCachePattern: vi.fn(),
		onDeleteMetadataCacheEntry: vi.fn(),
		onClearMetadataSearchHistory: vi.fn(),
		onLoadMoreMetadataCache: vi.fn(),
		onLoadMoreMetadataSearchHistory: vi.fn()
	});
}
