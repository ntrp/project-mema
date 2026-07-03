<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import MovieDiscoverFilterSheet from '$lib/components/app/discovery/movies/MovieDiscoverFilterSheet.svelte';
	import MovieDiscoverResults from '$lib/components/app/discovery/movies/MovieDiscoverResults.svelte';
	import MovieDiscoverToolbar from '$lib/components/app/discovery/movies/MovieDiscoverToolbar.svelte';
	import { autocompleteDiscoverMovieFacet, searchDiscoverMovies } from '$lib/settings/api';
	import type { DiscoverMovieFacetOption, MediaSearchResult } from '$lib/settings/types';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import {
		activeMovieFilterCount,
		filtersFromParams,
		movieFilterUrl,
		movieQuery,
		type DiscoverMovieFilters
	} from '$lib/components/app/discovery/movies/discoverMovieFilters';

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
	const filterCount = $derived(activeMovieFilterCount(filters));
	type MovieFacet = 'genres' | 'studios' | 'keywords';

	$effect(() => {
		const key = searchKey;
		const nextFilters = filters;
		if (loadedSearch === key && hasSearched) return;
		loadedSearch = key;
		void searchFirstPage(nextFilters);
	});

	async function searchFirstPage(nextFilters: DiscoverMovieFilters) {
		loading = true;
		hasSearched = true;
		pageNumber = 1;
		try {
			const response = await searchDiscoverMovies(movieQuery(nextFilters, 1));
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
			const response = await searchDiscoverMovies(movieQuery(filters, nextPage));
			results = [...results, ...(response.results ?? [])];
			hasMore = response.hasMore;
			pageNumber = nextPage;
		} finally {
			loadingMore = false;
		}
	}

	function updateFilters(next: DiscoverMovieFilters) {
		void goto(movieFilterUrl(next), {
			replaceState: true,
			noScroll: true,
			keepFocus: true
		});
	}

	function updateSort(sort: string) {
		updateFilters({ ...filters, sort });
	}

	async function loadFacet(facet: MovieFacet, query: string) {
		const cleaned = query.trim();
		if (cleaned.length < 2) {
			setFacetOptions(facet, []);
			setFacetLoading(facet, false);
			return;
		}
		setFacetLoading(facet, true);
		try {
			setFacetOptions(facet, await autocompleteDiscoverMovieFacet(facet, cleaned));
		} finally {
			setFacetLoading(facet, false);
		}
	}

	function setFacetOptions(facet: MovieFacet, options: DiscoverMovieFacetOption[]) {
		if (facet === 'genres') genreOptions = options;
		if (facet === 'studios') studioOptions = options;
		if (facet === 'keywords') keywordOptions = options;
	}

	function setFacetLoading(facet: MovieFacet, loading: boolean) {
		if (facet === 'genres') loadingGenres = loading;
		if (facet === 'studios') loadingStudios = loading;
		if (facet === 'keywords') loadingKeywords = loading;
	}

	function trackLoadMore(node: HTMLDivElement) {
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries.some((entry) => entry.isIntersecting)) {
					void loadMore();
				}
			},
			{ rootMargin: '700px 0px' }
		);
		observer.observe(node);
		return {
			destroy() {
				observer.disconnect();
			}
		};
	}
</script>

<section class="grid min-w-0 gap-[18px]" aria-labelledby="discover-movies-title">
	<PageHeading eyebrow="Discover" title="Movies" titleId="discover-movies-title" class="w-full">
		{#snippet actions()}
			<div class="mt-5">
				<MovieDiscoverToolbar
					sort={filters.sort}
					{filterCount}
					{filtersOpen}
					onSortChange={updateSort}
					onToggleFilters={() => (filtersOpen = !filtersOpen)}
				/>
			</div>
		{/snippet}
	</PageHeading>

	<MovieDiscoverFilterSheet
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
				<span class="sr-only">Loading more movies when visible</span>
			</div>
		{/if}
	</div>
</section>
