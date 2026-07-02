<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { IndexerSearchResponse } from '$lib/settings/types';

	interface Props {
		search: IndexerSearchResponse;
	}

	let { search }: Props = $props();
	const stats = $derived([
		{ label: 'Displayed entries', value: search.historyEntries.length },
		{ label: 'Cache hits', value: search.historyEntries.filter((entry) => entry.cacheHit).length },
		{
			label: 'Cache misses',
			value: search.historyEntries.filter((entry) => !entry.cacheHit).length
		},
		{ label: 'Failures', value: search.historyEntries.filter((entry) => !entry.success).length }
	]);
</script>

<div class="grid gap-4">
	<dl class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4" aria-label="Indexer history stats">
		{#each stats as stat (stat.label)}
			<div class="grid gap-1 px-3 py-2.5">
				<dt class="text-xs font-semibold text-muted-foreground">{stat.label}</dt>
				<dd class="m-0 text-xl leading-none font-bold text-foreground">{stat.value}</dd>
			</div>
		{/each}
	</dl>

	<div class="max-h-96 overflow-auto rounded-md border border-border">
		{#if search.historyEntries.length > 0}
			<Table.Root class="min-w-[980px]">
				<Table.Header class="sticky top-0 bg-card">
					<Table.Row>
						<Table.Head>Time</Table.Head>
						<Table.Head>Indexer</Table.Head>
						<Table.Head>Media</Table.Head>
						<Table.Head>Query</Table.Head>
						<Table.Head>Cache</Table.Head>
						<Table.Head>Items</Table.Head>
						<Table.Head>Response</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each search.historyEntries as entry (`${entry.createdAt}:${entry.indexerName}:${entry.query}`)}
						<Table.Row>
							<Table.Cell>{formatDateTime(entry.createdAt)}</Table.Cell>
							<Table.Cell>
								<strong class="block">{entry.indexerName}</strong>
								<span class="block text-xs text-muted-foreground">{entry.indexerType}</span>
							</Table.Cell>
							<Table.Cell>{entry.mediaType}</Table.Cell>
							<Table.Cell><code class="text-xs">{entry.query}</code></Table.Cell>
							<Table.Cell>
								<Badge variant={entry.cacheHit ? 'secondary' : 'outline'}>
									{entry.cacheHit ? 'Hit' : 'Miss'}
								</Badge>
							</Table.Cell>
							<Table.Cell>{entry.resultCount}</Table.Cell>
							<Table.Cell class="max-w-80">
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<code
												{...props}
												class={entry.success
													? 'block truncate text-xs'
													: 'block truncate text-xs text-destructive'}
											>
												{entry.error ?? entry.response}
											</code>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content class="max-w-160 whitespace-pre-wrap">
										{entry.error ?? entry.response}
									</Tooltip.Content>
								</Tooltip.Root>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		{:else}
			<p class="m-0 p-3 text-sm leading-6 text-muted-foreground">No indexer query history yet.</p>
		{/if}
	</div>
</div>
