<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import XIcon from '@lucide/svelte/icons/x';
	import MediaActionOptions from '$lib/components/app/media/actions/MediaActionOptions.svelte';
	import MediaTagSelector from '$lib/components/app/media/actions/MediaTagSelector.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import {
		preselectLibraryFolderId,
		preselectQualityProfileId
	} from '$lib/components/app/media/actions/mediaActionDefaults';
	import type { MediaActionSelection } from '$lib/components/app/media/actions/mediaActionTypes';
	import type {
		LibraryFolder,
		MediaSearchResult,
		MediaMonitorMode,
		MinimumAvailability,
		QualityProfileOption,
		SeriesType,
		Tag
	} from '$lib/settings/types';

	interface Props {
		candidate: MediaSearchResult;
		isAdmin: boolean;
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		tags: Tag[];
		saving: boolean;
		onClose: () => void;
		onConfirm: (_selection: MediaActionSelection) => void;
	}

	let {
		candidate,
		isAdmin,
		libraryFolders,
		qualityProfiles,
		tags,
		saving,
		onClose,
		onConfirm
	}: Props = $props();

	let qualityProfileId = $state('');
	let libraryFolderId = $state('');
	let selectedMonitorMode = $state<MediaMonitorMode | undefined>();
	let selectedSeriesType = $state<SeriesType | undefined>();
	let minimumAvailability = $state<MinimumAvailability>('released');
	let startSearch = $state(true);
	let selectedTags = $state<string[]>([]);

	const canConfirm = $derived(!isAdmin || (qualityProfileId !== '' && libraryFolderId !== ''));
	const monitorMode = $derived(selectedMonitorMode ?? defaultMonitorMode(candidate.type));
	const seriesType = $derived(selectedSeriesType ?? 'standard');

	$effect(() => {
		if (qualityProfileId === '') {
			qualityProfileId = preselectQualityProfileId(candidate, qualityProfiles);
		}
		if (libraryFolderId === '') {
			libraryFolderId = preselectLibraryFolderId(candidate, libraryFolders);
		}
	});

	function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!canConfirm || saving) {
			return;
		}
		onConfirm({
			qualityProfileId,
			libraryFolderId,
			tags: selectedTags,
			monitorMode,
			seriesType,
			minimumAvailability,
			startSearch
		});
	}

	function defaultMonitorMode(type: MediaSearchResult['type']): MediaMonitorMode {
		return type === 'series' ? 'all_episodes' : 'only_media';
	}
</script>

<SettingsFormModal title={candidate.title} {onClose}>
	<form class="-m-5 grid gap-3.5 overflow-hidden rounded-md bg-card p-5" onsubmit={submit}>
		<p class="m-0 mb-1.5 text-xs font-extrabold text-muted-foreground uppercase">
			{isAdmin ? 'Add media' : 'Request media'}
		</p>
		<p class="m-0 text-sm leading-6 text-muted-foreground">
			{candidate.type}{candidate.year ? ` · ${candidate.year}` : ''}
		</p>

		{#if isAdmin}
			<MediaActionOptions
				{isAdmin}
				mediaType={candidate.type}
				{libraryFolders}
				{qualityProfiles}
				bind:qualityProfileId
				bind:libraryFolderId
				{monitorMode}
				{seriesType}
				bind:minimumAvailability
				onMonitorModeChange={(mode) => (selectedMonitorMode = mode)}
				onSeriesTypeChange={(type) => (selectedSeriesType = type)}
			/>
			{#if libraryFolders.length === 0}
				<p class="m-0 text-sm leading-6 font-semibold text-destructive">
					Add a library folder in Settings before adding monitored media.
				</p>
			{/if}
		{:else}
			<MediaActionOptions
				{isAdmin}
				mediaType={candidate.type}
				{libraryFolders}
				{qualityProfiles}
				bind:qualityProfileId
				bind:libraryFolderId
				{monitorMode}
				{seriesType}
				bind:minimumAvailability
				onMonitorModeChange={(mode) => (selectedMonitorMode = mode)}
				onSeriesTypeChange={(type) => (selectedSeriesType = type)}
			/>
			<p class="m-0 text-sm leading-6 text-muted-foreground">
				Your request will be visible under Requests. An admin will choose the folder and quality
				profile before approval.
			</p>
		{/if}

		<MediaTagSelector {tags} bind:selectedTags />

		<div class="flex flex-wrap items-center justify-end gap-3">
			{#if isAdmin}
				<div class="mr-auto flex items-center gap-2">
					<Checkbox id="start-search" bind:checked={startSearch} />
					<Label for="start-search">Start searching</Label>
				</div>
			{/if}
			<Button type="button" variant="outline" onclick={onClose}>
				<XIcon aria-hidden="true" />
				<span>Cancel</span>
			</Button>
			<Button
				type="submit"
				disabled={!canConfirm || saving || (isAdmin && libraryFolders.length === 0)}
			>
				{#if saving}
					<LoaderCircleIcon class="animate-spin" aria-hidden="true" />
				{:else}
					<PlusIcon aria-hidden="true" />
				{/if}
				<span>
					{#if saving}
						{isAdmin ? 'Adding' : 'Requesting'}
					{:else}
						{isAdmin ? 'Add' : 'Request'}
					{/if}
				</span>
			</Button>
		</div>
	</form>
</SettingsFormModal>
