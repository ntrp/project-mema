<script lang="ts">
	interface Props {
		message: string;
		errorMessage: string;
	}

	let { message, errorMessage }: Props = $props();
</script>

{#if message || errorMessage}
	<div class="toast-stack" aria-live="polite" aria-atomic="true">
		{#if errorMessage}
			<p class="toast-notice error" role="alert">{errorMessage}</p>
		{/if}
		{#if message}
			<p class="toast-notice success">{message}</p>
		{/if}
	</div>
{/if}
