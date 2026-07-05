<script lang="ts">
	import LibraryScanImportFooter from '$lib/components/settings/library/scan/LibraryScanImportFooter.svelte';
	import LibraryScanImportTableBody from '$lib/components/settings/library/scan/LibraryScanImportTableBody.svelte';
	import LibraryScanImportTableHead from '$lib/components/settings/library/scan/LibraryScanImportTableHead.svelte';
	import LibraryScanImportToolbar from '$lib/components/settings/library/scan/LibraryScanImportToolbar.svelte';
	import * as Table from '$lib/components/ui/table';
	import {
		applyMovieOptions,
		applyQualityProfile,
		applySeriesOptions,
		matchedRowsByKind
	} from './libraryScanBulk';
	import { initialMatchDraft, matchFromScanItem } from './libraryScanDrafts';
	import { duplicateDraftStatesForRows, normalizeDuplicateDrafts } from './libraryScanDuplicates';
	import {
		canImportRows,
		defaultQualityProfileId,
		defaultMonitorModeForMatch,
		importPayloadForRows,
		searchCacheKey,
		sortedScanItems,
		type MatchDraft
	} from './libraryScanImport';
	import type { LibraryMediaKind, LibraryScan, LibraryScanImportRequest, LibraryScanItem, MediaMonitorMode, MediaSearchResult, MetadataProvider, MinimumAvailability, QualityProfileOption, SeriesType } from '$lib/settings/types';

	interface Props {
		scan: LibraryScan;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		loading: boolean;
		onSearchMatch: (_kind: LibraryMediaKind, _query: string, _providerId?: string) => Promise<MediaSearchResult[]>;
		onImport: (_scan: LibraryScan, _request: LibraryScanImportRequest) => Promise<void>;
	}
	let { scan, qualityProfiles, metadataProviders, loading, onSearchMatch, onImport }: Props = $props();
	let drafts = $state<Record<string, MatchDraft>>({});
	let showImported = $state(false);
	let bulkQualityProfileId = $state('');
	let defaultProfileApplied = $state(false);
	let movieMonitorMode = $state<MediaMonitorMode>('only_media');
	let movieMinimumAvailability = $state<MinimumAvailability>('released');
	let seriesMonitorMode = $state<MediaMonitorMode>('all_episodes');
	let bulkSeriesType = $state<SeriesType>('standard');
	let importing = $state(false);
	let searchTimers: Record<string, ReturnType<typeof globalThis.setTimeout>> = {};
	const searchCache: Record<string, MediaSearchResult[] | undefined> = {};
	const allRows = $derived(sortedScanItems(scan.items));
	const rows = $derived(allRows.filter((item) => showImported || !item.imported));
	const importableRows = $derived(rows.filter((item) => item.status === 'pending' && !item.imported && drafts[item.id]?.matched));
	const checkedRows = $derived(importableRows.filter((item) => drafts[item.id]?.selected));
	const checkedMatchedMovies = $derived(matchedRowsByKind(checkedRows, drafts, 'movie'));
	const checkedMatchedSeries = $derived(matchedRowsByKind(checkedRows, drafts, 'series'));
	const checkedRowsMatched = $derived(checkedRows.length > 0 && checkedRows.every((item) => drafts[item.id]?.matched));
	const canImport = $derived(canImportRows(checkedRows, drafts, bulkQualityProfileId));
	const allVisibleChecked = $derived(importableRows.length > 0 && importableRows.every((item) => drafts[item.id]?.selected));
	const matchedCount = $derived(scan.items.filter((item) => matchFromScanItem(item)).length);
	const importedCount = $derived(scan.items.filter((item) => item.imported).length);
	const duplicateCount = $derived(scan.items.filter((item) => item.duplicateGroupId).length);
	const noMatchCount = $derived(scan.items.length - matchedCount);
	const duplicateStates = $derived(duplicateDraftStatesForRows(allRows, drafts));
	const defaultProfileId = $derived(defaultQualityProfileId(qualityProfiles));
	const duplicateRemovalCount = $derived(Object.entries(drafts).filter(([id, draft]) => draft.removeDuplicate && (!draft.matched || duplicateStates[id]?.duplicate)).length);
	const showSeriesControls = $derived(scan.folderKind !== 'movie');
	$effect(() => {
		if (!bulkQualityProfileId && !defaultProfileApplied && defaultProfileId) bulkQualityProfileId = defaultProfileId;
		if (defaultProfileId) defaultProfileApplied = true;
	});
	$effect(() => {
		for (const item of scan.items) {
			if (drafts[item.id]) continue;
			drafts[item.id] = initialMatchDraft(item, metadataProviders, {
				qualityProfileId: bulkQualityProfileId,
				monitorMode: movieMonitorMode,
				minimumAvailability: movieMinimumAvailability,
				seriesType: bulkSeriesType
			});
		}
	});
	function toggleVisibleRows() {
		const next = !allVisibleChecked;
		for (const item of importableRows) {
			const draft = drafts[item.id];
			if (draft) draft.selected = next;
		}
	}
	function scheduleSearch(item: LibraryScanItem) {
		const draft = drafts[item.id];
		if (!draft) return;
		globalThis.clearTimeout(searchTimers[item.id]);
		if (draft.query.trim().length < 2) {
			draft.results = [];
			draft.searched = false;
			draft.searching = false;
			return;
		}
		searchTimers[item.id] = globalThis.setTimeout(() => void search(item), 1000);
	}
	async function search(item: LibraryScanItem) {
		const draft = drafts[item.id];
		if (!draft || draft.query.trim().length < 2) return;
		const query = draft.query;
		const key = searchCacheKey(draft.mediaKind, draft.metadataProviderId, query);
		draft.searching = true;
		try {
			const results =
				searchCache[key] ?? (await onSearchMatch(draft.mediaKind, query, draft.metadataProviderId));
			if (draft.query !== query) return;
			searchCache[key] = results;
			draft.results = results;
			draft.searched = true;
		} catch {
			if (draft.query !== query) return;
			draft.results = [];
			draft.searched = true;
		} finally {
			if (draft.query === query) draft.searching = false;
		}
	}
	function selectResult(item: LibraryScanItem, result: MediaSearchResult) {
		const draft = drafts[item.id];
		if (!draft) return;
		draft.matched = result;
		draft.query = result.title;
		draft.results = [];
		draft.selected = true;
		draft.monitorMode = defaultMonitorModeForMatch(result);
		normalizeDuplicateDrafts(allRows, drafts, item.duplicateGroupId);
	}
	function applyBulk() { applyQualityProfile(checkedRows, drafts, bulkQualityProfileId); }
	function applyMovieBulk() { applyMovieOptions(checkedMatchedMovies, drafts, movieMonitorMode, movieMinimumAvailability); }
	function applySeriesBulk() { applySeriesOptions(checkedMatchedSeries, drafts, seriesMonitorMode, bulkSeriesType); }
	async function importChecked() {
		if (!canImport) return;
		importing = true;
		try {
			await onImport(
				scan,
				importPayloadForRows(checkedRows, allRows, drafts, {
					qualityProfileId: bulkQualityProfileId,
					monitorMode: movieMonitorMode,
					minimumAvailability: movieMinimumAvailability,
					seriesType: bulkSeriesType
				})
			);
		} finally {
			importing = false;
		}
	}
