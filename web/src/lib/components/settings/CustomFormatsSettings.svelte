<script lang="ts">
	import ListChecksIcon from '@lucide/svelte/icons/list-checks';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import UploadIcon from '@lucide/svelte/icons/upload';
	import CustomFormatCard from '$lib/components/settings/CustomFormatCard.svelte';
	import CustomFormatForm from '$lib/components/settings/CustomFormatForm.svelte';
	import CustomFormatImportModal from '$lib/components/settings/CustomFormatImportModal.svelte';
	import CustomFormatTestParsingModal from '$lib/components/settings/CustomFormatTestParsingModal.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
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

<Card class="gap-4 p-5" aria-label="Custom formats">
	<SectionHeading>
		{#snippet actions()}
			<div class="flex flex-wrap justify-end gap-2.5">
				<Button type="button" variant="outline" onclick={openImport}>
					<UploadIcon aria-hidden="true" />
					Import
				</Button>
				<Button type="button" variant="outline" onclick={openTestParsing}>
					<ListChecksIcon aria-hidden="true" />
					Test parsing
				</Button>
				<Button type="button" onclick={openModal}>
					<PlusIcon aria-hidden="true" />
					<span>Add custom format</span>
				</Button>
			</div>
		{/snippet}
	</SectionHeading>

	<div class="grid gap-3 [grid-template-columns:repeat(auto-fit,minmax(min(100%,360px),1fr))]">
		{#each formats as format (format.id)}
			<CustomFormatCard
				{format}
				deleting={deletingId === format.id}
				onEdit={editFormat}
				{onDelete}
			/>
		{:else}
			<p class="col-span-full m-0 text-sm leading-6 text-muted-foreground">
				No custom formats configured
			</p>
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
</Card>
