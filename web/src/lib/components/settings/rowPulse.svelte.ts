const PULSE_MS = 1200;

export function createRowPulse() {
	let knownKeys = new Set<string>();
	let pulsingKeys = $state<Set<string>>(new Set());

	function update(keys: string[]) {
		const nextKeys = new Set(keys);
		if (knownKeys.size > 0) {
			const addedKeys = keys.filter((key) => !knownKeys.has(key));
			if (addedKeys.length > 0) {
				pulsingKeys = new Set([...pulsingKeys, ...addedKeys]);
				for (const key of addedKeys) {
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
		const nextKeys = new Set(pulsingKeys);
		nextKeys.delete(key);
		pulsingKeys = nextKeys;
	}

	return { update, classFor };
}
