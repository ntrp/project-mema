import type { SystemEvent } from '$lib/settings/types';
import { parseSystemEvent } from './systemEventStream';

interface SystemEventSubscription {
	onOpen: () => void;
	onError: () => void;
	onCreated: (event: SystemEvent) => void;
	onDeleted: (id: string) => void;
	onCleared: () => void;
}

export function subscribeSystemEvents(handlers: SystemEventSubscription) {
	const source = new EventSource('/api/events', { withCredentials: true });
	source.addEventListener('open', handlers.onOpen);
	source.addEventListener('error', handlers.onError);
	source.addEventListener('system.event.created', (event) => {
		const nextEvent = parseSystemEvent<SystemEvent>(event);
		if (nextEvent) {
			handlers.onCreated(nextEvent);
		}
	});
	source.addEventListener('system.event.deleted', (event) => {
		const deleted = parseSystemEvent<{ id: string }>(event);
		if (deleted?.id) {
			handlers.onDeleted(deleted.id);
		}
	});
	source.addEventListener('system.events.cleared', handlers.onCleared);
	return () => source.close();
}
