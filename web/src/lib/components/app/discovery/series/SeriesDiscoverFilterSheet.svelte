<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';
	import type { DiscoverMovieFacetOption } from '$lib/settings/types';
	import SeriesDiscoverFilters from './SeriesDiscoverFilters.svelte';
	import type { DiscoverSeriesFilters } from './discoverSeriesFilters';

	interface Props {
		open: boolean;
		filters: DiscoverSeriesFilters;
		genreOptions: DiscoverMovieFacetOption[];
		studioOptions: DiscoverMovieFacetOption[];
		keywordOptions: DiscoverMovieFacetOption[];
		loadingGenres?: boolean;
		loadingStudios?: boolean;
		loadingKeywords?: boolean;
		onOpenChange: (_open: boolean) => void;
		onChange: (_filters: DiscoverSeriesFilters) => void;
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
			<Sheet.Title>Series filters</Sheet.Title>
			<Sheet.Description>Filter series discovery results.</Sheet.Description>
		</Sheet.Header>
		<SeriesDiscoverFilters
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
