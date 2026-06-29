<script lang="ts">
	import DownloadClientForm from '$lib/components/settings/DownloadClientForm.svelte';
	import DownloadClientTable from '$lib/components/settings/DownloadClientTable.svelte';
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import { emptyDownloadClientForm } from '$lib/settings/forms';
	import type {
		DownloadClient,
		DownloadClientForm as DownloadClientFormValue,
		DownloadClientType,
		IntegrationTestResponse
	} from '$lib/settings/types';

	interface Props {
		clients: DownloadClient[];
		form: DownloadClientFormValue;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onTestConfig: (_form: DownloadClientFormValue) => Promise<IntegrationTestResponse>;
		onCancel: () => void;
		onEdit: (_client: DownloadClient) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let {
		clients,
		form = $bindable(),
		saving,
		onSave,
		onTestConfig,
		onCancel,
		onEdit,
		onDelete
	}: Props = $props();

	let modalOpen = $state(false);
	let typeSelected = $state(false);
	let testingConfig = $state(false);
	let testResult = $state<IntegrationTestResponse | undefined>();

	function openModal() {
		form = emptyDownloadClientForm();
		typeSelected = false;
		testResult = undefined;
		modalOpen = true;
	}

	function editClient(client: DownloadClient) {
		onEdit(client);
		typeSelected = true;
		testResult = undefined;
		modalOpen = true;
	}

	function closeModal() {
		onCancel();
		modalOpen = false;
		typeSelected = false;
		testResult = undefined;
	}

	function selectType(type: DownloadClientType) {
		form = { ...emptyDownloadClientForm(), type };
		typeSelected = true;
		testResult = undefined;
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();
		const passed = await testConfig();
		if (!passed) {
			return;
		}
		await onSave(event);
		if (isEmptyForm(form)) {
			modalOpen = false;
		}
	}

	async function testConfig() {
		testingConfig = true;
		testResult = undefined;
		try {
			const result = await onTestConfig(form);
			testResult = result;
			return result.success;
		} finally {
			testingConfig = false;
		}
	}

	function isEmptyForm(value: DownloadClientFormValue) {
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
</script>

<div class="page-heading">
	<p>Settings</p>
	<h1 id="settings-title">Download clients</h1>
</div>
<div class="settings-stack">
	<div class="settings-toolbar">
		<button type="button" onclick={openModal}>Add download client</button>
	</div>
	<DownloadClientTable {clients} onEdit={editClient} {onDelete} />
	{#if modalOpen}
		<SettingsFormModal
			title={form.id ? 'Edit download client' : 'Add download client'}
			onClose={closeModal}
		>
			{#if typeSelected}
				<DownloadClientForm
					bind:form
					{saving}
					onSave={save}
					onCancel={closeModal}
					onTest={testConfig}
					showTypeSelect={Boolean(form.id)}
					testing={testingConfig}
					{testResult}
				/>
			{:else}
				<div class="download-client-picker" aria-label="Download client type">
					<button type="button" onclick={() => selectType('transmission')}>
						<span class="app-icon" aria-hidden="true">sync_alt</span>
						<strong>Transmission</strong>
						<small>Torrent download client</small>
					</button>
					<button type="button" onclick={() => selectType('sabnzbd')}>
						<span class="app-icon" aria-hidden="true">cloud_download</span>
						<strong>SABnzbd</strong>
						<small>Usenet download client</small>
					</button>
				</div>
			{/if}
		</SettingsFormModal>
	{/if}
</div>
