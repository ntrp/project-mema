<script lang="ts">
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import {
		minimumAvailabilityOptions,
		movieMonitorModeOptions,
		seriesMonitorModeOptions,
		seriesTypeOptions
	} from '$lib/components/settings/library/scan/libraryScanImport';
	import type {
		MediaMonitorMode,
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
		qualityProfiles: QualityProfileOption[];
		hasMatchedMovies: boolean;
		hasMatchedSeries: boolean;
		showSeriesControls?: boolean;
		qualityProfileId: string;
		movieMonitorMode: MediaMonitorMode;
		movieMinimumAvailability: MinimumAvailability;
		seriesMonitorMode: MediaMonitorMode;
		seriesType: SeriesType;
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
		qualityProfiles,
		hasMatchedMovies,
		hasMatchedSeries,
		showSeriesControls = true,
		qualityProfileId = $bindable(),
		movieMonitorMode = $bindable(),
		movieMinimumAvailability = $bindable(),
		seriesMonitorMode = $bindable(),
		seriesType = $bindable(),
		onApplyQualityProfile,
		onApplyMovie,
		onApplySeries,
		onImport
	}: Props = $props();

	let monitorSelection = $state('');
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

	function applyMovieMinimumAvailability(value: string) {
		movieMinimumAvailability = value as MinimumAvailability;
		onApplyMovie();
	}

	function applyMonitorMode(value: string) {
		monitorSelection = value;
		const [section, monitorMode] = value.split(':') as ['movie' | 'series', MediaMonitorMode];
		if (section === 'movie') {
			movieMonitorMode = monitorMode;
			onApplyMovie();
		} else if (section === 'series') {
			seriesMonitorMode = monitorMode;
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
			{#if duplicateRemovalCount > 0}
				<ConfirmActionButton
					label="Import selected"
					title="Remove files"
					description={`Import selected rows and remove ${duplicateRemovalCount} file${duplicateRemovalCount === 1 ? '' : 's'}?`}
					confirmLabel="Import selected"
					confirmingLabel="Importing"
					variant="default"
					class="whitespace-nowrap"
					disabled={!canImport || loading || importing}
					onConfirm={onImport}
				>
					{importing ? 'Importing' : 'Import Selected'}
				</ConfirmActionButton>
			{:else}
				<Button
					type="button"
					class="whitespace-nowrap"
					disabled={!canImport || loading || importing}
					onclick={onImport}
				>
					{importing ? 'Importing' : 'Import Selected'}
				</Button>
			{/if}
		</Table.Cell>
		<Table.Cell class="w-px align-top"></Table.Cell>
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
