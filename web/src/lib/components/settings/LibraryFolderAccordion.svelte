<script lang="ts">
	import LibraryScanImportTable from '$lib/components/settings/LibraryScanImportTable.svelte';
	import type { LibraryScanImportRow } from '$lib/components/settings/libraryScanImport';
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

<div class="folder-accordion-list">
	{#each folders as folder (folder.id)}
		{@const scan = scansByFolder[folder.id]}
		<section class="folder-accordion">
			<div class="folder-accordion-header">
				<button type="button" class="folder-toggle" onclick={() => toggle(folder.id)}>
					<span class="app-icon" aria-hidden="true"
						>{opened[folder.id] ? 'expand_less' : 'expand_more'}</span
					>
					<span>{folder.path}</span>
				</button>
				<div class="row-actions">
					<button
						type="button"
						class="secondary icon-button"
						aria-label={`Scan ${folder.path}`}
						disabled={scanningFolderId === folder.id}
						onclick={() => onScan(folder.id)}
					>
						<span class="app-icon" aria-hidden="true">sync</span>
					</button>
					<button
						type="button"
						class="danger icon-button"
						aria-label={`Delete ${folder.path}`}
						onclick={() => onDelete(folder.id)}
					>
						<span class="app-icon" aria-hidden="true">delete</span>
					</button>
				</div>
			</div>
			{#if opened[folder.id]}
				<div class="folder-accordion-body">
					{#if scanningFolderId === folder.id && !scan}
						<p class="muted">Scanning folder</p>
					{:else if scan}
						<LibraryScanImportTable
							{scan}
							{qualityProfiles}
							loading={scanningFolderId === folder.id}
							{onSearchMatch}
							{onImport}
						/>
					{:else}
						<p class="empty">Open or scan this folder to review media.</p>
					{/if}
				</div>
			{/if}
		</section>
	{:else}
		<p class="empty">No library folders configured</p>
	{/each}
</div>
