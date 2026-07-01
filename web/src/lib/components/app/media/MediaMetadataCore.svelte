<script lang="ts">
	/* global HTMLDivElement */
	import { resolve } from '$app/paths';
	import MediaKeywordsSection from './MediaKeywordsSection.svelte';
	import MediaOverviewInfoCard from './MediaOverviewInfoCard.svelte';
	import { crewRolePreviews } from './mediaPeople';
	import { formatDate } from '$lib/settings/dateFormat';
	import type { Snippet } from 'svelte';
	import type { MediaMetadataDetails } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		castHref?: string;
		crewHref?: string;
		seasonsContent?: Snippet;
		beforeCastContent?: Snippet;
	}

	let { detail, castHref, crewHref, seasonsContent, beforeCastContent }: Props = $props();

	const facts = $derived(detail.facts ?? []);
	const keywords = $derived(detail.keywords ?? []);
	const crewRoles = $derived(crewRolePreviews(facts));
	const seasons = $derived(detail.seasons ?? []);
	const cast = $derived(detail.cast ?? []);
	const resolvedCastHref = $derived(castHref ?? resolvedPeopleHref(detail, 'cast'));
	const resolvedCrewHref = $derived(crewHref ?? resolvedPeopleHref(detail, 'crew'));
	let castRow = $state<HTMLDivElement | undefined>();

	function imageUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function episodeTitle(episode: { episodeNumber: number; name: string }) {
		return `${episode.episodeNumber} - ${episode.name}`;
	}

	function resolvedPeopleHref(details: MediaMetadataDetails, kind: 'cast' | 'crew') {
		if (!details.externalProvider || !details.externalId) {
			return undefined;
		}
		const params = {
			provider: details.externalProvider,
			type: details.type,
			externalId: details.externalId
		};
		return kind === 'cast'
			? resolve('/media/[provider]/[type]/[externalId]/cast', params)
			: resolve('/media/[provider]/[type]/[externalId]/crew', params);
	}

	function trackCastRow(node: HTMLDivElement) {
		castRow = node;
		return {
			destroy() {
				castRow = undefined;
			}
		};
	}

	function scrollCast(direction: number) {
		if (!castRow) {
			return;
		}
		castRow.scrollBy({
			left: direction * Math.max(castRow.clientWidth - 140, 220),
			behavior: 'smooth'
		});
	}
</script>

<section class="metadata-overview-layout" aria-labelledby="metadata-overview-title">
	<div class="metadata-overview-copy">
		<h2 id="metadata-overview-title">Overview</h2>
		<p>{detail.overview ?? 'No overview available.'}</p>
		{#if crewRoles.length > 0}
			<h3 class="metadata-crew-title">
				{#if resolvedCrewHref}
					<!-- eslint-disable svelte/no-navigation-without-resolve -->
					<a class="metadata-crew-title-link" href={resolvedCrewHref}>
						<span>Crew</span>
						<span class="app-icon" aria-hidden="true">arrow_forward</span>
					</a>
					<!-- eslint-enable svelte/no-navigation-without-resolve -->
				{:else}
					Crew
				{/if}
			</h3>
			<div class="metadata-crew-grid" aria-label="Crew">
				{#each crewRoles as role (role.role)}
					<div>
						<strong>{role.role}</strong>
						<span>{role.names.join(', ')}</span>
					</div>
				{/each}
			</div>
		{/if}
		<MediaKeywordsSection {keywords} />
	</div>
	<MediaOverviewInfoCard {detail} {facts} />
</section>

{#if seasonsContent}
	{@render seasonsContent()}
{:else if seasons.length > 0}
	<section aria-labelledby="metadata-seasons-title">
		<h2 id="metadata-seasons-title">Seasons</h2>
		<div class="metadata-season-list">
			{#each seasons as season (season.name)}
				<details class="metadata-season-panel">
					<summary>
						<span>{season.name}</span>
						<small
							>{season.episodeCount ? `${season.episodeCount} episodes` : 'Episodes unknown'}</small
						>
					</summary>
					{#if season.episodes && season.episodes.length > 0}
						<div class="metadata-episode-list">
							{#each season.episodes as episode (episode.episodeNumber)}
								<article class="metadata-episode-row">
									<div class="metadata-episode-copy">
										<h3>
											{episodeTitle(episode)}
											{#if episode.airDate}
												<span>{formatDate(episode.airDate)}</span>
											{/if}
										</h3>
										<p>{episode.overview ?? 'No episode overview available.'}</p>
									</div>
									{#if imageUrl(episode.stillPath, 'w300')}
										<img src={imageUrl(episode.stillPath, 'w300')} alt="" loading="lazy" />
									{/if}
								</article>
							{/each}
						</div>
					{:else}
						<p class="metadata-season-empty">No episode details available.</p>
					{/if}
				</details>
			{/each}
		</div>
	</section>
{/if}

{@render beforeCastContent?.()}

{#if cast.length > 0}
	<section aria-labelledby="metadata-cast-title">
		<div class="section-heading">
			{#if resolvedCastHref}
				<!-- eslint-disable svelte/no-navigation-without-resolve -->
				<a class="section-title-link" href={resolvedCastHref}>
					<h2 id="metadata-cast-title">Cast</h2>
					<span class="app-icon" aria-hidden="true">arrow_forward</span>
				</a>
				<!-- eslint-enable svelte/no-navigation-without-resolve -->
			{:else}
				<h2 id="metadata-cast-title">Cast</h2>
			{/if}
			<div class="section-heading-actions">
				<span>{cast.length}</span>
				<div class="poster-row-controls" aria-label="Cast carousel controls">
					<button type="button" aria-label="Scroll cast left" onclick={() => scrollCast(-1)}>
						<span class="app-icon" aria-hidden="true">chevron_left</span>
					</button>
					<button type="button" aria-label="Scroll cast right" onclick={() => scrollCast(1)}>
						<span class="app-icon" aria-hidden="true">chevron_right</span>
					</button>
				</div>
			</div>
		</div>
		<div class="metadata-cast-slider" use:trackCastRow>
			{#each cast as person (`${person.name}:${person.role ?? ''}`)}
				<article class="metadata-cast-card">
					<div>
						{#if imageUrl(person.profilePath, 'w185')}
							<img src={imageUrl(person.profilePath, 'w185')} alt="" loading="lazy" />
						{:else}
							<span>{person.name.slice(0, 1)}</span>
						{/if}
					</div>
					<strong>{person.name}</strong>
					{#if person.role}
						<p>{person.role}</p>
					{/if}
				</article>
			{/each}
		</div>
	</section>
{/if}
