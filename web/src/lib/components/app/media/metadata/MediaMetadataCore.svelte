<script lang="ts">
	import ArrowRightIcon from '@lucide/svelte/icons/arrow-right';
	import { resolve } from '$app/paths';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import MediaKeywordsSection from '$lib/components/app/media/metadata/MediaKeywordsSection.svelte';
	import MediaOverviewInfoCard from '$lib/components/app/media/detail/MediaOverviewInfoCard.svelte';
	import MediaEpisodeRow from '$lib/components/app/media/series/MediaEpisodeRow.svelte';
	import MediaPersonCard from '$lib/components/app/media/people/MediaPersonCard.svelte';
	import MediaSeasonPanel from '$lib/components/app/media/series/MediaSeasonPanel.svelte';
	import PosterRowControls from '$lib/components/app/media/posters/PosterRowControls.svelte';
	import { createPosterRowScroller } from '$lib/components/app/media/posters/posterRowScroller.svelte';
	import { crewRolePreviews } from '$lib/components/app/media/people/mediaPeople';
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
	const castRowKey = 'cast';
	const castScroller = createPosterRowScroller();

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
		return castScroller.trackRow(node, castRowKey);
	}

	function scrollCast(direction: -1 | 1) {
		castScroller.scrollRow(castRowKey, direction);
	}
</script>

<section
	class="grid items-start gap-5.5 min-[981px]:grid-cols-[minmax(0,1fr)_minmax(280px,390px)]"
	aria-labelledby="metadata-overview-title"
>
	<div class="grid min-w-0 gap-3">
		<h2 id="metadata-overview-title" class="m-0 text-3xl font-semibold text-foreground">
			Overview
		</h2>
		<p class="m-0 text-sm leading-6 text-muted-foreground">
			{detail.overview ?? 'No overview available.'}
		</p>
		{#if crewRoles.length > 0}
			<h3 class="mt-1 mb-0 text-xl text-foreground">
				{#if resolvedCrewHref}
					<!-- eslint-disable svelte/no-navigation-without-resolve -->
					<a
						class="inline-flex items-center gap-2 text-inherit no-underline hover:text-primary-hover focus-visible:text-primary-hover focus-visible:outline-none"
						href={resolvedCrewHref}
					>
						<span>Crew</span>
						<ArrowRightIcon aria-hidden="true" />
					</a>
					<!-- eslint-enable svelte/no-navigation-without-resolve -->
				{:else}
					Crew
				{/if}
			</h3>
			<div class="grid items-start gap-x-7 gap-y-[18px] md:grid-cols-3" aria-label="Crew">
				{#each crewRoles as role (role.role)}
					<div class="grid min-w-0 content-start gap-1">
						<strong class="[overflow-wrap:anywhere] text-foreground">{role.role}</strong>
						<span class="[overflow-wrap:anywhere] text-muted-foreground">
							{role.names.join(', ')}
						</span>
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
		<h2 id="metadata-seasons-title" class="m-0 text-3xl font-semibold text-foreground">Seasons</h2>
		<div class="grid gap-2.5">
			{#each seasons as season (season.name)}
				<MediaSeasonPanel
					summary={season.episodeCount ? `${season.episodeCount} episodes` : 'Episodes unknown'}
					size="-"
					tone="neutral"
				>
					{#snippet title()}
						<span>{season.name}</span>
					{/snippet}
					{#if season.episodes && season.episodes.length > 0}
						<div class="grid px-4.5">
							{#each season.episodes as episode (episode.episodeNumber)}
								<MediaEpisodeRow {episode} title={episodeTitle(episode)} />
							{/each}
						</div>
					{:else}
						<p class="p-4.5 text-sm text-muted-foreground">No episode details available.</p>
					{/if}
				</MediaSeasonPanel>
			{/each}
		</div>
	</section>
{/if}

{@render beforeCastContent?.()}

{#if cast.length > 0}
	{@const castEdges = castScroller.edgeState(castRowKey)}
	<section aria-labelledby="metadata-cast-title">
		<SectionHeading title="Cast" titleId="metadata-cast-title" href={resolvedCastHref}>
			{#if resolvedCastHref}
				<ArrowRightIcon aria-hidden="true" />
			{/if}
			{#snippet actions()}
				<span>{cast.length}</span>
				<PosterRowControls
					ariaLabel="Cast carousel controls"
					leftLabel="Scroll cast left"
					rightLabel="Scroll cast right"
					canScrollLeft={castEdges.canScrollLeft}
					canScrollRight={castEdges.canScrollRight}
					onScroll={scrollCast}
				/>
			{/snippet}
		</SectionHeading>
		<div
			class="-mx-3.5 mt-[-12px] grid max-w-full auto-cols-[minmax(190px,220px)] grid-flow-col gap-5 overflow-x-auto overflow-y-hidden overscroll-x-contain snap-x snap-proximity scroll-px-3.5 px-3.5 pt-[18px] pb-5 [scrollbar-width:none] max-[980px]:auto-cols-[minmax(160px,180px)] max-[980px]:gap-3.5 max-sm:mx-0 max-sm:auto-cols-[minmax(128px,150px)] max-sm:gap-3 max-sm:px-0 max-sm:pt-3.5 max-sm:pb-4 [&::-webkit-scrollbar]:hidden"
			use:trackCastRow
		>
			{#each cast as person (`${person.name}:${person.role ?? ''}`)}
				<MediaPersonCard name={person.name} role={person.role} image={person.profilePath} />
			{/each}
		</div>
	</section>
{/if}
