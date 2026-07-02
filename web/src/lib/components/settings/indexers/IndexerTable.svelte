<script lang="ts">
	import IndexerHealthStatus from './IndexerHealthStatus.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import SettingsRowActionButton from '../shared/SettingsRowActionButton.svelte';
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

<Card class="p-0" aria-label="Indexers">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head>Type</Table.Head>
				<Table.Head>Base URL</Table.Head>
				<Table.Head>Categories</Table.Head>
				<Table.Head>Priority</Table.Head>
				<Table.Head>Status</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each indexers as item (item.id)}
				<Table.Row>
					<Table.Cell>{item.name}</Table.Cell>
					<Table.Cell>{item.type}</Table.Cell>
					<Table.Cell class="max-w-[320px] truncate">{item.baseUrl}</Table.Cell>
					<Table.Cell>{(item.categories ?? []).join(', ') || '-'}</Table.Cell>
					<Table.Cell>{item.priority}</Table.Cell>
					<Table.Cell>
						<IndexerHealthStatus
							indexer={item}
							result={testResults[item.id]}
							checking={testingId === item.id}
						/>
					</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
							<Button
								type="button"
								variant="outline"
								size="sm"
								disabled={testingId === item.id}
								onclick={() => onTest(item.id)}
							>
								{testingId === item.id ? 'Checking' : 'Check'}
							</Button>
							<SettingsRowActionButton
								label={`Edit ${item.name}`}
								icon="edit"
								onclick={() => onEdit(item)}
							/>
							<SettingsRowActionButton
								label={`Delete ${item.name}`}
								icon="delete"
								variant="destructive"
								onclick={() => onDelete(item.id)}
							/>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={7} class="py-8 text-center text-muted-foreground">
						No indexers configured
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
