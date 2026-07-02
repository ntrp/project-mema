<script lang="ts">
	import DownloadClientsSettingsSection from '$lib/components/settings/download-clients/DownloadClientsSettingsSection.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<DownloadClientsSettingsSection
		clients={app.downloadClients}
		bind:form={app.downloadForm}
		saving={app.savingDownloadClient}
		onSave={app.saveDownloadClient}
		onTestConfig={app.testDownloadClientConfig}
		onCancel={app.cancelDownloadClient}
		onEdit={app.editDownloadClient}
		onDelete={app.deleteDownloadClient}
	/>
</SettingsRouteFrame>
