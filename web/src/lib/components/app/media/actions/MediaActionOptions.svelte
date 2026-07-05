<script lang="ts">
	import { matchingLibraryFolders } from '$lib/components/app/media/actions/mediaActionDefaults';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Label } from '$lib/components/ui/label';
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

	let matchingFolders = $derived(matchingLibraryFolders(mediaType, libraryFolders));
	let libraryFolderOptions = $derived([
		{ value: '', label: 'Select folder' },
		...matchingFolders.map((folder) => ({ value: folder.id, label: folder.path }))
	]);
	let qualityProfileOptions = $derived([
		{ value: '', label: 'Select profile' },
		...qualityProfiles.map((profile) => ({ value: profile.id, label: profile.name }))
	]);
	let monitorModeOptions = $derived([
		...(mediaType === 'serie'
			? [
					{ value: 'all_episodes', label: 'All episodes' },
					{ value: 'future_episodes', label: 'Future episodes' },
					{ value: 'missing_episodes', label: 'Missing episodes' },
					{ value: 'existing_episodes', label: 'Existing episodes' },
					{ value: 'no_specials', label: 'No specials' }
				]
			: [
					{ value: 'only_media', label: 'Only this media' },
					{ value: 'collection', label: 'Entire collection' }
				]),
		{ value: 'none', label: 'None' }
	]);
	const seriesTypeOptions: { value: SeriesType; label: string }[] = [
		{ value: 'standard', label: 'Standard' },
		{ value: 'daily', label: 'Daily / Date' },
		{ value: 'absolute', label: 'Absolute' }
	];
	const availabilityOptions: { value: MinimumAvailability; label: string }[] = [
		{ value: 'released', label: 'Released' },
		{ value: 'in_cinema', label: 'In cinema' },
		{ value: 'announced', label: 'Announced' }
	];
</script>

<div class="grid gap-4">
	{#if isAdmin}
		<div class="grid gap-2">
			<Label>Library folder</Label>
			<SettingsSelect
				value={libraryFolderId}
				options={libraryFolderOptions}
				onValueChange={(value) => (libraryFolderId = value)}
			/>
		</div>
		<div class="grid gap-2">
			<Label>Quality profile</Label>
			<SettingsSelect
				value={qualityProfileId}
				options={qualityProfileOptions}
				onValueChange={(value) => (qualityProfileId = value)}
			/>
		</div>
	{/if}

	<div class="grid gap-2">
		<Label>Monitor</Label>
		<SettingsSelect
			value={monitorMode}
			options={monitorModeOptions}
			onValueChange={(value) => onMonitorModeChange(value as MediaMonitorMode)}
		/>
	</div>

	{#if mediaType === 'serie'}
		<div class="grid gap-2">
			<Label>Series type</Label>
			<SettingsSelect
				value={seriesType}
				options={seriesTypeOptions}
				onValueChange={(value) => onSeriesTypeChange(value as SeriesType)}
			/>
		</div>
	{/if}

	<div class="grid gap-2">
		<Label>Minimum availability</Label>
		<SettingsSelect
			value={minimumAvailability}
			options={availabilityOptions}
			onValueChange={(value) => (minimumAvailability = value as MinimumAvailability)}
		/>
	</div>
</div>
