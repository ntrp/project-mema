<script lang="ts">
	import { Tooltip as TooltipPrimitive } from 'bits-ui';
	import { onMount } from 'svelte';
	import { openTouchTooltipRoot } from './tooltip-touch';

	let {
		delayDuration = 0,
		disableHoverableContent = true,
		...restProps
	}: TooltipPrimitive.ProviderProps = $props();

	onMount(() => {
		let pendingTouchRoot: HTMLElement | null = null;

		function handlePointerDown(event: PointerEvent) {
			pendingTouchRoot = null;
			if (event.pointerType !== 'touch' && event.pointerType !== 'pen') return;
			const trigger = (event.target as Element | null)?.closest?.('[data-slot="tooltip-trigger"]');
			pendingTouchRoot = trigger?.closest('[data-slot="tooltip-root"]') as HTMLElement | null;
		}

		function handlePointerCancel(event: PointerEvent) {
			if (event.pointerType === 'touch' || event.pointerType === 'pen') pendingTouchRoot = null;
		}

		function handleClick() {
			if (!pendingTouchRoot) return;
			const root = pendingTouchRoot;
			pendingTouchRoot = null;
			globalThis.setTimeout(() => openTouchTooltipRoot(root), 0);
		}

		document.addEventListener('pointerdown', handlePointerDown, true);
		document.addEventListener('pointercancel', handlePointerCancel, true);
		document.addEventListener('click', handleClick, true);
		return () => {
			document.removeEventListener('pointerdown', handlePointerDown, true);
			document.removeEventListener('pointercancel', handlePointerCancel, true);
			document.removeEventListener('click', handleClick, true);
		};
	});
</script>

<TooltipPrimitive.Provider {delayDuration} {disableHoverableContent} {...restProps} />
