<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import type { UserForm, UserRole } from '$lib/settings/types';

	interface Props {
		form: UserForm;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), saving, onSave, onCancel }: Props = $props();
	const roles: { value: UserRole; label: string }[] = [
		{ value: 'admin', label: 'admin' },
		{ value: 'user', label: 'user' }
	];
</script>

<Card.Root aria-labelledby="user-form-title">
	<Card.Header>
		<Card.Title id="user-form-title">{form.id ? 'Edit user' : 'Add user'}</Card.Title>
		{#if form.id}
			<Card.Action>
				<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
			</Card.Action>
		{/if}
	</Card.Header>

	<Card.Content>
		<form class="grid gap-4 sm:grid-cols-2" onsubmit={onSave}>
			<div class="space-y-2">
				<Label for="user-username">Username</Label>
				<Input
					id="user-username"
					bind:value={form.username}
					required
					maxlength={200}
					autocomplete="username"
				/>
			</div>
			<div class="space-y-2">
				<Label>Role</Label>
				<SettingsSelect
					value={form.role}
					options={roles}
					onValueChange={(value) => (form.role = value as UserRole)}
				/>
			</div>
			<div class="space-y-2 sm:col-span-2">
				<Label for="user-password">{form.id ? 'New password' : 'Password'}</Label>
				<Input
					id="user-password"
					bind:value={form.password}
					required={!form.id}
					minlength={8}
					maxlength={1024}
					autocomplete={form.id ? 'new-password' : 'new-password'}
					type="password"
					placeholder={form.id ? 'Leave blank to keep current password' : ''}
				/>
			</div>
			<Button class="w-fit" type="submit" disabled={saving}
				>{saving ? 'Saving' : 'Save user'}</Button
			>
		</form>
	</Card.Content>
</Card.Root>
