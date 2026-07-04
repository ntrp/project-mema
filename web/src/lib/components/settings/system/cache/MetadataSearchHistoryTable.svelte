<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { MetadataCacheResponse } from '$lib/settings/types';
	import HistoryResponseActions from './HistoryResponseActions.svelte';
	import InfiniteTableFrame from './InfiniteTableFrame.svelte';
	import { createRowPulse } from './rowPulse.svelte';

	interface Props {
		cache: MetadataCacheResponse;
		loading: boolean;
		onLoadMore: () => void | Promise<void>;
	}

	let { cache, loading, onLoadMore }: Props = $props();
	const rowPulse = createRowPulse();
	const hasMore = $derived(cache.historyEntries.length < cache.historyTotalEntries);
	const rowKeys = $derived(cache.historyEntries.map(historyEntryKey));
	const stats = $derived([
		{ label: 'Total queries', value: cache.historyStats.totalEntries },
		{ label: 'Cache hits', value: cache.historyStats.cacheHits },
		{ label: 'Cache misses', value: cache.historyStats.cacheMisses },
		{ label: 'Failures', value: cache.historyStats.failures }
	]);

	$effect(() => rowPulse.update(rowKeys));

	function historyEntryKey(entry: MetadataCacheResponse['historyEntries'][number]) {
		return `${entry.createdAt}:${entry.providerName}:${entry.query}`;
	}
</script>

<div class="grid gap-4">
	<dl class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4" aria-label="Metadata history stats">
		{#each stats as stat (stat.label)}
			<div class="grid gap-1 px-3 py-2.5">
				<dt class="text-xs font-semibold text-muted-foreground">{stat.label}</dt>
				<dd class="m-0 text-xl leading-none font-bold text-foreground">{stat.value}</dd>
			</div>
		{/each}
	</dl>

	<InfiniteTableFrame {hasMore} {loading} {onLoadMore}>
		{#if cache.historyEntries.length > 0}
			<Table.Root class="min-w-full table-auto border-collapse">
				<Table.Header class="sticky top-0 bg-card">
					<Table.Row>
						<Table.Head class="w-px">Time</Table.Head>
						<Table.Head class="w-px">Provider</Table.Head>
						<Table.Head class="w-px">Kind</Table.Head>
						<Table.Head class="w-px">Media</Table.Head>
						<Table.Head>Query</Table.Head>
						<Table.Head class="w-px text-right">Cache</Table.Head>
						<Table.Head class="w-px text-right">Items</Table.Head>
						<Table.Head class="w-px text-right">Response</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each cache.historyEntries as entry (historyEntryKey(entry))}
						{@const rowKey = historyEntryKey(entry)}
						<Table.Row class={rowPulse.classFor(rowKey)}>
							<Table.Cell class="w-px">{formatDateTime(entry.createdAt)}</Table.Cell>
							<Table.Cell class="w-px">
								<strong class="block truncate">{entry.providerName}</strong>
							</Table.Cell>
							<Table.Cell class="w-px">{entry.cacheKind}</Table.Cell>
							<Table.Cell class="w-px">
								{entry.mediaType}{entry.year ? ` · ${entry.year}` : ''}
							</Table.Cell>
							<Table.Cell class="max-w-96">
								<code class="block truncate text-xs">{entry.query}</code>
							</Table.Cell>
							<Table.Cell class="w-px text-right">
								<Badge
									variant="outline"
									class={entry.cacheHit
										? 'border-emerald-500/50 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300'
										: 'border-amber-500/50 bg-amber-500/10 text-amber-700 dark:text-amber-300'}
								>
									{entry.cacheHit ? 'Hit' : 'Miss'}
								</Badge>
							</Table.Cell>
							<Table.Cell class="w-px text-right">{entry.itemCount}</Table.Cell>
							<Table.Cell class="w-px text-right">
								<HistoryResponseActions
									value={entry.error ?? entry.response}
									success={entry.success}
								/>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		{:else}
			<p class="m-0 p-3 text-sm leading-6 text-muted-foreground">No metadata query history yet.</p>
		{/if}
	</InfiniteTableFrame>
</div>
