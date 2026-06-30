<script lang="ts">
	interface Props {
		totalFiles: number;
		sortMode: 'folders' | 'mixed';
	}

	let { totalFiles, sortMode = $bindable() }: Props = $props();
</script>

<div class="scan-table-toolbar">
	<span>{totalFiles} files</span>
	<label>
		<span>Sort</span>
		<select bind:value={sortMode}>
			<option value="folders">Folders</option>
			<option value="mixed">Mixed</option>
		</select>
	</label>
</div>
