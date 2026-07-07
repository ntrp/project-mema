<script lang="ts">
	import LibraryScanImportTableView from '$lib/components/settings/library/scan/LibraryScanImportTableView.svelte';
	import {
		applyMovieOptions,
		applyQualityProfile,
		applySeriesOptions,
		matchedRowsByKind,
		setRowsSelected
	} from './libraryScanBulk';
	import { applyAutoMatch } from './libraryScanAutoMatch';
	import {
		defaultMetadataProviderId,
		ensureScanDrafts,
		matchFromScanItem
	} from './libraryScanDrafts';
	import { duplicateDraftStatesForRows } from './libraryScanDuplicates';
	import {
		canImportRows,
		defaultQualityProfileId,
		sortedScanItems,
		type MatchDraft
	} from './libraryScanImport';
	import {
		applyScanItemProvider,
		changeScanItemProvider,
		importCheckedScanRows,
		scheduleScanItemSearch,
		searchPendingScanItems,
		searchScanItem
	} from './libraryScanTableActions';
	import type { LibraryScanImportTableProps } from './libraryScanTableProps';
	import type {
		LibraryScanItem,
		MediaMonitorMode,
		MediaSearchResult,
		MinimumAvailability,
		SeriesType
	} from '$lib/settings/types';

	let {
		scan,
		qualityProfiles,
		metadataProviders,
		loading,
		onSearchMatch,
		onImport
	}: LibraryScanImportTableProps = $props();
	let drafts = $state<Record<string, MatchDraft>>({});
	let showImported = $state(false);
	let bulkQualityProfileId = $state('');
	let bulkMetadataProviderId = $state('');
	let defaultsApplied = $state({ profile: false, provider: false });
	let movieMonitorMode = $state<MediaMonitorMode>('only_media');
	let movieMinimumAvailability = $state<MinimumAvailability>('released');
	let seriesMonitorMode = $state<MediaMonitorMode>('all_episodes');
	let bulkSeriesType = $state<SeriesType>('standard');
	let importing = $state(false);
	let importingItemId = $state('');
	let searchTimers: Record<string, ReturnType<typeof globalThis.setTimeout>> = {};
	let autoSearchStarted: Record<string, boolean> = {};
	const searchCache: Record<string, MediaSearchResult[] | undefined> = {};
	const allRows = $derived(sortedScanItems(scan.items));
	const rows = $derived(allRows.filter((item) => showImported || !item.imported));
	const importableRows = $derived(
		rows.filter((item) => item.status === 'pending' && !item.imported && drafts[item.id]?.matched)
	);
	const checkedRows = $derived(importableRows.filter((item) => drafts[item.id]?.selected));
	const checkedMatchedMovies = $derived(matchedRowsByKind(checkedRows, drafts, 'movie'));
	const checkedMatchedSeries = $derived(matchedRowsByKind(checkedRows, drafts, 'series'));
	const checkedRowsMatched = $derived(
		checkedRows.length > 0 && checkedRows.every((item) => drafts[item.id]?.matched)
	);
	const canImport = $derived(canImportRows(checkedRows, drafts, bulkQualityProfileId));
	const allVisibleChecked = $derived(
		importableRows.length > 0 && importableRows.every((item) => drafts[item.id]?.selected)
	);
	const matchedCount = $derived(
		allRows.filter((item) => drafts[item.id]?.matched || matchFromScanItem(item)).length
	);
	const importedCount = $derived(scan.items.filter((item) => item.imported).length);
	const duplicateCount = $derived(scan.items.filter((item) => item.duplicateGroupId).length);
	const noMatchCount = $derived(scan.items.length - matchedCount);
	const duplicateStates = $derived(duplicateDraftStatesForRows(allRows, drafts));
	const defaultProfileId = $derived(defaultQualityProfileId(qualityProfiles));
	const defaultProviderId = $derived(defaultMetadataProviderId(metadataProviders, 'movie'));
	const duplicateRemovalCount = $derived(
		Object.entries(drafts).filter(
			([id, draft]) => draft.removeDuplicate && (!draft.matched || duplicateStates[id]?.duplicate)
		).length
	);
	$effect(() => {
		if (!bulkQualityProfileId && !defaultsApplied.profile && defaultProfileId)
			bulkQualityProfileId = defaultProfileId;
		if (defaultProfileId) defaultsApplied.profile = true;
		if (!bulkMetadataProviderId && !defaultsApplied.provider && defaultProviderId)
			bulkMetadataProviderId = defaultProviderId;
		if (defaultProviderId) defaultsApplied.provider = true;
	});
	$effect(() => {
		ensureScanDrafts(scan.items, drafts, metadataProviders, {
			qualityProfileId: bulkQualityProfileId,
			monitorMode: movieMonitorMode,
			minimumAvailability: movieMinimumAvailability,
			seriesType: bulkSeriesType
		});
	});
	$effect(() => {
		searchPendingScanItems({ rows: allRows, drafts, autoSearchStarted, search });
	});
	function toggleVisibleRows() {
		setRowsSelected(importableRows, drafts, !allVisibleChecked);
	}
	function scheduleSearch(item: LibraryScanItem) {
		scheduleScanItemSearch({ item, drafts, searchTimers, search });
	}
	async function search(item: LibraryScanItem, auto: boolean) {
		await searchScanItem({ item, allRows, drafts, searchCache, auto, onSearchMatch });
	}
	function changeProvider(item: LibraryScanItem, providerId: string) {
		changeScanItemProvider({ item, providerId, drafts, search });
	}
	function selectResult(item: LibraryScanItem, result: MediaSearchResult) {
		applyAutoMatch(item, result, allRows, drafts);
	}
	function applyProvider() {
		applyScanItemProvider({
			rows: [...checkedRows],
			drafts,
			providerId: bulkMetadataProviderId,
			search
		});
	}
	async function importChecked() {
		if (!canImport) return;
		importing = true;
		try {
			await importCheckedScanRows({
				canImport,
				checkedRows: [...checkedRows],
				allRows,
				drafts,
				scan,
				onProgress: (id) => (importingItemId = id),
				onImport,
				bulk: {
					qualityProfileId: bulkQualityProfileId,
					monitorMode: movieMonitorMode,
					minimumAvailability: movieMinimumAvailability,
					seriesType: bulkSeriesType
				}
			});
		} finally {
			importingItemId = '';
			importing = false;
		}
	}
