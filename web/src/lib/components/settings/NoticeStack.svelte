<script lang="ts">
	interface Props {
		message: string;
		errorMessage: string;
	}

	let { message, errorMessage }: Props = $props();
</script>

{#if message}
	<p class="notice success">{message}</p>
{/if}
{#if errorMessage}
	<p class="notice error">{errorMessage}</p>
{/if}
