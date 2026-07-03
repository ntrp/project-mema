<script lang="ts">
	import { resolve } from '$app/paths';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { releaseSearchQuery } from '$lib/components/app/media/release-search/releaseSearchQuery';
	import type { MediaItem } from '$lib/settings/types';

	interface Props {
		items: MediaItem[];
		searchingItemId?: string;
		canManage: boolean;
		onFindReleases: (_item: MediaItem, _query?: string) => void;
	}

	let { items, searchingItemId, canManage, onFindReleases }: Props = $props();

	function monitorLabel(item: MediaItem) {
		if (!item.monitored || item.monitorMode === 'none') {
			return 'None';
		}
		return item.monitorMode === 'collection' ? 'Entire collection' : 'This media only';
	}
</script>

<PageHeading eyebrow="Library" title="Wanted" titleId="home-title" />

{#if items.length}
	<Card class="overflow-hidden p-0">
		<Table.Root class="wanted-table">
			<Table.Header>
				<Table.Row>
					<Table.Head>Title</Table.Head>
					<Table.Head>Type</Table.Head>
					<Table.Head>Year</Table.Head>
					<Table.Head>Monitor</Table.Head>
					<Table.Head>Profile</Table.Head>
					<Table.Head>Availability</Table.Head>
					<Table.Head></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each items as item (item.id)}
					<Table.Row>
						<Table.Cell>
							<a
								href={item.type === 'movie'
									? resolve('/movies/[id]', { id: item.id })
									: resolve('/series/[id]', { id: item.id })}
							>
								{item.title}
							</a>
						</Table.Cell>
						<Table.Cell>{item.type}</Table.Cell>
						<Table.Cell>{item.year ?? '-'}</Table.Cell>
						<Table.Cell>{monitorLabel(item)}</Table.Cell>
						<Table.Cell>{item.qualityProfileName ?? '-'}</Table.Cell>
						<Table.Cell>{item.minimumAvailability}</Table.Cell>
						<Table.Cell class="text-right">
							{#if canManage}
								<Button
									type="button"
									variant="outline"
									size="sm"
									disabled={searchingItemId === item.id}
									onclick={() => onFindReleases(item, releaseSearchQuery(item))}
								>
									{searchingItemId === item.id ? 'Searching' : 'Search'}
								</Button>
							{/if}
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card>
{:else}
	<EmptyState
		class="my-[18px] grid min-h-60 w-full place-items-center content-center gap-[18px] text-center"
	>
		<p class="m-0 text-lg font-black text-foreground">No missing media.</p>
	</EmptyState>
{/if}
