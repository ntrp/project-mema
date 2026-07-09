<script lang="ts">
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import type { Component as SvelteComponent } from 'svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		component: SvelteComponent<Record<string, unknown>>;
		componentProps?: Record<string, unknown>;
	}

	let { component: Component, componentProps = {} }: Props = $props();
	const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
	<Tooltip.Provider>
		<Component {...componentProps} />
	</Tooltip.Provider>
</QueryClientProvider>
