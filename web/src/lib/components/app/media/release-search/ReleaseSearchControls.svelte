<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import SlidersHorizontalIcon from '@lucide/svelte/icons/sliders-horizontal';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import ReleaseSearchQueryInput from '$lib/components/app/media/release-search/ReleaseSearchQueryInput.svelte';

	interface Props {
		overrideQuery: boolean;
		customQuery: string;
		queryVariants: string[];
		disabled?: boolean;
		filtersOpen: boolean;
		filterCount: number;
		searching?: boolean;
		searchDisabled?: boolean;
		onToggleFilters: () => void;
		onSearch: () => void;
	}

	let {
		overrideQuery = $bindable(),
		customQuery = $bindable(),
		queryVariants,
		disabled = false,
		filtersOpen,
		filterCount,
		searching = false,
		searchDisabled = false,
		onToggleFilters,
		onSearch
	}: Props = $props();
</script>

<div class="grid gap-3 md:grid-cols-2 md:items-end">
	<ReleaseSearchQueryInput bind:overrideQuery bind:customQuery {queryVariants} {disabled} />
	<div class="flex items-center justify-end gap-2">
		<Button type="button" variant="outline" aria-pressed={filtersOpen} onclick={onToggleFilters}>
			<SlidersHorizontalIcon aria-hidden="true" />
			<span>Filters</span>
			{#if filterCount > 0}
				<Badge variant="secondary" class="ml-1 h-5 min-w-5 rounded-[3px] px-1">
					{filterCount}
				</Badge>
			{/if}
		</Button>
		<Button type="button" disabled={searchDisabled} onclick={onSearch}>
			<SearchIcon aria-hidden="true" />
			{searching ? 'Searching' : 'Search'}
		</Button>
	</div>
</div>
