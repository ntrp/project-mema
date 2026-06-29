<script lang="ts">
	import type {
		LibraryFolder,
		MediaSearchResult,
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
		onConfirm: (_qualityProfileId?: string, _libraryFolderId?: string, _tags?: string[]) => void;
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
	let tagInput = $state('');
	let selectedTags = $state<string[]>([]);

	const canConfirm = $derived(!isAdmin || (qualityProfileId !== '' && libraryFolderId !== ''));

	function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!canConfirm || saving) {
			return;
		}
		commitTagInput();
		onConfirm(qualityProfileId, libraryFolderId, selectedTags);
	}

	function toggleTag(name: string) {
		selectedTags = selectedTags.some((tag) => tag.toLowerCase() === name.toLowerCase())
			? selectedTags.filter((tag) => tag.toLowerCase() !== name.toLowerCase())
			: [...selectedTags, name];
	}

	function removeTag(name: string) {
		selectedTags = selectedTags.filter((tag) => tag.toLowerCase() !== name.toLowerCase());
	}

	function commitTagInput() {
		const name = normalizeTag(tagInput);
		if (!name || selectedTags.some((tag) => tag.toLowerCase() === name.toLowerCase())) {
			tagInput = '';
			return;
		}
		selectedTags = [...selectedTags, name];
		tagInput = '';
	}

	function handleTagKeydown(event: Event) {
		if (!(event instanceof globalThis.KeyboardEvent)) {
			return;
		}
		if (event.key !== 'Enter' && event.key !== ',') {
			return;
		}
		event.preventDefault();
		commitTagInput();
	}

	function normalizeTag(value: string) {
		return value.trim().replace(/\s+/g, ' ');
	}

	function imageUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w780${path}`;
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
			style:--modal-bg-url={imageUrl(candidate.posterPath)
				? `url("${imageUrl(candidate.posterPath)}")`
				: undefined}
		>
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

			<div class="tag-selector">
				<div class="tag-selector-header">
					<span>Tags</span>
					<input
						bind:value={tagInput}
						type="text"
						maxlength="80"
						placeholder="Add tag"
						onkeydown={handleTagKeydown}
						onblur={commitTagInput}
					/>
				</div>
				{#if tags.length > 0}
					<div class="tag-options" aria-label="Existing tags">
						{#each tags as tag (tag.id)}
							<button
								type="button"
								class:active-tag={selectedTags.some(
									(selected) => selected.toLowerCase() === tag.name.toLowerCase()
								)}
								onclick={() => toggleTag(tag.name)}
							>
								{tag.name}
							</button>
						{/each}
					</div>
				{/if}
				{#if selectedTags.length > 0}
					<div class="selected-tags" aria-label="Selected tags">
						{#each selectedTags as tag (tag.toLowerCase())}
							<button type="button" onclick={() => removeTag(tag)}>{tag}</button>
						{/each}
					</div>
				{/if}
			</div>

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
