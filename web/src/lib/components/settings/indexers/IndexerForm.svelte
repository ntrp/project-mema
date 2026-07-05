<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import IndexerCatalogPicker from '$lib/components/settings/indexers/IndexerCatalogPicker.svelte';
	import IndexerDynamicField from '$lib/components/settings/indexers/IndexerDynamicField.svelte';
	import IndexerScopeFields from '$lib/components/settings/indexers/IndexerScopeFields.svelte';
	import { flattenCategories } from '$lib/components/settings/indexers/indexerCatalogFilters';
	import IntegrationTestStatus from '../shared/IntegrationTestStatus.svelte';
	import type {
		IndexerCatalogEntry,
		IndexerForm,
		IntegrationTestResponse,
		Tag
	} from '$lib/settings/types';

	interface Props {
		form: IndexerForm;
		catalog: IndexerCatalogEntry[];
		tags: Tag[];
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onTest: () => boolean | void | Promise<boolean | void>;
		testing?: boolean;
		testResult?: IntegrationTestResponse;
	}

	let {
		form = $bindable(),
		catalog,
		tags,
		saving,
		onSave,
		onTest,
		testing = false,
		testResult
	}: Props = $props();
	let configuring = $state(Boolean(form.id));
	const selected = $derived(catalog.find((entry) => entry.definitionId === form.definitionId));

	function applyDefinition(entry: IndexerCatalogEntry) {
		form.definitionId = entry.definitionId;
		form.name = entry.name;
		form.implementation = entry.implementation;
		form.implementationName = entry.implementationName;
		form.baseUrl = entry.indexerUrls?.[0] ?? '';
		form.fields = entry.fields.map((field) => ({
			name: field.name,
			value: defaultFieldValue(field)
		}));
		form.categoriesText = flattenCategories(entry.capabilities.categories)
			.map((category) => category.id)
			.join(', ');
		form.mediaTypeScopes = entry.mediaTypeScopes ?? ['movie', 'serie', 'anime', 'audio', 'book'];
		form.tagScopes = [];
		form.redirect = entry.supportsRedirect;
		configuring = true;
	}

	function fieldValue(name: string) {
		return form.fields?.find((field) => field.name === name)?.value ?? '';
	}

	function updateField(name: string, value: unknown) {
		const fields = [...(form.fields ?? [])];
		const index = fields.findIndex((field) => field.name === name);
		if (index >= 0) {
			fields[index] = { name, value };
		} else {
			fields.push({ name, value });
		}
		form.fields = fields;
	}

	function defaultFieldValue(field: IndexerCatalogEntry['fields'][number]) {
		if (field.value != null) {
			return field.value;
		}
		if (field.type === 'checkbox') {
			return false;
		}
		if (field.type === 'select') {
			return field.selectOptions?.[0]?.value ?? '';
		}
		return '';
	}

	function usesBuiltInControl(name: string) {
		return name === 'baseUrl' || name === 'apiKey' || name === 'categories';
	}
</script>

{#if !form.id && !configuring}
	<IndexerCatalogPicker {catalog} onSelect={applyDefinition} />
{:else}
	<form class="grid gap-4" onsubmit={onSave}>
		{#if !form.id}
			<div class="flex items-start justify-between gap-3">
				{#if selected}
					<div>
						<div class="font-bold text-foreground">{selected.name}</div>
						<p class="m-0 text-sm leading-6 text-muted-foreground">{selected.description}</p>
					</div>
				{/if}
				<Button type="button" variant="outline" onclick={() => (configuring = false)}>
					Back to catalog
				</Button>
			</div>
		{/if}
		<div class="space-y-2">
			<Label for="indexer-name">Name</Label>
			<Input id="indexer-name" bind:value={form.name} required maxlength={200} />
		</div>
		<div class="space-y-2">
			<Label for="indexer-priority">Priority</Label>
			<Input id="indexer-priority" bind:value={form.priority} min="0" max="1000" type="number" />
		</div>
		<div class="space-y-2">
			<Label for="indexer-base-url">Base URL</Label>
			<Input id="indexer-base-url" bind:value={form.baseUrl} type="url" required />
		</div>
		<div class="space-y-2">
			<Label for="indexer-api-key">API key</Label>
			<Input
				id="indexer-api-key"
				bind:value={form.apiKey}
				autocomplete="off"
				placeholder={form.apiKeySet ? 'Saved API key' : ''}
				type="password"
			/>
		</div>
		<div class="space-y-2">
			<Label for="indexer-categories">Categories</Label>
			<Input id="indexer-categories" bind:value={form.categoriesText} />
		</div>
		{#each selected?.fields ?? [] as field (field.name)}
			{#if !usesBuiltInControl(field.name)}
				<IndexerDynamicField
					{field}
					value={fieldValue(field.name)}
					onValueChange={(value) => updateField(field.name, value)}
				/>
			{/if}
		{/each}
		<IndexerScopeFields bind:form {tags} />
		<div class="flex items-center gap-3 py-2">
			<Switch id="indexer-enabled" bind:checked={form.enabled} />
			<Label for="indexer-enabled">Enabled</Label>
		</div>
		<IntegrationTestStatus enabled={form.enabled} result={testResult} {testing} />
		<div class="flex flex-wrap justify-end gap-2">
			<Button type="button" variant="outline" disabled={saving || testing} onclick={onTest}>
				{testing ? 'Testing' : 'Test indexer'}
			</Button>
			<Button type="submit" disabled={saving || testing}>
				{testing ? 'Testing' : saving ? 'Saving' : 'Save indexer'}
			</Button>
		</div>
	</form>
{/if}
