<script lang="ts">
	import CustomFormatsSettings from '$lib/components/settings/custom-formats/CustomFormatsSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<PageHeading eyebrow="Settings" title="Custom formats" titleId="settings-title" />
	<div class="space-y-4">
		<CustomFormatsSettings
			formats={app.customFormats}
			bind:form={app.customFormatForm}
			saving={app.savingCustomFormat}
			deletingId={app.deletingCustomFormatId}
			onSave={app.saveCustomFormat}
			onImport={app.importCustomFormat}
			onCancel={app.cancelCustomFormat}
			onEdit={app.editCustomFormat}
			onDelete={app.deleteCustomFormat}
		/>
	</div>
</SettingsRouteFrame>
