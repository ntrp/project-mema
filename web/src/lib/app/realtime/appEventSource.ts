export interface AppEvent<T = unknown> {
	id?: string;
	type?: string;
	time?: string;
	data: T;
}

type AppEventHandler<T = unknown> = (event: AppEvent<T>) => void;
export type AppEventSourceStatus = 'idle' | 'connecting' | 'open' | 'error';
type StatusHandler = (status: AppEventSourceStatus) => void;

const listeners = new Map<string, Set<AppEventHandler>>();
const processedIds = new Set<string>();
const processedIdQueue: string[] = [];
const registeredEventTypes = new Set<string>();
const statusListeners = new Set<StatusHandler>();
const maxProcessedIds = 1_000;

let source: EventSource | undefined;
let status: AppEventSourceStatus = 'idle';

export function startAppEventSource() {
	if (source) return;
	source = new EventSource('/api/events', { withCredentials: true });
	setStatus('connecting');
	source.addEventListener('open', () => setStatus('open'));
	source.addEventListener('error', () => setStatus('error'));
	for (const eventType of listeners.keys()) addSourceListener(eventType);
}

export function stopAppEventSource() {
	source?.close();
	source = undefined;
	registeredEventTypes.clear();
	processedIds.clear();
	processedIdQueue.length = 0;
	setStatus('idle');
}

export function subscribeToAppEvent<T>(eventType: string, handler: AppEventHandler<T>) {
	const handlers = listeners.get(eventType) ?? new Set<AppEventHandler>();
	const needsSourceListener = handlers.size === 0;
	handlers.add(handler as AppEventHandler);
	listeners.set(eventType, handlers);
	if (source && needsSourceListener) addSourceListener(eventType);
	return () => {
		handlers.delete(handler as AppEventHandler);
		if (handlers.size === 0) listeners.delete(eventType);
	};
}

export function hasActiveAppEventSource() {
	return Boolean(source);
}

export function subscribeToAppEventSourceStatus(handler: StatusHandler) {
	statusListeners.add(handler);
	handler(status);
	return () => statusListeners.delete(handler);
}

function addSourceListener(eventType: string) {
	if (registeredEventTypes.has(eventType)) return;
	registeredEventTypes.add(eventType);
	source?.addEventListener(eventType, (message) => dispatch(eventType, message));
}

function dispatch(eventType: string, message: Event) {
	const event = parseAppEvent(message);
	if (!event || isDuplicate(event.id)) return;
	for (const handler of listeners.get(eventType) ?? []) handler(event);
}

function setStatus(nextStatus: AppEventSourceStatus) {
	status = nextStatus;
	for (const handler of statusListeners) handler(status);
}

function parseAppEvent(message: Event): AppEvent | undefined {
	try {
		return JSON.parse((message as MessageEvent<string>).data) as AppEvent;
	} catch {
		return undefined;
	}
}

function isDuplicate(id?: string) {
	if (!id) return false;
	if (processedIds.has(id)) return true;
	processedIds.add(id);
	processedIdQueue.push(id);
	if (processedIdQueue.length > maxProcessedIds) {
		processedIds.delete(processedIdQueue.shift() as string);
	}
	return false;
}
