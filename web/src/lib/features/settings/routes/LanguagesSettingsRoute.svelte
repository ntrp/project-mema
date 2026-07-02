<script lang="ts">
	import LanguagesSettingsSection from '$lib/components/settings/languages/LanguagesSettingsSection.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<LanguagesSettingsSection
		languages={app.languages}
		bind:form={app.languageForm}
		saving={app.savingLanguage}
		deletingCode={app.deletingLanguageCode}
		onSave={app.saveLanguage}
		onCancel={app.cancelLanguage}
		onEdit={app.editLanguage}
		onDelete={app.deleteLanguage}
	/>
</SettingsRouteFrame>
