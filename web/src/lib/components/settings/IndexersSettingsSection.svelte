<script lang="ts">
	import IndexerForm from '$lib/components/settings/IndexerForm.svelte';
	import IndexerTable from '$lib/components/settings/IndexerTable.svelte';
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import { emptyIndexerForm } from '$lib/settings/forms';
	import type {
		Indexer,
		IndexerForm as IndexerFormValue,
		IntegrationTestResults
	} from '$lib/settings/types';

	interface Props {
		indexers: Indexer[];
		form: IndexerFormValue;
		saving: boolean;
		testingId?: string;
		testResults: IntegrationTestResults;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onEdit: (_indexer: Indexer) => void;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
	}

	let {
		indexers,
		form = $bindable(),
		saving,
		testingId,
		testResults,
		onSave,
		onCancel,
		onEdit,
		onDelete,
		onTest
	}: Props = $props();

	let modalOpen = $state(false);

	function openModal() {
		form = emptyIndexerForm();
		modalOpen = true;
	}

	function editIndexer(indexer: Indexer) {
		onEdit(indexer);
		modalOpen = true;
	}

	function closeModal() {
		onCancel();
		modalOpen = false;
	}

	async function save(event: SubmitEvent) {
		await onSave(event);
		if (!form.id && form.name === '' && form.baseUrl === '' && form.apiKey === '') {
			modalOpen = false;
		}
	}
</script>

<div class="page-heading">
	<p>Settings</p>
	<h1 id="settings-title">Indexers</h1>
</div>
<div class="settings-stack">
	<div class="settings-toolbar">
		<button type="button" onclick={openModal}>Add indexer</button>
	</div>
	<IndexerTable {indexers} onEdit={editIndexer} {onDelete} {onTest} {testingId} {testResults} />
	{#if modalOpen}
		<SettingsFormModal title={form.id ? 'Edit indexer' : 'Add indexer'} onClose={closeModal}>
			<IndexerForm bind:form {saving} onSave={save} onCancel={closeModal} />
		</SettingsFormModal>
	{/if}
</div>
