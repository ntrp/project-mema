<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import UserIcon from '@lucide/svelte/icons/user';
	import { resolve } from '$app/paths';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import MediaFileSearchModal from '$lib/components/app/media/files/MediaFileSearchModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { releaseSearchQuery } from '$lib/components/app/media/release-search/releaseSearchQuery';
	import type {
		Language,
		MediaItem,
		ReleaseCandidate,
		ReleaseOverrideDetails
	} from '$lib/settings/types';

	interface Props {
		items: MediaItem[];
		languages: Language[];
		searchingItemId?: string;
		grabbingKey?: string;
		canManage: boolean;
		onFindReleases: (_item: MediaItem, _query?: string) => void;
		onGrabRelease: (
			_item: MediaItem,
			_release: ReleaseCandidate,
			_overrideMatch?: boolean,
			_details?: ReleaseOverrideDetails
		) => void;
	}

	let {
		items,
		languages,
		searchingItemId,
		grabbingKey,
		canManage,
		onFindReleases,
		onGrabRelease
	}: Props = $props();
	let manualSearchItem = $state<MediaItem | undefined>();

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
								<div class="flex items-center justify-end gap-2">
									<Tooltip.Root>
										<Tooltip.Trigger>
											{#snippet child({ props })}
												<Button
													{...props}
													type="button"
													variant="outline"
													size="icon-sm"
													aria-label={`Automatic search ${item.title}`}
													disabled={searchingItemId === item.id}
													onclick={() => onFindReleases(item, releaseSearchQuery(item))}
												>
													<SearchIcon aria-hidden="true" />
												</Button>
											{/snippet}
										</Tooltip.Trigger>
										<Tooltip.Content>Automatic search</Tooltip.Content>
									</Tooltip.Root>
									<Tooltip.Root>
										<Tooltip.Trigger>
											{#snippet child({ props })}
												<Button
													{...props}
													type="button"
													variant="outline"
													size="icon-sm"
													aria-label={`Manual search ${item.title}`}
													disabled={searchingItemId === item.id}
													onclick={() => (manualSearchItem = item)}
												>
													<UserIcon aria-hidden="true" />
												</Button>
											{/snippet}
										</Tooltip.Trigger>
										<Tooltip.Content>Manual search</Tooltip.Content>
									</Tooltip.Root>
								</div>
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

{#if manualSearchItem}
	<MediaFileSearchModal
		item={manualSearchItem}
		{languages}
		searchContext={{ type: 'title' }}
		{grabbingKey}
		{canManage}
		onGrab={onGrabRelease}
		onClose={() => (manualSearchItem = undefined)}
	/>
{/if}
