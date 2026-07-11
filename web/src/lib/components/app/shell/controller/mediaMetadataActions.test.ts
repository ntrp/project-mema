import { beforeEach, describe, expect, it, vi } from 'vitest';

const refreshMediaItemMetadataMock = vi.hoisted(() => vi.fn());

vi.mock('$lib/settings/api', () => ({
	refreshMediaItemMetadata: refreshMediaItemMetadataMock
}));

import type { MediaItem } from '$lib/settings/types';
import { createMediaMetadataActions } from './mediaMetadataActions';
import type { AppShellState } from './state.svelte';

describe('media metadata actions (SCN-MEDIA-004)', () => {
	beforeEach(() => {
		refreshMediaItemMetadataMock.mockReset();
	});

	it('refreshes metadata, moves the updated item to the top, and resets busy state', async () => {
		const clearNotice = vi.fn();
		const state = testState([
			mediaItem({ id: 'media-1', title: 'Original Movie' }),
			mediaItem({ id: 'media-2', title: 'Other Movie' })
		]);
		refreshMediaItemMetadataMock.mockResolvedValue(
			mediaItem({ id: 'media-1', title: 'Updated Movie' })
		);

		await createMediaMetadataActions(state, {
			clearNotice,
			upsertMediaItem: state.upsertMediaItem
		}).refreshMediaMetadata(state.mediaItems[0]);

		expect(clearNotice).toHaveBeenCalledTimes(1);
		expect(refreshMediaItemMetadataMock).toHaveBeenCalledWith('media-1');
		expect(state.mediaItems.map((item) => item.title)).toEqual(['Updated Movie', 'Other Movie']);
		expect(state.message).toBe('Media metadata refreshed');
		expect(state.refreshingMetadataItemId).toBeUndefined();
	});

	it('keeps existing items and surfaces refresh failures', async () => {
		const state = testState([mediaItem({ id: 'media-1', title: 'Original Movie' })]);
		refreshMediaItemMetadataMock.mockRejectedValue(new Error('Provider unavailable'));

		await createMediaMetadataActions(state, {
			clearNotice: vi.fn(),
			upsertMediaItem: state.upsertMediaItem
		}).refreshMediaMetadata(state.mediaItems[0]);

		expect(state.mediaItems[0].title).toBe('Original Movie');
		expect(state.errorMessage).toBe('Provider unavailable');
		expect(state.refreshingMetadataItemId).toBeUndefined();
	});
});

function testState(mediaItems: MediaItem[]) {
	const value = {
		mediaItems,
		message: '',
		errorMessage: '',
		refreshingMetadataItemId: undefined,
		upsertMediaItem(item: MediaItem) {
			value.mediaItems = [item, ...value.mediaItems.filter((entry) => entry.id !== item.id)];
		}
	};
	return value as unknown as AppShellState & typeof value;
}

function mediaItem(overrides: Partial<MediaItem>): MediaItem {
	return {
		id: 'media-1',
		type: 'movie',
		title: 'Scenario Movie',
		year: 2026,
		status: 'missing',
		monitored: true,
		monitorMode: 'onlyMedia',
		minimumAvailability: 'released',
		...overrides
	} as MediaItem;
}
