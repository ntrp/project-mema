import { describe, expect, it, vi } from 'vitest';

const gotoMock = vi.hoisted(() => vi.fn());

vi.mock('$app/navigation', () => ({ goto: gotoMock }));
vi.mock('$app/paths', () => ({ resolve: (path: string) => path }));

import { createSettingsEditActions } from './settingsEditActions';
import type { AppShellState } from './state.svelte';

describe('settings edit actions (SCN-SETTINGS-002)', () => {
	it('opens integration and user edit forms on their settings routes', () => {
		const state = {} as AppShellState;
		const actions = createSettingsEditActions(state);

		actions.editDownloadClient({
			id: 'client-1',
			name: 'Transmission',
			type: 'transmission',
			baseUrl: 'http://client.local',
			username: 'user',
			password: 'secret',
			apiKey: 'key',
			category: 'movies',
			enabled: true,
			priority: 1
		} as never);
		expect(state.activeSettingsSection).toBe('download-clients');
		expect(state.downloadForm).toMatchObject({ id: 'client-1', name: 'Transmission' });
		expect(gotoMock).toHaveBeenLastCalledWith('/settings/download-clients');

		actions.editIndexer({
			id: 'indexer-1',
			name: 'Torznab',
			type: 'torznab',
			baseUrl: 'http://indexer.local',
			apiKey: 'key',
			categories: [2000, 2040],
			enabled: false,
			priority: 2
		} as never);
		expect(state.activeSettingsSection).toBe('indexers');
		expect(state.indexerForm).toMatchObject({ id: 'indexer-1', categoriesText: '2000, 2040' });
		expect(gotoMock).toHaveBeenLastCalledWith('/settings/indexers');

		actions.editUser({ id: 'user-1', username: 'viewer', role: 'user' } as never);
		expect(state.activeSettingsSection).toBe('users');
		expect(state.userForm).toMatchObject({ id: 'user-1', username: 'viewer', password: '' });
		expect(gotoMock).toHaveBeenLastCalledWith('/settings/users');

		actions.editTag({
			id: 'tag-1',
			name: 'Action',
			createdAt: '2026-07-03T00:00:00Z',
			updatedAt: '2026-07-03T00:00:00Z'
		});
		expect(state.tagForm).toEqual({ id: 'tag-1', name: 'Action' });
		expect(gotoMock).toHaveBeenLastCalledWith('/settings/tags');
	});

	it('opens language, profile, and custom-format edit forms', () => {
		const state = {} as AppShellState;
		const actions = createSettingsEditActions(state);

		actions.editLanguage({
			code: 'de',
			displayName: 'German',
			aliases: ['deu', 'ger']
		} as never);
		expect(state.activeSettingsSection).toBe('languages');
		expect(state.languageForm).toMatchObject({
			originalCode: 'de',
			code: 'de',
			aliasesText: 'deu, ger'
		});

		actions.editMediaProfile({
			id: 'profile-1',
			name: 'HD',
			qualityIds: ['q-1080p'],
			upgradesAllowed: true,
			minimumCustomFormatScore: 0,
			upgradeUntilCustomFormatScore: 0,
			minimumCustomFormatScoreIncrement: 1,
			removeNonEnabledLanguages: false,
			preferredProtocol: 'any',
			seriesPackPreference: 'auto',
			targetLanguages: ['english'],
			customFormatScores: []
		} as never);
		expect(state.activeSettingsSection).toBe('profiles');
		expect(state.mediaProfileForm).toMatchObject({ id: 'profile-1', name: 'HD' });

		actions.editCustomFormat({
			id: 'format-1',
			name: 'HDR',
			includeInRenameTemplate: true,
			includeSpecs: [],
			excludeSpecs: [],
			createdAt: '2026-07-03T00:00:00Z',
			updatedAt: '2026-07-03T00:00:00Z'
		});
		expect(state.activeSettingsSection).toBe('custom-formats');
		expect(state.customFormatForm).toMatchObject({ id: 'format-1', name: 'HDR' });
		expect(gotoMock).toHaveBeenLastCalledWith('/settings/custom-formats');
	});
});
