<script lang="ts">
	import SystemArea from '$lib/components/app/system/SystemArea.svelte';
	import type { SystemSection } from '$lib/settings/types';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		section: SystemSection;
	}

	let { section }: Props = $props();
	const app = getAppShellContext();
</script>

{#if app.isAdmin}
	<SystemArea activeSection={section} />
{/if}
