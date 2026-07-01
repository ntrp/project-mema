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

	function dismiss() {
		dismissedNoticeKey = noticeKey;
		onDismiss?.();
	}
</script>

{#if visible && (message || errorMessage)}
	<div
		class="fixed top-4 right-4 z-50 grid w-[min(360px,calc(100vw-2rem))] gap-2"
		aria-live="polite"
		aria-atomic="true"
	>
		{#if errorMessage}
			<Button
				type="button"
				variant="destructive"
				class="relative h-auto justify-start overflow-hidden whitespace-normal px-4 py-3 text-left shadow-lg"
				onclick={dismiss}
			>
				<span>{errorMessage}</span>
			</Button>
		{/if}
		{#if message}
			<Button
				type="button"
				variant="outline"
				class="relative h-auto justify-start overflow-hidden border-primary/30 bg-primary/10 px-4 py-3 text-left text-primary shadow-lg whitespace-normal hover:bg-primary/15"
				onclick={dismiss}
			>
				<span>{message}</span>
			</Button>
		{/if}
	</div>
{/if}
