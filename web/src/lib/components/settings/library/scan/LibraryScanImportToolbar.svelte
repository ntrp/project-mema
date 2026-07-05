<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';

	interface Props {
		totalFiles: number;
		matchedCount: number;
		noMatchCount: number;
		importedCount: number;
		duplicateCount: number;
		showImported: boolean;
	}

	let {
		totalFiles,
		matchedCount,
		noMatchCount,
		importedCount,
		duplicateCount,
		showImported = $bindable()
	}: Props = $props();
</script>

<div
	class="flex flex-wrap items-center justify-between gap-3 text-[13px] font-extrabold text-muted-foreground"
>
	<span>
		{totalFiles} files · {matchedCount} matched · {noMatchCount} no match · {importedCount} imported ·
		{duplicateCount} duplicates
	</span>
	<div class="flex items-center gap-4">
		<label class="flex items-center gap-2">
			<Checkbox bind:checked={showImported} />
			<span>Show Imported</span>
		</label>
	</div>
</div>
