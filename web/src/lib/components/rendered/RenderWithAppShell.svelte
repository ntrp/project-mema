<script lang="ts">
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import { createAppQueryClient } from '$lib/app/query/queryClient';
	import type { Component as SvelteComponent } from 'svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { setAppShellContext, type AppShellController } from '$lib/features/app/appShellContext';

	interface Props {
		app: AppShellController;
		component: SvelteComponent<Record<string, unknown>>;
		componentProps?: Record<string, unknown>;
	}

	let { app, component: Component, componentProps = {} }: Props = $props();
	const queryClient = createAppQueryClient();

	// svelte-ignore state_referenced_locally
	setAppShellContext(app);
</script>

<QueryClientProvider client={queryClient}>
	<Tooltip.Provider>
		<Component {...componentProps} />
	</Tooltip.Provider>
</QueryClientProvider>
