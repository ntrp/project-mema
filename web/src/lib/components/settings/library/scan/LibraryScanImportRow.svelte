<script lang="ts">
	import LibraryScanMatchCell from '$lib/components/settings/library/scan/LibraryScanMatchCell.svelte';
	import LibraryScanResetImportButton from '$lib/components/settings/library/scan/LibraryScanResetImportButton.svelte';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Table from '$lib/components/ui/table';
	import {
		isSeriesKind,
		minimumAvailabilityOptions,
		relativeScanPath,
		movieMonitorModeOptions,
		seriesMonitorModeOptions,
		seriesTypeOptions,
		type MatchDraft
	} from '$lib/components/settings/library/scan/libraryScanImport';
	import type { DuplicateDraftState } from './libraryScanDuplicates';
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
		item: LibraryScanItem;
		folderPath: string;
		draft: MatchDraft;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		duplicateState?: DuplicateDraftState;
		importing?: boolean;
		resetting?: boolean;
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
	const qualityProfileOptions = $derived([
		{ value: '', label: 'Select profile' },
		...qualityProfiles.map((profile) => ({ value: profile.id, label: profile.name }))
	]);
	const providerOptions = $derived(
		metadataProviders
			.filter((provider) => provider.enabled)
			.map((provider) => ({ value: provider.id, label: provider.name }))
	);
	const monitorOptions = $derived(series ? seriesMonitorModeOptions : movieMonitorModeOptions);

	function setRemoveFile(checked: boolean) {
		draft.removeDuplicate = checked;
		if (checked) draft.selected = false;
	}
</script>

<Table.Row class={duplicateState?.duplicate && !series ? 'bg-amber-500/5' : undefined}>
	<Table.Cell class="w-px align-middle">
		{#if !series}<Checkbox bind:checked={draft.selected} disabled={!importable || locked} />{/if}
	</Table.Cell>
	<Table.Cell class="min-w-120 align-top">
		<div class="flex items-end justify-between gap-3">
			<div class="min-w-0 flex-1">
				<div class="flex min-w-0 items-start justify-between gap-3">
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
							<span class="block truncate text-xs text-muted-foreground">{displayPath}</span>
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
			</div>
		</div>
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		<SettingsSelect
			value={draft.metadataProviderId}
			options={providerOptions}
			disabled={!providerOptions.length || locked}
			onValueChange={(value) => onProviderChange(item, value)}
		/>
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		<SettingsSelect
			value={draft.qualityProfileId}
			options={qualityProfileOptions}
			disabled={(!series && !draft.selected) || !draft.matched || locked}
			onValueChange={(value) => (draft.qualityProfileId = value)}
		/>
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		<SettingsSelect
			value={draft.monitorMode}
			options={monitorOptions}
			disabled={(!series && !draft.selected) || !draft.matched || locked}
			onValueChange={(value) => (draft.monitorMode = value as MediaMonitorMode)}
		/>
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		{#if series}
			<span class="text-xs text-muted-foreground">-</span>
		{:else}
			<SettingsSelect
				value={draft.minimumAvailability}
				options={minimumAvailabilityOptions}
				disabled={!draft.selected || !draft.matched || locked}
				onValueChange={(value) => (draft.minimumAvailability = value as MinimumAvailability)}
			/>
		{/if}
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		{#if series}
			<SettingsSelect
				value={draft.seriesType}
				options={seriesTypeOptions}
				disabled={!draft.matched || locked}
				onValueChange={(value) => (draft.seriesType = value as SeriesType)}
			/>
		{/if}
	</Table.Cell>
</Table.Row>
