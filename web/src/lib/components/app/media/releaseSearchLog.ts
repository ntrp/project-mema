import type { ReleaseSearchStreamStatus } from './releaseSearchStream';

export interface ReleaseSearchLogEntry {
	id: string;
	timestamp: string;
	message: string;
	resultMessage?: string;
	cacheHit?: boolean;
	durationMs?: number;
}

let nextLogId = 0;

export function placeholderLogEntry(): ReleaseSearchLogEntry {
	return {
		id: 'placeholder',
		timestamp: '',
		message: 'Press search to start'
	};
}

export function createLogEntry(message: string): ReleaseSearchLogEntry {
	return {
		id: nextId(),
		timestamp: timestampLabel(),
		message
	};
}

export function applyStatusToLog(
	entries: ReleaseSearchLogEntry[],
	status: ReleaseSearchStreamStatus
) {
	if (status.kind === 'indexer_finish' && status.indexerName && status.query) {
		const key = indexerKey(status.indexerName, status.query);
		const nextEntries = entries.map((entry) => {
			if (entry.id !== key) {
				return entry;
			}
			return {
				...entry,
				resultMessage: status.message,
				cacheHit: status.cacheHit,
				durationMs: status.durationMs
			};
		});
		return nextEntries.some((entry) => entry.id === key)
			? nextEntries
			: [...entries, statusLogEntry(status)];
	}
	return [...entries, statusLogEntry(status)];
}

function statusLogEntry(status: ReleaseSearchStreamStatus): ReleaseSearchLogEntry {
	return {
		id:
			status.kind === 'indexer_start' && status.indexerName && status.query
				? indexerKey(status.indexerName, status.query)
				: nextId(),
		timestamp: timestampLabel(),
		message: status.message,
		cacheHit: status.cacheHit,
		durationMs: status.durationMs
	};
}

function nextId() {
	nextLogId += 1;
	return `log:${nextLogId}`;
}

function indexerKey(indexerName: string, query: string) {
	return `indexer:${indexerName}:${query}`;
}

function timestampLabel() {
	const now = new Date();
	const hours = String(now.getHours()).padStart(2, '0');
	const minutes = String(now.getMinutes()).padStart(2, '0');
	const seconds = String(now.getSeconds()).padStart(2, '0');
	const milliseconds = String(now.getMilliseconds()).padStart(3, '0');
	return `${hours}:${minutes}:${seconds}.${milliseconds}`;
}
