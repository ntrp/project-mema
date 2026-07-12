<script lang="ts">
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import * as Card from '$lib/components/ui/card';
	import { onMount } from 'svelte';
	import { listSubtitleProviderCatalog } from '$lib/settings/api';
	import { emptySubtitleProviderForm, subtitleProviderFormFromProvider } from '$lib/settings/forms';
	import type {
		IntegrationTestResponse,
		IntegrationTestResults,
		SubtitleProvider,
		SubtitleProviderCatalogEntry,
		SubtitleProviderForm as FormValue
	} from '$lib/settings/types';
	import SubtitleProviderCatalogPicker from './catalog/SubtitleProviderCatalogPicker.svelte';
	import SubtitleProviderForm from './form/SubtitleProviderForm.svelte';
	import SubtitleProviderTable from './list/SubtitleProviderTable.svelte';

	interface Props {
		providers: SubtitleProvider[];
		onSave: (_form: FormValue) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		onTestConfig: (_form: FormValue) => Promise<IntegrationTestResponse>;
		testingId?: string;
		savingId?: string;
		testResults: IntegrationTestResults;
	}

	let {
		providers,
		onSave,
		onDelete,
		onTest,
		onTestConfig,
		testingId,
		savingId,
		testResults
	}: Props = $props();
	let catalog = $state.raw<SubtitleProviderCatalogEntry[]>([]);
	let catalogError = $state('');
	let modalOpen = $state(false);

	onMount(() => {
		let cancelled = false;
		void listSubtitleProviderCatalog()
			.then((response) => {
				if (!cancelled) catalog = response.providers;
			})
			.catch((error: unknown) => {
				if (!cancelled) {
					catalogError = error instanceof Error ? error.message : 'Could not load provider catalog';
				}
			});
		return () => {
			cancelled = true;
		};
	});
	let picking = $state(false);
	let form = $state<FormValue>(emptySubtitleProviderForm());
	let selectedEntry = $state<SubtitleProviderCatalogEntry | undefined>();
	let testingConfig = $state(false);
	let testResult = $state<IntegrationTestResponse | undefined>();

	function openPicker() {
		picking = true;
		modalOpen = true;
		testResult = undefined;
	}

	function selectEntry(entry: SubtitleProviderCatalogEntry) {
		selectedEntry = entry;
		form = emptySubtitleProviderForm(entry.key as FormValue['type'], entry);
		if (entry.runtimeStatus !== 'supported') form = { ...form, enabled: false };
		picking = false;
	}

	function editProvider(provider: SubtitleProvider) {
		selectedEntry = catalog.find(
			(entry) => entry.key === provider.catalogKey || entry.key === provider.type
		);
		form = subtitleProviderFormFromProvider(provider);
		if (form.runtimeStatus !== 'supported') form = { ...form, enabled: false };
		picking = false;
		modalOpen = true;
		testResult = undefined;
	}

	function closeModal() {
		modalOpen = false;
		picking = false;
		selectedEntry = undefined;
		testResult = undefined;
		form = emptySubtitleProviderForm();
	}

	async function saveForm(event: SubmitEvent) {
		event.preventDefault();
		if (!canSave()) return;
		if (canTest()) {
			const passed = await testConfig();
			if (!passed) return;
		}
		await onSave(form);
		closeModal();
	}

	async function testConfig() {
		if (!canTest()) return false;
		testingConfig = true;
		testResult = undefined;
		try {
			testResult = await onTestConfig(form);
			return testResult.success;
		} finally {
			testingConfig = false;
		}
	}

	function runTestConfig() {
		void testConfig();
	}

	function canTest() {
		return (selectedEntry?.runtimeStatus ?? form.runtimeStatus ?? 'supported') === 'supported';
	}

	function canSave() {
		return canTest() || !form.enabled;
	}
</script>

<div class="space-y-4">
	<div class="flex justify-end">
		<SettingsAddButton label="Add subtitle provider" onclick={openPicker} />
	</div>
	{#if catalogError}
		<Card.Root>
			<Card.Content class="pt-6">
				<p class="m-0 text-sm text-destructive">{catalogError}</p>
			</Card.Content>
		</Card.Root>
	{/if}
	<SubtitleProviderTable
		{providers}
		{catalog}
		{onDelete}
		{onTest}
		{testingId}
		{testResults}
		onEdit={editProvider}
	/>
	{#if modalOpen}
		<SettingsFormModal
			title={picking
				? 'Add subtitle provider'
				: form.id
					? 'Edit subtitle provider'
					: 'Configure subtitle provider'}
			onClose={closeModal}
			modalClass="w-[min(1200px,calc(100vw-32px))]"
		>
			{#if picking}
				<SubtitleProviderCatalogPicker {catalog} onSelect={selectEntry} />
			{:else}
				<SubtitleProviderForm
					bind:form
					entry={selectedEntry}
					saving={Boolean(form.id && savingId === form.id)}
					testing={testingConfig}
					{testResult}
					onSave={saveForm}
					onCancel={closeModal}
					onTest={runTestConfig}
				/>
			{/if}
		</SettingsFormModal>
	{/if}
</div>
