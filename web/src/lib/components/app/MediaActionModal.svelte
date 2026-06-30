<script lang="ts">
	import MediaActionOptions from '$lib/components/app/MediaActionOptions.svelte';
	import MediaTagSelector from '$lib/components/app/MediaTagSelector.svelte';
	import {
		mediaPosterUrl,
		preselectLibraryFolderId,
		preselectQualityProfileId
	} from '$lib/components/app/mediaActionDefaults';
	import type { MediaActionSelection } from '$lib/components/app/mediaActionTypes';
	import type {
		LibraryFolder,
		MediaSearchResult,
		MediaMonitorMode,
		MinimumAvailability,
		QualityProfileOption,
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
	let monitorMode = $state<MediaMonitorMode>('only_media');
	let minimumAvailability = $state<MinimumAvailability>('released');
	let startSearch = $state(true);
	let selectedTags = $state<string[]>([]);

	const canConfirm = $derived(!isAdmin || (qualityProfileId !== '' && libraryFolderId !== ''));
	const posterUrl = $derived(mediaPosterUrl(candidate.posterPath));

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
			minimumAvailability,
			startSearch
		});
	}
</script>

<div class="modal-backdrop" role="presentation" onclick={onClose}>
	<div
		class="modal-shell"
		aria-labelledby="media-action-title"
		role="dialog"
		aria-modal="true"
		onclick={(event) => event.stopPropagation()}
		onkeydown={(event) => event.stopPropagation()}
		tabindex="-1"
	>
		<form
			class="media-action-modal"
			onsubmit={submit}
			style:--modal-bg-url={posterUrl ? `url("${posterUrl}")` : undefined}
		>
			<div class="section-heading">
				<div>
					<p class="section-kicker">{isAdmin ? 'Add media' : 'Request media'}</p>
					<h2 id="media-action-title">{candidate.title}</h2>
				</div>
				<button type="button" class="secondary icon-button" aria-label="Close" onclick={onClose}>
					<span class="app-icon" aria-hidden="true">close</span>
				</button>
			</div>

			<p class="muted">
				{candidate.type}{candidate.year ? ` · ${candidate.year}` : ''}
			</p>

			{#if isAdmin}
				<MediaActionOptions
					{isAdmin}
					{libraryFolders}
					{qualityProfiles}
					bind:qualityProfileId
					bind:libraryFolderId
					bind:monitorMode
					bind:minimumAvailability
				/>
				{#if libraryFolders.length === 0}
					<p class="error">Add a library folder in Settings before adding monitored media.</p>
				{/if}
			{:else}
				<MediaActionOptions
					{isAdmin}
					{libraryFolders}
					{qualityProfiles}
					bind:qualityProfileId
					bind:libraryFolderId
					bind:monitorMode
					bind:minimumAvailability
				/>
				<p class="muted">
					Your request will be visible under Requests. An admin will choose the folder and quality
					profile before approval.
				</p>
			{/if}

			<MediaTagSelector {tags} bind:selectedTags />

			<div class="form-actions media-action-actions">
				{#if isAdmin}
					<label class="inline-check">
						<input type="checkbox" bind:checked={startSearch} />
						<span>Start searching</span>
					</label>
				{/if}
				<button type="button" class="secondary media-action-command" onclick={onClose}>
					<span class="app-icon" aria-hidden="true">close</span>
					<span>Cancel</span>
				</button>
				<button
					type="submit"
					class="media-action-command add-action-button"
					disabled={!canConfirm || saving || (isAdmin && libraryFolders.length === 0)}
				>
					<span class="app-icon" aria-hidden="true">add</span>
					<span>
						{#if saving}
							{isAdmin ? 'Adding' : 'Requesting'}
						{:else}
							{isAdmin ? 'Add' : 'Request'}
						{/if}
					</span>
				</button>
			</div>
		</form>
	</div>
</div>
