export const activeStatuses = ['available', 'scheduled', 'retryable', 'running'];
export const finalStatuses = ['completed', 'cancelled', 'discarded'];
export const allStatuses = [...activeStatuses, 'pending', ...finalStatuses];

export function canAbortStatus(status: string) {
	return !finalStatuses.includes(status);
}

export function statusClass(status: string) {
	if (status === 'running') return 'border-sky-500/50 bg-sky-500/10 text-sky-300';
	if (status === 'completed') return 'border-emerald-500/50 bg-emerald-500/10 text-emerald-300';
	if (status === 'cancelled' || status === 'discarded') {
		return 'border-destructive/50 bg-destructive/10 text-destructive';
	}
	return 'border-amber-500/50 bg-amber-500/10 text-amber-300';
}

export function formatInterval(seconds: number) {
	if (seconds < 60) return `${seconds}s`;
	if (seconds < 3600) return `${Math.round(seconds / 60)}m`;
	if (seconds < 86400) return `${Math.round(seconds / 3600)}h`;
	return `${Math.round(seconds / 86400)}d`;
}

export function scheduleCategoryLabel(category: string) {
	switch (category) {
		case 'release_search':
			return 'Release search';
		case 'download_import':
			return 'Download import';
		case 'subtitle_fulfillment':
			return 'Subtitles';
		default:
			return 'Maintenance';
	}
}

export function progressStyle(progress?: number) {
	if (progress === undefined) return '';
	return `width: ${Math.max(0, Math.min(100, progress))}%`;
}

export function executionMessage(progressLabel: string, infoMessage: string, status: string) {
	return progressLabel || infoMessage || status;
}
