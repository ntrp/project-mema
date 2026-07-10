const touchTooltipRoots = new Map<string, () => void>();

export function registerTouchTooltipRoot(root: HTMLElement, open: () => void) {
	const id = globalThis.crypto?.randomUUID?.() ?? Math.random().toString(36);
	root.dataset.touchTooltipId = id;
	touchTooltipRoots.set(id, open);
	return () => touchTooltipRoots.delete(id);
}

export function openTouchTooltipRoot(root: Element | null) {
	if (!(root instanceof HTMLElement)) return;
	const id = root.dataset.touchTooltipId;
	if (!id) return;
	const open = touchTooltipRoots.get(id);
	open?.();
}
