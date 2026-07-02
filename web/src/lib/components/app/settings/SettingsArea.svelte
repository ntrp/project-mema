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
		indexerSearch,
		metadataProviders,
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
		customFormatForm = $bindable(),
		tagForm = $bindable(),
		userForm = $bindable(),
		savingDownloadClient,
		savingIndexer,
		clearingIndexerSearchCache,
		savingIndexerSearchSettings,
		savingMetadataProviderId,
		savingLibraryFolder,
		savingPathMapping,
		deletingPathMappingId,
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
		onClearIndexerSearchCache,
		onSaveIndexerSearchSettings,
		onSaveMetadataProvider,
		onSaveLibraryFolder,
		onScanLibraryFolder,
		onSavePathMapping,
		onSaveCustomFormat,
		onImportCustomFormat,
		onSaveTag,
		onSaveUser,
		onCancelDownloadClient,
		onCancelIndexer,
		onCancelCustomFormat,
		onCancelTag,
		onCancelUser,
		onEditDownloadClient,
		onEditIndexer,
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

<section class="grid min-w-0 gap-[18px]" aria-labelledby="settings-title">
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
			{indexerSearch}
			bind:form={indexerForm}
			saving={savingIndexer}
			{clearingIndexerSearchCache}
			{savingIndexerSearchSettings}
			testingId={testingIndexerId}
			testResults={indexerTests}
			onSave={onSaveIndexer}
			onCancel={onCancelIndexer}
			onEdit={onEditIndexer}
			onDelete={onDeleteIndexer}
			onTest={onTestIndexer}
			{onClearIndexerSearchCache}
			{onSaveIndexerSearchSettings}
		/>
	{:else if isStaticSettingsSection(activeSection)}
		<SettingsStaticPanels
			{activeSection}
			{metadataProviders}
			{mediaProfiles}
			{customFormats}
			{tags}
			bind:customFormatForm
			bind:tagForm
			{savingMetadataProviderId}
			{testingMetadataProviderId}
			{deletingMediaProfileId}
			{savingCustomFormat}
			{deletingCustomFormatId}
			{savingTag}
			{deletingTagId}
			{metadataProviderTests}
			{onSaveMetadataProvider}
			{onTestMetadataProvider}
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
