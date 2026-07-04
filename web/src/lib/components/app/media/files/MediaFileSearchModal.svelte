<script lang="ts">
	import { onDestroy } from 'svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import type { MediaItem, ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';
	import {
		activeFilterCount,
		defaultReleaseFilters,
		filteredSortedReleases,
		releaseQualityOptions,
		type ReleaseFilters,
		type ReleaseSort,
		type ReleaseSortKey
	} from '$lib/components/app/media/release-display/releaseSearchResults';
	import {
		releaseSearchQuery,
		releaseSearchQueryVariants
	} from '$lib/components/app/media/release-search/releaseSearchQuery';
	import ReleaseOverrideDetailsStep from '$lib/components/app/media/release-override/ReleaseOverrideDetailsStep.svelte';
	import ReleaseSearchControls from '$lib/components/app/media/release-search/ReleaseSearchControls.svelte';
	import ReleaseSearchFilters from '$lib/components/app/media/release-search/ReleaseSearchFilters.svelte';
	import ReleaseSearchResultsTable from '$lib/components/app/media/release-search/ReleaseSearchResultsTable.svelte';
	import ReleaseSearchStatusLog from '$lib/components/app/media/release-search/ReleaseSearchStatusLog.svelte';
	import {
		applyStatusToLog,
		createLogEntry,
		placeholderLogEntry,
		type ReleaseSearchLogEntry
	} from '$lib/components/app/media/release-search/releaseSearchLog';
	import {
		subscribeReleaseSearchStream,
		type ReleaseSearchStreamResult
	} from '$lib/components/app/media/release-search/releaseSearchStream';
	import type { MediaFileSearchModalProps } from '$lib/components/app/media/files/mediaFileSearchModalTypes';

	let {
		item,
		grabbingKey,
		searchContext = { type: 'title' },
		languages,
		canManage,
		onGrab,
		onClose
	}: MediaFileSearchModalProps = $props();

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
	const currentResults = $derived(localResults);
	const releases = $derived(currentResults?.releases ?? []);
	const qualityOptions = $derived(releaseQualityOptions(releases));
	const visibleReleases = $derived(filteredSortedReleases(item, releases, filters, sort));
	const filterCount = $derived(activeFilterCount(filters));
	const filterResetKey = $derived(Object.values(filters).join('\0'));
	const showResultsTable = $derived(searching || releases.length > 0);

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
				: { key, direction: defaultSortDirection(key) };
	}

	function defaultSortDirection(key: ReleaseSortKey): ReleaseSort['direction'] {
		return ['score', 'quality', 'size', 'peers'].includes(key) ? 'desc' : 'asc';
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
	preventScroll={false}
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
			{#if showResultsTable}
				<ReleaseSearchResultsTable
					{item}
					releases={visibleReleases}
					{searching}
					{sort}
					resetKey={filterResetKey}
					{grabbingKey}
					{canManage}
					onSort={updateSort}
					onGrab={handleGrab}
				/>
			{/if}
		</div>
	{/if}
</SettingsFormModal>
