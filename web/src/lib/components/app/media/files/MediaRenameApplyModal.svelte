<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';

	interface Props {
		safeCount: number;
		onCancel: () => void;
		onConfirm: () => void;
	}

	let { safeCount, onCancel, onConfirm }: Props = $props();
	let open = $state(true);

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) {
			onCancel();
		}
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="w-[min(500px,calc(100vw-32px))]">
		<Dialog.Header>
			<Dialog.Title>Apply rename preview</Dialog.Title>
			<Dialog.Description>
				This moves {safeCount} file{safeCount === 1 ? '' : 's'} on disk and records file history.
			</Dialog.Description>
		</Dialog.Header>
		<p class="rounded-md border bg-muted/40 px-3 py-2 text-sm text-muted-foreground">
			The latest preview will be checked again before anything is moved.
		</p>
		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
			<Button type="button" onclick={onConfirm}>Apply</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
