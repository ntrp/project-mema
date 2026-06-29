<script lang="ts">
	import DownloadClientsSettingsSection from '$lib/components/settings/DownloadClientsSettingsSection.svelte';
	import FileNamingSettings from '$lib/components/settings/FileNamingSettings.svelte';
	import IndexersSettingsSection from '$lib/components/settings/IndexersSettingsSection.svelte';
	import LibrarySettingsSection from '$lib/components/settings/LibrarySettingsSection.svelte';
	import MetadataCacheSettings from '$lib/components/settings/MetadataCacheSettings.svelte';
	import MetadataProviderSettings from '$lib/components/settings/MetadataProviderSettings.svelte';
	import MediaProfilesSettings from '$lib/components/settings/MediaProfilesSettings.svelte';
	import QualitySizeSettings from '$lib/components/settings/QualitySizeSettings.svelte';
	import SystemLogsSettings from '$lib/components/settings/SystemLogsSettings.svelte';
	import TagSettings from '$lib/components/settings/TagSettings.svelte';
	import UsersSettingsSection from '$lib/components/settings/UsersSettingsSection.svelte';
	import type {
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		Indexer,
		IndexerForm as IndexerFormValue,
		IntegrationTestResponse,
		IntegrationTestResults,
		LibraryFolder,
		LibraryFolderForm as LibraryFolderFormValue,
		LibraryMediaKind,
		LibraryScan,
		LibraryScanItem,
		LibraryScanItemMatchRequest,
		ManagedUser,
		MediaProfile,
		MediaProfileForm as MediaProfileFormValue,
		MediaSearchResult,
		MetadataCacheResponse,
		MetadataProvider,
		MetadataProviderForm as MetadataProviderFormValue,
		PathMapping,
		PathMappingForm,
		SettingsSection,
		Tag,
		TagForm,
		UserForm as UserFormValue,
		UserSummary
	} from '$lib/settings/types';

	interface Props {
		activeSection: SettingsSection;
		downloadClients: DownloadClient[];
		indexers: Indexer[];
		metadataProviders: MetadataProvider[];
		metadataCache: MetadataCacheResponse;
		libraryFolders: LibraryFolder[];
		pathMappings: PathMapping[];
		mediaProfiles: MediaProfile[];
		users: ManagedUser[];
		tags: Tag[];
		currentUser?: UserSummary;
		activeLibraryScan?: LibraryScan;
		downloadForm: DownloadClientFormValue;
		indexerForm: IndexerFormValue;
		libraryFolderForm: LibraryFolderFormValue;
		pathMappingForm: PathMappingForm;
		mediaProfileForm: MediaProfileFormValue;
		tagForm: TagForm;
		userForm: UserFormValue;
		savingDownloadClient: boolean;
		savingIndexer: boolean;
		savingMetadataProviderId?: string;
		loadingMetadataCache: boolean;
		clearingMetadataCache: boolean;
		metadataCachePattern: string;
		savingLibraryFolder: boolean;
		savingPathMapping: boolean;
		deletingPathMappingId?: string;
		savingMediaProfile: boolean;
		deletingMediaProfileId?: string;
		savingTag: boolean;
		deletingTagId?: string;
		savingUser: boolean;
		loadingLibraryScan: boolean;
		testingIndexerId?: string;
		testingMetadataProviderId?: string;
		indexerTests: IntegrationTestResults;
		metadataProviderTests: IntegrationTestResults;
		onSaveDownloadClient: (_event: SubmitEvent) => void | Promise<void>;
		onTestDownloadClientConfig: (
			_form: DownloadClientFormValue
		) => Promise<IntegrationTestResponse>;
		onSaveIndexer: (_event: SubmitEvent) => void | Promise<void>;
		onSaveMetadataProvider: (_form: MetadataProviderFormValue) => void | Promise<void>;
		onRefreshMetadataCache: () => void | Promise<void>;
		onClearMetadataCache: () => void | Promise<void>;
		onClearMetadataCachePattern: (_event: SubmitEvent) => void | Promise<void>;
		onSaveLibraryFolder: (_event: SubmitEvent) => void | Promise<void>;
		onSavePathMapping: (_event: SubmitEvent) => void | Promise<void>;
		onSaveMediaProfile: (_event: SubmitEvent) => void | Promise<void>;
		onSaveTag: (_event: SubmitEvent) => void | Promise<void>;
		onSaveUser: (_event: SubmitEvent) => void | Promise<void>;
		onCancelDownloadClient: () => void;
		onCancelIndexer: () => void;
		onCancelMediaProfile: () => void;
		onCancelTag: () => void;
		onCancelUser: () => void;
		onEditDownloadClient: (_client: DownloadClient) => void;
		onEditIndexer: (_indexer: Indexer) => void;
		onEditMediaProfile: (_profile: MediaProfile) => void;
		onEditTag: (_tag: Tag) => void;
		onEditUser: (_user: ManagedUser) => void;
		onDeleteDownloadClient: (_id: string) => void | Promise<void>;
		onDeleteIndexer: (_id: string) => void | Promise<void>;
		onDeleteLibraryFolder: (_id: string) => void | Promise<void>;
		onDeletePathMapping: (_id: string) => void | Promise<void>;
		onDeleteMediaProfile: (_id: string) => void | Promise<void>;
		onDeleteTag: (_id: string) => void | Promise<void>;
		onDeleteUser: (_id: string) => void | Promise<void>;
		onTestIndexer: (_id: string) => void | Promise<void>;
		onTestMetadataProvider: (_id: string) => void | Promise<void>;
		onSearchLibraryMatch: (_kind: LibraryMediaKind, _query: string) => Promise<MediaSearchResult[]>;
		onMatchLibraryScanItem: (
			_item: LibraryScanItem,
			_request: LibraryScanItemMatchRequest
		) => Promise<void>;
	}

	let {
		activeSection,
		downloadClients,
		indexers,
		metadataProviders,
		metadataCache,
		libraryFolders,
		pathMappings,
		mediaProfiles,
		users,
		tags,
		currentUser,
		activeLibraryScan,
		downloadForm = $bindable(),
		indexerForm = $bindable(),
		libraryFolderForm = $bindable(),
		pathMappingForm = $bindable(),
		mediaProfileForm = $bindable(),
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
		savingTag,
		deletingTagId,
		savingUser,
		loadingLibraryScan,
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
		onSavePathMapping,
		onSaveMediaProfile,
		onSaveTag,
		onSaveUser,
		onCancelDownloadClient,
		onCancelIndexer,
		onCancelMediaProfile,
		onCancelTag,
		onCancelUser,
		onEditDownloadClient,
		onEditIndexer,
		onEditMediaProfile,
		onEditTag,
		onEditUser,
		onDeleteDownloadClient,
		onDeleteIndexer,
		onDeleteLibraryFolder,
		onDeletePathMapping,
		onDeleteMediaProfile,
		onDeleteTag,
		onDeleteUser,
		onTestIndexer,
		onTestMetadataProvider,
		onSearchLibraryMatch,
		onMatchLibraryScanItem
	}: Props = $props();
