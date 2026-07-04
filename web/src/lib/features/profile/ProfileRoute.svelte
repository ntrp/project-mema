<script lang="ts">
	import { onMount } from 'svelte';
	import ProfilePage from '$lib/components/profile/ProfilePage.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	const app = getAppShellContext();

	onMount(() => {
		if (!app.profile && !app.loadingProfile) {
			void app.loadProfile();
		}
	});
</script>

<ProfilePage
	profile={app.profile}
	loading={app.loadingProfile}
	saving={app.savingProfile}
	errorMessage={app.profileErrorMessage}
	onRetry={app.loadProfile}
	onSave={app.saveProfile}
/>
