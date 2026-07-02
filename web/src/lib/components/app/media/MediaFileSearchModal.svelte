<script lang="ts">
	import { onDestroy } from 'svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import type {
		MediaItem,
		Language,
		ReleaseCandidate,
		ReleaseOverrideDetails,
		ReleaseSearchState
	} from '$lib/settings/types';
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
	import ReleaseOverrideDetailsStep from './ReleaseOverrideDetailsStep.svelte';
	import ReleaseSearchControls from './ReleaseSearchControls.svelte';
	import ReleaseSearchFilters from './ReleaseSearchFilters.svelte';
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
		languages: Language[];
		canManage: boolean;
		onGrab: (
			_item: MediaItem,
			_release: ReleaseCandidate,
			_overrideMatch?: boolean,
			_details?: ReleaseOverrideDetails
		) => void;
		onClose: () => void;
	}

	let {
		item,
		releaseResults,
		grabbingKey,
		searchContext = { type: 'title' },
		languages,
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
	let sort = $state<ReleaseSort>({ key: 'score', direction: 'desc' });
	let filtersOpen = $state(false);
	let overrideRelease = $state<ReleaseCandidate | undefined>();
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

	function handleGrab(
		grabItem: MediaItem,
		release: ReleaseCandidate,
		overrideMatch = false,
		details?: ReleaseOverrideDetails
	) {
		if (overrideMatch && !details) {
			overrideRelease = release;
			return;
		}
		onGrab(grabItem, release, overrideMatch, details);
	}
</script>

<SettingsFormModal
	title={overrideRelease ? 'Grab with override' : 'Manual search'}
	modalClass="max-h-[calc(100vh-32px)] w-[min(1280px,calc(100vw-32px))]"
	{onClose}
>
	{#if overrideRelease}
		<ReleaseOverrideDetailsStep
			{item}
			release={overrideRelease}
			{languages}
			{qualityOptions}
			grabbing={grabbingKey === `${item.id}:${overrideRelease.id}`}
			onBack={() => (overrideRelease = undefined)}
			onConfirm={handleGrab}
		/>
	{:else}
		<div class="grid gap-5">
			<ReleaseSearchControls
				bind:overrideQuery
				bind:customQuery
				queryVariants={systemQueryVariants}
				disabled={!canManage || searching}
				{filtersOpen}
				{filterCount}
				{searching}
				searchDisabled={!canManage || searching || !searchQuery}
				onToggleFilters={() => (filtersOpen = !filtersOpen)}
				onSearch={submitSearch}
			/>
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
				onGrab={handleGrab}
			/>
		</div>
	{/if}
</SettingsFormModal>
