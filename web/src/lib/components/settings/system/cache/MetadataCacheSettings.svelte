<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import LivePulseDot from '$lib/components/shared/LivePulseDot.svelte';
	import * as Card from '$lib/components/ui/card';
	import type { MetadataCacheEntry, MetadataCacheResponse } from '$lib/settings/types';
	import MetadataCacheControls from './MetadataCacheControls.svelte';
	import MetadataCacheTable from './MetadataCacheTable.svelte';
	import MetadataSearchHistoryTable from './MetadataSearchHistoryTable.svelte';

	interface Props {
		cache: MetadataCacheResponse;
		pattern: string;
		clearing: boolean;
		loading: boolean;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_pattern: string) => void | Promise<void>;
		onDeleteEntry: (_entry: MetadataCacheEntry) => void | Promise<void>;
		onClearHistory: () => void | Promise<void>;
		onLoadMoreCache: () => void | Promise<void>;
		onLoadMoreHistory: () => void | Promise<void>;
	}

	let {
		cache,
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
		{ label: 'Total entries', value: cache.stats.totalEntries },
		{ label: 'Active', value: cache.stats.activeEntries },
		{ label: 'Expired', value: cache.stats.expiredEntries },
		{ label: 'Providers', value: cache.stats.providerCount }
	]);
</script>

<div class="grid gap-4">
	<Card.Root aria-labelledby="metadata-cache-title">
		<Card.Header>
			<div>
				<Card.Description class="flex items-center gap-2">
					<LivePulseDot />
					<span>Live</span>
				</Card.Description>
				<Card.Title id="metadata-cache-title">Metadata Cache</Card.Title>
			</div>
			<Card.Action>
				<ConfirmActionButton
					label="Reset metadata cache"
					title="Reset metadata cache"
					description="Delete every cached metadata result?"
					confirmLabel="Reset cache"
					confirmingLabel="Resetting"
					size="icon-sm"
					disabled={clearing}
					tooltip="Reset metadata cache"
					onConfirm={onClearAll}
				>
					<TrashIcon aria-hidden="true" />
				</ConfirmActionButton>
			</Card.Action>
		</Card.Header>

		<Card.Content class="grid gap-4">
			<dl class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4" aria-label="Metadata cache stats">
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

			<MetadataCacheTable
				{cache}
				{clearing}
				{loading}
				{onDeleteEntry}
				onLoadMore={onLoadMoreCache}
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root aria-labelledby="metadata-history-title">
		<Card.Header>
			<div>
				<Card.Description class="flex items-center gap-2">
					<LivePulseDot />
					<span>Live</span>
				</Card.Description>
				<Card.Title id="metadata-history-title">Query History</Card.Title>
			</div>
			<Card.Action>
				<ConfirmActionButton
					label="Clear metadata query history"
					title="Clear metadata query history"
					description="Delete every recorded metadata query history entry?"
					confirmLabel="Clear history"
					confirmingLabel="Clearing"
					size="icon-sm"
					disabled={clearing || cache.historyEntries.length === 0}
					tooltip="Clear metadata query history"
					onConfirm={onClearHistory}
				>
					<TrashIcon aria-hidden="true" />
				</ConfirmActionButton>
			</Card.Action>
		</Card.Header>
		<Card.Content>
			<MetadataSearchHistoryTable {cache} {loading} onLoadMore={onLoadMoreHistory} />
		</Card.Content>
	</Card.Root>
</div>
