<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import IndexerCatalogPicker from '$lib/components/settings/indexers/IndexerCatalogPicker.svelte';
	import { flattenCategories } from '$lib/components/settings/indexers/indexerCatalogFilters';
	import type { IndexerCatalogEntry, IndexerForm } from '$lib/settings/types';

	interface Props {
		form: IndexerForm;
		catalog: IndexerCatalogEntry[];
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), catalog, saving, onSave, onCancel }: Props = $props();
	let configuring = $state(Boolean(form.id));
	const selected = $derived(catalog.find((entry) => entry.definitionId === form.definitionId));

	function applyDefinition(entry: IndexerCatalogEntry) {
		form.definitionId = entry.definitionId;
		form.name = entry.name;
		form.implementation = entry.implementation;
		form.implementationName = entry.implementationName;
		form.baseUrl = entry.indexerUrls?.[0] ?? '';
		form.fields = entry.fields.map((field) => ({ name: field.name, value: field.value ?? '' }));
		form.categoriesText = flattenCategories(entry.capabilities.categories)
			.map((category) => category.id)
			.join(', ');
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
</script>

<Card.Root aria-labelledby="indexer-form-title">
	<Card.Header>
		<Card.Title id="indexer-form-title">
			{form.id ? 'Edit indexer' : configuring ? 'Configure indexer' : 'Add indexer'}
		</Card.Title>
		<Card.Action>
			<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
		</Card.Action>
	</Card.Header>

	<Card.Content>
		{#if !form.id && !configuring}
			<IndexerCatalogPicker {catalog} onSelect={applyDefinition} />
		{:else}
			<form class="grid gap-4 sm:grid-cols-2" onsubmit={onSave}>
				{#if !form.id}
					<div class="flex items-start justify-between gap-3 sm:col-span-2">
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
					<Input
						id="indexer-priority"
						bind:value={form.priority}
						min="0"
						max="1000"
						type="number"
					/>
				</div>
				<div class="space-y-2 sm:col-span-2">
					<Label for="indexer-base-url">Base URL</Label>
					<Input id="indexer-base-url" bind:value={form.baseUrl} required />
				</div>
				<div class="space-y-2">
					<Label for="indexer-api-key">API key</Label>
					<Input id="indexer-api-key" bind:value={form.apiKey} autocomplete="off" />
				</div>
				<div class="space-y-2">
					<Label for="indexer-categories">Categories</Label>
					<Input id="indexer-categories" bind:value={form.categoriesText} />
				</div>
				{#each selected?.fields ?? [] as field (field.name)}
					{#if field.name !== 'baseUrl' && field.name !== 'apiKey' && field.name !== 'categories' && field.type !== 'info'}
						<div class="space-y-2">
							<Label for={`indexer-field-${field.name}`}>{field.label}</Label>
							<Input
								id={`indexer-field-${field.name}`}
								value={String(fieldValue(field.name))}
								oninput={(event) => updateField(field.name, event.currentTarget.value)}
							/>
						</div>
					{/if}
				{/each}
				<div class="flex items-center gap-3 self-end py-2">
					<Switch id="indexer-enabled" bind:checked={form.enabled} />
					<Label for="indexer-enabled">Enabled</Label>
				</div>
				<Button class="w-fit" type="submit" disabled={saving}
					>{saving ? 'Saving' : 'Save indexer'}</Button
				>
			</form>
		{/if}
	</Card.Content>
</Card.Root>
