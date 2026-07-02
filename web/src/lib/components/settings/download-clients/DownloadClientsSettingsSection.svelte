<script lang="ts">
	import ArrowLeftRightIcon from '@lucide/svelte/icons/arrow-left-right';
	import CloudDownloadIcon from '@lucide/svelte/icons/cloud-download';
	import DownloadClientForm from '$lib/components/settings/download-clients/DownloadClientForm.svelte';
	import DownloadClientTable from '$lib/components/settings/download-clients/DownloadClientTable.svelte';
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
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

<PageHeading eyebrow="Settings" title="Download clients" titleId="settings-title" />
<div class="space-y-4">
	<div class="flex justify-end">
		<SettingsAddButton label="Add download client" onclick={openModal} />
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
				<div class="grid gap-3 sm:grid-cols-2" aria-label="Download client type">
					<Button
						type="button"
						variant="outline"
						class="h-auto flex-col items-start gap-2 p-4 text-left"
						onclick={() => selectType('transmission')}
					>
						<ArrowLeftRightIcon aria-hidden="true" />
						<strong>Transmission</strong>
						<small>Torrent download client</small>
					</Button>
					<Button
						type="button"
						variant="outline"
						class="h-auto flex-col items-start gap-2 p-4 text-left"
						onclick={() => selectType('sabnzbd')}
					>
						<CloudDownloadIcon aria-hidden="true" />
						<strong>SABnzbd</strong>
						<small>Usenet download client</small>
					</Button>
				</div>
			{/if}
		</SettingsFormModal>
	{/if}
</div>
