<script lang="ts">
	import { resolve } from '$app/paths';
	import type { Snippet } from 'svelte';
	import MediaMonitorBookmark from './MediaMonitorBookmark.svelte';
	import type { MediaMetadataDetails } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		titleId: string;
		actions?: Snippet;
	}

	let { detail, titleId, actions }: Props = $props();

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

<div class="metadata-hero">
	<div class="metadata-poster">
		{#if imageUrl(detail.posterPath, 'w342')}
			<img src={imageUrl(detail.posterPath, 'w342')} alt="" />
		{:else}
			<div class="poster-placeholder">{detail.type}</div>
		{/if}
	</div>
	<div class="metadata-title-block">
		<h1 id={titleId} class="metadata-title">
			<MediaMonitorBookmark monitored={detail.monitored === true} />
			<span>{detail.title}</span>
		</h1>
		{#if hasSubtitle(detail)}
			<p class="metadata-subtitle">
				{#if certification}
					<span class="metadata-certification">{certification}</span>
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
			<div class="metadata-info-bar" aria-label="Media information">
				{#each topInfo(detail) as [label, value] (`${label}:${value}`)}
					<span><strong>{label}</strong>{value}</span>
				{/each}
			</div>
		{/if}
		{#if genres.length > 0}
			<div class="metadata-tags" aria-label="Genres">
				{#each genres as genre (genre)}
					<span><span class="app-icon" aria-hidden="true">sell</span>{genre}</span>
				{/each}
			</div>
		{/if}
		<div class="metadata-actions">
			{#if detail.collectionId && detail.collectionName}
				<a
					class="external-link"
					href={resolve('/media/collections/[provider]/[collectionId]', {
						provider: detail.externalProvider,
						collectionId: detail.collectionId
					})}
				>
					<span class="app-icon" aria-hidden="true">collections_bookmark</span>
					<span>{detail.collectionName}</span>
				</a>
			{/if}
			{@render actions?.()}
		</div>
	</div>
</div>
