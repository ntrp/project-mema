<script lang="ts">
	import { activityDisplay, releaseGroupFromTitle } from './activityDisplay';
	import type { DownloadActivity, ManualImportRequest } from '$lib/settings/types';

	interface Props {
		activity: DownloadActivity;
		importing: boolean;
		error?: string;
		onImport: (_request: ManualImportRequest) => void;
		onClose: () => void;
	}

	let { activity, importing, error, onImport, onClose }: Props = $props();

	const summary = $derived(activityDisplay(activity));
	let initializedActivityId = $state('');
	let sourcePath = $state('');
	let targetFileName = $state('');
	let movieTitle = $state('');
	let year = $state<number | undefined>();
	let seasonNumber = $state<number | undefined>();
	let episodeNumber = $state<number | undefined>();
	let episodeTitle = $state('');
	let releaseGroup = $state('');
	let edition = $state('');
	let quality = $state('');
	let languagesText = $state('');

	$effect(() => {
		if (initializedActivityId === activity.id) return;
		initializedActivityId = activity.id;
		sourcePath = '';
		targetFileName = '';
		movieTitle = activity.mediaTitle;
		year = activity.mediaYear ?? undefined;
		seasonNumber = activity.mediaType === 'series' ? 1 : undefined;
		episodeNumber = activity.mediaType === 'series' ? 1 : undefined;
		episodeTitle = '';
		releaseGroup = releaseGroupFromTitle(activity.releaseTitle);
		edition = '';
		quality = summary.quality === '-' ? '' : summary.quality;
		languagesText = summary.languages.join(', ');
	});

	function submit() {
		onImport({
			sourcePath,
			targetFileName: optional(targetFileName),
			movieTitle: optional(movieTitle),
			year,
			seasonNumber,
			episodeNumber,
			episodeTitle: optional(episodeTitle),
			releaseGroup: optional(releaseGroup),
			edition: optional(edition),
			quality: optional(quality),
			languages: languagesText
				.split(',')
				.map((value) => value.trim())
				.filter(Boolean)
		});
	}

	function optional(value: string) {
		value = value.trim();
		return value === '' ? undefined : value;
	}

	function closeFromBackdrop(event: Event) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}
</script>

<div class="modal-backdrop" role="presentation" onclick={closeFromBackdrop}>
	<div
		class="modal-shell settings-modal manual-import-modal"
		role="dialog"
		aria-modal="true"
		aria-labelledby="manual-import-title"
	>
		<div class="modal-header">
			<div>
				<p>Manual import</p>
				<h2 id="manual-import-title">{activity.mediaTitle}</h2>
				<span>{activity.releaseTitle}</span>
			</div>
			<button type="button" class="secondary" onclick={onClose}>Close</button>
		</div>

		<form class="settings-form" onsubmit={(event) => (event.preventDefault(), submit())}>
			<label class="wide">
				<span>Source path</span>
				<input bind:value={sourcePath} required placeholder="/downloads/release/file.mkv" />
			</label>
			<label class="wide">
				<span>Target filename override</span>
				<input bind:value={targetFileName} placeholder="Leave empty to build from params" />
			</label>
			<label>
				<span>{activity.mediaType === 'series' ? 'Series' : 'Movie'}</span>
				<input bind:value={movieTitle} />
			</label>
			<label>
				<span>Year</span>
				<input bind:value={year} min="0" type="number" />
			</label>
			{#if activity.mediaType === 'series'}
				<label>
					<span>Season</span>
					<input bind:value={seasonNumber} min="0" type="number" required />
				</label>
				<label>
					<span>Episode</span>
					<input bind:value={episodeNumber} min="0" type="number" required />
				</label>
				<label class="wide">
					<span>Episode title</span>
					<input bind:value={episodeTitle} />
				</label>
			{/if}
			<label>
				<span>Release group</span>
				<input bind:value={releaseGroup} />
			</label>
			<label>
				<span>Edition</span>
				<input bind:value={edition} />
			</label>
			<label>
				<span>Quality</span>
				<input bind:value={quality} />
			</label>
			<label>
				<span>Languages</span>
				<input bind:value={languagesText} placeholder="English, German" />
			</label>
			{#if error}
				<p class="form-status error wide">{error}</p>
			{/if}
			<div class="modal-actions wide">
				<button type="button" class="secondary" onclick={onClose}>Cancel</button>
				<button type="submit" disabled={importing}>{importing ? 'Importing' : 'Import'}</button>
			</div>
		</form>
	</div>
</div>
