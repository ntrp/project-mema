<script lang="ts">
	import type { MediaItem } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		deleting: boolean;
		onClose: () => void;
		onDelete: (_keepFiles: boolean) => void;
	}

	let { item, deleting, onClose, onDelete }: Props = $props();
	let keepFiles = $state(false);
	const fileCount = $derived((item.filePaths?.length ?? 0) + (item.metadataFilePaths?.length ?? 0));
	const fileLabel = $derived(fileCount === 1 ? '1 file' : `${fileCount} files`);
</script>

<div class="modal-backdrop" role="presentation" onclick={onClose}>
	<div
		class="modal-shell media-delete-modal"
		aria-labelledby="media-delete-title"
		role="dialog"
		aria-modal="true"
		onclick={(event) => event.stopPropagation()}
		onkeydown={(event) => event.stopPropagation()}
		tabindex="-1"
	>
		<div class="section-heading">
			<div>
				<p class="section-kicker">Remove media</p>
				<h2 id="media-delete-title">{item.title}</h2>
			</div>
			<button type="button" class="secondary icon-button" aria-label="Close" onclick={onClose}>
				<span class="app-icon" aria-hidden="true">close</span>
			</button>
		</div>

		<p class="muted">
			This media item has {fileLabel}. Deleting removes it from the app and, unless kept, deletes
			its media folder from disk.
		</p>
		{#if item.mediaFolderPath}
			<p class="path-preview">{item.mediaFolderPath}</p>
		{/if}

		<label class="inline-check">
			<input type="checkbox" bind:checked={keepFiles} />
			<span>Keep media files</span>
		</label>

		<div class="form-actions media-delete-actions">
			<button type="button" class="secondary" onclick={onClose} disabled={deleting}>Cancel</button>
			<button type="button" class="danger" onclick={() => onDelete(keepFiles)} disabled={deleting}>
				{deleting ? 'Deleting' : 'Delete'}
			</button>
		</div>
	</div>
</div>
