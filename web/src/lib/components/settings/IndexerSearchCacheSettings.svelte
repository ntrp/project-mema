<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { IndexerSearchResponse } from '$lib/settings/types';
	import IndexerSearchHistoryTable from './IndexerSearchHistoryTable.svelte';
	import MetadataCacheControls from './MetadataCacheControls.svelte';

	interface Props {
		search: IndexerSearchResponse;
		pattern: string;
		loading: boolean;
		clearing: boolean;
		onRefresh: () => void | Promise<void>;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_pattern: string) => void | Promise<void>;
	}

	let {
		search,
		pattern = $bindable(),
		loading,
		clearing,
		onRefresh,
		onClearAll,
		onClearPattern
	}: Props = $props();

	const stats = $derived([
		{ label: 'Total entries', value: search.stats.totalEntries },
		{ label: 'Active', value: search.stats.activeEntries },
		{ label: 'Expired', value: search.stats.expiredEntries },
		{ label: 'Indexers', value: search.stats.indexerCount }
	]);
</script>

<Card.Root aria-labelledby="indexer-cache-title">
	<Card.Header>
		<div>
			<Card.Description class="flex items-center gap-2">
				<span class="relative flex size-2.5">
					<span
						class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75"
					></span>
					<span class="relative inline-flex size-2.5 rounded-full bg-emerald-500"></span>
				</span>
				<span>Live</span>
			</Card.Description>
			<Card.Title id="indexer-cache-title">Query Cache</Card.Title>
		</div>
		<Card.Action>
			<Button type="button" variant="outline" disabled={loading} onclick={() => void onRefresh()}>
				{loading ? 'Refreshing' : 'Refresh'}
			</Button>
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
			{onClearAll}
			{onClearPattern}
		/>

		<div class="max-h-80 overflow-auto rounded-md border border-border">
			{#if search.cacheEntries.length > 0}
				<Table.Root>
					<Table.Header class="sticky top-0 bg-card">
						<Table.Row>
							<Table.Head>Indexer</Table.Head>
							<Table.Head>Media</Table.Head>
							<Table.Head>Query</Table.Head>
							<Table.Head>Items</Table.Head>
							<Table.Head>Expires</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each search.cacheEntries as entry (`${entry.indexerName}:${entry.mediaType}:${entry.query}`)}
							<Table.Row>
								<Table.Cell>
									<strong class="block">{entry.indexerName}</strong>
									<span class="block text-xs text-muted-foreground">{entry.indexerType}</span>
								</Table.Cell>
								<Table.Cell>{entry.mediaType}</Table.Cell>
								<Table.Cell><code class="text-xs">{entry.query}</code></Table.Cell>
								<Table.Cell>{entry.resultCount}</Table.Cell>
								<Table.Cell>
									<span class={entry.expired ? 'text-muted-foreground' : 'text-primary'}>
										{entry.expired ? 'Expired' : formatDateTime(entry.expiresAt)}
									</span>
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			{:else}
				<p class="m-0 p-3 text-sm leading-6 text-muted-foreground">No indexer cache entries yet.</p>
			{/if}
		</div>
	</Card.Content>
</Card.Root>

<Card.Root aria-labelledby="indexer-history-title">
	<Card.Header>
		<div>
			<Card.Description class="flex items-center gap-2">
				<span class="relative flex size-2.5">
					<span
						class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75"
					></span>
					<span class="relative inline-flex size-2.5 rounded-full bg-emerald-500"></span>
				</span>
				<span>Live</span>
			</Card.Description>
			<Card.Title id="indexer-history-title">Query History</Card.Title>
		</div>
	</Card.Header>
	<Card.Content>
		<IndexerSearchHistoryTable {search} />
	</Card.Content>
</Card.Root>
