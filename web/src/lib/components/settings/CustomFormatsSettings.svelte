<script lang="ts">
	import CustomFormatForm from '$lib/components/settings/CustomFormatForm.svelte';
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
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
		onEdit,
		onDelete
	}: Props = $props();

	let modalOpen = $state(false);

	function openModal() {
		form = emptyCustomFormatForm();
		modalOpen = true;
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

<div class="panel" aria-labelledby="custom-format-settings-title">
	<div class="section-heading">
		<div>
			<p class="section-kicker">Release scoring</p>
			<h2 id="custom-format-settings-title">Custom formats</h2>
		</div>
		<button type="button" onclick={openModal}>Add custom format</button>
	</div>

	<div class="custom-format-grid">
		{#each formats as format (format.id)}
			<article class="custom-format-card">
				<div class="custom-format-card-header">
					<h3>{format.name}</h3>
					<div class="row-actions">
						<button type="button" class="secondary" onclick={() => editFormat(format)}>Edit</button>
						<button
							type="button"
							class="danger"
							disabled={deletingId === format.id}
							onclick={() => onDelete(format.id)}
						>
							{deletingId === format.id ? 'Deleting' : 'Delete'}
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
</div>
