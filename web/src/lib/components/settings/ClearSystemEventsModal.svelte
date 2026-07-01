<script lang="ts">
	import SettingsFormModal from './shared/SettingsFormModal.svelte';

	interface Props {
		clearing: boolean;
		onCancel: () => void;
		onConfirm: () => void;
	}

	let { clearing, onCancel, onConfirm }: Props = $props();
</script>

<SettingsFormModal title="Clear all events" modalClass="clear-events-modal" onClose={onCancel}>
	<div class="modal-confirmation">
		<p>This permanently removes every recorded system event.</p>
		<div class="modal-actions">
			<button type="button" class="secondary" disabled={clearing} onclick={onCancel}>Cancel</button>
			<button type="button" class="danger" disabled={clearing} onclick={onConfirm}>
				{clearing ? 'Clearing' : 'Clear all'}
			</button>
		</div>
	</div>
</SettingsFormModal>
