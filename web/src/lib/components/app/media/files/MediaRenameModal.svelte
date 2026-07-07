<script lang="ts">
	import { onMount } from 'svelte';
	import CheckIcon from '@lucide/svelte/icons/check';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Table from '$lib/components/ui/table';
	import { applyMediaRename, getFileNamingSettings, previewMediaRename } from '$lib/settings/api';
	import { defaultFileNamingTemplates } from '$lib/settings/fileNamingTemplates';
	import { relativePath } from '$lib/components/app/media/files/mediaFilePath';
	import type { FileNamingSettings, MediaItem, MediaRenamePreviewRow } from '$lib/settings/types';

	type Props = { item: MediaItem; onClose: () => void; onApplied?: () => void };

	let { item, onClose, onApplied = () => {} }: Props = $props();
	let open = $state(true);
	let rows = $state<MediaRenamePreviewRow[]>([]);
	let selected = $state<Record<string, boolean>>({});
	let settings = $state<FileNamingSettings>();
	let loading = $state(false);
	let applying = $state(false);
	let errorMessage = $state<string | undefined>();

	const rootPath = $derived(item.mediaFolderPath ?? '-');
	const template = $derived(
		item.type === 'movie'
			? (settings?.movieFileFormat ?? defaultFileNamingTemplates.movieFileFormat)
			: (settings?.seriesEpisodeFormat ?? defaultFileNamingTemplates.seriesEpisodeFormat)
	);
	const safeRows = $derived(rows.filter((row) => row.status === 'safe'));
	const selectedPaths = $derived(safeRows.filter((row) => selected[row.currentPath]));
	const allSafeSelected = $derived(safeRows.length > 0 && selectedPaths.length === safeRows.length);

	onMount(() => {
		void load();
	});

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) onClose();
	}

	async function load() {
		loading = true;
		errorMessage = undefined;
		try {
			const [preview, fileNaming] = await Promise.all([
				previewMediaRename(item.id),
				getFileNamingSettings()
			]);
			rows = preview.rows;
			settings = fileNaming;
			selected = Object.fromEntries(
				preview.rows.filter((row) => row.status === 'safe').map((row) => [row.currentPath, true])
			);
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not preview rename';
		} finally {
			loading = false;
		}
	}

	async function applySelected() {
		applying = true;
		errorMessage = undefined;
		try {
			const result = await applyMediaRename(
				item.id,
				selectedPaths.map((row) => row.currentPath)
			);
			rows = result.rows;
			selected = {};
			if (result.appliedCount > 0) onApplied();
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : 'Could not apply rename';
		} finally {
			applying = false;
		}
	}

	function setSelected(path: string, checked: boolean) {
		selected = { ...selected, [path]: checked };
	}

	function toggleAll() {
		const next = !allSafeSelected;
		selected = Object.fromEntries(safeRows.map((row) => [row.currentPath, next]));
	}

	function statusVariant(status: MediaRenamePreviewRow['status']) {
		if (status === 'safe' || status === 'applied') return 'secondary';
		if (status === 'unchanged' || status === 'skipped') return 'outline';
		return 'destructive';
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="w-[min(980px,calc(100vw-32px))] gap-4">
		<Dialog.Header>
			<Dialog.Title>Rename files</Dialog.Title>
		</Dialog.Header>

		<div class="grid gap-2 rounded-md border bg-muted/30 p-3 text-sm">
			<div class="grid gap-1 md:grid-cols-[9rem_minmax(0,1fr)]">
				<strong class="text-muted-foreground">Root path</strong>
				<span class="break-anywhere font-mono">{rootPath}</span>
			</div>
			<div class="grid gap-1 md:grid-cols-[9rem_minmax(0,1fr)]">
				<strong class="text-muted-foreground">Naming template</strong>
				<span class="break-anywhere font-mono">{template}</span>
			</div>
		</div>

		{#if errorMessage}
			<p class="m-0 text-sm font-medium text-destructive">{errorMessage}</p>
		{/if}

		<div class="max-h-[52vh] overflow-auto rounded-md border">
			<Table.Root class="min-w-180 text-sm">
					<Table.Header>
						<Table.Row>
							<Table.Head class="w-12">
							<Checkbox
								aria-label="Select all safe rename rows"
								checked={allSafeSelected}
								disabled={safeRows.length === 0 || loading || applying}
									onCheckedChange={toggleAll}
								/>
							</Table.Head>
							<Table.Head>Relative path</Table.Head>
							<Table.Head class="w-28">Status</Table.Head>
							<Table.Head>Reason</Table.Head>
						</Table.Row>
					</Table.Header>
				<Table.Body>
					{#each rows as row (row.currentPath)}
						<Table.Row>
							<Table.Cell>
								<Checkbox
									aria-label={`Rename ${relativePath(item.mediaFolderPath, row.currentPath)}`}
									checked={!!selected[row.currentPath]}
									disabled={row.status !== 'safe' || loading || applying}
									onCheckedChange={(checked) => setSelected(row.currentPath, checked === true)}
								/>
							</Table.Cell>
							<Table.Cell class="whitespace-normal">
								<div class="grid gap-1">
									<span
										class={`break-anywhere text-muted-foreground ${row.proposedPath && row.proposedPath !== row.currentPath ? 'line-through' : ''}`}
									>
										{relativePath(item.mediaFolderPath, row.currentPath)}
									</span>
									{#if row.proposedPath && row.proposedPath !== row.currentPath}
										<span class="break-anywhere font-medium">
											{relativePath(item.mediaFolderPath, row.proposedPath)}
										</span>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell>
								<Badge variant={statusVariant(row.status)}>
									{row.status === 'missing' ? 'skipped' : row.status}
								</Badge>
							</Table.Cell>
							<Table.Cell class="break-anywhere whitespace-normal text-muted-foreground">
								{row.messages.length > 0 ? row.messages.join(' ') : '-'}
							</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={4} class="text-muted-foreground">
								{loading ? 'Loading rename preview...' : 'No files available for rename.'}
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>

		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={load} disabled={loading || applying}>
				<RefreshCwIcon class={loading ? 'animate-spin' : undefined} aria-hidden="true" />
				{loading ? 'Previewing' : 'Refresh'}
			</Button>
			<Button type="button" variant="outline" onclick={onClose} disabled={applying}>Cancel</Button>
			<Button
				type="button"
				onclick={applySelected}
				disabled={selectedPaths.length === 0 || loading || applying}
			>
				<CheckIcon aria-hidden="true" />
				{applying ? 'Applying' : `Apply ${selectedPaths.length}`}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
