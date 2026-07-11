import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { MediaItem } from '$lib/settings/types';
import type { AppShellState } from './state.svelte';

const requests = vi.hoisted(() => ({
	search: vi.fn(),
	remove: vi.fn(),
	update: vi.fn(),
	grab: vi.fn()
}));
vi.mock('$lib/features/library/filesApi', () => ({
	deleteMediaItemSubtitle: requests.remove,
	updateMediaItemSubtitle: requests.update
}));
vi.mock('$lib/features/releases/api', () => ({
	enqueueMediaSubtitleSearch: requests.search,
	grabMediaSubtitle: requests.grab
}));
import { createMediaSubtitleActions } from './mediaSubtitleActions';

describe('media subtitle actions', () => {
	const item = { id: 'media-1' } as never;
	let state: AppShellState;
	let upsertMediaItem: ReturnType<typeof vi.fn<(item: MediaItem) => void>>;

	beforeEach(() => {
		Object.values(requests).forEach((request) => request.mockReset());
		state = { message: '', errorMessage: '' } as AppShellState;
		upsertMediaItem = vi.fn<(item: MediaItem) => void>();
	});

	it('searches, deletes, updates, and grabs subtitles', async () => {
		requests.search.mockResolvedValue({ jobId: 3, message: 'Queued' });
		requests.remove.mockResolvedValue({ id: 'removed' });
		requests.update.mockResolvedValue({ id: 'updated' });
		requests.grab.mockResolvedValue({ id: 'grabbed' });
		const actions = createMediaSubtitleActions(state, { clearNotice: vi.fn(), upsertMediaItem });
		await actions.searchMediaSubtitle(item);
		expect(state.message).toBe('Queued (#3)');
		await actions.deleteMediaSubtitle(item, 'subtitle-1');
		await actions.updateMediaSubtitle(item, 'subtitle-1', {} as never);
		await actions.grabMediaSubtitle(item, {} as never);
		expect(upsertMediaItem.mock.calls.map(([value]) => value.id)).toEqual([
			'removed',
			'updated',
			'grabbed'
		]);
		expect(state.message).toBe('Subtitle grabbed');
	});

	it('reports operation errors and rethrows grab failures', async () => {
		const actions = createMediaSubtitleActions(state, { clearNotice: vi.fn(), upsertMediaItem });
		requests.search.mockRejectedValue(new Error('search failed'));
		await actions.searchMediaSubtitle(item);
		expect(state.errorMessage).toBe('search failed');
		requests.remove.mockRejectedValue(new Error('delete failed'));
		await actions.deleteMediaSubtitle(item, 'subtitle-1');
		expect(state.errorMessage).toBe('delete failed');
		requests.update.mockRejectedValue(new Error('update failed'));
		await actions.updateMediaSubtitle(item, 'subtitle-1', {} as never);
		expect(state.errorMessage).toBe('update failed');
		requests.grab.mockRejectedValue(new Error('grab failed'));
		await expect(actions.grabMediaSubtitle(item, {} as never)).rejects.toThrow('grab failed');
		expect(state.errorMessage).toBe('grab failed');
	});
});
