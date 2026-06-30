<script lang="ts">
	import IntegrationTestStatus from './IntegrationTestStatus.svelte';
	import type { Indexer, IntegrationTestResults } from '$lib/settings/types';

	interface Props {
		indexers: Indexer[];
		onEdit: (_indexer: Indexer) => void;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		testingId?: string;
		testResults: IntegrationTestResults;
	}

	let { indexers, onEdit, onDelete, onTest, testingId, testResults }: Props = $props();
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
						<td>
							<IntegrationTestStatus
								enabled={item.enabled}
								result={testResults[item.id]}
								testing={testingId === item.id}
							/>
						</td>
						<td class="row-actions">
							<button
								type="button"
								class="secondary"
								disabled={testingId === item.id}
								onclick={() => onTest(item.id)}
							>
								{testingId === item.id ? 'Testing' : 'Test'}
							</button>
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
						<td colspan="7" class="empty">No indexers configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
