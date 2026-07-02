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
	<SystemArea
		activeSection={section}
		indexerSearch={app.indexerSearch}
		metadataCache={app.metadataCache}
		loadingIndexerSearch={app.loadingIndexerSearch}
		loadingMetadataCache={app.loadingMetadataCache}
		clearingIndexerSearchCache={app.clearingIndexerSearchCache}
		clearingMetadataCache={app.clearingMetadataCache}
		onRefreshIndexerSearch={app.refreshIndexerSearch}
		onRefreshMetadataCache={app.refreshMetadataCache}
		onClearIndexerSearchCache={app.clearIndexerSearchCache}
		onClearIndexerSearchCachePattern={app.clearIndexerSearchCachePattern}
		onClearMetadataCache={app.clearMetadataCache}
		onClearMetadataCachePattern={app.clearMetadataCachePattern}
	/>
{/if}
