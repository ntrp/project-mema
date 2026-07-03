<script lang="ts">
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
	import SlidersHorizontalIcon from '@lucide/svelte/icons/sliders-horizontal';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import {
		movieSortDirection,
		movieSortKey,
		movieSortOptions,
		nextMovieSort,
		type MovieSortKey
	} from './discoverMovieFilters';

	interface Props {
		sort: string;
		filterCount: number;
		filtersOpen: boolean;
		onSortChange: (_sort: string) => void;
		onToggleFilters: () => void;
	}

	let { sort, filterCount, filtersOpen, onSortChange, onToggleFilters }: Props = $props();

	const selectedSortKey = $derived(movieSortKey(sort));
	const selectedDirection = $derived(movieSortDirection(sort));
	const selectedSort = $derived(
		movieSortOptions.find((option) => option.key === selectedSortKey) ?? movieSortOptions[0]
	);

	function chooseSort(key: MovieSortKey) {
		onSortChange(nextMovieSort(sort, key));
	}
</script>

<div class="relative z-20 flex flex-wrap items-center justify-end gap-2">
	<Select.Root type="single" value={selectedSortKey}>
		<Select.Trigger class="min-w-[190px] justify-between" data-testid="movie-discover-sort">
			<span>{selectedSort.label}</span>
			{#if selectedDirection === 'desc'}
				<ArrowDownIcon aria-label="Descending" class="size-4 text-muted-foreground" />
			{:else}
				<ArrowUpIcon aria-label="Ascending" class="size-4 text-muted-foreground" />
			{/if}
		</Select.Trigger>
		<Select.Content>
			{#each movieSortOptions as option (option.key)}
				<Select.Item
					value={option.key}
					label={option.label}
					class="pr-2 [&>span:first-child]:hidden"
					onclick={() => chooseSort(option.key)}
				>
					<span class="flex min-w-0 flex-1 items-center justify-between gap-6">
						<span>{option.label}</span>
						{#if option.key === selectedSortKey}
							{#if selectedDirection === 'desc'}
								<ArrowDownIcon aria-label="Descending" class="size-4 text-muted-foreground" />
							{:else}
								<ArrowUpIcon aria-label="Ascending" class="size-4 text-muted-foreground" />
							{/if}
						{/if}
					</span>
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>

	<Button
		type="button"
		variant="outline"
		aria-pressed={filtersOpen}
		data-testid="movie-discover-filters"
		onclick={onToggleFilters}
	>
		<SlidersHorizontalIcon aria-hidden="true" />
		<span>Filters</span>
		{#if filterCount > 0}
			<Badge variant="secondary" class="ml-1 h-5 min-w-5 rounded-[3px] px-1">
				{filterCount}
			</Badge>
		{/if}
	</Button>
</div>
