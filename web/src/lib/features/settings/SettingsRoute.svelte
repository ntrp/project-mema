<script lang="ts">
	import SettingsArea from '$lib/components/app/settings/SettingsArea.svelte';
	import type { SettingsSection } from '$lib/settings/types';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		section: SettingsSection;
	}

	let { section }: Props = $props();
	const app = getAppShellContext();
</script>

{#if app.isAdmin}
	<SettingsArea
		bind:downloadForm={app.downloadForm}
		bind:indexerForm={app.indexerForm}
		bind:libraryFolderForm={app.libraryFolderForm}
		bind:pathMappingForm={app.pathMappingForm}
		bind:customFormatForm={app.customFormatForm}
		bind:tagForm={app.tagForm}
		bind:userForm={app.userForm}
		activeSection={section}
		downloadClients={app.downloadClients}
		indexers={app.indexers}
		metadataProviders={app.metadataProviders}
		metadataCache={app.metadataCache}
		libraryFolders={app.libraryFolders}
		pathMappings={app.pathMappings}
		mediaProfiles={app.mediaProfiles}
		customFormats={app.customFormats}
		users={app.users}
		tags={app.tags}
		currentUser={app.currentUser}
		libraryScansByFolder={app.libraryScansByFolder}
		openLibraryFolderId={app.openLibraryFolderId}
		savingDownloadClient={app.savingDownloadClient}
		savingIndexer={app.savingIndexer}
		savingMetadataProviderId={app.savingMetadataProviderId}
		loadingMetadataCache={app.loadingMetadataCache}
		clearingMetadataCache={app.clearingMetadataCache}
		savingLibraryFolder={app.savingLibraryFolder}
		scanningLibraryFolderId={app.scanningLibraryFolderId}
		savingPathMapping={app.savingPathMapping}
		deletingPathMappingId={app.deletingPathMappingId}
		deletingMediaProfileId={app.deletingMediaProfileId}
		savingCustomFormat={app.savingCustomFormat}
		deletingCustomFormatId={app.deletingCustomFormatId}
		savingTag={app.savingTag}
		deletingTagId={app.deletingTagId}
		savingUser={app.savingUser}
		bind:metadataCachePattern={app.metadataCachePattern}
		testingIndexerId={app.testingIndexerId}
		testingMetadataProviderId={app.testingMetadataProviderId}
		indexerTests={app.indexerTests}
		metadataProviderTests={app.metadataProviderTests}
		onSaveDownloadClient={app.saveDownloadClient}
		onTestDownloadClientConfig={app.testDownloadClientConfig}
		onSaveIndexer={app.saveIndexer}
		onSaveMetadataProvider={app.saveMetadataProvider}
		onRefreshMetadataCache={app.refreshMetadataCache}
		onClearMetadataCache={app.clearMetadataCache}
		onClearMetadataCachePattern={app.clearMetadataCachePattern}
		onSaveLibraryFolder={app.saveLibraryFolder}
		onScanLibraryFolder={app.scanLibraryFolder}
		onSavePathMapping={app.savePathMapping}
		onSaveCustomFormat={app.saveCustomFormat}
		onImportCustomFormat={app.importCustomFormat}
		onSaveTag={app.saveTag}
		onSaveUser={app.saveUser}
		onCancelDownloadClient={app.cancelDownloadClient}
		onCancelIndexer={app.cancelIndexer}
		onCancelCustomFormat={app.cancelCustomFormat}
		onCancelTag={app.cancelTag}
		onCancelUser={app.cancelUser}
		onEditDownloadClient={app.editDownloadClient}
		onEditIndexer={app.editIndexer}
		onEditCustomFormat={app.editCustomFormat}
		onEditUser={app.editUser}
		onEditTag={app.editTag}
		onDeleteDownloadClient={app.deleteDownloadClient}
		onDeleteIndexer={app.deleteIndexer}
		onDeleteLibraryFolder={app.deleteLibraryFolder}
		onDeletePathMapping={app.deletePathMapping}
		onDeleteMediaProfile={app.deleteMediaProfile}
		onDeleteCustomFormat={app.deleteCustomFormat}
		onDeleteTag={app.deleteTag}
		onDeleteUser={app.deleteUser}
		onTestIndexer={app.testIndexer}
		onTestMetadataProvider={app.testMetadataProvider}
		onSearchLibraryMatch={app.searchLibraryMatch}
		onImportLibraryScanRows={app.importLibraryScanRows}
	/>
{/if}
