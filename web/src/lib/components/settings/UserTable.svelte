<script lang="ts">
	import { formatDate } from '$lib/settings/dateFormat';
	import type { ManagedUser } from '$lib/settings/types';

	interface Props {
		users: ManagedUser[];
		currentUserId?: string;
		onEdit: (_user: ManagedUser) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { users, currentUserId, onEdit, onDelete }: Props = $props();
</script>

<div class="panel" aria-label="Users">
	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th>Username</th>
					<th>Role</th>
					<th>Created</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each users as user (user.id)}
					<tr>
						<td>
							{user.username}
							{#if user.id === currentUserId}
								<span class="status-pill">Current</span>
							{/if}
						</td>
						<td>{user.role}</td>
						<td>{formatDate(user.createdAt)}</td>
						<td class="row-actions">
							<button
								type="button"
								class="secondary icon-button"
								aria-label={`Edit ${user.username}`}
								onclick={() => onEdit(user)}
							>
								<span class="app-icon" aria-hidden="true">edit</span>
							</button>
							<button
								type="button"
								class="danger icon-button"
								disabled={user.id === currentUserId}
								aria-label={`Delete ${user.username}`}
								onclick={() => onDelete(user.id)}
							>
								<span class="app-icon" aria-hidden="true">delete</span>
							</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="4" class="empty">No users configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
