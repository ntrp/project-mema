import { describe, expect, it } from 'vitest';

import { AppShellState } from './state.svelte';
import { routeStateFromPath } from './routeState';
import type { MediaItem, MediaMetadataDetails } from '$lib/settings/types';

const noQuery = { get: () => null };

describe('app shell state (SCN-MEDIA-004)', () => {
	it('initializes active sections from route state', () => {
		const state = new AppShellState(routeStateFromPath('/settings/quality', {}, noQuery));

		expect(state.activeView).toBe('settings');
		expect(state.activeSettingsSection).toBe('quality');
		expect(state.activePrimarySection).toBe('settings');
		expect(state.activeSubmenuSection).toBe('quality');
	});

	it('hides admin-only primary items for regular users', () => {
		const state = new AppShellState();

		state.currentUser = { id: 'user-1', username: 'viewer', role: 'user' };
		expect(state.isAdmin).toBe(false);
		expect(state.primaryItems.map((item) => item.value)).not.toContain('settings');
		expect(state.primaryItems.map((item) => item.value)).not.toContain('blacklist');

		state.currentUser = { id: 'admin-1', username: 'admin', role: 'admin' };
		expect(state.isAdmin).toBe(true);
		expect(state.primaryItems.map((item) => item.value)).toContain('settings');
		expect(state.primaryItems.map((item) => item.value)).toContain('system');
	});

	it('derives people and related detail from selected media or metadata detail', () => {
		const state = new AppShellState(
			routeStateFromPath('/movies/movie-1/cast', { id: 'movie-1' }, noQuery)
		);
		state.mediaItems = [
			{
				id: 'movie-1',
				type: 'movie',
				title: 'Scenario Movie',
				monitored: true,
				recommendations: [{ title: 'Related Movie' }]
			} as MediaItem
		];

		expect(state.mediaPeopleMetadataDetail?.title).toBe('Scenario Movie');
		expect(state.relatedMediaSection).toBeUndefined();

		state.metadataDetail = {
			type: 'movie',
			title: 'Metadata Movie',
			externalProvider: 'tmdb',
			recommendations: [{ type: 'movie', title: 'Related Movie' }]
		} as MediaMetadataDetails;
		expect(state.mediaPeopleMetadataDetail?.title).toBe('Metadata Movie');
		expect(state.relatedMediaSection?.results[0].title).toBe('Related Movie');
	});
});
