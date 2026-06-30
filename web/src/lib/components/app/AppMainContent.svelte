<script lang="ts">
	import AdvancedSearchArea from '$lib/components/app/AdvancedSearchArea.svelte';
	import DiscoverSectionArea from '$lib/components/app/DiscoverSectionArea.svelte';
	import HomeAreaRoute from '$lib/components/app/HomeAreaRoute.svelte';
	import MediaCollectionArea from '$lib/components/app/MediaCollectionArea.svelte';
	import MediaPeopleArea from '$lib/components/app/MediaPeopleArea.svelte';
	import MetadataDetailArea from '$lib/components/app/MetadataDetailArea.svelte';
	import SettingsArea from '$lib/components/app/SettingsArea.svelte';
	import SystemArea from '$lib/components/app/SystemArea.svelte';
	import type {
		AppShellOptions,
		createAppShellController
	} from '$lib/components/app/appShellController.svelte';

	type AppShellController = ReturnType<typeof createAppShellController>;

	interface Props {
		app: AppShellController;
		options: AppShellOptions;
	}

	let { app = $bindable(), options }: Props = $props();
</script>

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
		initialQuery={options.initialAdvancedQuery ?? ''}
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
{:else if app.activeView === 'media-people'}
	<MediaPeopleArea
		detail={app.mediaPeopleDetail}
		loading={app.loadingMetadataDetail && !app.mediaPeopleDetail}
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
{:else if app.activeView === 'discover-section'}
	<DiscoverSectionArea
		section={app.discoverSection}
		mediaItems={app.mediaItems}
		loading={app.loadingDiscoverSection}
		loadingMore={app.loadingMoreDiscoverSection}
		hasMore={app.discoverSectionHasMore}
		addingKey={app.addingKey}
		blacklistingKey={app.blacklistingKey}
		actionLabel={app.isAdmin ? 'Add' : 'Request'}
		canManage={app.isAdmin}
		blacklist={app.discoverBlacklist}
		onAdd={app.addMedia}
		onBlacklist={app.blacklistDiscoverMedia}
		onLoadMore={app.loadMoreDiscoverSection}
	/>
{:else}
	<HomeAreaRoute {app} />
{/if}
