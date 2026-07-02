<script lang="ts">
	import SystemEventsSettings from '$lib/components/settings/system/events/SystemEventsSettings.svelte';
	import IndexerSearchCacheSettings from '$lib/components/settings/system/cache/IndexerSearchCacheSettings.svelte';
	import MetadataCacheSettings from '$lib/components/settings/system/cache/MetadataCacheSettings.svelte';
	import SystemJobsSettings from '$lib/components/settings/system/jobs/SystemJobsSettings.svelte';
	import SystemLogFilesSettings from '$lib/components/settings/system/logs/SystemLogFilesSettings.svelte';
	import SystemLogsSettings from '$lib/components/settings/system/logs/SystemLogsSettings.svelte';
	import SystemStatusSettings from '$lib/components/settings/system/SystemStatusSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import type {
		IndexerSearchCacheEntry,
		IndexerSearchResponse,
		MetadataCacheEntry,
		MetadataCacheResponse,
		SystemSection
	} from '$lib/settings/types';

	interface Props {
		activeSection: SystemSection;
		indexerSearch: IndexerSearchResponse;
		metadataCache: MetadataCacheResponse;
		loadingIndexerSearch: boolean;
		loadingMetadataCache: boolean;
		clearingIndexerSearchCache: boolean;
		clearingMetadataCache: boolean;
		onClearIndexerSearchCache: () => void | Promise<void>;
		onClearIndexerSearchCachePattern: (_pattern: string) => void | Promise<void>;
		onDeleteIndexerSearchCacheEntry: (_entry: IndexerSearchCacheEntry) => void | Promise<void>;
		onClearIndexerSearchHistory: () => void | Promise<void>;
		onLoadMoreIndexerSearchCache: () => void | Promise<void>;
		onLoadMoreIndexerSearchHistory: () => void | Promise<void>;
		onClearMetadataCache: () => void | Promise<void>;
		onClearMetadataCachePattern: (_pattern: string) => void | Promise<void>;
		onDeleteMetadataCacheEntry: (_entry: MetadataCacheEntry) => void | Promise<void>;
		onClearMetadataSearchHistory: () => void | Promise<void>;
		onLoadMoreMetadataCache: () => void | Promise<void>;
		onLoadMoreMetadataSearchHistory: () => void | Promise<void>;
	}

	let {
		activeSection,
		indexerSearch,
		metadataCache,
		loadingIndexerSearch,
		loadingMetadataCache,
		clearingIndexerSearchCache,
		clearingMetadataCache,
		onClearIndexerSearchCache,
		onClearIndexerSearchCachePattern,
		onDeleteIndexerSearchCacheEntry,
		onClearIndexerSearchHistory,
		onLoadMoreIndexerSearchCache,
		onLoadMoreIndexerSearchHistory,
		onClearMetadataCache,
		onClearMetadataCachePattern,
		onDeleteMetadataCacheEntry,
		onClearMetadataSearchHistory,
		onLoadMoreMetadataCache,
		onLoadMoreMetadataSearchHistory
	}: Props = $props();
	let logsConnected = $state(false);
	let indexerCachePattern = $state('');
	let metadataCachePattern = $state('');

	function connectionDotClass(connected: boolean) {
		const base = 'ml-2 inline-block size-3 translate-y-[-2px] rounded-full';
		return connected
			? `${base} animate-pulse bg-primary ring-4 ring-primary/20`
			: `${base} bg-muted-foreground ring-4 ring-muted`;
	}
</script>

<section aria-labelledby="system-title">
	{#if activeSection === 'status'}
		<PageHeading eyebrow="System" title="Status" titleId="system-title" />
		<div class="space-y-4">
			<SystemStatusSettings />
		</div>
	{:else if activeSection === 'indexing'}
		<PageHeading eyebrow="System" title="Indexing" titleId="system-title" />
		<div class="space-y-4">
			<IndexerSearchCacheSettings
				search={indexerSearch}
				bind:pattern={indexerCachePattern}
				clearing={clearingIndexerSearchCache}
				loading={loadingIndexerSearch}
				onClearAll={onClearIndexerSearchCache}
				onClearPattern={onClearIndexerSearchCachePattern}
				onDeleteEntry={onDeleteIndexerSearchCacheEntry}
				onClearHistory={onClearIndexerSearchHistory}
				onLoadMoreCache={onLoadMoreIndexerSearchCache}
				onLoadMoreHistory={onLoadMoreIndexerSearchHistory}
			/>
		</div>
	{:else if activeSection === 'metadata'}
		<PageHeading eyebrow="System" title="Metadata" titleId="system-title" />
		<div class="space-y-4">
			<MetadataCacheSettings
				cache={metadataCache}
				bind:pattern={metadataCachePattern}
				clearing={clearingMetadataCache}
				loading={loadingMetadataCache}
				onClearAll={onClearMetadataCache}
				onClearPattern={onClearMetadataCachePattern}
				onDeleteEntry={onDeleteMetadataCacheEntry}
				onClearHistory={onClearMetadataSearchHistory}
				onLoadMoreCache={onLoadMoreMetadataCache}
				onLoadMoreHistory={onLoadMoreMetadataSearchHistory}
			/>
		</div>
	{:else if activeSection === 'jobs'}
		<PageHeading eyebrow="System" title="Jobs" titleId="system-title" />
		<div class="space-y-4">
			<SystemJobsSettings />
		</div>
	{:else if activeSection === 'events'}
		<PageHeading eyebrow="System" title="Events" titleId="system-title" />
		<div class="space-y-4">
			<SystemEventsSettings />
		</div>
	{:else}
		<PageHeading eyebrow="System" title="Logs" titleId="system-title">
			<span
				class={connectionDotClass(logsConnected)}
				aria-label={logsConnected ? 'Log stream connected' : 'Log stream reconnecting'}
			></span>
		</PageHeading>
		<div class="space-y-4">
			<SystemLogsSettings onConnectionChange={(connected) => (logsConnected = connected)} />
			<SystemLogFilesSettings />
		</div>
	{/if}
</section>
