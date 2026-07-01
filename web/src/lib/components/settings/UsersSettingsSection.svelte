<script lang="ts">
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import UserForm from '$lib/components/settings/UserForm.svelte';
	import UserTable from '$lib/components/settings/UserTable.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { emptyUserForm } from '$lib/settings/forms';
	import type { ManagedUser, UserForm as UserFormValue, UserSummary } from '$lib/settings/types';

	interface Props {
		users: ManagedUser[];
		currentUser?: UserSummary;
		form: UserFormValue;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onEdit: (_user: ManagedUser) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let {
		users,
		currentUser,
		form = $bindable(),
		saving,
		onSave,
		onCancel,
		onEdit,
		onDelete
	}: Props = $props();

	let modalOpen = $state(false);

	function openModal() {
		form = emptyUserForm();
		modalOpen = true;
	}

	function editUser(user: ManagedUser) {
		onEdit(user);
		modalOpen = true;
	}

	function closeModal() {
		onCancel();
		modalOpen = false;
	}

	async function save(event: SubmitEvent) {
		await onSave(event);
		if (!form.id && form.username === '' && form.password === '') {
			modalOpen = false;
		}
	}
</script>

<PageHeading eyebrow="Settings" title="Users" titleId="settings-title" />
<div class="space-y-4">
	<div class="flex justify-end">
		<SettingsAddButton label="Add user" onclick={openModal} />
	</div>
	<UserTable {users} currentUserId={currentUser?.id} onEdit={editUser} {onDelete} />
	{#if modalOpen}
		<SettingsFormModal title={form.id ? 'Edit user' : 'Add user'} onClose={closeModal}>
			<UserForm bind:form {saving} onSave={save} onCancel={closeModal} />
		</SettingsFormModal>
	{/if}
</div>
