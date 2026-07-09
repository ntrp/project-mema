import { describe, expect, it, vi } from 'vitest';

vi.mock('$lib/components/settings/system/cache/rowPulse.svelte', () => ({
	createRowPulse: () => ({
		update: vi.fn(),
		classFor: () => undefined
	})
}));

import SystemEventsSettings from '$lib/components/settings/system/events/SystemEventsSettings.svelte';
import SystemJobsSettings from '$lib/components/settings/system/jobs/SystemJobsSettings.svelte';
import SystemLogsSettings from '$lib/components/settings/system/logs/SystemLogsSettings.svelte';
import { renderWithTooltip } from './renderHelpers';

describe('rendered system settings views (SCN-SYSTEM-008)', () => {
	it('renders event observability controls and loading state', () => {
		const events = renderWithTooltip(SystemEventsSettings, {
			onConnectionChange: vi.fn()
		});
		expect(events.body).toContain('Events');
		expect(events.body).toContain('Live');
		expect(events.body).toContain('Info');
		expect(events.body).toContain('Loading events');
		expect(events.body).toContain('Clear all events');
	});

	it('renders job observability filters and empty results guidance', () => {
		const { body } = renderWithTooltip(SystemJobsSettings, {});
		expect(body).toContain('Jobs');
		expect(body).toContain('Fixed Scheduled Jobs');
		expect(body).toContain('Refresh');
		expect(body).toContain('Status');
		expect(body).toContain('Queue');
		expect(body).toContain('Kind');
		expect(body).toContain('Search');
		expect(body).toContain('No one-shot jobs are running or scheduled.');
		expect(body).not.toContain('Manual Fulfillment Actions');
	});

	it('renders live log controls and waiting state before events arrive', () => {
		const { body } = renderWithTooltip(SystemLogsSettings, {
			onConnectionChange: vi.fn()
		});
		expect(body).toContain('Clear logs');
		expect(body).toContain('Follow logs');
		expect(body).toContain('Verbosity');
		expect(body).toContain('INFO');
		expect(body).toContain('Application logs');
		expect(body).toContain('Waiting for log entries');
	});
});
