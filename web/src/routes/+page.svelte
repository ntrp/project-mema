<script lang="ts">
	import { onMount } from 'svelte';

	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import DownloadClientForm from '$lib/components/settings/DownloadClientForm.svelte';
	import DownloadClientTable from '$lib/components/settings/DownloadClientTable.svelte';
	import IndexerForm from '$lib/components/settings/IndexerForm.svelte';
	import IndexerTable from '$lib/components/settings/IndexerTable.svelte';
	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import SettingsHeader from '$lib/components/settings/SettingsHeader.svelte';
	import {
		currentSessionAuthenticated,
		deleteDownloadClient as deleteDownloadClientRequest,
		deleteIndexer as deleteIndexerRequest,
		loadSettings as loadSettingsRequest,
		login as loginRequest,
		saveDownloadClient as saveDownloadClientRequest,
		saveIndexer as saveIndexerRequest
	} from '$lib/settings/api';
	import {
		downloadClientFormFromClient,
		emptyDownloadClientForm,
		emptyIndexerForm,
		indexerFormFromIndexer
	} from '$lib/settings/forms';
	import '$lib/settings/styles.css';
	import type {
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		Indexer,
		IndexerForm as IndexerFormValue
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

	const enabledDownloadClients = $derived(downloadClients.filter((item) => item.enabled).length);
	const enabledIndexers = $derived(indexers.filter((item) => item.enabled).length);

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
			message = 'Indexer deleted';
			await loadSettings();
		} catch (error) {
			errorMessage = errorMessageFrom(error, 'Could not delete indexer');
		}
	}

	function clearNotice() {
		errorMessage = '';
		message = '';
	}

	function errorMessageFrom(error: unknown, fallback: string) {
		return error instanceof Error ? error.message : fallback;
	}
</script>

<svelte:head>
	<title>Media Manager Settings</title>
	<meta
		name="description"
		content="Configure download clients and indexers for the self-hosted media manager"
	/>
</svelte:head>

<main class="shell">
	<SettingsHeader
		clientCount={downloadClients.length}
		enabledClientCount={enabledDownloadClients}
		indexerCount={indexers.length}
		enabledIndexerCount={enabledIndexers}
	/>

	{#if loading}
		<section class="panel">
			<p class="muted">Loading settings</p>
		</section>
	{:else if !authenticated}
		<AuthPanel bind:username bind:password onLogin={login} />
	{:else}
		<NoticeStack {message} {errorMessage} />

		<section class="settings-grid">
			<DownloadClientForm
				bind:form={downloadForm}
				saving={savingDownloadClient}
				onSave={saveDownloadClient}
				onCancel={() => (downloadForm = emptyDownloadClientForm())}
			/>
			<IndexerForm
				bind:form={indexerForm}
				saving={savingIndexer}
				onSave={saveIndexer}
				onCancel={() => (indexerForm = emptyIndexerForm())}
			/>
		</section>

		<section class="list-grid">
			<DownloadClientTable
				clients={downloadClients}
				onEdit={(client) => (downloadForm = downloadClientFormFromClient(client))}
				onDelete={deleteDownloadClient}
			/>
			<IndexerTable
				{indexers}
				onEdit={(indexer) => (indexerForm = indexerFormFromIndexer(indexer))}
				onDelete={deleteIndexer}
			/>
		</section>
	{/if}
</main>