</script>

<LibraryScanImportTableView
	totalFiles={scan.totalFiles}
	{matchedCount}
	{noMatchCount}
	{importedCount}
	{duplicateCount}
	{rows}
	folderPath={scan.folderPath}
	bind:drafts
	{duplicateStates}
	{qualityProfiles}
	{metadataProviders}
	{importingItemId}
	bind:showImported
	{allVisibleChecked}
	importableCount={importableRows.length}
	showSeriesControls={scan.folderKind !== 'movie'}
	{checkedRowsMatched}
	{canImport}
	{loading}
	{importing}
	{duplicateRemovalCount}
	hasMatchedMovies={checkedMatchedMovies.length > 0}
	hasMatchedSeries={checkedMatchedSeries.length > 0}
	bind:metadataProviderId={bulkMetadataProviderId}
	bind:qualityProfileId={bulkQualityProfileId}
	bind:movieMonitorMode
	bind:movieMinimumAvailability
	bind:seriesMonitorMode
	bind:seriesType={bulkSeriesType}
	onToggleRows={toggleVisibleRows}
	onSearch={scheduleSearch}
	onSelect={selectResult}
	onProviderChange={changeProvider}
	onApplyProvider={applyProvider}
	onApplyQualityProfile={() => applyQualityProfile(checkedRows, drafts, bulkQualityProfileId)}
	onApplyMovie={() =>
		applyMovieOptions(checkedMatchedMovies, drafts, movieMonitorMode, movieMinimumAvailability)}
	onApplySeries={() =>
		applySeriesOptions(checkedMatchedSeries, drafts, seriesMonitorMode, bulkSeriesType)}
	onImport={importChecked}
/>
