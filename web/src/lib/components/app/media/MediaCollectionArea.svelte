<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { resolve } from '$app/paths';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { MediaCollection, MediaItem, MediaSearchResult } from '$lib/settings/types';
	import MediaBadge from './MediaBadge.svelte';
	import PosterPlaceholder from './PosterPlaceholder.svelte';

	interface Props {
		collection?: MediaCollection;
		mediaItems: MediaItem[];
		loading: boolean;
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { collection, mediaItems, loading, addingKey, actionLabel, onAdd }: Props = $props();

	const libraryExternalKeys = $derived(
		new Set(
			(mediaItems ?? [])
				.map((item) => externalKey(item))
				.filter((key): key is string => Boolean(key))
		)
	);
	const libraryTitleKeys = $derived(new Set((mediaItems ?? []).map((item) => titleKey(item))));

	function resultKey(result: MediaSearchResult) {
		return `${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}

	function isInLibrary(result: MediaSearchResult) {
		const key = externalKey(result);
		return Boolean(key && libraryExternalKeys.has(key)) || libraryTitleKeys.has(titleKey(result));
	}

	function externalKey(item: MediaItem | MediaSearchResult) {
		if (!item.externalProvider || !item.externalId) {
			return undefined;
		}
		return `${item.type}:${clean(item.externalProvider)}:${clean(item.externalId)}`;
	}

	function titleKey(item: MediaItem | MediaSearchResult) {
		return `${item.type}:${clean(item.title)}:${item.year ?? ''}`;
	}

	function clean(value: string) {
		return value.trim().toLowerCase();
	}

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

{#if loading}
	<EmptyState
		title="Loading collection"
		description="Fetching collection media from the metadata provider."
	/>
{:else if !collection}
	<EmptyState title="Collection not available" description="Could not load this collection." />
{:else}
	<PageHeading
		eyebrow={collection.provider}
		title={collection.name}
		description={collection.overview}
	/>

	<section class="min-w-0" aria-labelledby="collection-results-title">
		<SectionHeading title="Collection media" titleId="collection-results-title">
			{#snippet actions()}
				<span>{collection.results.length} titles</span>
			{/snippet}
		</SectionHeading>
		<div
			class="-mx-3.5 mt-[-12px] grid auto-cols-[minmax(190px,220px)] grid-flow-col gap-5 overflow-x-auto overflow-y-hidden overscroll-x-contain snap-x snap-proximity scroll-px-3.5 px-3.5 pt-[18px] pb-5 [scrollbar-width:none] max-sm:mx-0 max-sm:auto-cols-[minmax(128px,150px)] max-sm:gap-3 max-sm:px-0 max-sm:pt-3.5 max-sm:pb-4 [&::-webkit-scrollbar]:hidden"
		>
			{#each collection.results as result (resultKey(result))}
				<article class="group/poster min-w-0 snap-start">
					<div
						class="relative aspect-[2/3] overflow-hidden rounded-md border border-border bg-card transition-[transform,border-color,box-shadow] duration-150 group-hover/poster:z-[2] group-hover/poster:-translate-y-1.5 group-hover/poster:scale-105 group-hover/poster:border-primary/50 group-hover/poster:shadow-xl group-focus-within/poster:z-[2] group-focus-within/poster:-translate-y-1.5 group-focus-within/poster:scale-105 group-focus-within/poster:border-primary/50 group-focus-within/poster:shadow-xl"
					>
						{#if posterUrl(result.posterPath)}
							<img
								class="block h-full w-full object-cover"
								src={posterUrl(result.posterPath)}
								alt=""
								loading="lazy"
							/>
						{:else}
							<PosterPlaceholder label={result.type} />
						{/if}
						{#if result.externalProvider && result.externalId}
							<a
								class="absolute inset-0 z-[1] rounded-md"
								href={resolve('/media/[provider]/[type]/[externalId]', {
									provider: result.externalProvider,
									type: result.type,
									externalId: result.externalId
								})}
								aria-label={`Open ${result.title} details`}
							></a>
						{/if}
						<MediaBadge type={result.type} />
						{#if isInLibrary(result)}
							<StatusPill
								class="absolute top-2 right-2 z-[3] bg-primary text-primary-foreground"
								tone="success">In library</StatusPill
							>
						{/if}
						<div
							class="pointer-events-none absolute inset-0 z-[2] flex flex-col justify-end gap-1.5 bg-card/70 px-[13px] pt-[58px] pb-[13px] opacity-0 transition-opacity duration-150 group-hover/poster:opacity-100 group-focus-within/poster:opacity-100"
						>
							<span class="text-sm leading-none text-primary-foreground"
								>{result.year ?? 'Unknown'}</span
							>
							<h3 class="m-0 text-xl leading-tight text-primary-foreground">{result.title}</h3>
							<p class="line-clamp-4 m-0 text-[13px] leading-tight text-primary-foreground">
								{result.overview ?? 'No overview available.'}
							</p>
							{#if isInLibrary(result)}
								<StatusPill class="bg-primary text-primary-foreground" tone="success"
									>In library</StatusPill
								>
							{:else}
								<Button
									type="button"
									class="pointer-events-auto mt-0.5 min-h-[30px] self-start px-3 text-[13px]"
									disabled={addingKey === resultKey(result)}
									onclick={() => onAdd(result)}
								>
									<PlusIcon aria-hidden="true" />
									<span>{addingKey === resultKey(result) ? 'Working' : actionLabel}</span>
								</Button>
							{/if}
						</div>
					</div>
				</article>
			{/each}
		</div>
	</section>
{/if}
