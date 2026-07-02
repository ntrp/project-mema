import { SvelteSet } from 'svelte/reactivity';

const PULSE_MS = 1200;

export function createRowPulse() {
	let knownKeys = new SvelteSet<string>();
	let pulsingKeys = new SvelteSet<string>();

	function update(keys: string[]) {
		const nextKeys = new SvelteSet(keys);
		if (knownKeys.size > 0) {
			const addedKeys = keys.filter((key) => !knownKeys.has(key));
			if (addedKeys.length > 0) {
				for (const key of addedKeys) {
					pulsingKeys.add(key);
					globalThis.setTimeout(() => removePulse(key), PULSE_MS);
				}
			}
		}
		knownKeys = nextKeys;
	}

	function classFor(key: string) {
		return pulsingKeys.has(key) ? 'live-row-pulse' : undefined;
	}

	function removePulse(key: string) {
		if (!pulsingKeys.has(key)) return;
		pulsingKeys.delete(key);
	}

	return { update, classFor };
}
