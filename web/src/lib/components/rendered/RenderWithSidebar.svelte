<script lang="ts">
	import type { Component as SvelteComponent } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar';

	interface Props {
		component: SvelteComponent<Record<string, unknown>>;
		componentProps?: Record<string, unknown>;
	}

	let { component: Component, componentProps = {} }: Props = $props();
</script>

<Sidebar.Provider>
	<Component {...componentProps} />
</Sidebar.Provider>
