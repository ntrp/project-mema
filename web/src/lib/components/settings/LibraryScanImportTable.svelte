<script lang="ts">
	import LibraryScanImportFooter from '$lib/components/settings/LibraryScanImportFooter.svelte';
	import LibraryScanImportRow from '$lib/components/settings/LibraryScanImportRow.svelte';
	import LibraryScanImportTableHead from '$lib/components/settings/LibraryScanImportTableHead.svelte';
	import LibraryScanImportToolbar from '$lib/components/settings/LibraryScanImportToolbar.svelte';
	import {
		folderName,
		scanMediaKind,
		searchCacheKey,
		sortedScanItems,
		wait,
		type LibraryScanImportRow as ImportRow,
		type MatchDraft
	} from '$lib/components/settings/libraryScanImport';
	import type {
		LibraryMediaKind,
		LibraryScan,
		LibraryScanItem,
		MediaMonitorMode,
		MediaSearchResult,
		MinimumAvailability,
		QualityProfileOption
	} from '$lib/settings/types';
	interface Props {
		scan: LibraryScan;
		qualityProfiles: QualityProfileOption[];
		loading: boolean;
		onSearchMatch: (_kind: LibraryMediaKind, _query: string) => Promise<MediaSearchResult[]>;
		onImport: (_scan: LibraryScan, _rows: ImportRow[]) => Promise<void>;
	}
	let { scan, qualityProfiles, loading, onSearchMatch, onImport }: Props = $props();
	let drafts = $state<Record<string, MatchDraft>>({});
	let sortMode = $state<'folders' | 'mixed'>('folders');
	let bulkQualityProfileId = $state('');
	let bulkMonitorMode = $state<MediaMonitorMode>('only_media');
	let bulkMinimumAvailability = $state<MinimumAvailability>('released');
	let importing = $state(false);
	let searchTimers: Record<string, number> = {};
	const searchCache: Record<string, MediaSearchResult[] | undefined> = {};
	const rows = $derived(sortedScanItems(scan.items, sortMode));
	const checkedRows = $derived(
		rows.filter((item) => item.status === 'pending' && drafts[item.id]?.selected)
	);
	const checkedRowsMatched = $derived(
		checkedRows.length > 0 && checkedRows.every((item) => drafts[item.id]?.matched)
	);
	const canImport = $derived(
		checkedRows.length > 0 &&
			checkedRows.every((item) => drafts[item.id]?.matched && drafts[item.id]?.qualityProfileId)
	);

	$effect(() => {
		for (const item of scan.items) {
			if (drafts[item.id]) continue;
			drafts[item.id] = {
				selected: false,
				query: item.matchedTitle ?? item.detectedTitle,
				mediaKind: scanMediaKind(item),
				matched: item.mediaItemId
					? {
							id: item.mediaItemId,
							title: item.matchedTitle ?? item.detectedTitle,
							type:
								item.detectedMediaKind === 'series' || item.detectedMediaKind === 'anime_series'
									? 'series'
									: 'movie',
							year: item.matchedYear ?? item.detectedYear
						}
					: undefined,
				results: [],
				searching: item.status === 'pending',
				searched: item.status !== 'pending',
				qualityProfileId: bulkQualityProfileId,
				monitorMode: bulkMonitorMode,
				minimumAvailability: bulkMinimumAvailability
			};
		}
		void autoSearchPendingRows();
	});
	async function autoSearchPendingRows() {
		for (const item of rows) {
			const draft = drafts[item.id];
			if (!draft || draft.searched || draft.matched || draft.query.trim().length < 2) continue;
			await search(item, true);
			await wait(1100);
		}
	}
	function scheduleSearch(item: LibraryScanItem) {
		window.clearTimeout(searchTimers[item.id]);
		searchTimers[item.id] = window.setTimeout(() => void search(item, false), 1000);
	}
	async function search(item: LibraryScanItem, auto: boolean) {
		const draft = drafts[item.id];
		if (!draft || draft.query.trim().length < 2) return;
		const key = searchCacheKey(draft.mediaKind, draft.query);
		draft.searching = true;
		try {
			const results = searchCache[key] ?? (await onSearchMatch(draft.mediaKind, draft.query));
			searchCache[key] = results;
			draft.results = results;
			draft.searched = true;
			if (auto && results.length > 0) {
				selectResult(item, results[0]);
			}
		} catch {
			draft.results = [];
			draft.searched = true;
		} finally {
			draft.searching = false;
		}
	}
	function selectResult(item: LibraryScanItem, result: MediaSearchResult) {
		const draft = drafts[item.id];
		if (!draft) return;
		draft.matched = result;
		draft.query = result.title;
		draft.results = [];
		draft.selected = true;
	}
	function applyBulk() {
		for (const item of checkedRows) {
			const draft = drafts[item.id];
			if (!draft) continue;
			draft.qualityProfileId = bulkQualityProfileId;
			draft.monitorMode = bulkMonitorMode;
			draft.minimumAvailability = bulkMinimumAvailability;
		}
	}
	async function importChecked() {
		if (!canImport) return;
		applyBulk();
		importing = true;
		try {
			await onImport(
				scan,
				checkedRows.map((item) => {
					const draft = drafts[item.id];
					const match = draft!.matched!;
					return {
						item,
						request: {
							mediaKind: draft!.mediaKind,
							title: match.title,
							year: match.year,
							monitored: true,
							qualityProfileId: draft!.qualityProfileId,
							monitorMode: draft!.monitorMode,
							minimumAvailability: draft!.minimumAvailability,
							externalProvider: match.externalProvider,
							externalId: match.externalId,
							overview: match.overview,
							posterPath: match.posterPath
						}
					};
				})
			);
		} finally {
			importing = false;
		}
	}
</script>

<LibraryScanImportToolbar totalFiles={scan.totalFiles} bind:sortMode />
<div class="table-wrap scan-import-table">
	<table>
		<LibraryScanImportTableHead />
		<tbody>
			{#each rows as item (item.id)}
				{@const draft = drafts[item.id]}
				{#if draft}
					<LibraryScanImportRow
						{item}
						{draft}
						{sortMode}
						{qualityProfiles}
						folderLabel={folderName(item.path)}
						onSearch={scheduleSearch}
						onSelect={selectResult}
					/>
				{/if}
			{/each}
		</tbody>
		<tfoot>
			<LibraryScanImportFooter
				checkedCount={checkedRows.length}
				{checkedRowsMatched}
				{canImport}
				{loading}
				{importing}
				{qualityProfiles}
				bind:qualityProfileId={bulkQualityProfileId}
				bind:monitorMode={bulkMonitorMode}
				bind:minimumAvailability={bulkMinimumAvailability}
				onApply={applyBulk}
				onImport={importChecked}
			/>
		</tfoot>
	</table>
</div>
