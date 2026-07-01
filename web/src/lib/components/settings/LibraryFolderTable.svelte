<script lang="ts">
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { formatDateTime } from '$lib/settings/dateFormat';
	import SettingsRowActionButton from './shared/SettingsRowActionButton.svelte';
	import type { LibraryFolder } from '$lib/settings/types';

	interface Props {
		folders: LibraryFolder[];
		scanningLibraryFolderId?: string;
		onScan: (_id: string) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { folders, scanningLibraryFolderId, onScan, onDelete }: Props = $props();
</script>

<Card class="gap-0 p-0" aria-labelledby="library-folder-list-title">
	<div class="border-b px-4 py-3">
		<h2 id="library-folder-list-title" class="text-lg font-semibold">Library folders</h2>
	</div>
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Path</Table.Head>
				<Table.Head>Added</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each folders as folder (folder.id)}
				<Table.Row>
					<Table.Cell class="max-w-[520px] truncate">{folder.path}</Table.Cell>
					<Table.Cell>{formatDateTime(folder.createdAt)}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
							<SettingsRowActionButton
								label={`Scan ${folder.path}`}
								icon="sync"
								disabled={scanningLibraryFolderId === folder.id}
								onclick={() => onScan(folder.id)}
							/>
							<SettingsRowActionButton
								label={`Delete ${folder.path}`}
								icon="delete"
								variant="destructive"
								onclick={() => onDelete(folder.id)}
							/>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={3} class="py-8 text-center text-muted-foreground">
						No library folders configured
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
