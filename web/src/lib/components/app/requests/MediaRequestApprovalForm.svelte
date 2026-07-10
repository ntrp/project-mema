<script lang="ts">
	import MediaActionOptions from '$lib/components/app/media/actions/MediaActionOptions.svelte';
	import MediaTagSelector from '$lib/components/app/media/actions/MediaTagSelector.svelte';
	import {
		matchingLibraryFolders,
		preselectLibraryFolderId,
		preselectQualityProfileId
	} from '$lib/components/app/media/actions/mediaActionDefaults';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import type {
		LibraryFolder,
		MediaMonitorMode,
		MediaRequest,
		MediaRequestApproveRequest,
		MinimumAvailability,
		QualityProfileOption,
		SeriesType,
		Tag
	} from '$lib/settings/types';

	interface Props {
		request: MediaRequest;
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		tags: Tag[];
		approvingRequestId?: string;
		onApprove: (_request: MediaRequest, _approval: MediaRequestApproveRequest) => void;
	}

	let { request, libraryFolders, qualityProfiles, tags, approvingRequestId, onApprove }: Props =
		$props();

	let qualityProfileId = $derived(preselectQualityProfileId(request, qualityProfiles));
	let libraryFolderId = $derived(preselectLibraryFolderId(request, libraryFolders));
	let selectedMonitorMode = $state<MediaMonitorMode | undefined>();
	let selectedSeriesType = $state<SeriesType | undefined>();
	let minimumAvailability = $state<MinimumAvailability>('released');
	let startSearch = $state(true);
	let selectedTags = $state<string[]>([]);
	let matchingFolders = $derived(matchingLibraryFolders(request.type, libraryFolders));
	let monitorMode = $derived(selectedMonitorMode ?? defaultMonitorMode(request.type));
	let seriesType = $derived(selectedSeriesType ?? 'standard');

	function approve(event: SubmitEvent) {
		event.preventDefault();
		if (!qualityProfileId || !libraryFolderId) {
			return;
		}
		onApprove(request, {
			qualityProfileId,
			libraryFolderId,
			monitorMode,
			seriesType: request.type === 'serie' ? seriesType : undefined,
			minimumAvailability,
			startSearch,
			tags: selectedTags
		});
	}

	function defaultMonitorMode(type: MediaRequest['type']): MediaMonitorMode {
		return type === 'serie' ? 'all_episodes' : 'only_media';
	}
</script>

<form class="grid gap-4" onsubmit={approve}>
	<MediaActionOptions
		isAdmin={true}
		mediaType={request.type}
		libraryFolders={matchingFolders}
		{qualityProfiles}
		bind:qualityProfileId
		bind:libraryFolderId
		{monitorMode}
		{seriesType}
		bind:minimumAvailability
		onMonitorModeChange={(mode) => (selectedMonitorMode = mode)}
		onSeriesTypeChange={(type) => (selectedSeriesType = type)}
	/>
	<MediaTagSelector {tags} bind:selectedTags />
	<div class="flex flex-wrap items-center gap-3">
		<div class="mr-auto flex items-center gap-2">
			<Checkbox id="request-start-search" bind:checked={startSearch} />
			<Label for="request-start-search">Start searching</Label>
		</div>
		<Button
			type="submit"
			disabled={approvingRequestId === request.id || matchingFolders.length === 0}
		>
			{approvingRequestId === request.id ? 'Approving' : 'Approve'}
		</Button>
	</div>
</form>
