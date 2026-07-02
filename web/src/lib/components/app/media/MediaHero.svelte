<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TagIcon from '@lucide/svelte/icons/tag';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { MediaItem, MediaItemStatus, MediaType } from '$lib/settings/types';
	import PosterPlaceholder from './PosterPlaceholder.svelte';
	import { releaseSearchQuery } from './releaseSearchQuery';

	interface Props {
		mediaType: MediaType;
		item: MediaItem;
		qualityProfileLabel: string;
		canManage: boolean;
		searchingItemId?: string;
		scanningMediaItemId?: string;
		deletingMediaItemId?: string;
		onFindReleases: (_item: MediaItem, _query?: string) => void;
		onRescanMediaFiles: (_item: MediaItem) => void;
		onDeleteMedia: (_item: MediaItem) => void;
	}

	let {
		mediaType,
		item,
		qualityProfileLabel,
		canManage,
		searchingItemId,
		scanningMediaItemId,
		deletingMediaItemId,
		onFindReleases,
		onRescanMediaFiles,
		onDeleteMedia
	}: Props = $props();

	function posterUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function statusLabel(status: MediaItemStatus) {
		switch (status) {
			case 'downloaded':
				return 'Downloaded';
			case 'downloading':
				return 'Downloading';
			default:
				return 'Missing';
		}
	}
</script>

<section
	class="grid items-end gap-6.5 min-[641px]:grid-cols-[clamp(190px,18vw,270px)_minmax(0,1fr)]"
	aria-labelledby="home-title"
>
	<div
		class="aspect-[2/3] overflow-hidden rounded-md border border-border bg-card max-sm:w-[170px]"
	>
		{#if posterUrl(item.posterPath, 'w500')}
			<img class="block size-full object-cover" src={posterUrl(item.posterPath, 'w500')} alt="" />
		{:else}
			<PosterPlaceholder
				label={mediaType === 'movie' ? 'Movie' : 'Series'}
				class="h-full min-h-0"
			/>
		{/if}
	</div>
	<div class="grid gap-3">
		<h1 id="home-title" class="text-[clamp(42px,4.6vw,68px)] leading-none">{item.title}</h1>
		<p>{mediaType === 'movie' ? 'Movie' : 'Series'}</p>
		<div class="flex flex-wrap gap-2" aria-label="Library media information">
			<StatusPill class="inline-flex items-center gap-1.5">
				<strong class="text-foreground">Year</strong>{item.year ?? 'Unknown'}
			</StatusPill>
			<StatusPill class="inline-flex items-center gap-1.5">
				<strong class="text-foreground">Type</strong>{item.type}
			</StatusPill>
			<StatusPill class="inline-flex items-center gap-1.5">
				<strong class="text-foreground">Status</strong>{statusLabel(item.status)}
			</StatusPill>
			<StatusPill class="inline-flex items-center gap-1.5">
				<strong class="text-foreground">Profile</strong>{qualityProfileLabel}
			</StatusPill>
			<StatusPill class="inline-flex items-center gap-1.5">
				<strong class="text-foreground">Monitor</strong>{item.monitored ? 'Monitored' : 'None'}
			</StatusPill>
		</div>
		{#if item.tags?.length}
			<div class="flex flex-wrap gap-[7px]" aria-label="Tags">
				{#each item.tags as tag (tag)}
					<StatusPill class="inline-flex items-center gap-1.5">
						<TagIcon aria-hidden="true" />{tag}
					</StatusPill>
				{/each}
			</div>
		{/if}
		{#if canManage}
			<div class="mt-1 flex flex-wrap items-center gap-2.5">
				<Button
					type="button"
					disabled={searchingItemId === item.id}
					onclick={() => onFindReleases(item, releaseSearchQuery(item))}
				>
					{searchingItemId === item.id ? 'Queued' : 'Find releases'}
				</Button>
				<Button
					type="button"
					variant="outline"
					disabled={scanningMediaItemId === item.id || !item.mediaFolderPath}
					onclick={() => onRescanMediaFiles(item)}
				>
					<RefreshCwIcon aria-hidden="true" />
					<span>{scanningMediaItemId === item.id ? 'Scanning' : 'Rescan files'}</span>
				</Button>
				<Button
					type="button"
					variant="destructive"
					disabled={deletingMediaItemId === item.id}
					onclick={() => onDeleteMedia(item)}
				>
					{deletingMediaItemId === item.id ? 'Removing' : 'Remove'}
				</Button>
			</div>
		{/if}
	</div>
</section>
