<script lang="ts">
	import CustomFormatForm from '$lib/components/settings/CustomFormatForm.svelte';
	import CustomFormatImportModal from '$lib/components/settings/CustomFormatImportModal.svelte';
	import CustomFormatTestParsingModal from '$lib/components/settings/CustomFormatTestParsingModal.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { emptyCustomFormatForm } from '$lib/settings/forms';
	import type {
		CustomFormat,
		CustomFormatForm as CustomFormatFormValue
	} from '$lib/settings/types';

	interface Props {
		formats: CustomFormat[];
		form: CustomFormatFormValue;
		saving: boolean;
		deletingId?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onImport: (_format: CustomFormatFormValue) => void | Promise<void>;
		onEdit: (_format: CustomFormat) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let {
		formats,
		form = $bindable(),
		saving,
		deletingId,
		onSave,
		onCancel,
		onImport,
		onEdit,
		onDelete
	}: Props = $props();

	let modalOpen = $state(false);
	let testParsingOpen = $state(false);
	let importOpen = $state(false);

	function openModal() {
		form = emptyCustomFormatForm();
		modalOpen = true;
	}

	function openTestParsing() {
		testParsingOpen = true;
	}

	function openImport() {
		importOpen = true;
	}

	function editFormat(format: CustomFormat) {
		onEdit(format);
		modalOpen = true;
	}

	function closeModal() {
		onCancel();
		modalOpen = false;
	}

	async function saveFormat(event: SubmitEvent) {
		await onSave(event);
		if (!form.id && form.name === '') {
			modalOpen = false;
		}
	}
</script>

<div class="panel" aria-label="Custom formats">
	<div class="section-heading">
		<div class="custom-format-topbar">
			<button type="button" class="secondary" onclick={openImport}>
				<span class="app-icon" aria-hidden="true">upload_file</span>
				Import
			</button>
			<button type="button" class="secondary" onclick={openTestParsing}>
				<span class="app-icon" aria-hidden="true">rule</span>
				Test parsing
			</button>
			<button type="button" class="add-action-button" onclick={openModal}>
				<span class="app-icon" aria-hidden="true">add</span>
				<span>Add custom format</span>
			</button>
		</div>
	</div>

	<div class="custom-format-grid">
		{#each formats as format (format.id)}
			<article class="custom-format-card">
				<div class="custom-format-card-header">
					<h3>{format.name}</h3>
					<div class="row-actions">
						<button
							type="button"
							class="secondary icon-button"
							aria-label={`Edit ${format.name}`}
							onclick={() => editFormat(format)}
						>
							<span class="app-icon" aria-hidden="true">edit</span>
						</button>
						<button
							type="button"
							class="danger icon-button"
							disabled={deletingId === format.id}
							aria-label={`${deletingId === format.id ? 'Deleting' : 'Delete'} ${format.name}`}
							onclick={() => onDelete(format.id)}
						>
							<span class="app-icon" aria-hidden="true">delete</span>
						</button>
					</div>
				</div>
				<div class="custom-format-tags">
					{#each format.includeSpecs as spec (spec.id)}
						<span
							class:include={spec.required}
							class:optional={!spec.required}
							title={`${spec.type}: ${spec.value}`}
						>
							{spec.name}
						</span>
					{/each}
					{#each format.excludeSpecs as spec (spec.id)}
						<span class="exclude" title={`${spec.type}: ${spec.value}`}>{spec.name}</span>
					{/each}
				</div>
			</article>
		{:else}
			<p class="empty custom-format-empty">No custom formats configured</p>
		{/each}
	</div>

	{#if modalOpen}
		<SettingsFormModal
			title={form.id ? 'Edit custom format' : 'Add custom format'}
			onClose={closeModal}
		>
			<CustomFormatForm bind:form {saving} onSave={saveFormat} onCancel={closeModal} />
		</SettingsFormModal>
	{/if}

	{#if testParsingOpen}
		<CustomFormatTestParsingModal onClose={() => (testParsingOpen = false)} />
	{/if}

	{#if importOpen}
		<CustomFormatImportModal onClose={() => (importOpen = false)} {onImport} />
	{/if}
</div>
