<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';
	import type { DiscoverMovieFacetOption } from '$lib/settings/types';
	import MovieDiscoverFilters from './MovieDiscoverFilters.svelte';
	import type { DiscoverMovieFilters } from './discoverMovieFilters';

	interface Props {
		open: boolean;
		filters: DiscoverMovieFilters;
		genreOptions: DiscoverMovieFacetOption[];
		studioOptions: DiscoverMovieFacetOption[];
		keywordOptions: DiscoverMovieFacetOption[];
		loadingGenres?: boolean;
		loadingStudios?: boolean;
		loadingKeywords?: boolean;
		onOpenChange: (_open: boolean) => void;
		onChange: (_filters: DiscoverMovieFilters) => void;
		onGenreQuery: (_query: string) => void;
		onStudioQuery: (_query: string) => void;
		onKeywordQuery: (_query: string) => void;
	}

	let {
		open,
		filters,
		genreOptions,
		studioOptions,
		keywordOptions,
		loadingGenres = false,
		loadingStudios = false,
		loadingKeywords = false,
		onOpenChange,
		onChange,
		onGenreQuery,
		onStudioQuery,
		onKeywordQuery
	}: Props = $props();
</script>

<Sheet.Root {open} {onOpenChange}>
	<Sheet.Content
		side="right"
		class="w-[min(420px,calc(100vw-24px))] overflow-y-auto p-4 sm:max-w-[420px]"
	>
		<Sheet.Header class="sr-only">
			<Sheet.Title>Movie filters</Sheet.Title>
			<Sheet.Description>Filter movie discovery results.</Sheet.Description>
		</Sheet.Header>
		<MovieDiscoverFilters
			{filters}
			{genreOptions}
			{studioOptions}
			{keywordOptions}
			{loadingGenres}
			{loadingStudios}
			{loadingKeywords}
			{onChange}
			{onGenreQuery}
			{onStudioQuery}
			{onKeywordQuery}
		/>
	</Sheet.Content>
</Sheet.Root>
