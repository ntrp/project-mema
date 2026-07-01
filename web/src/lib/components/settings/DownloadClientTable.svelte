<script lang="ts">
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import SettingsRowActionButton from './shared/SettingsRowActionButton.svelte';
	import type { DownloadClient } from '$lib/settings/types';

	interface Props {
		clients: DownloadClient[];
		onEdit: (_client: DownloadClient) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { clients, onEdit, onDelete }: Props = $props();
</script>

<Card class="p-0" aria-label="Download clients">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Name</Table.Head>
				<Table.Head>Type</Table.Head>
				<Table.Head>Base URL</Table.Head>
				<Table.Head>Priority</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each clients as item (item.id)}
				<Table.Row>
					<Table.Cell>{item.name}</Table.Cell>
					<Table.Cell>{item.type}</Table.Cell>
					<Table.Cell class="max-w-[320px] truncate">{item.baseUrl}</Table.Cell>
					<Table.Cell>{item.priority}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
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
					<Table.Cell colspan={5} class="py-8 text-center text-muted-foreground">
						No download clients configured
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
