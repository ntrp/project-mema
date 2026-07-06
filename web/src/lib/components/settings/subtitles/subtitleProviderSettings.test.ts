import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import SubtitleProviderSettings from './SubtitleProviderSettings.svelte';
import { subtitleProvider } from '$lib/components/rendered/appShellTestValues';

describe('subtitle provider settings (SCN-SETTINGS-024)', () => {
	it('masks saved OpenSubtitles secrets', () => {
		const { body } = render(SubtitleProviderSettings, {
			props: {
				providers: [subtitleProvider()],
				onSave: vi.fn(),
				onDelete: vi.fn(),
				onTest: vi.fn(),
				testResults: {}
			}
		});

		expect(body).toContain('OpenSubtitles');
		expect(body).toContain('type="password"');
		expect(body).toContain('Show secret');
		expect(body).toContain('value="scenario-key"');
		expect(body).toContain('value="scenario-password"');
	});
});
