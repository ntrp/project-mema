<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
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

<Card class="gap-4 p-5" aria-label="Log files">
	{#if errorMessage}
		<p class="m-0 font-bold text-destructive">{errorMessage}</p>
	{/if}

	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head>Size</Table.Head>
				<Table.Head>Modified</Table.Head>
				<Table.Head></Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each files as file (file.name)}
				<Table.Row>
					<Table.Cell>{file.name}</Table.Cell>
					<Table.Cell>{formatSize(file.sizeBytes)}</Table.Cell>
					<Table.Cell>{formatLongDateTime(file.modifiedAt)}</Table.Cell>
					<Table.Cell class="text-right">
						<Button
							type="button"
							variant="outline"
							size="sm"
							disabled={downloadingName === file.name}
							onclick={() => download(file)}
						>
							{downloadingName === file.name ? 'Downloading' : 'Download'}
						</Button>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={4}>No log files retained.</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
