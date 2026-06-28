<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';

	import AppNav from '$lib/components/app/AppNav.svelte';
	import HomeArea from '$lib/components/app/HomeArea.svelte';
	import SettingsArea from '$lib/components/app/SettingsArea.svelte';
	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import {
		createMediaItem as createMediaItemRequest,
		currentSessionAuthenticated,
		deleteDownloadClient as deleteDownloadClientRequest,
		deleteIndexer as deleteIndexerRequest,
		deleteMediaItem as deleteMediaItemRequest,
		enqueueMediaReleaseSearch as enqueueMediaReleaseSearchRequest,
		grabMediaRelease as grabMediaReleaseRequest,
		listDownloadActivity as listDownloadActivityRequest,
		listMediaItems as listMediaItemsRequest,
		loadSettings as loadSettingsRequest,
		login as loginRequest,
		logout as logoutRequest,
		saveDownloadClient as saveDownloadClientRequest,
		saveIndexer as saveIndexerRequest,
		searchMedia as searchMediaRequest,
		searchMediaReleases as searchMediaReleasesRequest,
		testDownloadClient as testDownloadClientRequest,
		testIndexer as testIndexerRequest
	} from '$lib/settings/api';
	import {
		downloadClientFormFromClient,
		emptyDownloadClientForm,
		emptyIndexerForm,
		indexerFormFromIndexer
	} from '$lib/settings/forms';
	import '$lib/settings/styles.css';
	import type {
		AppView,
		DownloadActivity,
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		HomeSection,
		Indexer,
		IndexerForm as IndexerFormValue,
		IntegrationTestResults,
		MediaItem,
		MediaSearchRequest,
		MediaSearchResult,
		ReleaseCandidate,
		ReleaseSearchResults,
		SettingsSection
	} from '$lib/settings/types';

	interface Props {
		initialView?: AppView;
		initialHomeSection?: HomeSection;
		initialSettingsSection?: SettingsSection;
		initialSelectedMediaItemId?: string;
	}

	let {
		initialView = 'home',
		initialHomeSection = 'explore',
		initialSettingsSection = 'download-clients',
		initialSelectedMediaItemId
	}: Props = $props();
	const routeDefaults = (() => ({
		view: initialView,
		homeSection: initialHomeSection,
		settingsSection: initialSettingsSection,
		selectedMediaItemId: initialSelectedMediaItemId
	}))();

	let authenticated = $state(false);
	let loading = $state(true);
	let savingDownloadClient = $state(false);
	let savingIndexer = $state(false);
	let message = $state('');
	let errorMessage = $state('');
	let username = $state('admin');
	let password = $state('admin');
	let downloadClients = $state<DownloadClient[]>([]);
	let indexers = $state<Indexer[]>([]);
	let mediaItems = $state<MediaItem[]>([]);
	let mediaSearchResults = $state<MediaSearchResult[]>([]);
	let releaseResults = $state<ReleaseSearchResults>({});
	let activities = $state<DownloadActivity[]>([]);
	let downloadForm = $state<DownloadClientFormValue>(emptyDownloadClientForm());
	let indexerForm = $state<IndexerFormValue>(emptyIndexerForm());
	let testingDownloadClientId = $state<string | undefined>();
	let testingIndexerId = $state<string | undefined>();
	let searchingMedia = $state(false);
	let addingKey = $state<string | undefined>();
	let searchingItemId = $state<string | undefined>();
	let grabbingKey = $state<string | undefined>();
	let deletingMediaItemId = $state<string | undefined>();
	let loadingActivity = $state(false);
	let downloadClientTests = $state<IntegrationTestResults>({});
	let indexerTests = $state<IntegrationTestResults>({});
	let activeView = $state<AppView>(routeDefaults.view);
	let activeHomeSection = $state<HomeSection>(routeDefaults.homeSection);
	let activeSettingsSection = $state<SettingsSection>(routeDefaults.settingsSection);
	let selectedMediaItemId = $state<string | undefined>(routeDefaults.selectedMediaItemId);
	let searchQuery = $state('');

	onMount(() => {
		void initialise();
	});

	async function initialise() {
		loading = true;
		errorMessage = '';

		authenticated = await currentSessionAuthenticated();
		if (authenticated) {
			await loadSettings();
			await loadLibrary();
		}

		loading = false;
	}

	async function login(event: SubmitEvent) {
		event.preventDefault();
		clearNotice();

		try {
			await loginRequest(username, password);
			authenticated = true;
			await loadSettings();
			await loadLibrary();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Login failed');
		}
	}

	async function logout() {
		clearNotice();

		try {
			await logoutRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not log out');
		} finally {
			authenticated = false;
			activeView = 'home';
			downloadClients = [];
			indexers = [];
			mediaItems = [];
			mediaSearchResults = [];
			releaseResults = {};
			activities = [];
			downloadForm = emptyDownloadClientForm();
			indexerForm = emptyIndexerForm();
		}
	}

	function showProfile() {
		clearNotice();
		message = 'Profile settings are not implemented yet';
	}

	function openSettings() {
		activeView = 'settings';
		activeSettingsSection = 'download-clients';
		void goto(resolve('/settings/download-clients'));
	}

	function openHome() {
		activeView = 'home';
		activeHomeSection = 'explore';
		void goto(resolve('/explore'));
	}

	async function loadSettings() {
		try {
			const settings = await loadSettingsRequest();
			downloadClients = settings.downloadClients;
			indexers = settings.indexers;
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load settings');
		}
	}

	async function loadLibrary() {
		await Promise.all([loadMediaItems(), loadDownloadActivity()]);
	}

	async function loadMediaItems() {
		try {
			mediaItems = await listMediaItemsRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load media items');
		}
	}

	async function loadDownloadActivity() {
		loadingActivity = true;
		try {
			activities = await listDownloadActivityRequest();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load download activity');
		} finally {
			loadingActivity = false;
		}
	}

	async function searchMedia(request: MediaSearchRequest) {
		searchingMedia = true;
		clearNotice();

		try {
			mediaSearchResults = await searchMediaRequest(request);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not search media');
		} finally {
			searchingMedia = false;
		}
	}

	async function addMedia(candidate: MediaSearchResult) {
		addingKey = candidateKey(candidate);
		clearNotice();

		try {
			const item = await createMediaItemRequest({
				title: candidate.title,
				type: candidate.type,
				year: candidate.year,
				monitored: true
			});
			mediaItems = [item, ...mediaItems];
			message = 'Media item added to monitored';
			activeHomeSection = candidate.type === 'movie' ? 'movies' : 'series';
			void goto(resolve(candidate.type === 'movie' ? '/movies' : '/series'));
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not add media item');
		} finally {
			addingKey = undefined;
		}
	}

	async function findReleases(item: MediaItem) {
		searchingItemId = item.id;
		clearNotice();

		try {
			const job = await enqueueMediaReleaseSearchRequest(item.id);
			releaseResults = {
				...releaseResults,
				[item.id]: { loaded: false, releases: [], errors: [`${job.message} (#${job.jobId})`] }
			};
			message = job.message;
			window.setTimeout(() => void loadReleaseResults(item.id), 1200);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not enqueue release search');
		} finally {
			searchingItemId = undefined;
		}
	}

	async function deleteMediaItem(item: MediaItem) {
		deletingMediaItemId = item.id;
		clearNotice();

		try {
			await deleteMediaItemRequest(item.id);
			mediaItems = mediaItems.filter((mediaItem) => mediaItem.id !== item.id);
			releaseResults = omitResult(releaseResults, item.id);
			activities = activities.filter((activity) => activity.mediaItemId !== item.id);
			message = 'Media item removed';
			if (selectedMediaItemId === item.id) {
				selectedMediaItemId = undefined;
				void goto(resolve(item.type === 'movie' ? '/movies' : '/series'));
			}
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not remove media item');
		} finally {
			deletingMediaItemId = undefined;
		}
	}

	async function loadReleaseResults(id: string) {
		try {
			const results = await searchMediaReleasesRequest(id);
			releaseResults = {
				...releaseResults,
				[id]: { loaded: true, releases: results.releases, errors: results.errors }
			};
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not load release results');
		}
	}

	async function grabRelease(item: MediaItem, release: ReleaseCandidate) {
		grabbingKey = `${item.id}:${release.id}`;
		clearNotice();

		try {
			const result = await grabMediaReleaseRequest(item.id, release);
			activities = [
				result.activity,
				...activities.filter((activity) => activity.id !== result.activity.id)
			];
			message = `${result.message} (#${result.jobId})`;
			activeHomeSection = 'activity';
			void goto(resolve('/activity'));
			window.setTimeout(() => void loadDownloadActivity(), 1200);
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not enqueue download');
		} finally {
			grabbingKey = undefined;
		}
	}

	async function saveDownloadClient(event: SubmitEvent) {
		event.preventDefault();
		savingDownloadClient = true;
		clearNotice();

		try {
			await saveDownloadClientRequest(downloadForm);
			downloadForm = emptyDownloadClientForm();
			message = 'Download client saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save download client');
		} finally {
			savingDownloadClient = false;
		}
	}

	async function saveIndexer(event: SubmitEvent) {
		event.preventDefault();
		savingIndexer = true;
		clearNotice();

		try {
			await saveIndexerRequest(indexerForm);
			indexerForm = emptyIndexerForm();
			message = 'Indexer saved';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not save indexer');
		} finally {
			savingIndexer = false;
		}
	}

	async function deleteDownloadClient(id: string) {
		clearNotice();

		try {
			await deleteDownloadClientRequest(id);
			if (downloadForm.id === id) {
				downloadForm = emptyDownloadClientForm();
			}
			downloadClientTests = omitResult(downloadClientTests, id);
			message = 'Download client deleted';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete download client');
		}
	}

	async function deleteIndexer(id: string) {
		clearNotice();

		try {
			await deleteIndexerRequest(id);
			if (indexerForm.id === id) {
				indexerForm = emptyIndexerForm();
			}
			indexerTests = omitResult(indexerTests, id);
			message = 'Indexer deleted';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete indexer');
		}
	}

	async function testDownloadClient(id: string) {
		clearNotice();
		testingDownloadClientId = id;

		try {
			const result = await testDownloadClientRequest(id);
			downloadClientTests = { ...downloadClientTests, [id]: result };
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not test download client');
		} finally {
			testingDownloadClientId = undefined;
		}
	}

	async function testIndexer(id: string) {
		clearNotice();
		testingIndexerId = id;

		try {
			const result = await testIndexerRequest(id);
			indexerTests = { ...indexerTests, [id]: result };
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not test indexer');
		} finally {
			testingIndexerId = undefined;
		}
	}

	function clearNotice() {
		errorMessage = '';
		message = '';
	}

	function omitResult<TValue>(results: Record<string, TValue | undefined>, id: string) {
		const { [id]: _removed, ...remaining } = results;
		return remaining;
	}

	function candidateKey(candidate: MediaSearchResult) {
		return `${candidate.type}:${candidate.title}:${candidate.year ?? ''}`;
	}

	function selectHomeSection(section: HomeSection) {
		activeHomeSection = section;
		void goto(resolve(`/${section}`));
	}

	function selectSettingsSection(section: SettingsSection) {
		activeSettingsSection = section;
		void goto(resolve(`/settings/${section}`));
	}

	function errorMessageFrom(error: unknown, fallback: string) {
		return error instanceof Error ? error.message : fallback;
	}
