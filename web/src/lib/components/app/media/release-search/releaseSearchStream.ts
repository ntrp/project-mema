import type { ReleaseCandidate } from '$lib/settings/types';

interface StreamEnvelope<T> {
	data: T;
}

export interface ReleaseSearchStreamResult {
	releases: ReleaseCandidate[];
	errors: string[];
}

export interface ReleaseSearchStreamStatus {
	kind?: 'message' | 'indexer_start' | 'indexer_finish' | 'error';
	message: string;
	indexerName?: string;
	query?: string;
	resultCount?: number;
	cacheHit?: boolean;
	durationMs?: number;
}

interface ReleaseSearchStreamHandlers {
	onStatus: (status: ReleaseSearchStreamStatus) => void;
	onResult: (result: ReleaseSearchStreamResult) => void;
	onError: (message: string) => void;
}

export function subscribeReleaseSearchStream(
	itemId: string,
	query: string,
	handlers: ReleaseSearchStreamHandlers
) {
	const suffix = query ? `?query=${encodeURIComponent(query)}` : '';
	const source = new EventSource(`/api/media/items/${itemId}/release-searches/stream${suffix}`, {
		withCredentials: true
	});

	source.addEventListener('media.release_search.status', (event) => {
		const payload = parseEvent<ReleaseSearchStreamStatus>(event);
		if (payload?.message) {
			handlers.onStatus(payload);
		}
	});
	source.addEventListener('media.release_search.result', (event) => {
		const result = parseEvent<ReleaseSearchStreamResult>(event);
		if (result) {
			handlers.onResult(result);
		}
		source.close();
	});
	source.addEventListener('media.release_search.error', (event) => {
		const payload = parseEvent<{ message?: string }>(event);
		handlers.onError(payload?.message ?? 'Release search failed');
		source.close();
	});
	source.addEventListener('error', () => {
		handlers.onError('Release search stream disconnected');
		source.close();
	});

	return () => source.close();
}

function parseEvent<T>(event: Event) {
	try {
		const message = event as MessageEvent<string>;
		const envelope = JSON.parse(message.data) as StreamEnvelope<T>;
		return envelope.data;
	} catch {
		return undefined;
	}
}
