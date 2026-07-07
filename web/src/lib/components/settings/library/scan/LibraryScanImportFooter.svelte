<script lang="ts">
	import LibraryScanImportButton from './LibraryScanImportButton.svelte';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import * as Table from '$lib/components/ui/table';
	import {
		minimumAvailabilityOptions,
		movieMonitorModeOptions,
		seriesMonitorModeOptions,
		seriesTypeOptions
	} from '$lib/components/settings/library/scan/libraryScanImport';
	import type {
		MediaMonitorMode,
		MetadataProvider,
		MinimumAvailability,
		QualityProfileOption,
		SeriesType
	} from '$lib/settings/types';

	interface Props {
		checkedRowsMatched: boolean;
		canImport: boolean;
		loading: boolean;
		importing: boolean;
		duplicateRemovalCount: number;
		metadataProviders: MetadataProvider[];
		qualityProfiles: QualityProfileOption[];
		hasMatchedMovies: boolean;
		hasMatchedSeries: boolean;
		showSeriesControls?: boolean;
		metadataProviderId: string;
		qualityProfileId: string;
		movieMonitorMode: MediaMonitorMode;
		movieMinimumAvailability: MinimumAvailability;
		seriesMonitorMode: MediaMonitorMode;
		seriesType: SeriesType;
		onApplyProvider: () => void;
		onApplyQualityProfile: () => void;
		onApplyMovie: () => void;
		onApplySeries: () => void;
		onImport: () => void;
	}

	let {
		checkedRowsMatched,
		canImport,
		loading,
		importing,
		duplicateRemovalCount,
		metadataProviders,
		qualityProfiles,
		hasMatchedMovies,
		hasMatchedSeries,
		showSeriesControls = true,
		metadataProviderId = $bindable(),
		qualityProfileId = $bindable(),
		movieMonitorMode = $bindable(),
		movieMinimumAvailability = $bindable(),
		seriesMonitorMode = $bindable(),
		seriesType = $bindable(),
		onApplyProvider,
		onApplyQualityProfile,
		onApplyMovie,
		onApplySeries,
		onImport
	}: Props = $props();

	let monitorSelection = $state('');
	const providerOptions = $derived(
		metadataProviders
			.filter((provider) => provider.enabled)
			.map((provider) => ({ value: provider.id, label: provider.name }))
	);
	const qualityProfileOptions = $derived([
		{ value: '', label: 'Select profile' },
		...qualityProfiles.map((profile) => ({ value: profile.id, label: profile.name }))
	]);
	const monitorOptions = $derived(
		showSeriesControls
			? [
					{ value: 'movie-header', label: 'Movie', disabled: true },
					...movieMonitorModeOptions.map((option) => ({
						value: `movie:${option.value}`,
						label: option.label
					})),
					{ value: 'series-header', label: 'Series', disabled: true },
					...seriesMonitorModeOptions.map((option) => ({
						value: `series:${option.value}`,
						label: option.label
					}))
				]
			: movieMonitorModeOptions.map((option) => ({
					value: `movie:${option.value}`,
					label: option.label
				}))
	);

	function applyQualityProfile(value: string) {
		qualityProfileId = value;
		onApplyQualityProfile();
	}

	function applyMetadataProvider(value: string) {
		metadataProviderId = value;
		onApplyProvider();
	}

	function applyMovieMinimumAvailability(value: string) {
		movieMinimumAvailability = value as MinimumAvailability;
		onApplyMovie();
	}

	function applyMonitorMode(value: string) {
		monitorSelection = value;
		const [section, monitorMode] = value.split(':') as ['movie' | 'series', MediaMonitorMode];
		if (section === 'movie') {
			movieMonitorMode = monitorMode;
			void movieMonitorMode;
			onApplyMovie();
		} else if (section === 'series') {
			seriesMonitorMode = monitorMode;
			void seriesMonitorMode;
			onApplySeries();
		}
	}

	function applySeriesType(value: string) {
		seriesType = value as SeriesType;
		onApplySeries();
	}
</script>

<Table.Footer class="sticky bottom-0 bg-muted shadow-sm">
	<Table.Row>
		<Table.Cell class="align-top" colspan={2}>
			<LibraryScanImportButton
				{canImport}
				{loading}
				{importing}
				{duplicateRemovalCount}
				{onImport}
			/>
		</Table.Cell>
		<Table.Cell class="w-px align-top">
			<SettingsSelect
				value={metadataProviderId}
				options={providerOptions}
				disabled={!checkedRowsMatched || !providerOptions.length}
				placeholder="Apply provider"
				onValueChange={applyMetadataProvider}
			/>
		</Table.Cell>
		<Table.Cell class="w-px align-top">
			<SettingsSelect
				value={qualityProfileId}
				options={qualityProfileOptions}
				disabled={!checkedRowsMatched}
				onValueChange={applyQualityProfile}
			/>
		</Table.Cell>
		<Table.Cell class="w-px align-top">
			<SettingsSelect
				value={monitorSelection}
				options={monitorOptions}
				disabled={!checkedRowsMatched}
				placeholder="Apply monitor"
				onValueChange={applyMonitorMode}
			/>
		</Table.Cell>
		<Table.Cell class="w-px align-top">
			<div class="grid gap-1.5">
				<SettingsSelect
					value={movieMinimumAvailability}
					options={minimumAvailabilityOptions}
					disabled={!checkedRowsMatched || !hasMatchedMovies}
					onValueChange={applyMovieMinimumAvailability}
				/>
			</div>
		</Table.Cell>
		<Table.Cell class="w-px align-top">
			{#if showSeriesControls}
				<SettingsSelect
					value={seriesType}
					options={seriesTypeOptions}
					disabled={!checkedRowsMatched || !hasMatchedSeries}
					onValueChange={applySeriesType}
				/>
			{/if}
		</Table.Cell>
	</Table.Row>
</Table.Footer>
