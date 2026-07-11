import type { SystemEvent } from '$lib/settings/types';
import {
	subscribeToAppEvent,
	subscribeToAppEventSourceStatus
} from '$lib/app/realtime/appEventSource';

interface SystemEventSubscription {
	onOpen: () => void;
	onError: () => void;
	onCreated: (event: SystemEvent) => void;
	onDeleted: (id: string) => void;
	onCleared: () => void;
}

export function subscribeSystemEvents(handlers: SystemEventSubscription) {
	const unsubscribers = [
		subscribeToAppEventSourceStatus((status) => {
			if (status === 'open') handlers.onOpen();
			if (status === 'error') handlers.onError();
		}),
		subscribeToAppEvent<SystemEvent>('system.event.created', ({ data: nextEvent }) => {
			if (nextEvent) {
				handlers.onCreated(nextEvent);
			}
		}),
		subscribeToAppEvent<{ id: string }>('system.event.deleted', ({ data: deleted }) => {
			if (deleted?.id) {
				handlers.onDeleted(deleted.id);
			}
		}),
		subscribeToAppEvent('system.events.cleared', handlers.onCleared)
	];
	return () => unsubscribers.forEach((unsubscribe) => unsubscribe());
}
