<script lang="ts">
	import DownloadClientForm from '$lib/components/settings/DownloadClientForm.svelte';
	import DownloadClientTable from '$lib/components/settings/DownloadClientTable.svelte';
	import IndexerForm from '$lib/components/settings/IndexerForm.svelte';
	import IndexerTable from '$lib/components/settings/IndexerTable.svelte';
	import LibraryFolderForm from '$lib/components/settings/LibraryFolderForm.svelte';
	import LibraryFolderTable from '$lib/components/settings/LibraryFolderTable.svelte';
	import LibraryScanReview from '$lib/components/settings/LibraryScanReview.svelte';
	import MetadataCacheSettings from '$lib/components/settings/MetadataCacheSettings.svelte';
	import MetadataProviderSettings from '$lib/components/settings/MetadataProviderSettings.svelte';
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import SystemLogsSettings from '$lib/components/settings/SystemLogsSettings.svelte';
	import TagSettings from '$lib/components/settings/TagSettings.svelte';
	import UserForm from '$lib/components/settings/UserForm.svelte';
	import UserTable from '$lib/components/settings/UserTable.svelte';
	import {
		emptyDownloadClientForm,
		emptyIndexerForm,
		emptyLibraryFolderForm,
		emptyUserForm
	} from '$lib/settings/forms';
	import type {
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		DownloadClientType,
		Indexer,
		IndexerForm as IndexerFormValue,
		IntegrationTestResults,
		IntegrationTestResponse,
		LibraryFolder,
		LibraryFolderForm as LibraryFolderFormValue,
		LibraryMediaKind,
		LibraryScan,
		LibraryScanItem,
		LibraryScanItemMatchRequest,
		ManagedUser,
		MediaSearchResult,
		MetadataCacheResponse,
		MetadataProvider,
		MetadataProviderForm as MetadataProviderFormValue,
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
		users: ManagedUser[];
		tags: Tag[];
		currentUser?: UserSummary;
		activeLibraryScan?: LibraryScan;
		downloadForm: DownloadClientFormValue;
		indexerForm: IndexerFormValue;
		libraryFolderForm: LibraryFolderFormValue;
		tagForm: TagForm;
		userForm: UserFormValue;
		savingDownloadClient: boolean;
		savingIndexer: boolean;
		savingMetadataProviderId?: string;
		loadingMetadataCache: boolean;
		clearingMetadataCache: boolean;
		metadataCachePattern: string;
		savingLibraryFolder: boolean;
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
		onSaveTag: (_event: SubmitEvent) => void | Promise<void>;
		onSaveUser: (_event: SubmitEvent) => void | Promise<void>;
		onCancelDownloadClient: () => void;
		onCancelIndexer: () => void;
		onCancelTag: () => void;
		onCancelUser: () => void;
		onEditDownloadClient: (_client: DownloadClient) => void;
		onEditIndexer: (_indexer: Indexer) => void;
		onEditTag: (_tag: Tag) => void;
		onEditUser: (_user: ManagedUser) => void;
		onDeleteDownloadClient: (_id: string) => void | Promise<void>;
		onDeleteIndexer: (_id: string) => void | Promise<void>;
		onDeleteLibraryFolder: (_id: string) => void | Promise<void>;
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
		users,
		tags,
		currentUser,
		activeLibraryScan,
		downloadForm = $bindable(),
		indexerForm = $bindable(),
		libraryFolderForm = $bindable(),
		tagForm = $bindable(),
		userForm = $bindable(),
		savingDownloadClient,
		savingIndexer,
		savingMetadataProviderId,
		loadingMetadataCache,
		clearingMetadataCache,
		metadataCachePattern = $bindable(),
		savingLibraryFolder,
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
		onSaveTag,
		onSaveUser,
		onCancelDownloadClient,
		onCancelIndexer,
		onCancelTag,
		onCancelUser,
		onEditDownloadClient,
		onEditIndexer,
		onEditTag,
		onEditUser,
		onDeleteDownloadClient,
		onDeleteIndexer,
		onDeleteLibraryFolder,
		onDeleteTag,
		onDeleteUser,
		onTestIndexer,
		onTestMetadataProvider,
		onSearchLibraryMatch,
		onMatchLibraryScanItem
	}: Props = $props();

	let downloadClientModalOpen = $state(false);
	let downloadClientTypeSelected = $state(false);
	let testingDownloadClientConfig = $state(false);
	let downloadClientModalTestResult = $state<IntegrationTestResponse | undefined>();
	let indexerModalOpen = $state(false);
	let libraryFolderModalOpen = $state(false);
	let userModalOpen = $state(false);

	function openDownloadClientModal() {
		downloadForm = emptyDownloadClientForm();
		downloadClientTypeSelected = false;
		downloadClientModalTestResult = undefined;
		downloadClientModalOpen = true;
	}

	function editDownloadClient(client: DownloadClient) {
		onEditDownloadClient(client);
		downloadClientTypeSelected = true;
		downloadClientModalTestResult = undefined;
		downloadClientModalOpen = true;
	}

	function closeDownloadClientModal() {
		onCancelDownloadClient();
		downloadClientModalOpen = false;
		downloadClientTypeSelected = false;
		downloadClientModalTestResult = undefined;
	}

	function selectDownloadClientType(type: DownloadClientType) {
		downloadForm = {
			...emptyDownloadClientForm(),
			type
		};
		downloadClientTypeSelected = true;
		downloadClientModalTestResult = undefined;
	}

	function openIndexerModal() {
		indexerForm = emptyIndexerForm();
		indexerModalOpen = true;
	}

	function editIndexer(indexer: Indexer) {
		onEditIndexer(indexer);
		indexerModalOpen = true;
	}

	function closeIndexerModal() {
		onCancelIndexer();
		indexerModalOpen = false;
	}

	function openLibraryFolderModal() {
		libraryFolderForm = emptyLibraryFolderForm();
		libraryFolderModalOpen = true;
	}

	function closeLibraryFolderModal() {
		libraryFolderForm = emptyLibraryFolderForm();
		libraryFolderModalOpen = false;
	}

	function openUserModal() {
		userForm = emptyUserForm();
		userModalOpen = true;
	}

	function editUser(user: ManagedUser) {
		onEditUser(user);
		userModalOpen = true;
	}

	function closeUserModal() {
		onCancelUser();
		userModalOpen = false;
	}

	async function saveDownloadClient(event: SubmitEvent) {
		event.preventDefault();
		const passed = await testDownloadClientConfig();
		if (!passed) {
			return;
		}
		await onSaveDownloadClient(event);
		if (isEmptyDownloadClientForm(downloadForm)) {
			downloadClientModalOpen = false;
		}
	}

	async function testDownloadClientConfig() {
		testingDownloadClientConfig = true;
		downloadClientModalTestResult = undefined;
		try {
			const result = await onTestDownloadClientConfig(downloadForm);
			downloadClientModalTestResult = result;
			return result.success;
		} finally {
			testingDownloadClientConfig = false;
		}
	}

	async function saveIndexer(event: SubmitEvent) {
		await onSaveIndexer(event);
		if (isEmptyIndexerForm(indexerForm)) {
			indexerModalOpen = false;
		}
	}

	async function saveLibraryFolder(event: SubmitEvent) {
		await onSaveLibraryFolder(event);
		if (libraryFolderForm.path.trim() === '') {
			libraryFolderModalOpen = false;
		}
	}

	async function saveUser(event: SubmitEvent) {
		await onSaveUser(event);
		if (isEmptyUserForm(userForm)) {
			userModalOpen = false;
		}
	}

	function isEmptyDownloadClientForm(value: DownloadClientFormValue) {
		return (
			!value.id &&
			value.name === '' &&
			value.baseUrl === '' &&
			value.username === '' &&
			value.password === '' &&
			value.apiKey === '' &&
			value.category === ''
		);
	}

	function isEmptyIndexerForm(value: IndexerFormValue) {
		return !value.id && value.name === '' && value.baseUrl === '' && value.apiKey === '';
	}

	function isEmptyUserForm(value: UserFormValue) {
		return !value.id && value.username === '' && value.password === '';
	}
