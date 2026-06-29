<script lang="ts">
	import type { ManagedUser } from '$lib/settings/types';

	interface Props {
		users: ManagedUser[];
		currentUserId?: string;
		onEdit: (_user: ManagedUser) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { users, currentUserId, onEdit, onDelete }: Props = $props();
</script>

<div class="panel" aria-labelledby="user-list-title">
	<h2 id="user-list-title">Users</h2>
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
						<td>{new Date(user.createdAt).toLocaleDateString()}</td>
						<td class="row-actions">
							<button type="button" class="secondary" onclick={() => onEdit(user)}>Edit</button>
							<button
								type="button"
								class="danger"
								disabled={user.id === currentUserId}
								onclick={() => onDelete(user.id)}
							>
								Delete
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
