<script lang="ts">
	import { formatDateTime } from '$lib/settings/dateFormat';
	import type { LibraryFolder } from '$lib/settings/types';

	interface Props {
		folders: LibraryFolder[];
		scanningLibraryFolderId?: string;
		onScan: (_id: string) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { folders, scanningLibraryFolderId, onScan, onDelete }: Props = $props();
</script>

<div class="panel" aria-labelledby="library-folder-list-title">
	<h2 id="library-folder-list-title">Library folders</h2>
	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th>Path</th>
					<th>Added</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each folders as folder (folder.id)}
					<tr>
						<td>{folder.path}</td>
						<td>{formatDateTime(folder.createdAt)}</td>
						<td class="row-actions">
							<button
								type="button"
								class="secondary icon-button"
								aria-label={`Scan ${folder.path}`}
								disabled={scanningLibraryFolderId === folder.id}
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
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="3" class="empty">No library folders configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
