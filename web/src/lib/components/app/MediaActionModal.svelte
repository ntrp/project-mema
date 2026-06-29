<script lang="ts">
	import type { LibraryFolder, MediaSearchResult, QualityProfileOption } from '$lib/settings/types';

	interface Props {
		candidate: MediaSearchResult;
		isAdmin: boolean;
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		saving: boolean;
		onClose: () => void;
		onConfirm: (_qualityProfileId?: string, _libraryFolderId?: string) => void;
	}

	let { candidate, isAdmin, libraryFolders, qualityProfiles, saving, onClose, onConfirm }: Props =
		$props();

	let qualityProfileId = $state('');
	let libraryFolderId = $state('');

	const canConfirm = $derived(!isAdmin || (qualityProfileId !== '' && libraryFolderId !== ''));

	function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!canConfirm || saving) {
			return;
		}
		onConfirm(qualityProfileId, libraryFolderId);
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
		<form class="media-action-modal" onsubmit={submit}>
			<div class="section-heading">
				<div>
					<p class="section-kicker">{isAdmin ? 'Add media' : 'Request media'}</p>
					<h2 id="media-action-title">{candidate.title}</h2>
				</div>
				<button type="button" class="secondary" onclick={onClose}>Close</button>
			</div>

			<p class="muted">
				{candidate.type}{candidate.year ? ` · ${candidate.year}` : ''}
			</p>

			{#if isAdmin}
				<div class="settings-form compact-form">
					<label>
						<span>Quality profile</span>
						<select bind:value={qualityProfileId}>
							<option value="" disabled>Select profile</option>
							{#each qualityProfiles as profile (profile.id)}
								<option value={profile.id}>{profile.name}</option>
							{/each}
						</select>
					</label>
					<label>
						<span>Library folder</span>
						<select bind:value={libraryFolderId}>
							<option value="" disabled>Select folder</option>
							{#each libraryFolders as folder (folder.id)}
								<option value={folder.id}>{folder.path}</option>
							{/each}
						</select>
					</label>
				</div>
				{#if libraryFolders.length === 0}
					<p class="error">Add a library folder in Settings before adding monitored media.</p>
				{/if}
			{:else}
				<p class="muted">
					Your request will be visible under Requests. An admin will choose the folder and quality
					profile before approval.
				</p>
			{/if}

			<div class="form-actions">
				<button type="button" class="secondary" onclick={onClose}>Cancel</button>
				<button
					type="submit"
					disabled={!canConfirm || saving || (isAdmin && libraryFolders.length === 0)}
				>
					{#if saving}
						{isAdmin ? 'Adding' : 'Requesting'}
					{:else}
						{isAdmin ? 'Add' : 'Request'}
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
