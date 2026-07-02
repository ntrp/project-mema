<script lang="ts">
	import IndexersSettingsSection from '$lib/components/settings/indexers/IndexersSettingsSection.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<IndexersSettingsSection
		indexers={app.indexers}
		indexerSearch={app.indexerSearch}
		bind:form={app.indexerForm}
		saving={app.savingIndexer}
		clearingIndexerSearchCache={app.clearingIndexerSearchCache}
		savingIndexerSearchSettings={app.savingIndexerSearchSettings}
		testingId={app.testingIndexerId}
		testResults={app.indexerTests}
		onSave={app.saveIndexer}
		onCancel={app.cancelIndexer}
		onEdit={app.editIndexer}
		onDelete={app.deleteIndexer}
		onTest={app.testIndexer}
		onClearIndexerSearchCache={app.clearIndexerSearchCache}
		onSaveIndexerSearchSettings={app.saveIndexerSearchSettings}
	/>
</SettingsRouteFrame>
