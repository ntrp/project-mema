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
		onClearIndexerSearchCache={app.clearIndexerSearchCache}
		onClearIndexerSearchCachePattern={app.clearIndexerSearchCachePattern}
		onDeleteIndexerSearchCacheEntry={app.deleteIndexerSearchCacheEntry}
		onClearIndexerSearchHistory={app.clearIndexerSearchHistory}
		onLoadMoreIndexerSearchCache={app.loadMoreIndexerSearchCache}
		onLoadMoreIndexerSearchHistory={app.loadMoreIndexerSearchHistory}
		onClearMetadataCache={app.clearMetadataCache}
		onClearMetadataCachePattern={app.clearMetadataCachePattern}
		onDeleteMetadataCacheEntry={app.deleteMetadataCacheEntry}
		onClearMetadataSearchHistory={app.clearMetadataSearchHistory}
		onLoadMoreMetadataCache={app.loadMoreMetadataCache}
		onLoadMoreMetadataSearchHistory={app.loadMoreMetadataSearchHistory}
	/>
{/if}
