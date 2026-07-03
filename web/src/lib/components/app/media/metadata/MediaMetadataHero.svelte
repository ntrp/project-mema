<script lang="ts">
	import LibraryIcon from '@lucide/svelte/icons/library';
	import PlayIcon from '@lucide/svelte/icons/play';
	import TagIcon from '@lucide/svelte/icons/tag';
	import { resolve } from '$app/paths';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import type { Snippet } from 'svelte';
	import MediaMonitorBookmark from '$lib/components/app/media/detail/MediaMonitorBookmark.svelte';
	import MediaTrailerModal from '$lib/components/app/media/detail/MediaTrailerModal.svelte';
	import PosterPlaceholder from '$lib/components/app/media/posters/PosterPlaceholder.svelte';
	import {
		imageUrl,
		mediaHeroTopInfo,
		monitorHint,
		monitorStatus,
		runtimeText,
		statusBadgeClass,
		statusLabel
	} from '$lib/components/app/media/detail/mediaHeroDisplay';
	import type { MediaItemStatus, MediaMetadataDetails } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		titleId: string;
		showMonitorBookmark?: boolean;
		showTrailerButton?: boolean;
		mediaStatus?: MediaItemStatus;
		monitorMonitored?: boolean;
		monitorStatusText?: string;
		monitorHintText?: string;
		monitorDisabled?: boolean;
		onToggleMonitor?: () => void;
		actions?: Snippet;
	}

	let {
		detail,
		titleId,
		showMonitorBookmark = true,
		showTrailerButton = true,
		mediaStatus,
		monitorMonitored = detail.monitored === true,
		monitorStatusText = monitorStatus(detail),
		monitorHintText = monitorHint(detail),
		monitorDisabled = false,
		onToggleMonitor = () => {},
		actions
	}: Props = $props();
	let trailerOpen = $state(false);

	const genres = $derived(detail.genres ?? []);
	const factMap = $derived(new Map((detail.facts ?? []).map((fact) => [fact.label, fact.value])));
	const certification = $derived(certificationText());
	const duration = $derived(runtimeText(detail.runtimeMinutes));
	const trailerTitle = $derived(`${detail.title} trailer`);

	function hasSubtitle(details: MediaMetadataDetails) {
		return Boolean(mediaStatus || certification || details.year || duration);
	}

	function certificationText() {
		return factMap.get('Certification');
	}
</script>

<div
	class="grid items-end gap-6.5 min-[641px]:grid-cols-[clamp(190px,18vw,270px)_minmax(0,1fr)] mb-6"
>
	<div class="aspect-2/3 overflow-hidden rounded-md border border-border bg-card max-sm:w-42.5">
		{#if imageUrl(detail.posterPath, 'w342')}
			<img class="block size-full object-cover" src={imageUrl(detail.posterPath, 'w342')} alt="" />
		{:else}
			<PosterPlaceholder label={detail.type} class="h-full min-h-0" />
		{/if}
	</div>
	<div class="grid gap-3">
		<h1 id={titleId} class="flex items-center gap-3 text-[clamp(42px,4.6vw,68px)] leading-none">
			{#if showMonitorBookmark}
				<MediaMonitorBookmark
					monitored={monitorMonitored}
					status={monitorStatusText}
					hint={monitorHintText}
					disabled={monitorDisabled}
					size={64}
					onToggle={onToggleMonitor}
				/>
			{/if}
			<span>{detail.title}</span>
		</h1>
		{#if hasSubtitle(detail)}
			<p class="m-0 flex flex-wrap items-center gap-3.5 text-muted-foreground">
				{#if certification}
					<span
						class="inline-flex min-h-6 items-center rounded-md border border-white p-1 text-base leading-none font-extrabold text-foreground"
						>{certification}</span
					>
				{/if}
				{#if detail.year}
					<span>{detail.year}</span>
				{/if}
				{#if duration}
					<span>{duration}</span>
				{/if}
				{#if mediaStatus}
					<Badge
						variant="outline"
						class={`min-h-6 rounded-[3px] bg-popover px-2 py-1 text-base leading-none font-extrabold ${statusBadgeClass(mediaStatus)}`}
					>
						{statusLabel(mediaStatus)}
					</Badge>
				{/if}
			</p>
		{/if}
		{#if mediaHeroTopInfo(detail).length > 0}
			<div class="flex flex-wrap gap-2" aria-label="Media information">
				{#each mediaHeroTopInfo(detail) as [label, value] (`${label}:${value}`)}
					<StatusPill class="inline-flex items-center gap-1.5">
						<strong class="text-foreground">{label}</strong>{value}
					</StatusPill>
				{/each}
			</div>
		{/if}
		{#if genres.length > 0}
			<div class="flex flex-wrap gap-1.75" aria-label="Genres">
				{#each genres as genre (genre)}
					<StatusPill class="inline-flex items-center gap-1.5">
						<TagIcon aria-hidden="true" />{genre}
					</StatusPill>
				{/each}
			</div>
		{/if}
		<div class="mt-1 flex flex-wrap items-end gap-2.5">
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
			{#if showTrailerButton && detail.trailerUrl}
				<Button variant="outline" size="sm" onclick={() => (trailerOpen = true)}>
					<PlayIcon aria-hidden="true" />
					<span>Trailer</span>
				</Button>
			{/if}
			{@render actions?.()}
		</div>
	</div>
</div>

{#if showTrailerButton && trailerOpen && detail.trailerUrl}
	<MediaTrailerModal
		title={trailerTitle}
		url={detail.trailerUrl}
		onClose={() => (trailerOpen = false)}
	/>
{/if}
