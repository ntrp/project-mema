<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { getMediaMetadataDetails, searchMedia } from '$lib/settings/api';
	import type { MediaItem, MediaMetadataDetails, MediaSearchResult } from '$lib/settings/types';
	import { cn } from '$lib/utils';
	import type { ReleaseOverrideDraft } from '$lib/components/app/media/release-override/releaseOverrideDetails';
	import {
		episodeNumbers,
		seasonOptions,
		selectedSeason
	} from '$lib/components/app/media/release-override/releaseOverrideSeriesOptions';
	import ReleaseOverrideEpisodeSelect from '$lib/components/app/media/release-override/ReleaseOverrideEpisodeSelect.svelte';
	import ReleaseOverrideSeasonSelect from '$lib/components/app/media/release-override/ReleaseOverrideSeasonSelect.svelte';

	interface Props {
		item: MediaItem;
		draft: ReleaseOverrideDraft;
	}

	let { item, draft }: Props = $props();
	let results = $state<MediaSearchResult[]>([]);
	let details = $state<MediaMetadataDetails | undefined>();
	let loadingSearch = $state(false);
	let loadingDetails = $state(false);
	let open = $state(false);
	let selectedIndex = $state(-1);
	let requestNumber = 0;

	const trimmed = $derived(draft.seriesTitle.trim());
	const selectedResult = $derived(selectedIndex >= 0 ? results[selectedIndex] : undefined);
	const seasons = $derived(seasonOptions(details));
	const season = $derived(selectedSeason(seasons, draft.seasonNumber));
	const episodes = $derived(season?.season.episodes ?? []);

	onMount(() => {
		if (item.externalProvider === 'tmdb' && item.externalId) {
			void loadSeriesDetails(item.externalId, item.title);
		}
	});

	function handleInput() {
		open = true;
		selectedIndex = -1;
		void searchSeries(trimmed);
	}

	async function searchSeries(query: string) {
		const current = ++requestNumber;
		if (query.length < 2) {
			results = [];
			loadingSearch = false;
			return;
		}
		loadingSearch = true;
		try {
			const found = await searchMedia({ query, type: 'serie' });
			if (current === requestNumber) {
				results = found.filter((result) => result.externalProvider === 'tmdb').slice(0, 6);
			}
		} finally {
			if (current === requestNumber) loadingSearch = false;
		}
	}

	async function choose(result: MediaSearchResult) {
		setSeriesTitle(result.title);
		open = false;
		selectedIndex = -1;
		if (result.externalProvider !== 'tmdb' || !result.externalId) return;
		await loadSeriesDetails(result.externalId, result.title);
	}

	async function loadSeriesDetails(externalId: string, title: string) {
		loadingDetails = true;
		try {
			details = await getMediaMetadataDetails('tmdb', 'serie', externalId);
			setSeriesTitle(title);
			selectDefaultSeason();
		} finally {
			loadingDetails = false;
		}
	}

	function selectDefaultSeason() {
		const availableSeasons = seasonOptions(details);
		const match =
			availableSeasons.find((option) => option.value === draft.seasonNumber) ?? availableSeasons[0];
		if (!match) return;
		setSeason(match.value);
	}

	function setSeriesTitle(value: string) {
		draft.seriesTitle = value;
	}

	function setSeason(value: string) {
		draft.seasonNumber = value;
		const match = selectedSeason(seasons, value);
		const valid = new Set((match?.season.episodes ?? []).map((episode) => episode.episodeNumber));
		const kept = episodeNumbers(draft.episodeNumbers).filter((number) => valid.has(number));
		draft.episodeNumbers = kept.join(', ');
	}

	function setEpisodes(value: string) {
		draft.episodeNumbers = value;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
			if (results.length === 0) return;
			event.preventDefault();
			open = true;
			selectedIndex =
				event.key === 'ArrowDown'
					? Math.min(selectedIndex + 1, results.length - 1)
					: Math.max(selectedIndex - 1, 0);
		}
		if (event.key === 'Enter' && selectedResult) {
			event.preventDefault();
			void choose(selectedResult);
		}
		if (event.key === 'Escape') open = false;
	}

	function resultKey(result: MediaSearchResult) {
		return `${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}`;
	}
</script>

<div class="grid gap-3 md:grid-cols-2">
	<div class="relative grid gap-1.5">
		<Label for="override-series">Series</Label>
		<Input
			id="override-series"
			value={draft.seriesTitle}
			autocomplete="off"
			oninput={(event) => {
				setSeriesTitle(event.currentTarget.value);
				handleInput();
			}}
			onfocus={() => {
				open = true;
				if (trimmed.length >= 2 && results.length === 0) void searchSeries(trimmed);
			}}
			onkeydown={handleKeydown}
			onblur={() => {
				window.setTimeout(() => (open = false), 120);
			}}
		/>
		{#if open && trimmed.length >= 2}
			<div
				class="absolute inset-x-0 top-[calc(100%+6px)] z-30 max-h-72 overflow-auto rounded-md border border-border bg-popover p-1.5 text-popover-foreground shadow-md"
				role="listbox"
				aria-label="Series matches"
			>
				{#if results.length > 0}
					{#each results as result, index (resultKey(result))}
						<Button
							type="button"
							variant="ghost"
							role="option"
							aria-selected={index === selectedIndex}
							class={cn(
								'grid min-h-10 w-full justify-start rounded-sm px-2 py-1.5 text-left',
								index === selectedIndex && 'bg-accent text-accent-foreground'
							)}
							onpointerdown={(event) => event.preventDefault()}
							onclick={() => choose(result)}
						>
							<span class="truncate">{result.title}</span>
							{#if result.year}<small class="text-xs text-muted-foreground">{result.year}</small
								>{/if}
						</Button>
					{/each}
				{:else}
					<div class="px-2 py-1 text-xs font-bold text-muted-foreground uppercase">
						{loadingSearch ? 'Searching' : 'No TMDB matches'}
					</div>
				{/if}
			</div>
		{/if}
	</div>
	<ReleaseOverrideSeasonSelect
		value={draft.seasonNumber}
		label={season?.label ?? ''}
		{seasons}
		onChange={setSeason}
	/>
</div>
{#if loadingDetails}
	<span class="text-xs text-muted-foreground">Loading series metadata</span>
{/if}
<ReleaseOverrideEpisodeSelect value={draft.episodeNumbers} {episodes} onChange={setEpisodes} />
