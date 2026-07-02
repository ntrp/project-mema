<script lang="ts">
	import TagSettings from '$lib/components/settings/tags/TagSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<PageHeading eyebrow="Settings" title="Tags" titleId="settings-title" />
	<div class="space-y-4">
		<TagSettings
			tags={app.tags}
			bind:form={app.tagForm}
			saving={app.savingTag}
			deletingId={app.deletingTagId}
			onSave={app.saveTag}
			onCancel={app.cancelTag}
			onEdit={app.editTag}
			onDelete={app.deleteTag}
		/>
	</div>
</SettingsRouteFrame>
