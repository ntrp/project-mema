import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
import {
	clearSystemEvents,
	deleteSystemEvent,
	listSystemEvents
} from '$lib/components/settings/system/events/api';
import type { SystemEvent } from '$lib/settings/types';

const limit = 100;
type EventPage = Awaited<ReturnType<typeof listSystemEvents>>;

export const systemEventKeys = {
	list: () => ['system', 'events'] as const,
	page: (before: string) => ['system', 'events', 'page', before] as const
};

export function createSystemEventsResource() {
	const client = useQueryClient();
	const update = (fn: (page: EventPage) => EventPage) =>
		client.setQueryData<EventPage>(systemEventKeys.list(), (page) =>
			fn(page ?? { events: [], hasMore: false })
		);
	return {
		query: createQuery(() => ({
			queryKey: systemEventKeys.list(),
			queryFn: () => listSystemEvents({ limit })
		})),
		remove: createMutation(() => ({
			mutationFn: deleteSystemEvent,
			onSuccess: (_, id) =>
				update((page) => ({ ...page, events: page.events.filter((e) => e.id !== id) }))
		})),
		clear: createMutation(() => ({
			mutationFn: clearSystemEvents,
			onSuccess: () => update(() => ({ events: [], hasMore: false }))
		})),
		loadMore: async (before: string) => {
			const next = await client.fetchQuery({
				queryKey: systemEventKeys.page(before),
				queryFn: () => listSystemEvents({ before, limit })
			});
			update((page) => {
				return {
					events: [
						...page.events,
						...next.events.filter(
							(event) => !page.events.some((existing) => existing.id === event.id)
						)
					],
					hasMore: next.hasMore
				};
			});
		},
		created: (event: SystemEvent) =>
			update((page) => ({
				...page,
				events: [event, ...page.events.filter((item) => item.id !== event.id)]
			})),
		deleted: (id: string) =>
			update((page) => ({ ...page, events: page.events.filter((item) => item.id !== id) })),
		cleared: () => update(() => ({ events: [], hasMore: false }))
	};
}
