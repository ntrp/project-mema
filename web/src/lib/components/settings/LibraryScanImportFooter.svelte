<script lang="ts">
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Button } from '$lib/components/ui/button';
	import {
		minimumAvailabilityOptions,
		monitorModeOptions
	} from '$lib/components/settings/libraryScanImport';
	import type {
		MediaMonitorMode,
		MinimumAvailability,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		checkedCount: number;
		checkedRowsMatched: boolean;
		canImport: boolean;
		loading: boolean;
		importing: boolean;
		qualityProfiles: QualityProfileOption[];
		qualityProfileId: string;
		monitorMode: MediaMonitorMode;
		minimumAvailability: MinimumAvailability;
		onApply: () => void;
		onImport: () => void;
	}

	let {
		checkedCount,
		checkedRowsMatched,
		canImport,
		loading,
		importing,
		qualityProfiles,
		qualityProfileId = $bindable(),
		monitorMode = $bindable(),
		minimumAvailability = $bindable(),
		onApply,
		onImport
	}: Props = $props();

	let qualityProfileOptions = $derived([
		{ value: '', label: 'Select profile' },
		...qualityProfiles.map((profile) => ({ value: profile.id, label: profile.name }))
	]);

	function applyQualityProfile(value: string) {
		qualityProfileId = value;
		onApply();
	}

	function applyMonitorMode(value: string) {
		monitorMode = value as MediaMonitorMode;
		onApply();
	}

	function applyMinimumAvailability(value: string) {
		minimumAvailability = value as MinimumAvailability;
		onApply();
	}
</script>

<tr class="sticky bottom-0 bg-muted shadow-sm">
	<td class="align-top" colspan="3">{checkedCount} checked</td>
	<td class="align-top">
		<SettingsSelect
			value={qualityProfileId}
			options={qualityProfileOptions}
			disabled={!checkedRowsMatched}
			onValueChange={applyQualityProfile}
		/>
	</td>
	<td class="align-top">
		<SettingsSelect
			value={monitorMode}
			options={monitorModeOptions}
			disabled={!checkedRowsMatched}
			onValueChange={applyMonitorMode}
		/>
	</td>
	<td class="align-top">
		<SettingsSelect
			value={minimumAvailability}
			options={minimumAvailabilityOptions}
			disabled={!checkedRowsMatched}
			onValueChange={applyMinimumAvailability}
		/>
		<Button
			type="button"
			class="mt-2 w-full"
			disabled={!canImport || loading || importing}
			onclick={onImport}
		>
			{importing ? 'Importing' : 'Import'}
		</Button>
	</td>
</tr>
