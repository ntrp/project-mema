<script lang="ts">
	import { onMount } from 'svelte';
	import AdvancedSearchArea from '$lib/components/app/AdvancedSearchArea.svelte';
	import AppDocumentHead from '$lib/components/app/AppDocumentHead.svelte';
	import AppNav from '$lib/components/app/AppNav.svelte';
	import {
		createAppShellController,
		type AppShellOptions
	} from '$lib/components/app/appShellController.svelte';
	import HomeAreaRoute from '$lib/components/app/HomeAreaRoute.svelte';
	import MediaCollectionArea from '$lib/components/app/MediaCollectionArea.svelte';
	import MediaDeleteModal from '$lib/components/app/MediaDeleteModal.svelte';
	import MediaActionModal from '$lib/components/app/MediaActionModal.svelte';
	import MetadataDetailArea from '$lib/components/app/MetadataDetailArea.svelte';
	import SettingsArea from '$lib/components/app/SettingsArea.svelte';
	import SidebarMenu from '$lib/components/app/SidebarMenu.svelte';
	import SystemArea from '$lib/components/app/SystemArea.svelte';
	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import '$lib/settings/styles.css';

	let props: AppShellOptions = $props();
	// svelte-ignore state_referenced_locally
	const app = createAppShellController(props);

	onMount(() => {
		void app.initialise();
		return app.disconnectEvents;
	});
</script>

<AppDocumentHead />

{#if app.loading}
	<main class="shell">
		<section class="panel">
			<p class="muted">Loading app</p>
		</section>
	</main>
{:else if !app.authenticated}
	<main class="shell login-shell">
		<div class="login-brand">
			<span class="brand-mark large" aria-hidden="true">M</span>
			<h1>mema</h1>
		</div>
		<AuthPanel bind:username={app.username} bind:password={app.password} onLogin={app.login} />
	</main>
{:else}
	<div class="app-frame">
		<SidebarMenu
			title="mema"
			items={app.primaryItems}
			active={app.activePrimarySection}
			activeSubmenu={app.activeSubmenuSection}
			onSelect={app.selectPrimarySection}
			onSubmenuSelect={app.selectSubmenuSection}
		/>
		<div class="app-main">
			<AppNav
				bind:searchQuery={app.searchQuery}
				groups={app.autocompleteGroups}
				loading={app.loadingAutocomplete}
				onSearch={app.autocompleteMedia}
				onSelect={app.selectAutocompleteResult}
				onAdvancedSearch={app.openAdvancedSearch}
				onProfile={app.showProfile}
				onLogout={app.logout}
				showNotifications={app.isAdmin}
			/>
			<main class="app-content">
				{#if app.activeView === 'settings' && app.isAdmin}
					<SettingsArea
						bind:downloadForm={app.downloadForm}
						bind:indexerForm={app.indexerForm}
						bind:libraryFolderForm={app.libraryFolderForm}
						bind:pathMappingForm={app.pathMappingForm}
						bind:mediaProfileForm={app.mediaProfileForm}
						bind:customFormatForm={app.customFormatForm}
						bind:tagForm={app.tagForm}
						bind:userForm={app.userForm}
						activeSection={app.activeSettingsSection}
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
						savingMediaProfile={app.savingMediaProfile}
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
						onSaveMediaProfile={app.saveMediaProfile}
						onSaveCustomFormat={app.saveCustomFormat}
						onImportCustomFormat={app.importCustomFormat}
						onSaveTag={app.saveTag}
						onSaveUser={app.saveUser}
						onCancelDownloadClient={app.cancelDownloadClient}
						onCancelIndexer={app.cancelIndexer}
						onCancelMediaProfile={app.cancelMediaProfile}
						onCancelCustomFormat={app.cancelCustomFormat}
						onCancelTag={app.cancelTag}
						onCancelUser={app.cancelUser}
						onEditDownloadClient={app.editDownloadClient}
						onEditIndexer={app.editIndexer}
						onEditMediaProfile={app.editMediaProfile}
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
				{:else if app.activeView === 'system' && app.isAdmin}
					<SystemArea activeSection={app.activeSystemSection} />
				{:else if app.activeView === 'advanced-search'}
					<AdvancedSearchArea
						initialQuery={props.initialAdvancedQuery ?? ''}
						metadataProviders={app.metadataProviders}
						groups={app.advancedSearchGroups}
						searching={app.searchingAdvanced}
						addingKey={app.addingKey}
						actionLabel={app.isAdmin ? 'Add' : 'Request'}
						onSearch={app.advancedSearch}
						onAdd={app.addMedia}
					/>
				{:else if app.activeView === 'metadata-detail'}
					<MetadataDetailArea
						detail={app.metadataDetail}
						loading={app.loadingMetadataDetail}
						addingKey={app.addingKey}
						actionLabel={app.isAdmin ? 'Add' : 'Request'}
						onAdd={app.addMedia}
					/>
				{:else if app.activeView === 'media-collection'}
					<MediaCollectionArea
						collection={app.mediaCollection}
						mediaItems={app.mediaItems}
						loading={app.loadingMediaCollection}
						addingKey={app.addingKey}
						actionLabel={app.isAdmin ? 'Add' : 'Request'}
						onAdd={app.addMedia}
					/>
				{:else}
					<HomeAreaRoute {app} />
				{/if}
			</main>
		</div>
	</div>
	{#if app.activeMediaCandidate}
		<MediaActionModal
			candidate={app.activeMediaCandidate}
			isAdmin={app.isAdmin}
			libraryFolders={app.libraryFolders}
			qualityProfiles={app.mediaProfiles}
			tags={app.tags}
			saving={app.savingMediaAction}
			onClose={app.closeMediaAction}
			onConfirm={app.confirmMediaAction}
		/>
	{/if}
	{#if app.mediaDeleteCandidate}
		<MediaDeleteModal
			item={app.mediaDeleteCandidate}
			deleting={app.deletingMediaItemId === app.mediaDeleteCandidate.id}
			onClose={app.closeMediaDelete}
			onDelete={app.confirmMediaDelete}
		/>
	{/if}
{/if}

<NoticeStack message={app.message} errorMessage={app.errorMessage} onDismiss={app.clearNotice} />
