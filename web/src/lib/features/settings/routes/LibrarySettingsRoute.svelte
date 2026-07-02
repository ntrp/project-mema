<script lang="ts">
	import LibrarySettingsSection from '$lib/components/settings/library/LibrarySettingsSection.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import SettingsRouteFrame from '../SettingsRouteFrame.svelte';

	const app = getAppShellContext();
</script>

<SettingsRouteFrame>
	<LibrarySettingsSection
		folders={app.libraryFolders}
		bind:form={app.libraryFolderForm}
		pathMappings={app.pathMappings}
		bind:pathMappingForm={app.pathMappingForm}
		scansByFolder={app.libraryScansByFolder}
		openFolderId={app.openLibraryFolderId}
		qualityProfiles={app.mediaProfiles}
		saving={app.savingLibraryFolder}
		scanningLibraryFolderId={app.scanningLibraryFolderId}
		savingPathMapping={app.savingPathMapping}
		deletingPathMappingId={app.deletingPathMappingId}
		onSave={app.saveLibraryFolder}
		onScan={app.scanLibraryFolder}
		onDelete={app.deleteLibraryFolder}
		onSavePathMapping={app.savePathMapping}
		onDeletePathMapping={app.deletePathMapping}
		onSearchMatch={app.searchLibraryMatch}
		onImport={app.importLibraryScanRows}
	/>
</SettingsRouteFrame>
