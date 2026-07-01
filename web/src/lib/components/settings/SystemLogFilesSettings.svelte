<script lang="ts">
	import { onMount } from 'svelte';
	import { downloadSystemLogFile, listSystemLogFiles } from '$lib/settings/api';
	import { formatLongDateTime } from '$lib/settings/dateFormat';
	import type { SystemLogFile } from '$lib/settings/types';

	let files = $state<SystemLogFile[]>([]);
	let downloadingName = $state<string | undefined>();
	let errorMessage = $state('');

	onMount(() => {
		void load();
	});

	async function load() {
		errorMessage = '';
		try {
			files = await listSystemLogFiles();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not load log files';
		}
	}

	async function download(file: SystemLogFile) {
		downloadingName = file.name;
		errorMessage = '';
		try {
			const blob = await downloadSystemLogFile(file.name);
			const href = globalThis.URL.createObjectURL(blob);
			const link = globalThis.document.createElement('a');
			link.href = href;
			link.download = file.name;
			link.click();
			globalThis.URL.revokeObjectURL(href);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not download log file';
		} finally {
			downloadingName = undefined;
		}
	}

	function formatSize(bytes: number) {
		if (bytes < 1024) {
			return `${bytes} B`;
		}
		if (bytes < 1024 * 1024) {
			return `${(bytes / 1024).toFixed(1)} KB`;
		}
		return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
	}
</script>

<section class="panel log-settings-panel" aria-label="Log files">
	{#if errorMessage}
		<p class="inline-error">{errorMessage}</p>
	{/if}

	<div class="table-wrap">
		<table class="data-table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Size</th>
					<th>Modified</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each files as file (file.name)}
					<tr>
						<td>{file.name}</td>
						<td>{formatSize(file.sizeBytes)}</td>
						<td>{formatLongDateTime(file.modifiedAt)}</td>
						<td>
							<button
								type="button"
								class="secondary compact-action"
								disabled={downloadingName === file.name}
								onclick={() => download(file)}
							>
								{downloadingName === file.name ? 'Downloading' : 'Download'}
							</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="4">No log files retained.</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>