</script>

<svelte:head>
	<title>mema</title>
	<meta
		name="description"
		content="Self-hosted media manager for movies, series, books, music, clients, and indexers"
	/>
</svelte:head>

{#if loading}
	<main class="shell">
		<section class="panel">
			<p class="muted">Loading app</p>
		</section>
	</main>
{:else if !authenticated}
	<main class="shell login-shell">
		<div class="login-brand">
			<span class="brand-mark large" aria-hidden="true">M</span>
			<h1>mema</h1>
		</div>
		<NoticeStack {message} {errorMessage} />
		<AuthPanel bind:username bind:password onLogin={login} />
	</main>
{:else}
	<div class="app-frame">
		<AppNav
			{activeView}
			bind:searchQuery
			onHome={openHome}
			onSettings={openSettings}
			onProfile={showProfile}
			onLogout={logout}
		/>
		<main class="app-content">
			<NoticeStack {message} {errorMessage} />
			{#if activeView === 'settings'}
				<SettingsArea
					bind:downloadForm
					bind:indexerForm
					activeSection={activeSettingsSection}
					{downloadClients}
					{indexers}
					{savingDownloadClient}
					{savingIndexer}
					{testingDownloadClientId}
					{testingIndexerId}
					{downloadClientTests}
					{indexerTests}
					onSectionSelect={selectSettingsSection}
					onSaveDownloadClient={saveDownloadClient}
					onSaveIndexer={saveIndexer}
					onCancelDownloadClient={() => (downloadForm = emptyDownloadClientForm())}
					onCancelIndexer={() => (indexerForm = emptyIndexerForm())}
					onEditDownloadClient={(client) => {
						downloadForm = downloadClientFormFromClient(client);
						activeSettingsSection = 'download-clients';
						void goto(resolve('/settings/download-clients'));
					}}
					onEditIndexer={(indexer) => {
						indexerForm = indexerFormFromIndexer(indexer);
						activeSettingsSection = 'indexers';
						void goto(resolve('/settings/indexers'));
					}}
					onDeleteDownloadClient={deleteDownloadClient}
					onDeleteIndexer={deleteIndexer}
					onTestDownloadClient={testDownloadClient}
					onTestIndexer={testIndexer}
				/>
			{:else}
				<HomeArea
					activeSection={activeHomeSection}
					{selectedMediaItemId}
					{mediaItems}
					{mediaSearchResults}
					{releaseResults}
					{activities}
					{searchingMedia}
					{addingKey}
					{searchingItemId}
					{grabbingKey}
					{deletingMediaItemId}
					{loadingActivity}
					onSelect={selectHomeSection}
					onSearchMedia={searchMedia}
					onAddMedia={addMedia}
					onFindReleases={findReleases}
					onDeleteMedia={deleteMediaItem}
					onGrabRelease={grabRelease}
					onRefreshActivity={loadDownloadActivity}
				/>
			{/if}
		</main>
	</div>
{/if}
