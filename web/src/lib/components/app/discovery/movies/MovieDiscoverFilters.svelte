<script lang="ts">
	import RotateCcwIcon from '@lucide/svelte/icons/rotate-ccw';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { DiscoverMovieFacetOption } from '$lib/settings/types';
	import MovieFacetAutocomplete from './MovieFacetAutocomplete.svelte';
	import MovieFilterMultiSelect from './MovieFilterMultiSelect.svelte';
	import MovieRangeFilter from './MovieRangeFilter.svelte';
	import MovieSingleRangeFilter from './MovieSingleRangeFilter.svelte';
	import {
		contentRatingOptions,
		defaultMovieFilters,
		languageOptions,
		type DiscoverMovieFilters
	} from './discoverMovieFilters';

	interface Props {
		filters: DiscoverMovieFilters;
		genreOptions: DiscoverMovieFacetOption[];
		studioOptions: DiscoverMovieFacetOption[];
		keywordOptions: DiscoverMovieFacetOption[];
		loadingGenres?: boolean;
		loadingStudios?: boolean;
		loadingKeywords?: boolean;
		onChange: (_filters: DiscoverMovieFilters) => void;
		onGenreQuery: (_query: string) => void;
		onStudioQuery: (_query: string) => void;
		onKeywordQuery: (_query: string) => void;
	}

	let {
		filters,
		genreOptions,
		studioOptions,
		keywordOptions,
		loadingGenres = false,
		loadingStudios = false,
		loadingKeywords = false,
		onChange,
		onGenreQuery,
		onStudioQuery,
		onKeywordQuery
	}: Props = $props();

	function patch(next: Partial<DiscoverMovieFilters>) {
		onChange({ ...filters, ...next });
	}

	function resetFilters() {
		onChange({ ...defaultMovieFilters(), sort: filters.sort });
	}
</script>

<aside class="grid gap-3" aria-label="Movie filters">
	<div class="flex items-center justify-between gap-3">
		<h2 class="m-0 text-base font-semibold text-foreground">Filters</h2>
		<Button type="button" variant="ghost" size="sm" onclick={resetFilters}>
			<RotateCcwIcon aria-hidden="true" />
			Reset
		</Button>
	</div>

	<div class="grid grid-cols-2 gap-2">
		<div class="grid gap-2">
			<Label for="movie-release-from">Release from</Label>
			<Input
				id="movie-release-from"
				type="date"
				value={filters.releaseDateFrom}
				oninput={(event) => patch({ releaseDateFrom: event.currentTarget.value })}
			/>
		</div>
		<div class="grid gap-2">
			<Label for="movie-release-to">Release to</Label>
			<Input
				id="movie-release-to"
				type="date"
				value={filters.releaseDateTo}
				oninput={(event) => patch({ releaseDateTo: event.currentTarget.value })}
			/>
		</div>
	</div>

	<MovieFacetAutocomplete
		id="movie-studios"
		label="Studio"
		values={filters.studios}
		placeholder="Search studios"
		options={studioOptions}
		loading={loadingStudios}
		onQuery={onStudioQuery}
		onChange={(studios) => patch({ studios })}
	/>
	<MovieFacetAutocomplete
		id="movie-genres"
		label="Genres"
		values={filters.genres}
		excludedValues={filters.withoutGenres}
		placeholder="Search genres"
		options={genreOptions}
		loading={loadingGenres}
		onQuery={onGenreQuery}
		onChange={(genres) => patch({ genres })}
		onExcludedChange={(withoutGenres) => patch({ withoutGenres })}
		onSignedChange={(genres, withoutGenres) => patch({ genres, withoutGenres })}
	/>
	<MovieFacetAutocomplete
		id="movie-keywords"
		label="Keywords"
		values={filters.keywords}
		excludedValues={filters.withoutKeywords}
		placeholder="Search keywords"
		options={keywordOptions}
		loading={loadingKeywords}
		onQuery={onKeywordQuery}
		onChange={(keywords) => patch({ keywords })}
		onExcludedChange={(withoutKeywords) => patch({ withoutKeywords })}
		onSignedChange={(keywords, withoutKeywords) => patch({ keywords, withoutKeywords })}
	/>
	<MovieFilterMultiSelect
		id="movie-languages"
		label="Original language"
		values={filters.originalLanguages}
		options={languageOptions}
		placeholder="Select languages"
		onChange={(originalLanguages) => patch({ originalLanguages })}
	/>
	<MovieFilterMultiSelect
		id="movie-ratings"
		label="Content rating"
		values={filters.contentRatings}
		options={contentRatingOptions}
		placeholder="Select ratings"
		onChange={(contentRatings) => patch({ contentRatings })}
	/>
	<MovieRangeFilter
		label="Runtime"
		value={filters.runtime}
		min={0}
		max={400}
		unit="m"
		onChange={(runtime) => patch({ runtime })}
	/>
	<MovieRangeFilter
		label="TMDB score"
		value={filters.score}
		min={0}
		max={10}
		step={0.1}
		onChange={(score) => patch({ score })}
	/>
	<MovieSingleRangeFilter
		label="TMDB min vote count"
		value={filters.minVoteCount}
		min={0}
		max={1000}
		unit=" votes"
		onChange={(minVoteCount) => patch({ minVoteCount })}
	/>
</aside>
