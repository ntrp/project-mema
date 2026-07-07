<script lang="ts">
	import LibraryScanImportFooter from '$lib/components/settings/library/scan/LibraryScanImportFooter.svelte';
	import LibraryScanImportTableBody from '$lib/components/settings/library/scan/LibraryScanImportTableBody.svelte';
	import LibraryScanImportTableHead from '$lib/components/settings/library/scan/LibraryScanImportTableHead.svelte';
	import LibraryScanImportToolbar from '$lib/components/settings/library/scan/LibraryScanImportToolbar.svelte';
	import * as Table from '$lib/components/ui/table';
	import type { DuplicateDraftState } from './libraryScanDuplicates';
	import type { MatchDraft } from './libraryScanImport';
	import type {
		LibraryScanItem,
		MediaMonitorMode,
		MediaSearchResult,
		MetadataProvider,
		MinimumAvailability,
		QualityProfileOption,
		SeriesType
	} from '$lib/settings/types';

	interface Props {
		totalFiles: number;
		matchedCount: number;
		noMatchCount: number;
		importedCount: number;
		duplicateCount: number;
		rows: LibraryScanItem[];
		folderPath: string;
		drafts: Record<string, MatchDraft>;
		duplicateStates: Record<string, DuplicateDraftState>;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		importingItemId: string;
		showImported: boolean;
		allVisibleChecked: boolean;
		importableCount: number;
		showSeriesControls: boolean;
		checkedRowsMatched: boolean;
		canImport: boolean;
		loading: boolean;
		importing: boolean;
		duplicateRemovalCount: number;
		hasMatchedMovies: boolean;
		hasMatchedSeries: boolean;
		metadataProviderId: string;
		qualityProfileId: string;
		movieMonitorMode: MediaMonitorMode;
		movieMinimumAvailability: MinimumAvailability;
		seriesMonitorMode: MediaMonitorMode;
		seriesType: SeriesType;
		onToggleRows: () => void;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
		onProviderChange: (_item: LibraryScanItem, _providerId: string) => void;
		onApplyProvider: () => void;
		onApplyQualityProfile: () => void;
		onApplyMovie: () => void;
		onApplySeries: () => void;
		onImport: () => void;
	}

	let {
		totalFiles,
		matchedCount,
		noMatchCount,
		importedCount,
		duplicateCount,
		rows,
		folderPath,
		drafts = $bindable(),
		duplicateStates,
		qualityProfiles,
		metadataProviders,
		importingItemId,
		showImported = $bindable(),
		allVisibleChecked,
		importableCount,
		showSeriesControls,
		checkedRowsMatched,
		canImport,
		loading,
		importing,
		duplicateRemovalCount,
		hasMatchedMovies,
		hasMatchedSeries,
		metadataProviderId = $bindable(),
		qualityProfileId = $bindable(),
		movieMonitorMode = $bindable(),
		movieMinimumAvailability = $bindable(),
		seriesMonitorMode = $bindable(),
		seriesType = $bindable(),
		onToggleRows,
		onSearch,
		onSelect,
		onProviderChange,
		onApplyProvider,
		onApplyQualityProfile,
		onApplyMovie,
		onApplySeries,
		onImport
	}: Props = $props();
</script>

<LibraryScanImportToolbar
	{totalFiles}
	{matchedCount}
	{noMatchCount}
	{importedCount}
	{duplicateCount}
	bind:showImported
/>
<div class="mt-4 overflow-auto rounded-md border border-border">
	<Table.Root class="min-w-7xl table-auto border-collapse">
		<colgroup>
			<col class="w-[1%]" />
			<col class="w-full" />
			<col span={5} class="w-[1%]" />
		</colgroup>
		<LibraryScanImportTableHead
			checked={allVisibleChecked}
			disabled={importableCount === 0}
			{showSeriesControls}
			onToggle={onToggleRows}
		/>
		<LibraryScanImportTableBody
			{rows}
			{folderPath}
			bind:drafts
			{duplicateStates}
			{qualityProfiles}
			{metadataProviders}
			{importingItemId}
			{onSearch}
			{onSelect}
			{onProviderChange}
		/>
		<LibraryScanImportFooter
			{checkedRowsMatched}
			{canImport}
			{loading}
			{importing}
			{duplicateRemovalCount}
			{metadataProviders}
			{qualityProfiles}
			{hasMatchedMovies}
			{hasMatchedSeries}
			{showSeriesControls}
			bind:metadataProviderId
			bind:qualityProfileId
			bind:movieMonitorMode
			bind:movieMinimumAvailability
			bind:seriesMonitorMode
			bind:seriesType
			{onApplyProvider}
			{onApplyQualityProfile}
			{onApplyMovie}
			{onApplySeries}
			{onImport}
		/>
	</Table.Root>
</div>
