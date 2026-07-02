<script lang="ts">
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { emptyLanguageForm } from '$lib/settings/forms';
	import type { Language, LanguageForm as LanguageFormValue } from '$lib/settings/types';
	import LanguageDeleteModal from './LanguageDeleteModal.svelte';
	import LanguageForm from './LanguageForm.svelte';
	import LanguageTable from './LanguageTable.svelte';

	interface Props {
		languages: Language[];
		form: LanguageFormValue;
		saving: boolean;
		deletingCode?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onEdit: (_language: Language) => void;
		onDelete: (_code: string) => void | Promise<void>;
	}

	let {
		languages,
		form = $bindable(),
		saving,
		deletingCode,
		onSave,
		onCancel,
		onEdit,
		onDelete
	}: Props = $props();

	let formOpen = $state(false);
	let deleteCandidate = $state<Language | undefined>();

	function addLanguage() {
		form = emptyLanguageForm();
		formOpen = true;
	}

	function editLanguage(language: Language) {
		onEdit(language);
		formOpen = true;
	}

	function closeForm() {
		onCancel();
		formOpen = false;
	}

	async function save(event: SubmitEvent) {
		await onSave(event);
		if (!form.originalCode && form.code === '' && form.displayName === '') {
			formOpen = false;
		}
	}
</script>

<PageHeading eyebrow="Settings" title="Languages" titleId="settings-title" />
<div class="space-y-4">
	<div class="flex justify-end">
		<SettingsAddButton label="Add language" onclick={addLanguage} />
	</div>
	<LanguageTable
		{languages}
		{deletingCode}
		onEdit={editLanguage}
		onDelete={(language) => (deleteCandidate = language)}
	/>
	{#if formOpen}
		<SettingsFormModal
			title={form.originalCode ? 'Edit language' : 'Add language'}
			onClose={closeForm}
		>
			<LanguageForm bind:form {saving} onSave={save} onCancel={closeForm} />
		</SettingsFormModal>
	{/if}
	{#if deleteCandidate}
		<LanguageDeleteModal
			language={deleteCandidate}
			deleting={deletingCode === deleteCandidate.code}
			onCancel={() => (deleteCandidate = undefined)}
			onConfirm={async () => {
				await onDelete(deleteCandidate?.code ?? '');
				deleteCandidate = undefined;
			}}
		/>
	{/if}
</div>
