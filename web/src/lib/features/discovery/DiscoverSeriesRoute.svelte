<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SeriesDiscoverFilterSheet from '$lib/components/app/discovery/series/SeriesDiscoverFilterSheet.svelte';
	import SeriesDiscoverToolbar from '$lib/components/app/discovery/series/SeriesDiscoverToolbar.svelte';
	import MovieDiscoverResults from '$lib/components/app/discovery/movies/MovieDiscoverResults.svelte';
	import { autocompleteDiscoverSeriesFacet, searchDiscoverSeries } from '$lib/settings/api';
	import type { DiscoverMovieFacetOption, MediaSearchResult } from '$lib/settings/types';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import {
		activeSeriesFilterCount,
		filtersFromParams,
		seriesFilterUrl,
		seriesQuery,
		type DiscoverSeriesFilters
	} from '$lib/components/app/discovery/series/discoverSeriesFilters';

	const app = getAppShellContext();
	const filters = $derived(filtersFromParams(page.url.searchParams));
	let results = $state<MediaSearchResult[]>([]);
	let genreOptions = $state<DiscoverMovieFacetOption[]>([]);
	let studioOptions = $state<DiscoverMovieFacetOption[]>([]);
	let keywordOptions = $state<DiscoverMovieFacetOption[]>([]);
	let loading = $state(false);
	let loadingMore = $state(false);
	let loadingGenres = $state(false);
	let loadingStudios = $state(false);
	let loadingKeywords = $state(false);
	let hasMore = $state(false);
	let hasSearched = $state(false);
	let pageNumber = $state(1);
	let loadedSearch = $state('');
	let filtersOpen = $state(false);

	const searchKey = $derived(page.url.search);
	const filterCount = $derived(activeSeriesFilterCount(filters));
	type SeriesFacet = 'genres' | 'studios' | 'keywords';

	$effect(() => {
		const key = searchKey;
		const nextFilters = filters;
		if (loadedSearch === key && hasSearched) return;
		loadedSearch = key;
		void searchFirstPage(nextFilters);
	});

	async function searchFirstPage(nextFilters: DiscoverSeriesFilters) {
		loading = true;
		hasSearched = true;
		pageNumber = 1;
		try {
			const response = await searchDiscoverSeries(seriesQuery(nextFilters, 1));
			results = response.results ?? [];
			hasMore = response.hasMore;
		} finally {
			loading = false;
		}
	}

	async function loadMore() {
		if (loading || loadingMore || !hasMore) return;
		loadingMore = true;
		try {
			const nextPage = pageNumber + 1;
			const response = await searchDiscoverSeries(seriesQuery(filters, nextPage));
			results = [...results, ...(response.results ?? [])];
			hasMore = response.hasMore;
			pageNumber = nextPage;
		} finally {
			loadingMore = false;
		}
	}

	function updateFilters(next: DiscoverSeriesFilters) {
		void goto(seriesFilterUrl(next), { replaceState: true, noScroll: true, keepFocus: true });
	}

	function updateSort(sort: string) {
		updateFilters({ ...filters, sort });
	}

	async function loadFacet(facet: SeriesFacet, query: string) {
		const cleaned = query.trim();
		if (cleaned.length < 2) {
			setFacetOptions(facet, []);
			setFacetLoading(facet, false);
			return;
		}
		setFacetLoading(facet, true);
		try {
			setFacetOptions(facet, await autocompleteDiscoverSeriesFacet(facet, cleaned));
		} finally {
			setFacetLoading(facet, false);
		}
	}

	function setFacetOptions(facet: SeriesFacet, options: DiscoverMovieFacetOption[]) {
		if (facet === 'genres') genreOptions = options;
		if (facet === 'studios') studioOptions = options;
		if (facet === 'keywords') keywordOptions = options;
	}

	function setFacetLoading(facet: SeriesFacet, loadingValue: boolean) {
		if (facet === 'genres') loadingGenres = loadingValue;
		if (facet === 'studios') loadingStudios = loadingValue;
		if (facet === 'keywords') loadingKeywords = loadingValue;
	}

	function trackLoadMore(node: HTMLDivElement) {
		const observer = new IntersectionObserver(
			(entries) => entries.some((entry) => entry.isIntersecting) && void loadMore(),
			{ rootMargin: '700px 0px' }
		);
		observer.observe(node);
		return { destroy: () => observer.disconnect() };
	}
</script>

<section class="grid min-w-0 gap-[18px]" aria-labelledby="discover-series-title">
	<PageHeading eyebrow="Discover" title="Series" titleId="discover-series-title" class="w-full">
		{#snippet actions()}
			<div class="mt-5">
				<SeriesDiscoverToolbar
					sort={filters.sort}
					{filterCount}
					{filtersOpen}
					onSortChange={updateSort}
					onToggleFilters={() => (filtersOpen = !filtersOpen)}
				/>
			</div>
		{/snippet}
	</PageHeading>

	<SeriesDiscoverFilterSheet
		open={filtersOpen}
		{filters}
		{genreOptions}
		{studioOptions}
		{keywordOptions}
		{loadingGenres}
		{loadingStudios}
		{loadingKeywords}
		onOpenChange={(open) => (filtersOpen = open)}
		onChange={updateFilters}
		onGenreQuery={(query) => void loadFacet('genres', query)}
		onStudioQuery={(query) => void loadFacet('studios', query)}
		onKeywordQuery={(query) => void loadFacet('keywords', query)}
	/>

	<div class="grid min-w-0 content-start gap-4">
		<MovieDiscoverResults
			{results}
			mediaItems={app.mediaItems}
			blacklist={app.discoverBlacklist}
			{loading}
			{loadingMore}
			{hasSearched}
			addingKey={app.addingKey}
			blacklistingKey={app.blacklistingKey}
			actionLabel={app.isAdmin ? 'Add' : 'Request'}
			canManage={app.isAdmin}
			onAdd={app.addMedia}
			onBlacklist={app.blacklistDiscoverMedia}
		/>
		{#if hasMore}
			<div class="min-h-8" use:trackLoadMore>
				<span class="sr-only">Loading more series when visible</span>
			</div>
		{/if}
	</div>
</section>
