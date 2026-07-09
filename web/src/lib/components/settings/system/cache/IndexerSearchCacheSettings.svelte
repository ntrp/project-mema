<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import LivePulseDot from '$lib/components/shared/LivePulseDot.svelte';
	import * as Card from '$lib/components/ui/card';
	import type { IndexerSearchCacheEntry, IndexerSearchResponse } from '$lib/settings/types';
	import IndexerSearchCacheTable from './IndexerSearchCacheTable.svelte';
	import IndexerSearchHistoryTable from './IndexerSearchHistoryTable.svelte';
	import MetadataCacheControls from './MetadataCacheControls.svelte';

	interface Props {
		search: IndexerSearchResponse;
		pattern: string;
		clearing: boolean;
		loading: boolean;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_pattern: string) => void | Promise<void>;
		onDeleteEntry: (_entry: IndexerSearchCacheEntry) => void | Promise<void>;
		onClearHistory: () => void | Promise<void>;
		onLoadMoreCache: () => void | Promise<void>;
		onLoadMoreHistory: () => void | Promise<void>;
	}

	let {
		search,
		pattern = $bindable(),
		clearing,
		loading,
		onClearAll,
		onClearPattern,
		onDeleteEntry,
		onClearHistory,
		onLoadMoreCache,
		onLoadMoreHistory
	}: Props = $props();

	const stats = $derived([
		{ label: 'Total entries', value: search.stats.totalEntries },
		{ label: 'Active', value: search.stats.activeEntries },
		{ label: 'Expired', value: search.stats.expiredEntries },
		{ label: 'Indexers', value: search.stats.indexerCount }
	]);
</script>

<div class="grid gap-4">
	<Card.Root aria-labelledby="indexer-cache-title">
		<Card.Header>
			<div>
				<Card.Description class="flex items-center gap-2">
					<LivePulseDot />
					<span>Live</span>
				</Card.Description>
				<Card.Title id="indexer-cache-title">Query Cache</Card.Title>
			</div>
			<Card.Action>
				<ConfirmActionButton
					label="Reset indexer query cache"
					title="Reset indexer query cache"
					description="Delete every cached indexer query result?"
					confirmLabel="Reset cache"
					confirmingLabel="Resetting"
					size="icon-sm"
					disabled={clearing}
					tooltip="Reset indexer query cache"
					onConfirm={onClearAll}
				>
					<TrashIcon aria-hidden="true" />
				</ConfirmActionButton>
			</Card.Action>
		</Card.Header>

		<Card.Content class="grid gap-4">
			<dl class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4" aria-label="Indexer cache stats">
				{#each stats as stat (stat.label)}
					<div class="grid gap-1 px-3 py-2.5">
						<dt class="text-xs font-semibold text-muted-foreground">{stat.label}</dt>
						<dd class="m-0 text-xl leading-none font-bold text-foreground">{stat.value}</dd>
					</div>
				{/each}
			</dl>

			<MetadataCacheControls
				bind:pattern
				{clearing}
				showClearAll={false}
				{onClearAll}
				{onClearPattern}
			/>

			<IndexerSearchCacheTable
				{search}
				{clearing}
				{loading}
				{onDeleteEntry}
				onLoadMore={onLoadMoreCache}
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root aria-labelledby="indexer-history-title">
		<Card.Header>
			<div>
				<Card.Description class="flex items-center gap-2">
					<LivePulseDot />
					<span>Live</span>
				</Card.Description>
				<Card.Title id="indexer-history-title">Query History</Card.Title>
			</div>
			<Card.Action>
				<ConfirmActionButton
					label="Clear indexer query history"
					title="Clear indexer query history"
					description="Delete every recorded indexer query history entry?"
					confirmLabel="Clear history"
					confirmingLabel="Clearing"
					size="icon-sm"
					disabled={clearing || search.historyEntries.length === 0}
					tooltip="Clear indexer query history"
					onConfirm={onClearHistory}
				>
					<TrashIcon aria-hidden="true" />
				</ConfirmActionButton>
			</Card.Action>
		</Card.Header>
		<Card.Content>
			<IndexerSearchHistoryTable {search} {loading} onLoadMore={onLoadMoreHistory} />
		</Card.Content>
	</Card.Root>
</div>
