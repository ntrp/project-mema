<script lang="ts">
	import { Button } from '$lib/components/ui/button';

	interface Props {
		message: string;
		errorMessage: string;
		onDismiss?: () => void;
	}

	let { message, errorMessage, onDismiss }: Props = $props();
	let dismissedNoticeKey = $state('');
	let noticeKey = $derived(`${errorMessage}|${message}`);
	let visible = $derived(noticeKey !== '|' && dismissedNoticeKey !== noticeKey);

	$effect(() => {
		if (noticeKey === '|') {
			dismissedNoticeKey = '';
			return;
		}
		if (!visible) {
			return;
		}
		const key = noticeKey;
		const timeout = globalThis.setTimeout(() => {
			dismissedNoticeKey = key;
			onDismiss?.();
		}, 5000);
		return () => globalThis.clearTimeout(timeout);
	});

	function dismiss() {
		dismissedNoticeKey = noticeKey;
		onDismiss?.();
	}
</script>

{#if visible && (message || errorMessage)}
	<div
		class="fixed right-4 bottom-4 z-50 grid w-[min(380px,calc(100vw-2rem))] gap-2"
		aria-live="polite"
		aria-atomic="true"
	>
		{#if errorMessage}
			<Button
				type="button"
				variant="destructive"
				class="relative h-auto justify-start overflow-hidden bg-destructive px-4 py-3 text-left text-destructive-foreground shadow-lg whitespace-normal hover:bg-destructive hover:text-destructive-foreground"
				onclick={dismiss}
			>
				<span>{errorMessage}</span>
				<span
					class="notice-progress absolute right-0 bottom-0 left-0 h-0.5 bg-destructive-foreground/70"
					aria-hidden="true"
				></span>
			</Button>
		{/if}
		{#if message}
			<Button
				type="button"
				variant="outline"
				class="relative h-auto justify-start overflow-hidden border-primary bg-popover px-4 py-3 text-left text-popover-foreground shadow-lg whitespace-normal hover:bg-popover hover:text-popover-foreground dark:bg-popover dark:hover:bg-popover"
				onclick={dismiss}
			>
				<span>{message}</span>
				<span
					class="notice-progress absolute right-0 bottom-0 left-0 h-0.5 bg-primary"
					aria-hidden="true"
				></span>
			</Button>
		{/if}
	</div>
{/if}

<style>
	.notice-progress {
		transform-origin: left;
		animation: notice-progress 5s linear forwards;
	}

	@keyframes notice-progress {
		from {
			transform: scaleX(1);
		}
		to {
			transform: scaleX(0);
		}
	}
</style>
