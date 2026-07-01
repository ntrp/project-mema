<script lang="ts">
	import CompassIcon from '@lucide/svelte/icons/compass';
	import { resolve } from '$app/paths';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { MediaItem, MediaType } from '$lib/settings/types';
	import PosterPlaceholder from '../media/PosterPlaceholder.svelte';

	interface Props {
		mediaType: MediaType;
		items: MediaItem[];
	}

	let { mediaType, items }: Props = $props();

	const heading = $derived(mediaType === 'movie' ? 'Added movies' : 'Added series');
	const label = $derived(mediaType === 'movie' ? 'Movies' : 'Series');
	const emptyText = $derived(
		mediaType === 'movie' ? 'No movies added yet.' : 'No series added yet.'
	);

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w500${path}`;
	}

	function statusLabel(status: MediaItem['status']) {
		switch (status) {
			case 'downloaded':
				return 'Downloaded';
			case 'downloading':
				return 'Downloading';
			default:
				return 'Missing';
		}
	}

	function typeLabel(type: MediaType) {
		return type === 'movie' ? 'Movie' : 'Series';
	}

	function statusLineClass(status: MediaItem['status']) {
		if (status === 'downloaded') {
			return 'bg-primary';
		}
		if (status === 'downloading') {
			return 'bg-secondary';
		}
		return 'bg-destructive';
	}
</script>

<PageHeading eyebrow={label} title={heading} titleId="home-title" />

{#if items.length === 0}
	<EmptyState
		class="my-[18px] grid min-h-60 w-full place-items-center content-center gap-[18px] text-center"
	>
		<p class="m-0 text-lg font-black text-foreground">{emptyText}</p>
		<Button href={resolve('/discover')}>
			<CompassIcon aria-hidden="true" />
			<span>Discover</span>
		</Button>
	</EmptyState>
{:else}
	<div
		class="grid grid-cols-[repeat(auto-fill,minmax(132px,1fr))] gap-3 sm:grid-cols-[repeat(auto-fill,minmax(190px,220px))] sm:gap-5"
	>
		{#each items as item (item.id)}
			<a
				class="group/library relative block isolate overflow-hidden rounded-md border border-border bg-card text-foreground no-underline transition-[transform,border-color,box-shadow] duration-150 hover:z-[2] hover:-translate-y-1.5 hover:scale-[1.04] hover:border-primary/50 hover:shadow-xl focus-visible:z-[2] focus-visible:-translate-y-1.5 focus-visible:scale-[1.04] focus-visible:border-primary/50 focus-visible:shadow-xl focus-visible:outline-none max-sm:hover:-translate-y-0.5 max-sm:hover:scale-[1.02] max-sm:focus-visible:-translate-y-0.5 max-sm:focus-visible:scale-[1.02]"
				href={mediaType === 'movie'
					? resolve('/movies/[id]', { id: item.id })
					: resolve('/series/[id]', { id: item.id })}
				aria-label={`Open ${item.title} details`}
			>
				<div class="aspect-[2/3] overflow-hidden bg-card">
					{#if posterUrl(item.posterPath)}
						<img
							class="block h-full w-full object-cover"
							src={posterUrl(item.posterPath)}
							alt=""
							loading="lazy"
						/>
					{:else}
						<PosterPlaceholder label={typeLabel(mediaType)} class="h-full min-h-0" />
					{/if}
				</div>
				<div
					class="pointer-events-none absolute inset-0 z-[2] flex flex-col justify-end gap-1.5 bg-card/70 px-[13px] pt-[58px] pb-[13px] opacity-0 transition-opacity duration-150 group-hover/library:opacity-100 group-focus-visible/library:opacity-100"
				>
					<span class="text-sm leading-none text-primary-foreground">{item.year ?? 'Unknown'}</span>
					<h3 class="m-0 text-xl leading-tight text-primary-foreground">{item.title}</h3>
					<p class="line-clamp-4 m-0 text-[13px] leading-tight text-primary-foreground">
						{item.overview ?? 'No overview available.'}
					</p>
				</div>
				<span
					class={`absolute right-0 bottom-0 left-0 z-[4] h-[5px] ${statusLineClass(item.status)}`}
					aria-label={statusLabel(item.status)}
				></span>
			</a>
		{/each}
	</div>
{/if}