</script>

<section class="workspace-main" aria-labelledby="settings-title">
	{#if activeSection === 'download-clients'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Download clients</h1>
		</div>
		<div class="settings-stack">
			<div class="settings-toolbar">
				<button type="button" onclick={openDownloadClientModal}>Add download client</button>
			</div>
			<DownloadClientTable
				clients={downloadClients}
				onEdit={editDownloadClient}
				onDelete={onDeleteDownloadClient}
			/>
			{#if downloadClientModalOpen}
				<SettingsFormModal
					title={downloadForm.id ? 'Edit download client' : 'Add download client'}
					onClose={closeDownloadClientModal}
				>
					{#if downloadClientTypeSelected}
						<DownloadClientForm
							bind:form={downloadForm}
							saving={savingDownloadClient}
							onSave={saveDownloadClient}
							onCancel={closeDownloadClientModal}
							onTest={testDownloadClientConfig}
							showTypeSelect={Boolean(downloadForm.id)}
							testing={testingDownloadClientConfig}
							testResult={downloadClientModalTestResult}
						/>
					{:else}
						<div class="download-client-picker" aria-label="Download client type">
							<button type="button" onclick={() => selectDownloadClientType('transmission')}>
								<span class="app-icon" aria-hidden="true">sync_alt</span>
								<strong>Transmission</strong>
								<small>Torrent download client</small>
							</button>
							<button type="button" onclick={() => selectDownloadClientType('sabnzbd')}>
								<span class="app-icon" aria-hidden="true">cloud_download</span>
								<strong>SABnzbd</strong>
								<small>Usenet download client</small>
							</button>
						</div>
					{/if}
				</SettingsFormModal>
			{/if}
		</div>
	{:else if activeSection === 'indexers'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Indexers</h1>
		</div>
		<div class="settings-stack">
			<div class="settings-toolbar">
				<button type="button" onclick={openIndexerModal}>Add indexer</button>
			</div>
			<IndexerTable
				{indexers}
				onEdit={editIndexer}
				onDelete={onDeleteIndexer}
				onTest={onTestIndexer}
				testingId={testingIndexerId}
				testResults={indexerTests}
			/>
			{#if indexerModalOpen}
				<SettingsFormModal
					title={indexerForm.id ? 'Edit indexer' : 'Add indexer'}
					onClose={closeIndexerModal}
				>
					<IndexerForm
						bind:form={indexerForm}
						saving={savingIndexer}
						onSave={saveIndexer}
						onCancel={closeIndexerModal}
					/>
				</SettingsFormModal>
			{/if}
		</div>
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
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Users</h1>
		</div>
		<div class="settings-stack">
			<div class="settings-toolbar">
				<button type="button" onclick={openUserModal}>Add user</button>
			</div>
			<UserTable
				{users}
				currentUserId={currentUser?.id}
				onEdit={editUser}
				onDelete={onDeleteUser}
			/>
			{#if userModalOpen}
				<SettingsFormModal title={userForm.id ? 'Edit user' : 'Add user'} onClose={closeUserModal}>
					<UserForm
						bind:form={userForm}
						saving={savingUser}
						onSave={saveUser}
						onCancel={closeUserModal}
					/>
				</SettingsFormModal>
			{/if}
		</div>
	{:else if activeSection === 'system-logs'}
		<div class="page-heading">
			<p>Settings / System</p>
			<h1 id="settings-title">Logs</h1>
		</div>
		<div class="settings-stack">
			<SystemLogsSettings />
		</div>
	{:else}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Library</h1>
		</div>
		<div class="settings-stack">
			<div class="settings-toolbar">
				<button type="button" onclick={openLibraryFolderModal}>Add library folder</button>
			</div>
			<LibraryFolderTable folders={libraryFolders} onDelete={onDeleteLibraryFolder} />
			<LibraryScanReview
				scan={activeLibraryScan}
				loading={loadingLibraryScan}
				onSearchMatch={onSearchLibraryMatch}
				onMatch={onMatchLibraryScanItem}
			/>
			{#if libraryFolderModalOpen}
				<SettingsFormModal title="Add library folder" onClose={closeLibraryFolderModal}>
					<LibraryFolderForm
						bind:form={libraryFolderForm}
						saving={savingLibraryFolder}
						onSave={saveLibraryFolder}
					/>
				</SettingsFormModal>
			{/if}
		</div>
	{/if}
</section>
