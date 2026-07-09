<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';

	interface Props {
		track: MediaFileDetailRow;
		onCancel: () => void;
		onConfirm: () => void | Promise<void>;
	}

	let { track, onCancel, onConfirm }: Props = $props();

	function deleteDescription() {
		if (track.chapterSummary) return 'Delete all chapters from this file?';
		if (track.type === 'chapter') return `Delete chapter ${track.trackNumber} from this file?`;
		return `Delete ${track.type} track ${track.trackNumber} from this file?`;
	}
</script>

<Dialog.Root open onOpenChange={(open) => !open && onCancel()}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Delete embedded track</Dialog.Title>
			<Dialog.Description>{deleteDescription()}</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
			<Button type="button" variant="destructive" onclick={onConfirm}>Delete</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
