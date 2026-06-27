<script lang="ts">
	import type { Indexer } from '$lib/settings/types';

	interface Props {
		indexers: Indexer[];
		onEdit: (_indexer: Indexer) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { indexers, onEdit, onDelete }: Props = $props();
</script>

<div class="panel" aria-labelledby="indexer-list-title">
	<h2 id="indexer-list-title">Indexers</h2>
	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Type</th>
					<th>Base URL</th>
					<th>Categories</th>
					<th>Priority</th>
					<th>Status</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each indexers as item (item.id)}
					<tr>
						<td>{item.name}</td>
						<td>{item.type}</td>
						<td>{item.baseUrl}</td>
						<td>{(item.categories ?? []).join(', ') || '-'}</td>
						<td>{item.priority}</td>
						<td>{item.enabled ? 'Enabled' : 'Disabled'}</td>
						<td class="row-actions">
							<button type="button" class="secondary" onclick={() => onEdit(item)}>Edit</button>
							<button type="button" class="danger" onclick={() => onDelete(item.id)}>Delete</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="7" class="empty">No indexers configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
