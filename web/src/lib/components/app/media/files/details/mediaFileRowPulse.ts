import type { MediaFileDetailRow } from '../mediaFileDetails';
import { SvelteMap, SvelteSet } from 'svelte/reactivity';

export function createMediaFileRowPulse() {
	const pulsing = new SvelteSet<string>();
	let signatures = new SvelteMap<string, string>();
	const timers = new SvelteMap<string, number>();
	return (rows: MediaFileDetailRow[]) => {
		const next = new SvelteMap(rows.map((row) => [row.key, signature(row)]));
		for (const [key, value] of next) {
			if (signatures.has(key) && signatures.get(key) !== value) pulse(key);
		}
		signatures = next;
		return pulsing;
	};

	function pulse(key: string) {
		const active = timers.get(key);
		if (active) window.clearTimeout(active);
		pulsing.add(key);
		timers.set(
			key,
			window.setTimeout(() => pulsing.delete(key), 1200)
		);
	}
}

function signature(row: MediaFileDetailRow) {
	return [
		row.description,
		row.visualState,
		row.statusLabel,
		row.operationLabel,
		...(row.details ?? [])
	].join('\u001f');
}
