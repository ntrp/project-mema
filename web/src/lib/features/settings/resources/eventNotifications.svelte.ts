import { createQuery, useQueryClient } from '@tanstack/svelte-query';
import { listSystemEvents } from '$lib/components/settings/system/events/api';
import type { SystemEvent } from '$lib/settings/types';

const key = ['system', 'events', 'notifications'] as const;
type EventPage = Awaited<ReturnType<typeof listSystemEvents>>;

export function createEventNotificationsResource(enabled: () => boolean) {
	const client = useQueryClient();
	const update = (fn: (events: SystemEvent[]) => SystemEvent[]) =>
		client.setQueryData<EventPage>(key, (page) => ({
			events: fn(page?.events ?? []),
			hasMore: page?.hasMore ?? false
		}));
	return {
		query: createQuery(() => ({
			queryKey: key,
			queryFn: () => listSystemEvents(),
			enabled: enabled()
		})),
		created: (event: SystemEvent) =>
			update((events) => [event, ...events.filter((item) => item.id !== event.id)]),
		deleted: (id: string) => update((events) => events.filter((item) => item.id !== id)),
		cleared: () => update(() => [])
	};
}
