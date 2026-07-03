<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

	interface Props {
		row: MediaFileRow;
		onCancel: () => void;
		onConfirm: () => void;
	}

	let { row, onCancel, onConfirm }: Props = $props();
	let open = $state(true);

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) {
			onCancel();
		}
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="w-[min(460px,calc(100vw-32px))]">
		<Dialog.Header>
			<Dialog.Title>Delete file</Dialog.Title>
			<Dialog.Description>
				This removes the file from disk. The media item stays in the app.
			</Dialog.Description>
		</Dialog.Header>
		<p class="rounded-md border bg-muted/40 px-3 py-2 text-sm break-all text-muted-foreground">
			{row.relativePath}
		</p>
		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
			<Button type="button" variant="destructive" onclick={onConfirm}>Delete</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
