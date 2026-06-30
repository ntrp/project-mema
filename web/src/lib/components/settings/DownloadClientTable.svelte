<script lang="ts">
	import type { DownloadClient } from '$lib/settings/types';

	interface Props {
		clients: DownloadClient[];
		onEdit: (_client: DownloadClient) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { clients, onEdit, onDelete }: Props = $props();
</script>

<div class="panel" aria-label="Download clients">
	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Type</th>
					<th>Base URL</th>
					<th>Priority</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each clients as item (item.id)}
					<tr>
						<td>{item.name}</td>
						<td>{item.type}</td>
						<td>{item.baseUrl}</td>
						<td>{item.priority}</td>
						<td class="row-actions">
							<button
								type="button"
								class="secondary icon-button"
								aria-label={`Edit ${item.name}`}
								onclick={() => onEdit(item)}
							>
								<span class="app-icon" aria-hidden="true">edit</span>
							</button>
							<button
								type="button"
								class="danger icon-button"
								aria-label={`Delete ${item.name}`}
								onclick={() => onDelete(item.id)}
							>
								<span class="app-icon" aria-hidden="true">delete</span>
							</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="5" class="empty">No download clients configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
