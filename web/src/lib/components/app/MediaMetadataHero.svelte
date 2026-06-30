<script lang="ts">
	import { resolve } from '$app/paths';
	import type { Snippet } from 'svelte';
	import MediaMonitorBookmark from './MediaMonitorBookmark.svelte';
	import { formatDate } from '$lib/settings/dateFormat';
	import { providerDisplayName, providerPageUrl } from '$lib/settings/providerLinks';
	import type { MediaMetadataDetails } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		titleId: string;
		actions?: Snippet;
	}

	let { detail, titleId, actions }: Props = $props();

	const genres = $derived(detail.genres ?? []);

	function imageUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function titleWithYear(details: MediaMetadataDetails) {
		return `${details.title}${details.year ? ` (${details.year})` : ''}`;
	}

	function subtitle(details: MediaMetadataDetails) {
		const parts = [
			details.type === 'movie' ? 'Movie' : 'Series',
			details.year,
			details.voteAverage ? `${Math.round(details.voteAverage * 10)}%` : undefined
		].filter(Boolean);
		return parts.join(' | ');
	}

	function topInfo(details: MediaMetadataDetails) {
		return [
			details.type === 'movie' && details.releaseDate
				? ['Release', formatDate(details.releaseDate)]
				: undefined,
			details.runtimeMinutes ? ['Runtime', `${details.runtimeMinutes} min`] : undefined,
			details.seasonCount ? ['Seasons', `${details.seasonCount}`] : undefined,
			details.episodeCount ? ['Episodes', `${details.episodeCount}`] : undefined,
			details.originalLanguage ? ['Language', details.originalLanguage.toUpperCase()] : undefined
		].filter((item): item is [string, string] => Boolean(item));
	}

	function externalUrl(details: MediaMetadataDetails) {
		return providerPageUrl(details.externalProvider, details.type, details.externalId);
	}

	function externalLabel(details: MediaMetadataDetails) {
		return providerDisplayName(details.externalProvider);
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
			<span>{titleWithYear(detail)}</span>
		</h1>
		<p>{subtitle(detail)}</p>
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
			{#if externalUrl(detail)}
				<!-- eslint-disable svelte/no-navigation-without-resolve -->
				<a
					class="external-link"
					href={externalUrl(detail)}
					target="_blank"
					rel="noreferrer"
					aria-label={`Open ${externalLabel(detail)} page in a new tab`}
				>
					<span class="app-icon" aria-hidden="true">open_in_new</span>
					<span>{externalLabel(detail)}</span>
				</a>
				<!-- eslint-enable svelte/no-navigation-without-resolve -->
			{/if}
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