</script>

<section class="workspace-main" aria-labelledby="settings-title">
	{#if activeSection === 'download-clients'}
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
	{:else if activeSection === 'metadata'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Metadata</h1>
		</div>
		<div class="settings-stack">
			<MetadataProviderSettings
				{metadataProviders}
				onSave={onSaveMetadataProvider}
				onTest={onTestMetadataProvider}
				testingId={testingMetadataProviderId}
				savingId={savingMetadataProviderId}
				testResults={metadataProviderTests}
			/>
			<MetadataCacheSettings
				cache={metadataCache}
				bind:pattern={metadataCachePattern}
				loading={loadingMetadataCache}
				clearing={clearingMetadataCache}
				onRefresh={onRefreshMetadataCache}
				onClearAll={onClearMetadataCache}
				onClearPattern={onClearMetadataCachePattern}
			/>
		</div>
	{:else if activeSection === 'quality'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Quality</h1>
		</div>
		<div class="settings-stack">
			<QualitySizeSettings />
		</div>
	{:else if activeSection === 'profiles'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Profiles</h1>
		</div>
		<div class="settings-stack">
			<MediaProfilesSettings
				profiles={mediaProfiles}
				bind:form={mediaProfileForm}
				saving={savingMediaProfile}
				deletingId={deletingMediaProfileId}
				onSave={onSaveMediaProfile}
				onCancel={onCancelMediaProfile}
				onEdit={onEditMediaProfile}
				onDelete={onDeleteMediaProfile}
			/>
		</div>
	{:else if activeSection === 'file-naming'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">File naming</h1>
		</div>
		<div class="settings-stack">
			<FileNamingSettings />
		</div>
	{:else if activeSection === 'tags'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Tags</h1>
		</div>
		<div class="settings-stack">
			<TagSettings
				{tags}
				bind:form={tagForm}
				saving={savingTag}
				deletingId={deletingTagId}
				onSave={onSaveTag}
				onCancel={onCancelTag}
				onEdit={onEditTag}
				onDelete={onDeleteTag}
			/>
		</div>
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
	{:else if activeSection === 'system-logs'}
		<div class="page-heading">
			<p>Settings / System</p>
			<h1 id="settings-title">Logs</h1>
		</div>
		<div class="settings-stack">
			<SystemLogsSettings />
		</div>
	{:else}
		<LibrarySettingsSection
			folders={libraryFolders}
			bind:form={libraryFolderForm}
			{pathMappings}
			bind:pathMappingForm
			scan={activeLibraryScan}
			saving={savingLibraryFolder}
			{savingPathMapping}
			{deletingPathMappingId}
			loadingScan={loadingLibraryScan}
			onSave={onSaveLibraryFolder}
			onDelete={onDeleteLibraryFolder}
			{onSavePathMapping}
			{onDeletePathMapping}
			onSearchMatch={onSearchLibraryMatch}
			onMatch={onMatchLibraryScanItem}
		/>
	{/if}
</section>
