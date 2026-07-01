<script lang="ts">
	/* global HTMLSelectElement */
	import type {
		LibraryFolder,
		MediaMonitorMode,
		MediaType,
		MinimumAvailability,
		QualityProfileOption,
		SeriesType
	} from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
		isAdmin: boolean;
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		qualityProfileId: string;
		libraryFolderId: string;
		monitorMode: MediaMonitorMode;
		seriesType: SeriesType;
		minimumAvailability: MinimumAvailability;
		onMonitorModeChange: (_mode: MediaMonitorMode) => void;
		onSeriesTypeChange: (_type: SeriesType) => void;
	}

	let {
		mediaType,
		isAdmin,
		libraryFolders,
		qualityProfiles,
		qualityProfileId = $bindable(),
		libraryFolderId = $bindable(),
		monitorMode,
		seriesType,
		minimumAvailability = $bindable(),
		onMonitorModeChange,
		onSeriesTypeChange
	}: Props = $props();
</script>

<div class="settings-form compact-form media-action-fields">
	{#if isAdmin}
		<label>
			<span>Library folder</span>
			<select bind:value={libraryFolderId}>
				<option value="" disabled>Select folder</option>
				{#each libraryFolders as folder (folder.id)}
					<option value={folder.id}>{folder.path}</option>
				{/each}
			</select>
		</label>
		<label>
			<span>Quality profile</span>
			<select bind:value={qualityProfileId}>
				<option value="" disabled>Select profile</option>
				{#each qualityProfiles as profile (profile.id)}
					<option value={profile.id}>{profile.name}</option>
				{/each}
			</select>
		</label>
	{/if}

	<label>
		<span>Monitor</span>
		<select
			value={monitorMode}
			onchange={(event) =>
				onMonitorModeChange((event.currentTarget as HTMLSelectElement).value as MediaMonitorMode)}
		>
			{#if mediaType === 'series'}
				<option value="all_episodes">All episodes</option>
				<option value="future_episodes">Future episodes</option>
				<option value="missing_episodes">Missing episodes</option>
				<option value="existing_episodes">Existing episodes</option>
				<option value="no_specials">No specials</option>
			{:else}
				<option value="only_media">Only this media</option>
				<option value="collection">Entire collection</option>
			{/if}
			<option value="none">None</option>
		</select>
	</label>

	{#if mediaType === 'series'}
		<label>
			<span>Series type</span>
			<select
				value={seriesType}
				onchange={(event) =>
					onSeriesTypeChange((event.currentTarget as HTMLSelectElement).value as SeriesType)}
			>
				<option value="standard">Standard</option>
				<option value="daily">Daily / Date</option>
				<option value="absolute">Absolute</option>
			</select>
		</label>
	{/if}

	<label>
		<span>Minimum availability</span>
		<select bind:value={minimumAvailability}>
			<option value="released">Released</option>
			<option value="in_cinema">In cinema</option>
			<option value="announced">Announced</option>
		</select>
	</label>
</div>
