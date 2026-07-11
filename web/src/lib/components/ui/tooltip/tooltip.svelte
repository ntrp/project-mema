<script lang="ts" generics="T = never">
	import { Tooltip as TooltipPrimitive } from 'bits-ui';
	import { registerTouchTooltipRoot } from './tooltip-touch';

	let { open = $bindable(false), ...restProps }: TooltipPrimitive.RootProps<T> = $props();
	let rootOpen = $state(open);
	let touchCloseTimer: ReturnType<typeof globalThis.setTimeout> | undefined;
	let touchOpenActive = false;

	$effect(() => {
		if (!touchOpenActive) rootOpen = open;
	});

	function clearTouchCloseTimer() {
		if (!touchCloseTimer) return;
		globalThis.clearTimeout(touchCloseTimer);
		touchCloseTimer = undefined;
	}

	function openTouchTooltip() {
		clearTouchCloseTimer();
		touchOpenActive = true;
		rootOpen = true;
		open = true;
		touchCloseTimer = globalThis.setTimeout(() => {
			touchOpenActive = false;
			rootOpen = false;
			open = false;
			touchCloseTimer = undefined;
		}, 2500);
	}

	function handleOpenChange(value: boolean) {
		if (!value && touchOpenActive) return;
		rootOpen = value;
		open = value;
	}

	function touchTooltipRoot(node: HTMLSpanElement) {
		const unregister = registerTouchTooltipRoot(node, openTouchTooltip);
		return {
			destroy() {
				unregister();
			}
		};
	}
</script>

<span use:touchTooltipRoot data-slot="tooltip-root" style="display: contents;">
	<TooltipPrimitive.Root open={rootOpen} onOpenChange={handleOpenChange} {...restProps} />
</span>
