import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	setQueryData: vi.fn(),
	subscribe: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({
	useQueryClient: () => ({ setQueryData: mocks.setQueryData })
}));
vi.mock('$lib/app/realtime/appEventSource', () => ({
	subscribeToAppEvent: mocks.subscribe
}));

import { connectActivityQueryEvents } from './realtime.svelte';
import { activityKeys } from './queries.svelte';

describe('activity real-time cache updates', () => {
	beforeEach(() => vi.clearAllMocks());

	it('subscribes to download updates and returns the unsubscribe callback', () => {
		const unsubscribe = vi.fn();
		mocks.subscribe.mockReturnValue(unsubscribe);

		expect(connectActivityQueryEvents()).toBe(unsubscribe);
		expect(mocks.subscribe).toHaveBeenCalledWith('activity.download.updated', expect.any(Function));
	});

	it('moves an updated download to the front without duplicating it', () => {
		connectActivityQueryEvents();
		const handler = mocks.subscribe.mock.calls[0][1];
		const updated = { id: 'one', status: 'completed' };

		handler({ data: updated });

		expect(mocks.setQueryData).toHaveBeenCalledWith(activityKeys.downloads(), expect.any(Function));
		const update = mocks.setQueryData.mock.calls[0][1];
		expect(
			update({
				activities: [
					{ id: 'two', status: 'queued' },
					{ id: 'one', status: 'queued' }
				]
			})
		).toEqual({ activities: [updated, { id: 'two', status: 'queued' }] });
	});

	it('ignores empty data and initializes an empty cache', () => {
		connectActivityQueryEvents();
		const handler = mocks.subscribe.mock.calls[0][1];
		handler({ data: undefined });
		expect(mocks.setQueryData).not.toHaveBeenCalled();

		handler({ data: { id: 'new' } });
		const update = mocks.setQueryData.mock.calls[0][1];
		expect(update(undefined)).toEqual({ activities: [{ id: 'new' }] });
	});
});
