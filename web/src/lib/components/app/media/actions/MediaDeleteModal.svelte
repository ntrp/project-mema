<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
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

<SettingsFormModal title={item.title} modalClass="grid w-[min(560px,100%)] gap-4" {onClose}>
	<div class="grid gap-4">
		<p class="m-0 mb-1.5 text-xs font-extrabold text-muted-foreground uppercase">Remove media</p>
		<p class="m-0 text-sm leading-6 text-muted-foreground">
			This media item has {fileLabel}. Deleting removes it from the app and, unless kept, deletes
			its media folder from disk.
		</p>
		{#if item.mediaFolderPath}
			<p
				class="rounded-md border border-border bg-card px-3 py-2.5 font-mono text-sm break-anywhere text-muted-foreground"
			>
				{item.mediaFolderPath}
			</p>
		{/if}

		<div class="flex items-center gap-2">
			<Checkbox id="keep-media-files" bind:checked={keepFiles} />
			<Label for="keep-media-files">Keep media files</Label>
		</div>

		<div class="flex justify-end gap-2">
			<Button type="button" variant="outline" onclick={onClose} disabled={deleting}>Cancel</Button>
			<Button
				type="button"
				variant="destructive"
				onclick={() => onDelete(keepFiles)}
				disabled={deleting}
			>
				{deleting ? 'Deleting' : 'Delete'}
			</Button>
		</div>
	</div>
</SettingsFormModal>
