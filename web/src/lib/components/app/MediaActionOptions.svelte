<script lang="ts">
	import type {
		LibraryFolder,
		MediaMonitorMode,
		MinimumAvailability,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		isAdmin: boolean;
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		qualityProfileId: string;
		libraryFolderId: string;
		monitorMode: MediaMonitorMode;
		minimumAvailability: MinimumAvailability;
	}

	let {
		isAdmin,
		libraryFolders,
		qualityProfiles,
		qualityProfileId = $bindable(),
		libraryFolderId = $bindable(),
		monitorMode = $bindable(),
		minimumAvailability = $bindable()
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
		<select bind:value={monitorMode}>
			<option value="only_media">Only this media</option>
			<option value="collection">Entire collection</option>
		</select>
	</label>

	<label>
		<span>Minimum availability</span>
		<select bind:value={minimumAvailability}>
			<option value="released">Released</option>
			<option value="in_cinema">In cinema</option>
			<option value="announced">Announced</option>
		</select>
	</label>
</div>
