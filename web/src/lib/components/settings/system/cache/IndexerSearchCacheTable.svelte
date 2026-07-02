<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { IndexerSearchCacheEntry, IndexerSearchResponse } from '$lib/settings/types';
	import InfiniteTableFrame from './InfiniteTableFrame.svelte';
	import { createRowPulse } from './rowPulse.svelte';

	interface Props {
		search: IndexerSearchResponse;
		clearing: boolean;
		loading: boolean;
		onDeleteEntry: (_entry: IndexerSearchCacheEntry) => void | Promise<void>;
		onLoadMore: () => void | Promise<void>;
	}

	let { search, clearing, loading, onDeleteEntry, onLoadMore }: Props = $props();
	const rowPulse = createRowPulse();
	const hasMore = $derived(search.cacheEntries.length < search.stats.totalEntries);
	const rowKeys = $derived(search.cacheEntries.map(cacheEntryKey));

	$effect(() => rowPulse.update(rowKeys));

	function cacheEntryKey(entry: IndexerSearchCacheEntry) {
		return `${entry.indexerName}:${entry.mediaType}:${entry.query}`;
	}
</script>

<InfiniteTableFrame {hasMore} {loading} {onLoadMore}>
	{#if search.cacheEntries.length > 0}
		<Table.Root class="min-w-full table-auto border-collapse">
			<Table.Header class="sticky top-0 bg-card">
				<Table.Row>
					<Table.Head class="w-px">Expires</Table.Head>
					<Table.Head class="w-px">Indexer</Table.Head>
					<Table.Head class="w-px">Media</Table.Head>
					<Table.Head>Query</Table.Head>
					<Table.Head class="w-px">Items</Table.Head>
					<Table.Head class="w-px text-right">Actions</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each search.cacheEntries as entry (cacheEntryKey(entry))}
					{@const rowKey = cacheEntryKey(entry)}
					<Table.Row class={rowPulse.classFor(rowKey)}>
						<Table.Cell class="w-px">
							{entry.expired ? 'Expired' : formatDateTime(entry.expiresAt)}
						</Table.Cell>
						<Table.Cell class="max-w-48">
							<strong class="block truncate">{entry.indexerName}</strong>
						</Table.Cell>
						<Table.Cell class="w-px">{entry.mediaType}</Table.Cell>
						<Table.Cell class="max-w-96">
							<code class="block truncate text-xs">{entry.query}</code>
						</Table.Cell>
						<Table.Cell class="w-px">{entry.resultCount}</Table.Cell>
						<Table.Cell class="w-px text-right">
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="destructive"
											size="icon-sm"
											aria-label="Delete cache entry"
											disabled={clearing}
											onclick={() => void onDeleteEntry(entry)}
										>
											<TrashIcon aria-hidden="true" />
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>Delete cache entry</Tooltip.Content>
							</Tooltip.Root>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{:else}
		<p class="m-0 p-3 text-sm leading-6 text-muted-foreground">No indexer cache entries yet.</p>
	{/if}
</InfiniteTableFrame>
