<script lang="ts">
	import { onDestroy } from 'svelte';
	import SearchIcon from '@lucide/svelte/icons/search';
	import SlidersHorizontalIcon from '@lucide/svelte/icons/sliders-horizontal';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
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
	import {
		releaseSearchQuery,
		releaseSearchQueryVariants,
		type ReleaseSearchContext
	} from './releaseSearchQuery';
	import ReleaseSearchFilters from './ReleaseSearchFilters.svelte';
	import ReleaseSearchQueryInput from './ReleaseSearchQueryInput.svelte';
	import ReleaseSearchResultsTable from './ReleaseSearchResultsTable.svelte';
	import ReleaseSearchStatusLog from './ReleaseSearchStatusLog.svelte';
	import {
		applyStatusToLog,
		createLogEntry,
		placeholderLogEntry,
		type ReleaseSearchLogEntry
	} from './releaseSearchLog';
	import {
		subscribeReleaseSearchStream,
		type ReleaseSearchStreamResult
	} from './releaseSearchStream';

	interface Props {
		item: MediaItem;
		releaseResults?: ReleaseSearchState;
		grabbingKey?: string;
		searchContext?: ReleaseSearchContext;
		canManage: boolean;
		onGrab: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onClose: () => void;
	}

	let {
		item,
		releaseResults,
		grabbingKey,
		searchContext = { type: 'title' },
		canManage,
		onGrab,
		onClose
	}: Props = $props();

	let overrideQuery = $state(false);
	let customQuery = $state('');
	let localResults = $state<ReleaseSearchStreamResult | undefined>();
	let searching = $state(false);
	let statusMessages = $state<ReleaseSearchLogEntry[]>([placeholderLogEntry()]);
	let filters = $state<ReleaseFilters>(defaultReleaseFilters());
	let sort = $state<ReleaseSort>({ direction: 'desc' });
	let filtersOpen = $state(false);
	let unsubscribeSearch: (() => void) | undefined;
	const systemQuery = $derived(releaseSearchQuery(item, searchContext));
	const systemQueryVariants = $derived(releaseSearchQueryVariants(item, searchContext));
	const searchQuery = $derived(overrideQuery ? customQuery.trim() : systemQuery);
	const currentResults = $derived(localResults ?? releaseResults);
	const releases = $derived(currentResults?.releases ?? []);
	const qualityOptions = $derived(releaseQualityOptions(releases));
	const visibleReleases = $derived(filteredSortedReleases(item, releases, filters, sort));
	const filterCount = $derived(activeFilterCount(filters));

	$effect(() => {
		if (!overrideQuery) {
			customQuery = systemQuery;
		}
	});

	onDestroy(() => {
		unsubscribeSearch?.();
	});

	function submitSearch() {
		unsubscribeSearch?.();
		searching = true;
		localResults = { releases: [], errors: [] };
		statusMessages = [createLogEntry('Search started')];
		unsubscribeSearch = subscribeReleaseSearchStream(item.id, searchQuery, {
			onStatus: appendStatus,
			onResult: (result) => {
				localResults = result;
				statusMessages = [
					...statusMessages,
					createLogEntry(`Search finished: ${result.releases.length} releases`)
				].slice(-100);
				searching = false;
				unsubscribeSearch = undefined;
			},
			onError: (message) => {
				statusMessages = [...statusMessages, createLogEntry(message)].slice(-100);
				searching = false;
				unsubscribeSearch = undefined;
			}
		});
	}

	function updateSort(key: ReleaseSortKey) {
		sort =
			sort.key === key
				? { key, direction: sort.direction === 'asc' ? 'desc' : 'asc' }
				: { key, direction: 'asc' };
	}

	function appendStatus(status: Parameters<typeof applyStatusToLog>[1]) {
		statusMessages = applyStatusToLog(statusMessages, status).slice(-100);
	}
</script>

<SettingsFormModal
	title="Manual search"
	modalClass="max-h-[calc(100vh-32px)] w-[min(1280px,calc(100vw-32px))]"
	{onClose}
>
	<div class="grid gap-5">
		<div class="grid gap-3 md:grid-cols-2 md:items-end">
			<ReleaseSearchQueryInput
				bind:overrideQuery
				bind:customQuery
				queryVariants={systemQueryVariants}
				disabled={!canManage || searching}
			/>
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
		<ReleaseSearchStatusLog messages={statusMessages} />
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
			{searching}
			{sort}
			{grabbingKey}
			{canManage}
			onSort={updateSort}
			{onGrab}
		/>
	</div>
</SettingsFormModal>
