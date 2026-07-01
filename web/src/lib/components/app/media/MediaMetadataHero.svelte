<script lang="ts">
	import LibraryIcon from '@lucide/svelte/icons/library';
	import TagIcon from '@lucide/svelte/icons/tag';
	import { resolve } from '$app/paths';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { Snippet } from 'svelte';
	import MediaMonitorBookmark from './MediaMonitorBookmark.svelte';
	import PosterPlaceholder from './PosterPlaceholder.svelte';
	import type { MediaMetadataDetails } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		titleId: string;
		showMonitorBookmark?: boolean;
		actions?: Snippet;
	}

	let { detail, titleId, showMonitorBookmark = true, actions }: Props = $props();

	const genres = $derived(detail.genres ?? []);
	const factMap = $derived(new Map((detail.facts ?? []).map((fact) => [fact.label, fact.value])));
	const certification = $derived(certificationText());
	const duration = $derived(runtimeText(detail.runtimeMinutes));

	function imageUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function hasSubtitle(details: MediaMetadataDetails) {
		return Boolean(certification || details.year || duration);
	}

	function topInfo(details: MediaMetadataDetails) {
		return [
			details.seasonCount ? ['Seasons', `${details.seasonCount}`] : undefined,
			details.episodeCount ? ['Episodes', `${details.episodeCount}`] : undefined
		].filter((item): item is [string, string] => Boolean(item));
	}

	function certificationText() {
		return factMap.get('Certification');
	}

	function runtimeText(minutes?: number) {
		if (!minutes || minutes <= 0) {
			return undefined;
		}
		const hours = Math.floor(minutes / 60);
		const remainingMinutes = minutes % 60;
		if (hours > 0 && remainingMinutes > 0) {
			return `${hours}h ${remainingMinutes}m`;
		}
		if (hours > 0) {
			return `${hours}h`;
		}
		return `${remainingMinutes}m`;
	}
</script>

<div class="grid items-end gap-6.5 min-[641px]:grid-cols-[clamp(190px,18vw,270px)_minmax(0,1fr)]">
	<div
		class="aspect-[2/3] overflow-hidden rounded-md border border-border bg-card max-sm:w-[170px]"
	>
		{#if imageUrl(detail.posterPath, 'w342')}
			<img class="block size-full object-cover" src={imageUrl(detail.posterPath, 'w342')} alt="" />
		{:else}
			<PosterPlaceholder label={detail.type} class="h-full min-h-0" />
		{/if}
	</div>
	<div class="grid gap-3">
		<h1 id={titleId} class="flex items-center gap-3 text-[clamp(42px,4.6vw,68px)] leading-none">
			{#if showMonitorBookmark}
				<MediaMonitorBookmark monitored={detail.monitored === true} />
			{/if}
			<span>{detail.title}</span>
		</h1>
		{#if hasSubtitle(detail)}
			<p class="m-0 flex flex-wrap items-center gap-3.5 text-muted-foreground">
				{#if certification}
					<span
						class="inline-flex min-h-6 items-center rounded-md border-2 border-primary/50 px-2.5 py-0.5 text-base leading-none font-extrabold text-foreground"
						>{certification}</span
					>
				{/if}
				{#if detail.year}
					<span>{detail.year}</span>
				{/if}
				{#if duration}
					<span>{duration}</span>
				{/if}
			</p>
		{/if}
		{#if topInfo(detail).length > 0}
			<div class="flex flex-wrap gap-2" aria-label="Media information">
				{#each topInfo(detail) as [label, value] (`${label}:${value}`)}
					<StatusPill class="inline-flex items-center gap-1.5">
						<strong class="text-foreground">{label}</strong>{value}
					</StatusPill>
				{/each}
			</div>
		{/if}
		{#if genres.length > 0}
			<div class="flex flex-wrap gap-[7px]" aria-label="Genres">
				{#each genres as genre (genre)}
					<StatusPill class="inline-flex items-center gap-1.5">
						<TagIcon aria-hidden="true" />{genre}
					</StatusPill>
				{/each}
			</div>
		{/if}
		<div class="mt-1 flex flex-wrap items-center gap-2.5">
			{#if detail.collectionId && detail.collectionName}
				<Button
					variant="outline"
					size="sm"
					href={resolve('/media/collections/[provider]/[collectionId]', {
						provider: detail.externalProvider,
						collectionId: detail.collectionId
					})}
				>
					<LibraryIcon aria-hidden="true" />
					<span>{detail.collectionName}</span>
				</Button>
			{/if}
			{@render actions?.()}
		</div>
	</div>
</div>
