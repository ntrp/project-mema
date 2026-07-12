import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import SubtitleProviderCatalogPicker from './catalog/SubtitleProviderCatalogPicker.svelte';
import SubtitleProviderForm from './form/SubtitleProviderForm.svelte';
import SubtitleProviderSettings from './SubtitleProviderSettings.svelte';
import { subtitleProvider } from '$lib/components/rendered/appShellTestValues';
import { renderWithTooltip } from '$lib/components/rendered/renderHelpers';
import { emptySubtitleProviderForm } from '$lib/settings/forms';
import type { SubtitleProviderCatalogEntry } from '$lib/settings/types';

const openSubtitlesEntry: SubtitleProviderCatalogEntry = {
	key: 'opensubtitlescom',
	displayName: 'OpenSubtitles.com',
	provenanceCommit: 'e54edd769b7062280118a14aa0fef3808829714d',
	runtimeStatus: 'supported',
	runtimeMessage: 'Runtime supported through the OpenSubtitles.com HTTP API.',
	mediaTypes: ['movie', 'serie'],
	dependencies: {},
	outboundPolicy: {},
	fields: [
		{ key: 'baseUrl', label: 'Base URL', type: 'text', required: true, persisted: true },
		{ key: 'username', label: 'Username', type: 'text', persisted: true },
		{ key: 'apiKey', label: 'API key', type: 'password', secret: true, persisted: true },
		{ key: 'password', label: 'Password', type: 'password', secret: true, persisted: true }
	]
};

const whisperEntry: SubtitleProviderCatalogEntry = {
	key: 'whisperai',
	displayName: 'Whisper',
	runtimeStatus: 'unsupported',
	runtimeMessage: 'Requires a reviewed local Whisper service integration before runtime use.',
	mediaTypes: ['movie', 'serie'],
	dependencies: { local_http_endpoint: true },
	outboundPolicy: { allowLocalHosts: true },
	fields: []
};

describe('subtitle provider settings (SCN-SETTINGS-024)', () => {
	it('renders configured providers as a table without marking unsaved entries as testing', () => {
		const { body } = renderWithTooltip(SubtitleProviderSettings, {
			providers: [subtitleProvider()],
			onSave: vi.fn(),
			onDelete: vi.fn(),
			onTest: vi.fn(),
			onTestConfig: vi.fn(),
			testResults: {}
		});

		expect(body).toContain('Add subtitle provider');
		expect(body).toContain('OpenSubtitles');
		expect(body).toContain('Supported');
		expect(body).not.toContain('Testing');
	});

	it('shows catalog picker runtime states and dependencies', () => {
		const { body } = render(SubtitleProviderCatalogPicker, {
			props: { catalog: [openSubtitlesEntry, whisperEntry], onSelect: vi.fn() }
		});

		expect(body).toContain('OpenSubtitles.com');
		expect(body).toContain('Whisper');
		expect(body).toContain('Unsupported');
		expect(body).toContain('local http endpoint');
	});

	it('does not render saved subtitle provider secret values in dynamic fields', () => {
		const form = emptySubtitleProviderForm('opensubtitlescom', openSubtitlesEntry);
		form.apiKeySet = true;
		form.passwordSet = true;
		form.secretFieldsSet = ['apiKey', 'password'];
		const { body } = render(SubtitleProviderForm, {
			props: {
				form,
				entry: openSubtitlesEntry,
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onTest: vi.fn()
			}
		});

		expect(body).toContain('Saved secret');
		expect(body).toContain('type="password"');
		expect(body).not.toContain('scenario-key');
		expect(body).not.toContain('scenario-password');
	});

	it('preserves mock subtitle provider rows in the form', () => {
		const form = emptySubtitleProviderForm('mock', {
			...openSubtitlesEntry,
			key: 'mock',
			displayName: 'Mock',
			fields: []
		});
		form.mockSubtitles = [{ title: 'Scenario Movie', languageId: 'english', format: 'vtt' }];
		const { body } = render(SubtitleProviderForm, {
			props: { form, onSave: vi.fn(), onCancel: vi.fn(), onTest: vi.fn() }
		});

		expect(body).toContain('value="Scenario Movie"');
		expect(body).toContain('value="english"');
		expect(body).toContain('VTT');
	});
});
