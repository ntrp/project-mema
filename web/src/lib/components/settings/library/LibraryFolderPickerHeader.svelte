<script lang="ts">
	interface Props {
		selectedPath: string;
		onClose: () => void;
	}

	let { selectedPath, onClose }: Props = $props();
</script>

<div class="section-heading">
	<div>
		<p class="section-kicker">Folder picker</p>
		<h3 id="folder-picker-title">Select library folder</h3>
	</div>
	<button type="button" class="secondary" onclick={onClose}>Close</button>
</div>

<div class="folder-picker-selected">
	<span>Selected</span>
	<strong>{selectedPath || 'No folder selected'}</strong>
</div>
