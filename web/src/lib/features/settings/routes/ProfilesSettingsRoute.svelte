<script lang="ts">
	import MediaProfilesSettings from '$lib/components/settings/profiles/MediaProfilesSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<PageHeading eyebrow="Settings" title="Profiles" titleId="settings-title" />
	<div class="space-y-4">
		<MediaProfilesSettings
			profiles={app.mediaProfiles}
			deletingId={app.deletingMediaProfileId}
			onDelete={app.deleteMediaProfile}
		/>
	</div>
</SettingsRouteFrame>
