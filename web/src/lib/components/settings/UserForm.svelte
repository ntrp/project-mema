<script lang="ts">
	import type { UserForm, UserRole } from '$lib/settings/types';

	interface Props {
		form: UserForm;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), saving, onSave, onCancel }: Props = $props();
	const roles: UserRole[] = ['admin', 'user'];
</script>

<div class="panel" aria-labelledby="user-form-title">
	<div class="section-heading">
		<h2 id="user-form-title">{form.id ? 'Edit user' : 'Add user'}</h2>
		{#if form.id}
			<button type="button" class="secondary" onclick={onCancel}>Cancel</button>
		{/if}
	</div>

	<form class="settings-form" onsubmit={onSave}>
		<label>
			<span>Username</span>
			<input bind:value={form.username} required maxlength="200" autocomplete="username" />
		</label>
		<label>
			<span>Role</span>
			<select bind:value={form.role}>
				{#each roles as role (role)}
					<option value={role}>{role}</option>
				{/each}
			</select>
		</label>
		<label class="wide">
			<span>{form.id ? 'New password' : 'Password'}</span>
			<input
				bind:value={form.password}
				required={!form.id}
				minlength="8"
				maxlength="1024"
				autocomplete={form.id ? 'new-password' : 'new-password'}
				type="password"
				placeholder={form.id ? 'Leave blank to keep current password' : ''}
			/>
		</label>
		<button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save user'}</button>
	</form>
</div>
