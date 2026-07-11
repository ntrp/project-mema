import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { SystemJobExecution } from '$lib/settings/types';
import type { AppShellState } from './state.svelte';

const enqueue = vi.hoisted(() => vi.fn());
vi.mock('$lib/features/releases/api', () => ({ enqueueMediaFulfillmentAction: enqueue }));
import { createMediaFulfillmentActions } from './mediaFulfillmentActions';

describe('media fulfillment actions', () => {
	const item = { id: 'media-1' } as never;
	const request = { operation: 'video_transcode', filePath: '/movie.mkv' } as never;
	let state: AppShellState;
	let loadMediaItems: ReturnType<typeof vi.fn<() => Promise<void>>>;

	beforeEach(() => {
		enqueue.mockReset();
		state = { pendingFulfillmentActions: {}, message: '', errorMessage: '' } as AppShellState;
		loadMediaItems = vi.fn<() => Promise<void>>().mockResolvedValue(undefined);
	});

	it('tracks an enqueued job and refreshes media', async () => {
		enqueue.mockResolvedValue({ jobId: 42, message: 'Queued' });
		const actions = createMediaFulfillmentActions(state, { clearNotice: vi.fn(), loadMediaItems });
		await expect(actions.enqueueMediaFulfillment(item, request)).resolves.toEqual({
			jobId: 42,
			message: 'Queued'
		});
		expect(state.pendingFulfillmentActions).toEqual({ 'video_transcode|/movie.mkv||': 42 });
		expect(state.message).toBe('Queued (#42)');
		expect(loadMediaItems).toHaveBeenCalledOnce();
	});

	it('clears failed jobs and records the request error', async () => {
		enqueue.mockRejectedValue(new Error('offline'));
		const actions = createMediaFulfillmentActions(state, { clearNotice: vi.fn(), loadMediaItems });
		await actions.enqueueMediaFulfillment(item, request);
		expect(state.pendingFulfillmentActions).toEqual({});
		expect(state.errorMessage).toBe('offline');
	});

	it('reconciles running and final job executions', () => {
		const actions = createMediaFulfillmentActions(state, {
			clearNotice: vi.fn(),
			loadMediaItems
		});
		const execution = {
			kind: 'media.fulfillment.video_transcode',
			args: JSON.stringify({ file_path: '/movie.mkv' }),
			riverJobId: 81,
			status: 'running'
		} as SystemJobExecution;
		actions.updateFulfillmentJobExecution(execution);
		expect(state.pendingFulfillmentActions).toEqual({ 'video_transcode|/movie.mkv||': 81 });
		actions.updateFulfillmentJobExecution({ ...execution, status: 'completed' });
		expect(state.pendingFulfillmentActions).toEqual({});
	});

	it('ignores unrelated, unknown, and malformed executions', () => {
		const actions = createMediaFulfillmentActions(state, {
			clearNotice: vi.fn(),
			loadMediaItems
		});
		actions.updateFulfillmentJobExecution({
			kind: 'other',
			args: '{}',
			status: 'running'
		} as never);
		actions.updateFulfillmentJobExecution({
			kind: 'media.fulfillment.unknown',
			args: '{',
			status: 'running'
		} as never);
		expect(state.pendingFulfillmentActions).toEqual({});
	});
});