</script>

<LibraryScanImportToolbar totalFiles={scan.totalFiles} {matchedCount} {noMatchCount} {importedCount} {duplicateCount} bind:showImported />
<div class="mt-4 overflow-auto rounded-md border border-border">
	<Table.Root class="min-w-7xl table-auto border-collapse">
		<colgroup>
			<col class="w-[1%]" />
			<col class="w-full" />
			<col span={5} class="w-[1%]" />
		</colgroup>
		<LibraryScanImportTableHead checked={allVisibleChecked} disabled={importableRows.length === 0} {showSeriesControls} onToggle={toggleVisibleRows} />
		<LibraryScanImportTableBody {rows} bind:drafts {duplicateStates} {qualityProfiles} {metadataProviders} onSearch={scheduleSearch} onSelect={selectResult} />
		<LibraryScanImportFooter
			{checkedRowsMatched}
			{canImport}
			{loading}
			{importing}
			{duplicateRemovalCount}
			{qualityProfiles}
			hasMatchedMovies={checkedMatchedMovies.length > 0}
			hasMatchedSeries={checkedMatchedSeries.length > 0}
			{showSeriesControls}
			bind:qualityProfileId={bulkQualityProfileId}
			bind:movieMonitorMode
			bind:movieMinimumAvailability
			bind:seriesMonitorMode
			bind:seriesType={bulkSeriesType}
			onApplyQualityProfile={applyBulk}
			onApplyMovie={applyMovieBulk}
			onApplySeries={applySeriesBulk}
			onImport={importChecked}
		/>
	</Table.Root>
</div>
