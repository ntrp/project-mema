<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import SlidersHorizontalIcon from '@lucide/svelte/icons/sliders-horizontal';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { MediaItem, ReleaseCandidate, ReleaseSearchState } from '$lib/settings/types';
	import {
		activeFilterCount,
		defaultReleaseFilters,
		filteredSortedReleases,
		releaseQualityOptions,
		type ReleaseFilters,
		type ReleaseSort,
		type ReleaseSortKey
	} from './releaseSearchResults';
	import { releaseSearchQuery, type ReleaseSearchContext } from './releaseSearchQuery';
	import ReleaseSearchFilters from './ReleaseSearchFilters.svelte';
	import ReleaseSearchResultsTable from './ReleaseSearchResultsTable.svelte';

	interface Props {
		item: MediaItem;
		releaseResults?: ReleaseSearchState;
		searching?: boolean;
		grabbingKey?: string;
		searchContext?: ReleaseSearchContext;
		canManage: boolean;
		onSearch: (_item: MediaItem, _query?: string) => void;
		onGrab: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onClose: () => void;
	}

	let {
		item,
		releaseResults,
		searching = false,
		grabbingKey,
		searchContext = { type: 'title' },
		canManage,
		onSearch,
		onGrab,
		onClose
	}: Props = $props();

	let overrideQuery = $state(false);
	let customQuery = $state('');
	let filters = $state<ReleaseFilters>(defaultReleaseFilters());
	let sort = $state<ReleaseSort>({ direction: 'desc' });
	let filtersOpen = $state(false);
	const systemQuery = $derived(releaseSearchQuery(item, searchContext));
	const searchQuery = $derived(overrideQuery ? customQuery.trim() : systemQuery);
	const releases = $derived(releaseResults?.releases ?? []);
	const qualityOptions = $derived(releaseQualityOptions(releases));
	const visibleReleases = $derived(filteredSortedReleases(item, releases, filters, sort));
	const filterCount = $derived(activeFilterCount(filters));

	$effect(() => {
		if (!overrideQuery) {
			customQuery = systemQuery;
		}
	});

	function submitSearch() {
		onSearch(item, searchQuery);
	}

	function updateSort(key: ReleaseSortKey) {
		sort =
			sort.key === key
				? { key, direction: sort.direction === 'asc' ? 'desc' : 'asc' }
				: { key, direction: 'asc' };
	}
</script>

<SettingsFormModal
	title="Manual search"
	modalClass="max-h-[calc(100vh-32px)] w-[min(1280px,calc(100vw-32px))]"
	{onClose}
>
	<div class="grid gap-5">
		<div class="grid gap-3 md:grid-cols-2 md:items-end">
			<div class="flex min-w-0 flex-wrap items-end gap-3">
				<div class="grid min-w-72 flex-1 gap-2">
					<Label for="release-search-query">Search query</Label>
					<Input
						id="release-search-query"
						class={!overrideQuery ? 'bg-muted text-muted-foreground opacity-80' : ''}
						bind:value={customQuery}
						readonly={!overrideQuery}
						disabled={!canManage || searching}
						maxlength={500}
					/>
				</div>
				<div class="flex h-9 items-center gap-2">
					<Checkbox
						id="release-search-query-override"
						bind:checked={overrideQuery}
						disabled={!canManage || searching}
					/>
					<Label
						for="release-search-query-override"
						class={!canManage || searching ? 'text-muted-foreground opacity-70' : ''}
					>
						Override
					</Label>
				</div>
			</div>
			<div class="flex items-center justify-end gap-2">
				<Button
					type="button"
					variant="outline"
					aria-pressed={filtersOpen}
					onclick={() => (filtersOpen = !filtersOpen)}
				>
					<SlidersHorizontalIcon aria-hidden="true" />
					<span>Filters</span>
					{#if filterCount > 0}
						<Badge variant="secondary" class="ml-1 h-5 min-w-5 rounded-full px-1">
							{filterCount}
						</Badge>
					{/if}
				</Button>
				<Button
					type="button"
					disabled={!canManage || searching || !searchQuery}
					onclick={submitSearch}
				>
					<SearchIcon aria-hidden="true" />
					{searching ? 'Searching' : 'Search'}
				</Button>
			</div>
		</div>
		{#if releaseResults?.loaded && releaseResults.errors.length}
			<div
				class="grid gap-1 rounded-md bg-secondary px-3 py-2.5 font-bold text-secondary-foreground"
			>
				{#each releaseResults.errors as error (error)}
					<p class="m-0">{error}</p>
				{/each}
			</div>
		{/if}
		{#if filtersOpen}
			<ReleaseSearchFilters
				{filters}
				{qualityOptions}
				onChange={(nextFilters) => (filters = nextFilters)}
				onReset={() => (filters = defaultReleaseFilters())}
			/>
		{/if}
		<ReleaseSearchResultsTable
			{item}
			releases={visibleReleases}
			{releaseResults}
			{sort}
			{grabbingKey}
			{canManage}
			onSort={updateSort}
			{onGrab}
		/>
	</div>
</SettingsFormModal>
