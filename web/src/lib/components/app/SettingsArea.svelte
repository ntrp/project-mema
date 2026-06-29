<script lang="ts">
	import { resolve } from '$app/paths';

	import DownloadClientForm from '$lib/components/settings/DownloadClientForm.svelte';
	import DownloadClientTable from '$lib/components/settings/DownloadClientTable.svelte';
	import IndexerForm from '$lib/components/settings/IndexerForm.svelte';
	import IndexerTable from '$lib/components/settings/IndexerTable.svelte';
	import LibraryFolderForm from '$lib/components/settings/LibraryFolderForm.svelte';
	import LibraryFolderTable from '$lib/components/settings/LibraryFolderTable.svelte';
	import LibraryScanReview from '$lib/components/settings/LibraryScanReview.svelte';
	import MetadataProviderSettings from '$lib/components/settings/MetadataProviderSettings.svelte';
	import UserForm from '$lib/components/settings/UserForm.svelte';
	import UserTable from '$lib/components/settings/UserTable.svelte';
	import type {
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		Indexer,
		IndexerForm as IndexerFormValue,
		IntegrationTestResults,
		LibraryFolder,
		LibraryFolderForm as LibraryFolderFormValue,
		LibraryMediaKind,
		LibraryScan,
		LibraryScanItem,
		LibraryScanItemMatchRequest,
		ManagedUser,
		MediaSearchResult,
		MetadataProvider,
		MetadataProviderForm as MetadataProviderFormValue,
		SettingsSection,
		UserForm as UserFormValue,
		UserSummary
	} from '$lib/settings/types';

	type SettingsHref =
		| '/settings/library'
		| '/settings/download-clients'
		| '/settings/indexers'
		| '/settings/metadata'
		| '/settings/users';

	interface Props {
		activeSection: SettingsSection;
		downloadClients: DownloadClient[];
		indexers: Indexer[];
		metadataProviders: MetadataProvider[];
		libraryFolders: LibraryFolder[];
		users: ManagedUser[];
		currentUser?: UserSummary;
		activeLibraryScan?: LibraryScan;
		downloadForm: DownloadClientFormValue;
		indexerForm: IndexerFormValue;
		libraryFolderForm: LibraryFolderFormValue;
		userForm: UserFormValue;
		savingDownloadClient: boolean;
		savingIndexer: boolean;
		savingMetadataProviderId?: string;
		savingLibraryFolder: boolean;
		savingUser: boolean;
		loadingLibraryScan: boolean;
		testingDownloadClientId?: string;
		testingIndexerId?: string;
		testingMetadataProviderId?: string;
		downloadClientTests: IntegrationTestResults;
		indexerTests: IntegrationTestResults;
		metadataProviderTests: IntegrationTestResults;
		onSectionSelect: (_section: SettingsSection) => void;
		onSaveDownloadClient: (_event: SubmitEvent) => void | Promise<void>;
		onSaveIndexer: (_event: SubmitEvent) => void | Promise<void>;
		onSaveMetadataProvider: (_form: MetadataProviderFormValue) => void | Promise<void>;
		onSaveLibraryFolder: (_event: SubmitEvent) => void | Promise<void>;
		onSaveUser: (_event: SubmitEvent) => void | Promise<void>;
		onCancelDownloadClient: () => void;
		onCancelIndexer: () => void;
		onCancelUser: () => void;
		onEditDownloadClient: (_client: DownloadClient) => void;
		onEditIndexer: (_indexer: Indexer) => void;
		onEditUser: (_user: ManagedUser) => void;
		onDeleteDownloadClient: (_id: string) => void | Promise<void>;
		onDeleteIndexer: (_id: string) => void | Promise<void>;
		onDeleteLibraryFolder: (_id: string) => void | Promise<void>;
		onDeleteUser: (_id: string) => void | Promise<void>;
		onTestDownloadClient: (_id: string) => void | Promise<void>;
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
		libraryFolders,
		users,
		currentUser,
		activeLibraryScan,
		downloadForm = $bindable(),
		indexerForm = $bindable(),
		libraryFolderForm = $bindable(),
		userForm = $bindable(),
		savingDownloadClient,
		savingIndexer,
		savingMetadataProviderId,
		savingLibraryFolder,
		savingUser,
		loadingLibraryScan,
		testingDownloadClientId,
		testingIndexerId,
		testingMetadataProviderId,
		downloadClientTests,
		indexerTests,
		metadataProviderTests,
		onSectionSelect,
		onSaveDownloadClient,
		onSaveIndexer,
		onSaveMetadataProvider,
		onSaveLibraryFolder,
		onSaveUser,
		onCancelDownloadClient,
		onCancelIndexer,
		onCancelUser,
		onEditDownloadClient,
		onEditIndexer,
		onEditUser,
		onDeleteDownloadClient,
		onDeleteIndexer,
		onDeleteLibraryFolder,
		onDeleteUser,
		onTestDownloadClient,
		onTestIndexer,
		onTestMetadataProvider,
		onSearchLibraryMatch,
		onMatchLibraryScanItem
	}: Props = $props();

	const settingsItems = [
		{
			value: 'library',
			label: 'Library',
			href: '/settings/library'
		},
		{
			value: 'download-clients',
			label: 'Download clients',
			href: '/settings/download-clients'
		},
		{
			value: 'indexers',
			label: 'Indexers',
			href: '/settings/indexers'
		},
		{
			value: 'metadata',
			label: 'Metadata',
			href: '/settings/metadata'
		},
		{
			value: 'users',
			label: 'Users',
			href: '/settings/users'
		}
	] satisfies {
		value: SettingsSection;
		label: string;
		href: SettingsHref;
	}[];
