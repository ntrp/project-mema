<script lang="ts">
	import IntegrationTestStatus from './IntegrationTestStatus.svelte';
	import type { DownloadClient, IntegrationTestResults } from '$lib/settings/types';

	interface Props {
		clients: DownloadClient[];
		onEdit: (_client: DownloadClient) => void;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		testingId?: string;
		testResults: IntegrationTestResults;
	}

	let { clients, onEdit, onDelete, onTest, testingId, testResults }: Props = $props();
</script>

<div class="panel" aria-labelledby="download-client-list-title">
	<h2 id="download-client-list-title">Download clients</h2>
	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Type</th>
					<th>Base URL</th>
					<th>Priority</th>
					<th>Status</th>
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
							<button type="button" class="secondary" onclick={() => onEdit(item)}>Edit</button>
							<button type="button" class="danger" onclick={() => onDelete(item.id)}>Delete</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="6" class="empty">No download clients configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
