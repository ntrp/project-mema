<script lang="ts">
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import * as Table from '$lib/components/ui/table';
	import {
		isSeriesKind,
		minimumAvailabilityOptions,
		movieMonitorModeOptions,
		seriesMonitorModeOptions,
		seriesTypeOptions,
		type MatchDraft
	} from '../libraryScanImport';
	import type {
		LibraryScanItem,
		MediaMonitorMode,
		MetadataProvider,
		MinimumAvailability,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		item: LibraryScanItem;
		draft: MatchDraft;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		locked: boolean;
		selectionActive: boolean;
		onProviderChange: (_item: LibraryScanItem, _providerId: string) => void;
		onSeriesSettingChange?: (_field: 'quality' | 'monitor' | 'type', _value: string) => void;
	}

	let {
		item,
		draft = $bindable(),
		qualityProfiles,
		metadataProviders,
		locked,
		selectionActive,
		onProviderChange,
		onSeriesSettingChange
	}: Props = $props();
	const series = $derived(isSeriesKind(draft.mediaKind));
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
</script>

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
		disabled={!selectionActive || !draft.matched || locked}
		onValueChange={(value) =>
			series ? onSeriesSettingChange?.('quality', value) : (draft.qualityProfileId = value)}
	/>
</Table.Cell>
<Table.Cell class="w-px align-top">
	<SettingsSelect
		value={draft.monitorMode}
		options={monitorOptions}
		disabled={!selectionActive || !draft.matched || locked}
		onValueChange={(value) =>
			series
				? onSeriesSettingChange?.('monitor', value)
				: (draft.monitorMode = value as MediaMonitorMode)}
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
			disabled={!selectionActive || !draft.matched || locked}
			onValueChange={(value) => onSeriesSettingChange?.('type', value)}
		/>
	{/if}
</Table.Cell>