</script>

<section class="workspace-main" aria-labelledby="settings-title">
	<nav class="settings-tabs" aria-label="Settings sections">
		{#each settingsItems as item (item.value)}
			<a
				href={resolve(item.href)}
				class:active-tab={activeSection === item.value}
				aria-current={activeSection === item.value ? 'page' : undefined}
				onclick={() => onSectionSelect(item.value)}
			>
				{item.label}
			</a>
		{/each}
	</nav>

	{#if activeSection === 'download-clients'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Download clients</h1>
		</div>
		<div class="settings-stack">
			<DownloadClientForm
				bind:form={downloadForm}
				saving={savingDownloadClient}
				onSave={onSaveDownloadClient}
				onCancel={onCancelDownloadClient}
			/>
			<DownloadClientTable
				clients={downloadClients}
				onEdit={onEditDownloadClient}
				onDelete={onDeleteDownloadClient}
				onTest={onTestDownloadClient}
				testingId={testingDownloadClientId}
				testResults={downloadClientTests}
			/>
		</div>
	{:else if activeSection === 'indexers'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Indexers</h1>
		</div>
		<div class="settings-stack">
			<IndexerForm
				bind:form={indexerForm}
				saving={savingIndexer}
				onSave={onSaveIndexer}
				onCancel={onCancelIndexer}
			/>
			<IndexerTable
				{indexers}
				onEdit={onEditIndexer}
				onDelete={onDeleteIndexer}
				onTest={onTestIndexer}
				testingId={testingIndexerId}
				testResults={indexerTests}
			/>
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
		</div>
	{:else if activeSection === 'users'}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Users</h1>
		</div>
		<div class="settings-stack">
			<UserForm
				bind:form={userForm}
				saving={savingUser}
				onSave={onSaveUser}
				onCancel={onCancelUser}
			/>
			<UserTable
				{users}
				currentUserId={currentUser?.id}
				onEdit={onEditUser}
				onDelete={onDeleteUser}
			/>
		</div>
	{:else}
		<div class="page-heading">
			<p>Settings</p>
			<h1 id="settings-title">Library</h1>
		</div>
		<div class="settings-stack">
			<LibraryFolderForm
				bind:form={libraryFolderForm}
				saving={savingLibraryFolder}
				onSave={onSaveLibraryFolder}
			/>
			<LibraryFolderTable folders={libraryFolders} onDelete={onDeleteLibraryFolder} />
			<LibraryScanReview
				scan={activeLibraryScan}
				loading={loadingLibraryScan}
				onSearchMatch={onSearchLibraryMatch}
				onMatch={onMatchLibraryScanItem}
			/>
		</div>
	{/if}
</section>
