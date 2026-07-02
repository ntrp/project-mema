<script lang="ts">
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import {
		scanSortModeOptions,
		type LibraryScanSortMode
	} from '$lib/components/settings/library/scan/libraryScanImport';

	interface Props {
		totalFiles: number;
		sortMode: LibraryScanSortMode;
	}

	let { totalFiles, sortMode = $bindable() }: Props = $props();
</script>

<div
	class="flex items-center justify-between gap-3 text-[13px] font-extrabold text-muted-foreground"
>
	<span>{totalFiles} files</span>
	<label class="flex items-center gap-2">
		<span>Sort</span>
		<SettingsSelect
			value={sortMode}
			options={scanSortModeOptions}
			onValueChange={(value) => (sortMode = value as LibraryScanSortMode)}
		/>
	</label>
</div>
