<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import ChevronUpIcon from '@lucide/svelte/icons/chevron-up';
	import LibraryScanImportTable from '$lib/components/settings/library/scan/LibraryScanImportTable.svelte';
	import type { LibraryScanImportRow } from '$lib/components/settings/library/scan/libraryScanImport';
	import SettingsRowActionButton from '$lib/components/settings/shared/SettingsRowActionButton.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import type {
		LibraryFolder,
		LibraryMediaKind,
		LibraryScan,
		MediaSearchResult,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		folders: LibraryFolder[];
		scansByFolder: Record<string, LibraryScan>;
		openFolderId?: string;
		scanningFolderId?: string;
		qualityProfiles: QualityProfileOption[];
		onScan: (_id: string) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
		onSearchMatch: (_kind: LibraryMediaKind, _query: string) => Promise<MediaSearchResult[]>;
		onImport: (_scan: LibraryScan, _rows: LibraryScanImportRow[]) => Promise<void>;
	}

	let {
		folders,
		scansByFolder,
		openFolderId,
		scanningFolderId,
		qualityProfiles,
		onScan,
		onDelete,
		onSearchMatch,
		onImport
	}: Props = $props();

	let opened = $state<Record<string, boolean>>({});

	$effect(() => {
		if (openFolderId) {
			opened[openFolderId] = true;
		}
	});

	function toggle(id: string) {
		opened[id] = !opened[id];
		if (opened[id] && !scansByFolder[id]) {
			void onScan(id);
		}
	}
</script>

<div class="grid gap-3">
	{#each folders as folder (folder.id)}
		{@const scan = scansByFolder[folder.id]}
		<Card class="overflow-hidden p-0">
			<div class="flex items-center gap-2 border-b p-2">
				<Button
					type="button"
					variant="ghost"
					class="min-w-0 flex-1 justify-start px-2 text-left"
					onclick={() => toggle(folder.id)}
				>
					{#if opened[folder.id]}
						<ChevronUpIcon aria-hidden="true" />
					{:else}
						<ChevronDownIcon aria-hidden="true" />
					{/if}
					<span class="truncate">{folder.path}</span>
				</Button>
				<div class="flex shrink-0 items-center gap-2">
					<SettingsRowActionButton
						label={`Scan ${folder.path}`}
						icon="sync"
						disabled={scanningFolderId === folder.id}
						onclick={() => onScan(folder.id)}
					/>
					<SettingsRowActionButton
						label={`Delete ${folder.path}`}
						icon="delete"
						variant="destructive"
						confirmTitle="Delete library folder"
						confirmDescription={`Delete library folder "${folder.path}" from settings?`}
						confirmLabel="Delete folder"
						onclick={() => onDelete(folder.id)}
					/>
				</div>
			</div>
			{#if opened[folder.id]}
				<div class="grid gap-3 p-3">
					{#if scanningFolderId === folder.id && !scan}
						<p class="m-0 text-sm leading-6 text-muted-foreground">Scanning folder</p>
					{:else if scan}
						<LibraryScanImportTable
							{scan}
							{qualityProfiles}
							loading={scanningFolderId === folder.id}
							{onSearchMatch}
							{onImport}
						/>
					{:else}
						<p class="m-0 text-sm leading-6 text-muted-foreground">
							Open or scan this folder to review media.
						</p>
					{/if}
				</div>
			{/if}
		</Card>
	{:else}
		<p class="m-0 text-sm leading-6 text-muted-foreground">No library folders configured</p>
	{/each}
</div>
