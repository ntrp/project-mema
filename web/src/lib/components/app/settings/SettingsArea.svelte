<script lang="ts">
	import DownloadClientsSettingsSection from '$lib/components/settings/DownloadClientsSettingsSection.svelte';
	import IndexersSettingsSection from '$lib/components/settings/IndexersSettingsSection.svelte';
	import LibrarySettingsSection from '$lib/components/settings/LibrarySettingsSection.svelte';
	import SystemGeneralSettings from '$lib/components/settings/SystemGeneralSettings.svelte';
	import UsersSettingsSection from '$lib/components/settings/UsersSettingsSection.svelte';
	import SettingsStaticPanels from './SettingsStaticPanels.svelte';
	import { isStaticSettingsSection, type SettingsAreaProps } from './settingsAreaTypes';

	let {
		activeSection,
		downloadClients,
		indexers,
		metadataProviders,
		metadataCache,
		libraryFolders,
		pathMappings,
		mediaProfiles,
		customFormats,
		users,
		tags,
		currentUser,
		libraryScansByFolder,
		openLibraryFolderId,
		downloadForm = $bindable(),
		indexerForm = $bindable(),
		libraryFolderForm = $bindable(),
		pathMappingForm = $bindable(),
		mediaProfileForm = $bindable(),
		customFormatForm = $bindable(),
		tagForm = $bindable(),
		userForm = $bindable(),
		savingDownloadClient,
		savingIndexer,
		savingMetadataProviderId,
		loadingMetadataCache,
		clearingMetadataCache,
		metadataCachePattern = $bindable(),
		savingLibraryFolder,
		savingPathMapping,
		deletingPathMappingId,
		savingMediaProfile,
		deletingMediaProfileId,
		savingCustomFormat,
		deletingCustomFormatId,
		savingTag,
		deletingTagId,
		savingUser,
		scanningLibraryFolderId,
		testingIndexerId,
		testingMetadataProviderId,
		indexerTests,
		metadataProviderTests,
		onSaveDownloadClient,
		onTestDownloadClientConfig,
		onSaveIndexer,
		onSaveMetadataProvider,
		onRefreshMetadataCache,
		onClearMetadataCache,
		onClearMetadataCachePattern,
		onSaveLibraryFolder,
		onScanLibraryFolder,
		onSavePathMapping,
		onSaveMediaProfile,
		onSaveCustomFormat,
		onImportCustomFormat,
		onSaveTag,
		onSaveUser,
		onCancelDownloadClient,
		onCancelIndexer,
		onCancelMediaProfile,
		onCancelCustomFormat,
		onCancelTag,
		onCancelUser,
		onEditDownloadClient,
		onEditIndexer,
		onEditMediaProfile,
		onEditCustomFormat,
		onEditTag,
		onEditUser,
		onDeleteDownloadClient,
		onDeleteIndexer,
		onDeleteLibraryFolder,
		onDeletePathMapping,
		onDeleteMediaProfile,
		onDeleteCustomFormat,
		onDeleteTag,
		onDeleteUser,
		onTestIndexer,
		onTestMetadataProvider,
		onSearchLibraryMatch,
		onImportLibraryScanRows
	}: SettingsAreaProps = $props();
</script>

<section class="workspace-main" aria-labelledby="settings-title">
	{#if activeSection === 'general'}
		<SystemGeneralSettings />
	{:else if activeSection === 'download-clients'}
		<DownloadClientsSettingsSection
			clients={downloadClients}
			bind:form={downloadForm}
			saving={savingDownloadClient}
			onSave={onSaveDownloadClient}
			onTestConfig={onTestDownloadClientConfig}
			onCancel={onCancelDownloadClient}
			onEdit={onEditDownloadClient}
			onDelete={onDeleteDownloadClient}
		/>
	{:else if activeSection === 'indexers'}
		<IndexersSettingsSection
			{indexers}
			bind:form={indexerForm}
			saving={savingIndexer}
			testingId={testingIndexerId}
			testResults={indexerTests}
			onSave={onSaveIndexer}
			onCancel={onCancelIndexer}
			onEdit={onEditIndexer}
			onDelete={onDeleteIndexer}
			onTest={onTestIndexer}
		/>
	{:else if isStaticSettingsSection(activeSection)}
		<SettingsStaticPanels
			{activeSection}
			{metadataProviders}
			{metadataCache}
			{mediaProfiles}
			{customFormats}
			{tags}
			bind:metadataCachePattern
			bind:mediaProfileForm
			bind:customFormatForm
			bind:tagForm
			{savingMetadataProviderId}
			{testingMetadataProviderId}
			{loadingMetadataCache}
			{clearingMetadataCache}
			{savingMediaProfile}
			{deletingMediaProfileId}
			{savingCustomFormat}
			{deletingCustomFormatId}
			{savingTag}
			{deletingTagId}
			{metadataProviderTests}
			{onSaveMetadataProvider}
			{onTestMetadataProvider}
			{onRefreshMetadataCache}
			{onClearMetadataCache}
			{onClearMetadataCachePattern}
			{onSaveMediaProfile}
			{onCancelMediaProfile}
			{onEditMediaProfile}
			{onDeleteMediaProfile}
			{onSaveCustomFormat}
			{onImportCustomFormat}
			{onCancelCustomFormat}
			{onEditCustomFormat}
			{onDeleteCustomFormat}
			{onSaveTag}
			{onCancelTag}
			{onEditTag}
			{onDeleteTag}
		/>
	{:else if activeSection === 'users'}
		<UsersSettingsSection
			{users}
			{currentUser}
			bind:form={userForm}
			saving={savingUser}
			onSave={onSaveUser}
			onCancel={onCancelUser}
			onEdit={onEditUser}
			onDelete={onDeleteUser}
		/>
	{:else}
		<LibrarySettingsSection
			folders={libraryFolders}
			bind:form={libraryFolderForm}
			{pathMappings}
			bind:pathMappingForm
			scansByFolder={libraryScansByFolder}
			openFolderId={openLibraryFolderId}
			qualityProfiles={mediaProfiles}
			saving={savingLibraryFolder}
			{scanningLibraryFolderId}
			{savingPathMapping}
			{deletingPathMappingId}
			onSave={onSaveLibraryFolder}
			onScan={onScanLibraryFolder}
			onDelete={onDeleteLibraryFolder}
			{onSavePathMapping}
			{onDeletePathMapping}
			onSearchMatch={onSearchLibraryMatch}
			onImport={onImportLibraryScanRows}
		/>
	{/if}
</section>
