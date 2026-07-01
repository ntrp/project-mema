<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Table from '$lib/components/ui/table';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { MetadataCacheResponse } from '$lib/settings/types';

	interface Props {
		cache: MetadataCacheResponse;
		pattern: string;
		loading: boolean;
		clearing: boolean;
		onRefresh: () => void | Promise<void>;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_event: SubmitEvent) => void | Promise<void>;
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

<Card.Root aria-labelledby="metadata-cache-title">
	<Card.Header>
		<div>
			<Card.Description>Cache</Card.Description>
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

		<form class="grid items-end gap-3 md:grid-cols-[minmax(0,1fr)_auto]" onsubmit={onClearPattern}>
			<div class="grid gap-1.5">
				<Label>Reset by regex</Label>
				<Input bind:value={pattern} placeholder="discover:|details:123|matrix" autocomplete="off" />
			</div>
			<div class="flex flex-wrap justify-end gap-2">
				<Button
					type="submit"
					variant="destructive"
					disabled={clearing || pattern.trim().length === 0}
				>
					{clearing ? 'Resetting' : 'Reset matching'}
				</Button>
				<Button
					type="button"
					variant="destructive"
					disabled={clearing}
					onclick={() => void onClearAll()}
				>
					Reset all
				</Button>
			</div>
		</form>

		{#if cache.entries.length > 0}
			<Table.Root>
				<Table.Header>
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
			<p class="m-0 text-sm leading-6 text-muted-foreground">No metadata cache entries yet.</p>
		{/if}
	</Card.Content>
</Card.Root>
