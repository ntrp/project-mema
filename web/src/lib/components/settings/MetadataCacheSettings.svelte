<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { MetadataCacheResponse } from '$lib/settings/types';
	import MetadataCacheControls from './MetadataCacheControls.svelte';
	import MetadataSearchHistoryTable from './MetadataSearchHistoryTable.svelte';

	interface Props {
		cache: MetadataCacheResponse;
		pattern: string;
		loading: boolean;
		clearing: boolean;
		onRefresh: () => void | Promise<void>;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_pattern: string) => void | Promise<void>;
	}

	let {
		cache,
		pattern = $bindable(),
		loading,
		clearing,
		onRefresh,
		onClearAll,
		onClearPattern
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
					<span class="relative flex size-2.5">
						<span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75"></span>
						<span class="relative inline-flex size-2.5 rounded-full bg-emerald-500"></span>
					</span>
					<span>Live cache</span>
				</Card.Description>
				<Card.Title id="metadata-cache-title">Metadata provider cache</Card.Title>
			</div>
			<Card.Action>
				<Button type="button" variant="outline" disabled={loading} onclick={() => void onRefresh()}>
					{loading ? 'Refreshing' : 'Refresh'}
				</Button>
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
				{onClearAll}
				{onClearPattern}
			/>

			<div class="max-h-80 overflow-auto rounded-md border border-border">
				{#if cache.entries.length > 0}
					<Table.Root>
						<Table.Header class="sticky top-0 bg-card">
							<Table.Row>
								<Table.Head>Provider</Table.Head>
								<Table.Head>Kind</Table.Head>
								<Table.Head>Media</Table.Head>
								<Table.Head>Key</Table.Head>
								<Table.Head>Items</Table.Head>
								<Table.Head>Expires</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each cache.entries as entry (`${entry.providerName}:${entry.mediaType}:${entry.query}:${entry.year}`)}
								<Table.Row>
									<Table.Cell>
										<strong class="block">{entry.providerName}</strong>
										<span class="block text-xs text-muted-foreground">{entry.providerType}</span>
									</Table.Cell>
									<Table.Cell>{entry.cacheKind}</Table.Cell>
									<Table.Cell>{entry.mediaType}{entry.year ? ` · ${entry.year}` : ''}</Table.Cell>
									<Table.Cell><code class="text-xs">{entry.query}</code></Table.Cell>
									<Table.Cell>{entry.itemCount}</Table.Cell>
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
					<p class="m-0 p-3 text-sm leading-6 text-muted-foreground">
						No metadata cache entries yet.
					</p>
				{/if}
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root aria-labelledby="metadata-history-title">
		<Card.Header>
			<Card.Description class="flex items-center gap-2">
				<span class="relative flex size-2.5">
					<span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75"></span>
					<span class="relative inline-flex size-2.5 rounded-full bg-emerald-500"></span>
				</span>
				<span>Live history</span>
			</Card.Description>
			<Card.Title id="metadata-history-title">Metadata query history</Card.Title>
		</Card.Header>
		<Card.Content>
			<MetadataSearchHistoryTable {cache} />
		</Card.Content>
	</Card.Root>
</div>
