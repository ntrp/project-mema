<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Tooltip from '$lib/components/ui/tooltip';
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

<div
	class="grid min-h-[50rem] grid-rows-[minmax(24rem,1fr)_minmax(24rem,1fr)] gap-4 lg:h-[calc(100vh-12rem)]"
>
	<Card.Root class="min-h-0" aria-labelledby="metadata-cache-title">
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
				<Card.Title id="metadata-cache-title">Metadata Cache</Card.Title>
			</div>
			<Card.Action>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								{...props}
								type="button"
								variant="destructive"
								size="icon-sm"
								aria-label="Reset metadata cache"
								disabled={clearing}
								onclick={() => void onClearAll()}
							>
								<TrashIcon aria-hidden="true" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Reset metadata cache</Tooltip.Content>
				</Tooltip.Root>
			</Card.Action>
		</Card.Header>

		<Card.Content class="grid min-h-0 flex-1 grid-rows-[auto_auto_minmax(0,1fr)] gap-4">
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

	<Card.Root class="min-h-0" aria-labelledby="metadata-history-title">
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
				<Card.Title id="metadata-history-title">Query History</Card.Title>
			</div>
			<Card.Action>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Button
								{...props}
								type="button"
								variant="destructive"
								size="icon-sm"
								aria-label="Clear metadata query history"
								disabled={clearing || cache.historyEntries.length === 0}
								onclick={() => void onClearHistory()}
							>
								<TrashIcon aria-hidden="true" />
							</Button>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>Clear metadata query history</Tooltip.Content>
				</Tooltip.Root>
			</Card.Action>
		</Card.Header>
		<Card.Content class="min-h-0 flex-1">
			<MetadataSearchHistoryTable
				{cache}
				{loading}
				onLoadMore={onLoadMoreHistory}
			/>
		</Card.Content>
	</Card.Root>
</div>
