<script lang="ts">
	import LibraryScanEpisodeRow from '$lib/components/settings/library/scan/LibraryScanEpisodeRow.svelte';
	import LibraryScanImportRow from '$lib/components/settings/library/scan/LibraryScanImportRow.svelte';
	import {
		isDraftSelectableForImport,
		isScanItemGroupRepresentative,
		scanItemGroupRows,
		type MatchDraft
	} from './libraryScanImport';
	import type { DuplicateDraftState } from './libraryScanDuplicates';
	import * as Table from '$lib/components/ui/table';
	import type {
		LibraryScanItem,
		MediaSearchResult,
		MetadataProvider,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		rows: LibraryScanItem[];
		folderPath: string;
		drafts: Record<string, MatchDraft>;
		duplicateStates: Record<string, DuplicateDraftState>;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		importingItemId: string;
		resettingItemId: string;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
		onProviderChange: (_item: LibraryScanItem, _providerId: string) => void;
		onSeasonChange: (_item: LibraryScanItem, _value: number | undefined) => void;
		onEpisodeChange: (_item: LibraryScanItem, _value: number | undefined) => void;
		onResetImport: (_item: LibraryScanItem) => void | Promise<void>;
	}

	let {
		rows,
		folderPath,
		drafts = $bindable(),
		duplicateStates,
		qualityProfiles,
		metadataProviders,
		importingItemId,
		resettingItemId,
		onSearch,
		onSelect,
		onProviderChange,
		onSeasonChange,
		onEpisodeChange,
		onResetImport
	}: Props = $props();
	let collapsedFolders = $state<Record<string, boolean>>({});
	let dragSelecting = $state(false);
	let dragValue = false;
	let dragOrigin = '';
	let suppressedClick = '';

	function seriesCollapsed(item: LibraryScanItem) {
		return Boolean(item.seriesFolderPath && collapsedFolders[item.seriesFolderPath]);
	}

	function toggleSeries(item: LibraryScanItem) {
		if (!item.seriesFolderPath) return;
		collapsedFolders[item.seriesFolderPath] = !collapsedFolders[item.seriesFolderPath];
	}

	function selectable(item: LibraryScanItem) {
		return (
			isDraftSelectableForImport(item, drafts[item.id]) &&
			importingItemId !== item.id &&
			resettingItemId !== item.id
		);
	}

	function seriesSelection(item: LibraryScanItem) {
		const group = scanItemGroupRows(item, rows).filter(selectable);
		const selected = group.filter((row) => drafts[row.id]?.selected).length;
		return {
			checked: group.length > 0 && selected === group.length,
			indeterminate: selected > 0 && selected < group.length,
			disabled: group.length === 0
		};
	}

	function selectSeries(item: LibraryScanItem, selected: boolean) {
		for (const row of scanItemGroupRows(item, rows)) {
			if (selectable(row)) drafts[row.id].selected = selected;
		}
	}

	function updateSeriesSetting(
		item: LibraryScanItem,
		field: 'quality' | 'monitor' | 'type',
		value: string
	) {
		for (const row of scanItemGroupRows(item, rows)) {
			const draft = drafts[row.id];
			if (!draft) continue;
			if (field === 'quality') draft.qualityProfileId = value;
			else if (field === 'monitor') draft.monitorMode = value as MatchDraft['monitorMode'];
			else draft.seriesType = value as MatchDraft['seriesType'];
		}
	}

	function startSelection(item: LibraryScanItem, event: PointerEvent) {
		if (event.button !== 0 || !selectable(item)) return;
		event.preventDefault();
		dragSelecting = true;
		dragValue = !drafts[item.id].selected;
		dragOrigin = item.id;
		drafts[item.id].selected = dragValue;
	}

	function enterSelection(item: LibraryScanItem) {
		if (dragSelecting && selectable(item)) drafts[item.id].selected = dragValue;
	}

	function endSelection() {
		if (!dragSelecting) return;
		dragSelecting = false;
		suppressedClick = dragOrigin;
		globalThis.setTimeout(() => (suppressedClick = ''), 0);
	}

	function clickSelection(item: LibraryScanItem, event: MouseEvent) {
		if (suppressedClick === item.id) {
			event.preventDefault();
			suppressedClick = '';
			return;
		}
		if (selectable(item)) drafts[item.id].selected = !drafts[item.id].selected;
	}
</script>

<svelte:window onpointerup={endSelection} onpointercancel={endSelection} />

<Table.Body>
	{#each rows as item (item.id)}
		{#if drafts[item.id]}
			{#if isScanItemGroupRepresentative(item, rows)}
				{@const selection = seriesSelection(item)}
				<LibraryScanImportRow
					{item}
					{folderPath}
					bind:draft={drafts[item.id]}
					{qualityProfiles}
					{metadataProviders}
					duplicateState={duplicateStates[item.id]}
					importing={importingItemId === item.id}
					resetting={resettingItemId === item.id}
					selectionChecked={selection.checked}
					selectionIndeterminate={selection.indeterminate}
					selectionDisabled={selection.disabled}
					collapsed={seriesCollapsed(item)}
					onToggleSeries={() => toggleSeries(item)}
					onSelectionChange={(checked) => selectSeries(item, checked)}
					onSeriesSettingChange={(field, value) => updateSeriesSetting(item, field, value)}
					{onSearch}
					{onSelect}
					{onProviderChange}
					{onResetImport}
				/>
			{/if}
			{#if item.seriesFolderPath && !seriesCollapsed(item)}
				<LibraryScanEpisodeRow
					{item}
					{folderPath}
					bind:draft={drafts[item.id]}
					duplicateState={duplicateStates[item.id]}
					importing={importingItemId === item.id}
					resetting={resettingItemId === item.id}
					onSelectionStart={startSelection}
					onSelectionEnter={enterSelection}
					onSelectionClick={clickSelection}
					{onSeasonChange}
					{onEpisodeChange}
					{onResetImport}
				/>
			{/if}
		{/if}
	{:else}
		<Table.Row>
			<Table.Cell colspan={7} class="text-muted-foreground">No files to import.</Table.Cell>
		</Table.Row>
	{/each}
</Table.Body>
