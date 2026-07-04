import { describe, expect, it, vi } from 'vitest';

vi.mock('$lib/components/settings/system/cache/rowPulse.svelte', () => ({
	createRowPulse: () => ({
		update: vi.fn(),
		classFor: () => undefined
	})
}));

import MediaOverviewInfoCard from '$lib/components/app/media/detail/MediaOverviewInfoCard.svelte';
import SystemEventSeverityIcon from '$lib/components/settings/system/events/SystemEventSeverityIcon.svelte';
import SystemEventsControls from '$lib/components/settings/system/events/SystemEventsControls.svelte';
import SystemEventsTable from '$lib/components/settings/system/events/SystemEventsTable.svelte';
import SystemJobsTable from '$lib/components/settings/system/jobs/SystemJobsTable.svelte';
import SystemLogRow from '$lib/components/settings/system/logs/SystemLogRow.svelte';
import type { SystemJob } from '$lib/settings/types';
import { renderWithTooltip } from './renderHelpers';

describe('rendered system components (SCN-SYSTEM-006)', () => {
	it('renders event controls and disables destructive clearing when no events exist', () => {
		const { body } = renderWithTooltip(SystemEventsControls, {
			severityFilter: 'warning' as const,
			loading: false,
			clearing: true,
			eventsEmpty: true,
			onSeverityChange: vi.fn(),
			onClear: vi.fn()
		});

		expect(body).toContain('Warning');
		expect(body).toContain('Clearing events');
		expect(body).toContain('disabled');
		expect(body).toContain('data-size="sm"');
	});

	it('renders info event severity in blue', () => {
		const { body } = renderWithTooltip(SystemEventSeverityIcon, {
			severity: 'info' as const
		});

		expect(body).toContain('text-sky-600');
		expect(body).toContain('dark:text-sky-300');
	});

	it('renders warning event severity in yellow', () => {
		const { body } = renderWithTooltip(SystemEventSeverityIcon, {
			severity: 'warning' as const
		});

		expect(body).toContain('text-yellow-600');
		expect(body).toContain('dark:text-yellow-300');
	});

	it('renders event rows with severity, category, message, and error text', () => {
		const { body } = renderWithTooltip(SystemEventsTable, {
			events: [
				{
					id: 'event-1',
					createdAt: '2026-07-03T01:02:03Z',
					severity: 'error' as const,
					category: 'indexer',
					message: 'Indexer failed',
					data: { error: 'timeout' }
				}
			],
			loading: false,
			hasMore: true,
			loadingMore: true,
			onDelete: vi.fn(),
			onLoadMore: vi.fn()
		});

		expect(body).toContain('Indexer failed');
		expect(body).toContain('indexer');
		expect(body).toContain('timeout');
		expect(body).toContain('Loading more events');
	});

	it('renders job rows with abort affordances only for active jobs', () => {
		const { body } = renderWithTooltip(SystemJobsTable, {
			jobs: [
				systemJob({ id: 42, status: 'running', infoMessage: 'Searching releases' }),
				systemJob({ id: 43, status: 'completed', infoMessage: 'Done' })
			],
			abortingId: 42,
			onAbort: vi.fn()
		});

		expect(body).toContain('media.release_search');
		expect(body).toContain('Searching releases');
		expect(body).toContain('completed');
		expect(body.match(/Abort job/g)?.length).toBe(1);
	});

	it('renders log rows with levels, messages, and attributes', () => {
		const { body } = renderWithTooltip(SystemLogRow, {
			entry: {
				id: 'log-1',
				time: '2026-07-03T01:02:03Z',
				level: 'error' as const,
				message: 'Import failed',
				attributes: { path: '/downloads/movie.mkv' }
			}
		});

		expect(body).toContain('ERROR');
		expect(body).toContain('Import failed');
		expect(body).toContain('Log attributes');
	});
});

describe('rendered media components (SCN-MEDIA-003)', () => {
	it('renders media overview facts for movie releases', () => {
		const { body } = renderWithTooltip(MediaOverviewInfoCard, {
			detail: {
				title: 'Scenario Movie',
				type: 'movie' as const,
				externalProvider: 'tmdb',
				externalId: 'movie-1',
				status: 'Released',
				releaseDate: '2026-07-03',
				originalLanguage: 'en',
				voteAverage: 8.2
			},
			facts: [
				{ label: 'Revenue', value: '$100,000,000' },
				{ label: 'Budget', value: '$40,000,000' },
				{ label: 'Production Countries', value: 'United States, Germany' },
				{ label: 'Studios', value: 'Scenario Pictures\nExample Studios' },
				{ label: 'Digital Release Date', value: '2026-08-01' }
			]
		});

		expect(body).toContain('TMDb');
		expect(body).toContain('82%');
		expect(body).toContain('Released');
		expect(body).toContain('Cinema release date');
		expect(body).toContain('Digital release date');
		expect(body).toContain('Aug 1, 2026');
		expect(body).toContain('English');
		expect(body).toContain('Scenario Pictures');
	});

	it('renders series status and first-aired facts', () => {
		const { body } = renderWithTooltip(MediaOverviewInfoCard, {
			detail: {
				title: 'Scenario Series',
				type: 'series' as const,
				externalProvider: 'tmdb',
				externalId: 'series-1',
				status: 'Continuing',
				firstAirDate: '2025-02-14',
				originalLanguage: 'de'
			},
			facts: [{ label: 'Networks', value: 'Scenario Network' }]
		});

		expect(body).toContain('Continuing');
		expect(body).toContain('Feb 14, 2025');
		expect(body).toContain('German');
		expect(body).toContain('Scenario Network');
	});
});

function systemJob(overrides: Partial<SystemJob>): SystemJob {
	return {
		id: 1,
		status: 'running',
		kind: 'media.release_search',
		queue: 'media_search',
		attempt: 1,
		maxAttempts: 3,
		priority: 1,
		infoMessage: '',
		scheduledAt: '2026-07-03T01:02:03Z',
		createdAt: '2026-07-03T01:02:03Z',
		args: '',
		metadata: '',
		errors: '',
		...overrides
	};
}
