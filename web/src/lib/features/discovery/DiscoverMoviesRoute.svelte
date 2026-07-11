<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import MovieDiscoverFilterSheet from '$lib/components/app/discovery/movies/MovieDiscoverFilterSheet.svelte';
	import MovieDiscoverResults from '$lib/components/app/discovery/movies/MovieDiscoverResults.svelte';
	import MovieDiscoverToolbar from '$lib/components/app/discovery/movies/MovieDiscoverToolbar.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { createMediaItemsQuery } from '$lib/features/library/queries.svelte';
	import {
		activeMovieFilterCount,
		filtersFromParams,
		movieFilterUrl,
		movieQuery,
		type DiscoverMovieFilters
	} from '$lib/components/app/discovery/movies/discoverMovieFilters';
	import {
		createDiscoverFacetQuery,
		createMovieSearchQuery,
		type DiscoverFacet
	} from './search/queries.svelte';

	const app = getAppShellContext();
	const library = createMediaItemsQuery();
	const filters = $derived(filtersFromParams(page.url.searchParams));
	let filtersOpen = $state(false);
	let genreInput = $state('');
	let studioInput = $state('');
	let keywordInput = $state('');
	const search = createMovieSearchQuery(() => movieQuery(filters, 1));
	const genres = createDiscoverFacetQuery('movies', 'genres', () => genreInput);
	const studios = createDiscoverFacetQuery('movies', 'studios', () => studioInput);
	const keywords = createDiscoverFacetQuery('movies', 'keywords', () => keywordInput);
	const results = $derived(search.data?.pages.flatMap((result) => result.results ?? []) ?? []);
	const filterCount = $derived(activeMovieFilterCount(filters));

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

	function setFacetInput(facet: DiscoverFacet, query: string) {
		if (facet === 'genres') genreInput = query;
		if (facet === 'studios') studioInput = query;
		if (facet === 'keywords') keywordInput = query;
	}

	function trackLoadMore(node: HTMLDivElement) {
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries.some((entry) => entry.isIntersecting)) {
					void search.fetchNextPage();
				}
			},
			{ rootMargin: '700px 0px' }
		);
		observer.observe(node);
		return () => observer.disconnect();
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
		genreOptions={genres.data ?? []}
		studioOptions={studios.data ?? []}
		keywordOptions={keywords.data ?? []}
		loadingGenres={genres.isFetching}
		loadingStudios={studios.isFetching}
		loadingKeywords={keywords.isFetching}
		onOpenChange={(open) => (filtersOpen = open)}
		onChange={updateFilters}
		onGenreQuery={(query) => setFacetInput('genres', query)}
		onStudioQuery={(query) => setFacetInput('studios', query)}
		onKeywordQuery={(query) => setFacetInput('keywords', query)}
	/>

	<div class="grid min-w-0 content-start gap-4">
		<MovieDiscoverResults
			{results}
			mediaItems={library.data ?? []}
			blacklist={app.discoverBlacklist}
			loading={search.isPending}
			loadingMore={search.isFetchingNextPage}
			hasSearched={!search.isPending}
			addingKey={app.addingKey}
			blacklistingKey={app.blacklistingKey}
			actionLabel={app.isAdmin ? 'Add' : 'Request'}
			canManage={app.isAdmin}
			onAdd={app.addMedia}
			onBlacklist={app.blacklistDiscoverMedia}
		/>
		{#if search.hasNextPage}
			<div class="min-h-8" {@attach trackLoadMore}>
				<span class="sr-only">Loading more movies when visible</span>
			</div>
		{/if}
	</div>
</section>
