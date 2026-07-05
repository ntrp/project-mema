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
		expect(body).toContain('Saved API key');
		expect(body).toContain('Saved password');
		expect(body).not.toContain('scenario-key');
	});
});
