<script lang="ts">
	import type { Component as SvelteComponent } from 'svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { setAppShellContext, type AppShellController } from '$lib/features/app/appShellContext';

	interface Props {
		app: AppShellController;
		component: SvelteComponent<Record<string, unknown>>;
		componentProps?: Record<string, unknown>;
	}

	let { app, component: Component, componentProps = {} }: Props = $props();

	// svelte-ignore state_referenced_locally
	setAppShellContext(app);
</script>

<Tooltip.Provider>
	<Component {...componentProps} />
</Tooltip.Provider>
