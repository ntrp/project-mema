<script lang="ts">
	import type { LibraryFolder } from '$lib/settings/types';

	interface Props {
		folders: LibraryFolder[];
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { folders, onDelete }: Props = $props();
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
						<td>{new Date(folder.createdAt).toLocaleString()}</td>
						<td class="row-actions">
							<button type="button" class="danger" onclick={() => onDelete(folder.id)}
								>Delete</button
							>
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
