<script lang="ts">
	import type { Snippet } from 'svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();
	const app = getAppShellContext();
</script>

{#if app.isAdmin}
	<section class="grid min-w-0 gap-[18px]" aria-labelledby="settings-title">
		{@render children()}
	</section>
{/if}
