import { describe, expect, it } from 'vitest';

import EventsPage from './system/events/+page.svelte';
import IndexingPage from './system/indexing/+page.svelte';
import JobsPage from './system/jobs/+page.svelte';
import LogsPage from './system/logs/+page.svelte';
import MetadataPage from './system/metadata/+page.svelte';
import StatusPage from './system/status/+page.svelte';
import { renderPage } from './routeTestHelpers';

describe('system route pages (SCN-SYSTEM-005)', () => {
	it.each([
		[StatusPage, ['Status', 'Loading system status']],
		[IndexingPage, ['Indexing', 'Query Cache', 'Query History']],
		[MetadataPage, ['Metadata', 'Metadata Cache', 'Query History']],
		[
			JobsPage,
			['Jobs', 'No one-shot jobs are running or scheduled.', 'Manual Fulfillment Actions']
		],
		[EventsPage, ['Events', 'Loading events']],
		[LogsPage, ['Logs', 'Waiting for log entries', 'No log files retained.']]
	])('renders the system route section', (component, expectedText) => {
		const { body } = renderPage(component);

		for (const text of expectedText) {
			expect(body).toContain(text);
		}
	});
});
