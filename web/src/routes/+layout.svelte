<script lang="ts">
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import favicon from '$lib/assets/favicon.svg';
	import { createAppQueryClient } from '$lib/app/query/queryClient';
	import AppShell from '$lib/components/app/shell/AppShell.svelte';
	import '../styles.css';

	let { children } = $props();
	const queryClient = createAppQueryClient();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	<AppShell>
		{@render children()}
	</AppShell>
</QueryClientProvider>
