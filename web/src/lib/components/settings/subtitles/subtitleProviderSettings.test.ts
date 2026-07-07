import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import SubtitleProviderSettings from './SubtitleProviderSettings.svelte';
import { subtitleProvider } from '$lib/components/rendered/appShellTestValues';

describe('subtitle provider settings (SCN-SETTINGS-024)', () => {
	it('does not show unsaved providers as testing', () => {
		const { body } = render(SubtitleProviderSettings, {
			props: {
				providers: [],
				onSave: vi.fn(),
				onDelete: vi.fn(),
				onTest: vi.fn(),
				testResults: {}
			}
		});

		expect(body).toContain('OpenSubtitles');
		expect(body).toContain('Mock subtitles');
		expect(body).not.toContain('Testing');
		expect(body).toContain('Test');
	});

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

	it('renders mock subtitle provider rows', () => {
		const { body } = render(SubtitleProviderSettings, {
			props: {
				providers: [
					subtitleProvider({
						id: 'mock-subtitle-1',
						name: 'Mock Subtitles',
						type: 'mock',
						baseUrl: 'mock://subtitles',
						apiKey: undefined,
						password: undefined,
						apiKeySet: false,
						passwordSet: false,
						mockSubtitles: [
							{
								id: 'mock-row-1',
								title: 'Scenario Movie',
								languageId: 'english',
								format: 'vtt'
							}
						]
					})
				],
				onSave: vi.fn(),
				onDelete: vi.fn(),
				onTest: vi.fn(),
				testResults: {}
			}
		});

		expect(body).toContain('Mock subtitles');
		expect(body).toContain('value="Scenario Movie"');
		expect(body).toContain('value="english"');
		expect(body).toContain('VTT');
	});
});
