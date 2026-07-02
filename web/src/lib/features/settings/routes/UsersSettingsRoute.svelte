<script lang="ts">
	import UsersSettingsSection from '$lib/components/settings/users/UsersSettingsSection.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<UsersSettingsSection
		users={app.users}
		currentUser={app.currentUser}
		bind:form={app.userForm}
		saving={app.savingUser}
		onSave={app.saveUser}
		onCancel={app.cancelUser}
		onEdit={app.editUser}
		onDelete={app.deleteUser}
	/>
</SettingsRouteFrame>
