<script lang="ts">
	import SubtitleProviderSettings from '$lib/components/settings/subtitles/SubtitleProviderSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<PageHeading eyebrow="Settings" title="Subtitles" titleId="settings-title" />
	<div class="space-y-4">
		<SubtitleProviderSettings
			providers={app.subtitleProviders}
			onSave={app.saveSubtitleProvider}
			onDelete={app.deleteSubtitleProvider}
			onTest={app.testSubtitleProvider}
			testingId={app.testingSubtitleProviderId}
			savingId={app.savingSubtitleProviderId}
			testResults={app.subtitleProviderTests}
		/>
	</div>
</SettingsRouteFrame>
