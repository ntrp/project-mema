<script lang="ts">
	import IndexerForm from '$lib/components/settings/indexers/IndexerForm.svelte';
	import IndexerSearchSettingsPanel from '$lib/components/settings/indexers/IndexerSearchSettings.svelte';
	import IndexerTable from '$lib/components/settings/indexers/IndexerTable.svelte';
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { emptyIndexerForm } from '$lib/settings/forms';
	import type {
		Indexer,
		IndexerForm as IndexerFormValue,
		IndexerSearchResponse,
		IndexerSearchSettings,
		IntegrationTestResults
	} from '$lib/settings/types';

	interface Props {
		indexers: Indexer[];
		indexerSearch: IndexerSearchResponse;
		form: IndexerFormValue;
		saving: boolean;
		clearingIndexerSearchCache: boolean;
		savingIndexerSearchSettings: boolean;
		testingId?: string;
		testResults: IntegrationTestResults;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onEdit: (_indexer: Indexer) => void;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		onClearIndexerSearchCache: () => void | Promise<void>;
		onSaveIndexerSearchSettings: (_settings: IndexerSearchSettings) => void | Promise<void>;
	}

	let {
		indexers,
		indexerSearch,
		form = $bindable(),
		saving,
		clearingIndexerSearchCache,
		savingIndexerSearchSettings,
		testingId,
		testResults,
		onSave,
		onCancel,
		onEdit,
		onDelete,
		onTest,
		onClearIndexerSearchCache,
		onSaveIndexerSearchSettings
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

<PageHeading eyebrow="Settings" title="Indexers" titleId="settings-title" />
<div class="space-y-4">
	<div class="flex justify-end">
		<SettingsAddButton label="Add indexer" onclick={openModal} />
	</div>
	<IndexerTable {indexers} onEdit={editIndexer} {onDelete} {onTest} {testingId} {testResults} />
	<IndexerSearchSettingsPanel
		search={indexerSearch}
		clearing={clearingIndexerSearchCache}
		saving={savingIndexerSearchSettings}
		onClearCache={onClearIndexerSearchCache}
		onSaveSettings={onSaveIndexerSearchSettings}
	/>
	{#if modalOpen}
		<SettingsFormModal title={form.id ? 'Edit indexer' : 'Add indexer'} onClose={closeModal}>
			<IndexerForm bind:form {saving} onSave={save} onCancel={closeModal} />
		</SettingsFormModal>
	{/if}
</div>
