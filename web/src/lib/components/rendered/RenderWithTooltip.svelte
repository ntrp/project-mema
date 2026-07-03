<script lang="ts">
	import type { Component as SvelteComponent } from 'svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		component: SvelteComponent<Record<string, unknown>>;
		componentProps?: Record<string, unknown>;
	}

	let { component: Component, componentProps = {} }: Props = $props();
</script>

<Tooltip.Provider>
	<Component {...componentProps} />
</Tooltip.Provider>
