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

	type SmartMediaCandidate = MediaSearchResult & {
		genres?: string[];
		originalLanguage?: string;
	};

	const smartCandidate = $derived(candidate as SmartMediaCandidate);

	let qualityProfileId = $state(preselectQualityProfileId());
	let libraryFolderId = $state(preselectLibraryFolderId());
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

	function preselectQualityProfileId() {
		return bestScored(qualityProfiles, profileScore)?.id ?? '';
	}

	function profileScore(profile: QualityProfileOption) {
		const text = normalizedText(`${profile.id} ${profile.name}`);
		let score = 0;

		if (isAnimeCandidate()) {
			score += hasAny(text, ['anime']) ? 50 : 0;
		} else if (hasAny(text, ['anime'])) {
			score -= 20;
		}

		if (hasAny(text, ['1080', '1080p'])) {
			score += 40;
		}
		if (hasAny(text, ['2160', '2160p', '4k', 'uhd'])) {
			score += 25;
		}
		if (hasAny(text, ['default', 'best'])) {
			score += 15;
		}
		if (hasAny(text, ['any acceptable', 'any'])) {
			score -= 10;
		}

		return score;
	}

	function preselectLibraryFolderId() {
		const bestFolder = bestScored(libraryFolders, folderScore);
		return bestFolder?.id ?? '';
	}

	function folderScore(folder: LibraryFolder) {
		const path = normalizedText(folder.path);
		const hasAnime = hasAny(path, ['anime']);
		const hasMovie = hasAny(path, ['movie', 'movies', 'film', 'films']);
		const hasSeries = hasAny(path, ['series', 'tv', 'show', 'shows']);
		let score = 0;

		if (candidate.type === 'series') {
			score += hasSeries ? 100 : 0;
			score += hasAnime && isAnimeCandidate() ? 20 : 0;
			score -= hasMovie ? 25 : 0;
			return score;
		}

		if (isAnimeCandidate()) {
			score += hasAnime && hasMovie ? 120 : 0;
			score += hasAnime && !hasMovie ? 90 : 0;
			score += !hasAnime && hasMovie ? 60 : 0;
			score -= hasSeries ? 25 : 0;
			return score;
		}

		score += hasMovie ? 100 : 0;
		score -= hasAnime ? 15 : 0;
		score -= hasSeries ? 25 : 0;
		return score;
	}

	function isAnimeCandidate() {
		const genres = smartCandidate.genres ?? [];
		const text = normalizedText(
			[candidate.title, candidate.overview, genres.join(' '), smartCandidate.originalLanguage].join(
				' '
			)
		);
		const isJapaneseAnimation =
			smartCandidate.originalLanguage?.toLowerCase() === 'ja' &&
			genres.some((genre) => normalizedText(genre).includes('animation'));

		return text.includes('anime') || isJapaneseAnimation;
	}

	function bestScored<T>(items: T[], scoreItem: (item: T) => number) {
		return items
			.map((item, index) => ({ item, score: scoreItem(item), index }))
			.sort((left, right) => right.score - left.score || left.index - right.index)[0]?.item;
	}

	function normalizedText(value: string) {
		return value.toLowerCase();
	}

	function hasAny(text: string, needles: string[]) {
		return needles.some((needle) => text.includes(needle));
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
					<div class="tag-input-box">
						{#each selectedTags as tag (tag.toLowerCase())}
							<button type="button" onclick={() => removeTag(tag)}>{tag}</button>
						{/each}
						<input
							bind:value={tagInput}
							type="text"
							list="media-action-tag-options"
							maxlength="80"
							placeholder={selectedTags.length === 0 ? 'Add tag' : ''}
							autocomplete="off"
							onkeydown={handleTagKeydown}
							onblur={commitTagInput}
						/>
					</div>
					<datalist id="media-action-tag-options">
						{#each tags as tag (tag.id)}
							{#if !selectedTags.some((selected) => selected.toLowerCase() === tag.name.toLowerCase())}
								<option value={tag.name}></option>
							{/if}
						{/each}
					</datalist>
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
