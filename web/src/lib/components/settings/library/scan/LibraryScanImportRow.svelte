<script lang="ts">
	import LibraryScanMatchCell from '$lib/components/settings/library/scan/LibraryScanMatchCell.svelte';
	import LibraryScanResetImportButton from '$lib/components/settings/library/scan/LibraryScanResetImportButton.svelte';
	import LibraryScanRowSettings from './row/LibraryScanRowSettings.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Table from '$lib/components/ui/table';
	import {
		isSeriesKind,
		relativeScanPath,
		type MatchDraft
	} from '$lib/components/settings/library/scan/libraryScanImport';
	import type { DuplicateDraftState } from './libraryScanDuplicates';
	import type {
		LibraryScanItem,
		MediaSearchResult,
		MetadataProvider,
		QualityProfileOption
	} from '$lib/settings/types';
	interface Props {
		item: LibraryScanItem;
		folderPath: string;
		draft: MatchDraft;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		duplicateState?: DuplicateDraftState;
		importing?: boolean;
		resetting?: boolean;
		selectionChecked?: boolean;
		selectionIndeterminate?: boolean;
		selectionDisabled?: boolean;
		collapsed?: boolean;
		onToggleSeries?: () => void;
		onSelectionChange?: (_checked: boolean) => void;
		onSeriesSettingChange?: (_field: 'quality' | 'monitor' | 'type', _value: string) => void;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
		onProviderChange: (_item: LibraryScanItem, _providerId: string) => void;
		onResetImport: (_item: LibraryScanItem) => void | Promise<void>;
	}

	let {
		item,
		folderPath,
		draft = $bindable(),
		qualityProfiles,
		metadataProviders,
		duplicateState,
		importing = false,
		resetting = false,
		selectionChecked = false,
		selectionIndeterminate = false,
		selectionDisabled = false,
		collapsed = false,
		onToggleSeries,
		onSelectionChange,
		onSeriesSettingChange,
		onSearch,
		onSelect,
		onProviderChange,
		onResetImport
	}: Props = $props();
	const series = $derived(isSeriesKind(draft.mediaKind));
	const locked = $derived(item.imported || importing || resetting);
	const canRemoveFile = $derived(
		!item.imported && item.status === 'pending' && (!draft.matched || duplicateState?.duplicate)
	);
	const importable = $derived(
		Boolean(draft.matched) && item.status === 'pending' && !item.imported && !draft.removeDuplicate
	);
	const displayPath = $derived(
		relativeScanPath(
			series && item.seriesFolderPath ? item.seriesFolderPath : item.path,
			folderPath
		)
	);
	const selectionActive = $derived(
		series ? selectionChecked || selectionIndeterminate : draft.selected
	);
	function handleRowClick(event: MouseEvent) {
		const target = event.target as { closest?: (_selector: string) => unknown };
		if (series && !target.closest?.('button, input, [role="button"], [role="checkbox"]')) {
			onToggleSeries?.();
		}
	}
	function setRemoveFile(checked: boolean) {
		draft.removeDuplicate = checked;
		if (checked) draft.selected = false;
	}
</script>

<Table.Row
	class={duplicateState?.duplicate && !series ? 'bg-amber-500/5' : undefined}
	onclick={handleRowClick}
>
	<Table.Cell class="w-px align-middle">
		<Checkbox
			checked={series ? selectionChecked : draft.selected}
			indeterminate={series && selectionIndeterminate}
			disabled={series ? selectionDisabled : !importable || locked}
			aria-label={series ? `Select all episodes in ${item.detectedTitle}` : undefined}
			onCheckedChange={(checked) =>
				series ? onSelectionChange?.(checked === true) : (draft.selected = checked === true)}
		/>
	</Table.Cell>
	<Table.Cell class="min-w-120 align-top">
		<div class="flex min-w-0 flex-1 items-start justify-between gap-3">
			<div class="grid min-w-0 flex-1 gap-1">
				<div class="flex min-w-0 items-start gap-2">
					<div class="min-w-0 flex-1">
						<LibraryScanMatchCell {item} bind:draft disabled={locked} {onSearch} {onSelect} />
					</div>
					{#if !series && item.imported}
						<Badge
							variant="outline"
							class="mt-1 border-emerald-500/40 bg-emerald-500/10 text-emerald-400"
						>
							Imported
						</Badge>
					{:else if !series && importing}
						<span class="mt-1 shrink-0">
							<InlineSpinner label="Importing" />
						</span>
					{/if}
				</div>
				<div class="flex ml-3 min-w-0 items-start gap-3">
					{#if series}<button
							type="button"
							class="block truncate text-xs text-muted-foreground hover:text-foreground"
							aria-expanded={!collapsed}
							onclick={onToggleSeries}>{collapsed ? '▸' : '▾'} {displayPath}</button
						>{:else}<span class="block truncate text-xs text-muted-foreground">{displayPath}</span
						>{/if}
					{#if !series && canRemoveFile}
						<label class="flex shrink-0 items-center gap-2 text-xs text-amber-500">
							<Checkbox
								checked={draft.removeDuplicate}
								disabled={Boolean(
									importing ||
									(draft.matched && duplicateState?.duplicate && !duplicateState.removalAllowed)
								)}
								onCheckedChange={(checked) => setRemoveFile(checked === true)}
							/>
							<span>Remove file</span>
						</label>
					{/if}
				</div>
			</div>
			{#if !series && item.imported}
				<LibraryScanResetImportButton {item} {resetting} {onResetImport} />
			{/if}
		</div>
	</Table.Cell>
	<LibraryScanRowSettings
		{item}
		bind:draft
		{qualityProfiles}
		{metadataProviders}
		{locked}
		{selectionActive}
		{onProviderChange}
		{onSeriesSettingChange}
	/>
</Table.Row>
