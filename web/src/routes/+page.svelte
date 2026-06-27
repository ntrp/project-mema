<script lang="ts">
	import { onMount } from 'svelte';

	import AppNav from '$lib/components/app/AppNav.svelte';
	import HomeArea from '$lib/components/app/HomeArea.svelte';
	import SettingsArea from '$lib/components/app/SettingsArea.svelte';
	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import {
		currentSessionAuthenticated,
		deleteDownloadClient as deleteDownloadClientRequest,
		deleteIndexer as deleteIndexerRequest,
		loadSettings as loadSettingsRequest,
		login as loginRequest,
		logout as logoutRequest,
		saveDownloadClient as saveDownloadClientRequest,
		saveIndexer as saveIndexerRequest,
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
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		HomeSection,
		Indexer,
		IndexerForm as IndexerFormValue,
		IntegrationTestResults,
		SettingsSection
	} from '$lib/settings/types';

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
	let downloadForm = $state<DownloadClientFormValue>(emptyDownloadClientForm());
	let indexerForm = $state<IndexerFormValue>(emptyIndexerForm());
	let testingDownloadClientId = $state<string | undefined>();
	let testingIndexerId = $state<string | undefined>();
	let downloadClientTests = $state<IntegrationTestResults>({});
	let indexerTests = $state<IntegrationTestResults>({});
	let activeView = $state<AppView>('home');
	let activeHomeSection = $state<HomeSection>('explore');
	let activeSettingsSection = $state<SettingsSection>('download-clients');
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
	}

	function openHome() {
		activeView = 'home';
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

	function omitResult(results: IntegrationTestResults, id: string) {
		const { [id]: _removed, ...remaining } = results;
		return remaining;
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
					onSectionSelect={(section) => (activeSettingsSection = section)}
					onSaveDownloadClient={saveDownloadClient}
					onSaveIndexer={saveIndexer}
					onCancelDownloadClient={() => (downloadForm = emptyDownloadClientForm())}
					onCancelIndexer={() => (indexerForm = emptyIndexerForm())}
					onEditDownloadClient={(client) => {
						downloadForm = downloadClientFormFromClient(client);
						activeSettingsSection = 'download-clients';
					}}
					onEditIndexer={(indexer) => {
						indexerForm = indexerFormFromIndexer(indexer);
						activeSettingsSection = 'indexers';
					}}
					onDeleteDownloadClient={deleteDownloadClient}
					onDeleteIndexer={deleteIndexer}
					onTestDownloadClient={testDownloadClient}
					onTestIndexer={testIndexer}
				/>
			{:else}
				<HomeArea
					activeSection={activeHomeSection}
					onSelect={(section) => (activeHomeSection = section)}
				/>
			{/if}
		</main>
	</div>
{/if}
