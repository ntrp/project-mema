<script lang="ts">
	import EyeIcon from '@lucide/svelte/icons/eye';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { DiscoverBlacklistItem } from '$lib/settings/types';
	import MediaBadge from '../media/MediaBadge.svelte';
	import PosterPlaceholder from '../media/PosterPlaceholder.svelte';

	interface Props {
		items: DiscoverBlacklistItem[];
		loading: boolean;
		removingId?: string;
		onRemove: (_item: DiscoverBlacklistItem) => void;
	}

	let { items, loading, removingId, onRemove }: Props = $props();

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w342${path}`;
	}
</script>

<PageHeading
	eyebrow="Discover"
	title="Blacklist"
	titleId="home-title"
	description={`${items.length} hidden titles`}
/>

{#if loading}
	<div
		class="grid grid-cols-[repeat(auto-fill,minmax(132px,1fr))] gap-3 sm:grid-cols-[repeat(auto-fill,minmax(190px,220px))] sm:gap-5"
	>
		{#each Array.from({ length: 8 }) as _, index (index)}
			<div class="min-w-0 snap-start aspect-[2/3] rounded-md bg-card" aria-hidden="true"></div>
		{/each}
	</div>
{:else if items.length === 0}
	<EmptyState
		title="No blacklisted media"
		description="Use the hidden-eye action on discover cards to hide titles from discovery."
	/>
{:else}
	<div
		class="grid grid-cols-[repeat(auto-fill,minmax(132px,1fr))] gap-3 sm:grid-cols-[repeat(auto-fill,minmax(170px,220px))] sm:gap-5"
	>
		{#each items as item (item.id)}
			<article class="group/poster min-w-0 snap-start">
				<div
					class="relative aspect-[2/3] overflow-hidden rounded-md border border-border bg-card transition-[transform,border-color,box-shadow] duration-150 group-hover/poster:z-[2] group-hover/poster:-translate-y-1.5 group-hover/poster:scale-105 group-hover/poster:border-primary/50 group-hover/poster:shadow-xl group-focus-within/poster:z-[2] group-focus-within/poster:-translate-y-1.5 group-focus-within/poster:scale-105 group-focus-within/poster:border-primary/50 group-focus-within/poster:shadow-xl"
				>
					{#if posterUrl(item.posterPath)}
						<img
							class="block h-full w-full object-cover"
							src={posterUrl(item.posterPath)}
							alt=""
							loading="lazy"
						/>
					{:else}
						<PosterPlaceholder label={item.type} />
					{/if}
					<MediaBadge type={item.type} />
					<div
						class="pointer-events-none absolute inset-0 z-[2] flex flex-col justify-end gap-1.5 bg-card/70 px-[13px] pt-[58px] pb-[13px] opacity-0 transition-opacity duration-150 group-hover/poster:opacity-100 group-focus-within/poster:opacity-100"
					>
						<span class="text-sm leading-none text-primary-foreground"
							>{item.year ?? 'Unknown'}</span
						>
						<h3 class="m-0 text-xl leading-tight text-primary-foreground">{item.title}</h3>
						<p class="line-clamp-4 m-0 text-[13px] leading-tight text-primary-foreground">
							{item.overview ?? 'No overview available.'}
						</p>
						<Tooltip.Root>
							<Tooltip.Trigger>
								{#snippet child({ props })}
									<Button
										{...props}
										type="button"
										variant="outline"
										size="icon-sm"
										class="pointer-events-auto absolute top-2.5 right-2.5 size-[34px] min-h-[34px] min-w-[34px] border-border bg-card/80 p-0 text-foreground backdrop-blur-md hover:border-primary/50 hover:bg-muted hover:text-primary-foreground focus-visible:border-primary/50 focus-visible:bg-muted focus-visible:text-primary-foreground"
										disabled={removingId === item.id}
										aria-label={`Remove ${item.title} from blacklist`}
										onclick={() => onRemove(item)}
									>
										<EyeIcon aria-hidden="true" />
									</Button>
								{/snippet}
							</Tooltip.Trigger>
							<Tooltip.Content>Remove from blacklist</Tooltip.Content>
						</Tooltip.Root>
					</div>
				</div>
			</article>
		{/each}
	</div>
{/if}
