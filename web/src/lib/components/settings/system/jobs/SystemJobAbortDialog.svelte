<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';

	interface AbortJob {
		id: number;
		kind: string;
	}

	interface Props {
		job?: AbortJob;
		onClose: () => void;
		onAbort: () => void | Promise<void>;
	}

	let { job, onClose, onAbort }: Props = $props();
</script>

<Dialog.Root open={!!job} onOpenChange={(open) => !open && onClose()}>
	<Dialog.Content class="w-[min(300px,calc(100vw-32px))]">
		<Dialog.Header>
			<Dialog.Title>Abort job</Dialog.Title>
			<Dialog.Description>
				Abort job {job?.id} ({job?.kind})?
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={onClose}>Cancel</Button>
			<Button type="button" variant="destructive" onclick={() => void onAbort()}>Abort</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
