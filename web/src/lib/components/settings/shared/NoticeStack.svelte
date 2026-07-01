<script lang="ts">
	const dismissDelayMs = 5000;

	interface Props {
		message: string;
		errorMessage: string;
		onDismiss?: () => void;
	}

	let { message, errorMessage, onDismiss }: Props = $props();
	let visible = $state(false);
	let noticeKey = $state('');

	$effect(() => {
		const nextKey = `${errorMessage}|${message}`;
		if (!nextKey || nextKey === '|') {
			visible = false;
			noticeKey = '';
			return;
		}
		visible = true;
		noticeKey = nextKey;
		const timeout = window.setTimeout(dismiss, dismissDelayMs);
		return () => window.clearTimeout(timeout);
	});

	function dismiss() {
		visible = false;
		onDismiss?.();
	}
</script>

{#if visible && (message || errorMessage)}
	<div class="toast-stack" aria-live="polite" aria-atomic="true">
		{#if errorMessage}
			<button type="button" class="toast-notice error" onclick={dismiss}>
				<span>{errorMessage}</span>
				{#key noticeKey}
					<span class="toast-progress" aria-hidden="true"></span>
				{/key}
			</button>
		{/if}
		{#if message}
			<button type="button" class="toast-notice success" onclick={dismiss}>
				<span>{message}</span>
				{#key noticeKey}
					<span class="toast-progress" aria-hidden="true"></span>
				{/key}
			</button>
		{/if}
	</div>
{/if}
