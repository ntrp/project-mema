<script lang="ts">
	import LibraryScanMatchCell from '$lib/components/settings/library/scan/LibraryScanMatchCell.svelte';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Table from '$lib/components/ui/table';
	import {
		isSeriesKind,
		minimumAvailabilityOptions,
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
		draft: MatchDraft;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		duplicateState?: DuplicateDraftState;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
	}

	let {
		item,
		draft = $bindable(),
		qualityProfiles,
		metadataProviders,
		duplicateState,
		onSearch,
		onSelect
	}: Props = $props();
	const series = $derived(isSeriesKind(draft.mediaKind));
	const canRemoveFile = $derived(
		!item.imported && item.status === 'pending' && (!draft.matched || duplicateState?.duplicate)
	);
	const importable = $derived(
		Boolean(draft.matched) && item.status === 'pending' && !item.imported && !draft.removeDuplicate
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

<Table.Row class={duplicateState?.duplicate ? 'bg-amber-500/5' : undefined}>
	<Table.Cell class="w-px align-middle">
		<Checkbox bind:checked={draft.selected} disabled={!importable} />
	</Table.Cell>
	<Table.Cell class="min-w-120 align-top">
		<div class="grid gap-2">
			<div class="flex items-end justify-between gap-3">
				<div class="min-w-0">
					<LibraryScanMatchCell {item} bind:draft {onSearch} {onSelect} />
					<div class="flex ml-3 mt-1 min-w-0 items-start justify-end gap-3">
						<span class="block truncate text-xs text-muted-foreground">{item.path}</span>
						{#if canRemoveFile}
							<label class="flex shrink-0 items-center justify-end gap-2 text-xs text-amber-500">
								<Checkbox
									checked={draft.removeDuplicate}
									disabled={Boolean(draft.matched && duplicateState?.duplicate && !duplicateState.removalAllowed)}
									onCheckedChange={(checked) => setRemoveFile(checked === true)}
								/>
								<span>Remove file</span>
							</label>
						{/if}
					</div>
					{#if item.imported}
						<span class="text-xs font-bold text-muted-foreground">Imported</span>
					{/if}
				</div>
			</div>
		</div>
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		<SettingsSelect
			value={draft.metadataProviderId}
			options={providerOptions}
			disabled={!providerOptions.length || Boolean(draft.matched)}
			onValueChange={(value) => (draft.metadataProviderId = value)}
		/>
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		<SettingsSelect
			value={draft.qualityProfileId}
			options={qualityProfileOptions}
			disabled={!draft.selected || !draft.matched}
			onValueChange={(value) => (draft.qualityProfileId = value)}
		/>
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		<SettingsSelect
			value={draft.monitorMode}
			options={monitorOptions}
			disabled={!draft.selected || !draft.matched}
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
				disabled={!draft.selected || !draft.matched}
				onValueChange={(value) => (draft.minimumAvailability = value as MinimumAvailability)}
			/>
		{/if}
	</Table.Cell>
	<Table.Cell class="w-px align-top">
		{#if series}
			<SettingsSelect
				value={draft.seriesType}
				options={seriesTypeOptions}
				disabled={!draft.selected || !draft.matched}
				onValueChange={(value) => (draft.seriesType = value as SeriesType)}
			/>
		{/if}
	</Table.Cell>
</Table.Row>
