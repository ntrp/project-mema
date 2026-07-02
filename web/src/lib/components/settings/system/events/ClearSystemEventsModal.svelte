<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import SettingsFormModal from '../../shared/SettingsFormModal.svelte';

	interface Props {
		clearing: boolean;
		onCancel: () => void;
		onConfirm: () => void;
	}

	let { clearing, onCancel, onConfirm }: Props = $props();
</script>

<SettingsFormModal
	title="Clear all events"
	modalClass="grid w-[min(460px,100%)] gap-4"
	onClose={onCancel}
>
	<div class="grid gap-[18px]">
		<p class="m-0 leading-normal text-muted-foreground">
			This permanently removes every recorded system event.
		</p>
		<div class="flex flex-wrap justify-end gap-2.5">
			<Button type="button" variant="outline" disabled={clearing} onclick={onCancel}>Cancel</Button>
			<Button type="button" variant="destructive" disabled={clearing} onclick={onConfirm}>
				{clearing ? 'Clearing' : 'Clear all'}
			</Button>
		</div>
	</div>
</SettingsFormModal>
