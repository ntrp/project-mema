<script lang="ts">
	import MetadataProviderSettings from '$lib/components/settings/metadata/MetadataProviderSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<PageHeading eyebrow="Settings" title="Metadata" titleId="settings-title" />
	<div class="space-y-4">
		<MetadataProviderSettings
			metadataProviders={app.metadataProviders}
			onSave={app.saveMetadataProvider}
			onTest={app.testMetadataProvider}
			testingId={app.testingMetadataProviderId}
			savingId={app.savingMetadataProviderId}
			testResults={app.metadataProviderTests}
		/>
	</div>
</SettingsRouteFrame>
