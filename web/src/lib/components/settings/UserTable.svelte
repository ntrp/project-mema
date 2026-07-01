<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { formatDate } from '$lib/settings/dateFormat';
	import SettingsRowActionButton from './shared/SettingsRowActionButton.svelte';
	import type { ManagedUser } from '$lib/settings/types';

	interface Props {
		users: ManagedUser[];
		currentUserId?: string;
		onEdit: (_user: ManagedUser) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { users, currentUserId, onEdit, onDelete }: Props = $props();
</script>

<Card class="p-0" aria-label="Users">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Username</Table.Head>
				<Table.Head>Role</Table.Head>
				<Table.Head>Created</Table.Head>
				<Table.Head class="text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each users as user (user.id)}
				<Table.Row>
					<Table.Cell>
						<div class="flex items-center gap-2">
							{user.username}
							{#if user.id === currentUserId}
								<Badge variant="secondary">Current</Badge>
							{/if}
						</div>
					</Table.Cell>
					<Table.Cell>{user.role}</Table.Cell>
					<Table.Cell>{formatDate(user.createdAt)}</Table.Cell>
					<Table.Cell>
						<div class="flex justify-end gap-2">
							<SettingsRowActionButton
								label={`Edit ${user.username}`}
								icon="edit"
								onclick={() => onEdit(user)}
							/>
							<SettingsRowActionButton
								label={`Delete ${user.username}`}
								icon="delete"
								variant="destructive"
								disabled={user.id === currentUserId}
								onclick={() => onDelete(user.id)}
							/>
						</div>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={4} class="py-8 text-center text-muted-foreground">
						No users configured
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
