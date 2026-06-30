<script lang="ts">
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
</script>

<tr class="scan-import-row">
	<td colspan="3">{checkedCount} checked</td>
	<td>
		<select bind:value={qualityProfileId} disabled={!checkedRowsMatched} onchange={onApply}>
			<option value="">Select profile</option>
			{#each qualityProfiles as profile (profile.id)}
				<option value={profile.id}>{profile.name}</option>
			{/each}
		</select>
	</td>
	<td>
		<select bind:value={monitorMode} disabled={!checkedRowsMatched} onchange={onApply}>
			<option value="only_media">Only this media</option>
			<option value="collection">Entire collection</option>
		</select>
	</td>
	<td>
		<select bind:value={minimumAvailability} disabled={!checkedRowsMatched} onchange={onApply}>
			<option value="released">Released</option>
			<option value="in_cinema">In cinema</option>
			<option value="announced">Announced</option>
		</select>
		<button type="button" disabled={!canImport || loading || importing} onclick={onImport}>
			{importing ? 'Importing' : 'Import'}
		</button>
	</td>
</tr>
