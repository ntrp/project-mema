<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import type { IndexerForm, IndexerType } from '$lib/settings/types';

	interface Props {
		form: IndexerForm;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), saving, onSave, onCancel }: Props = $props();
	const indexerTypes: { value: IndexerType; label: string }[] = [
		{ value: 'torznab', label: 'torznab' },
		{ value: 'newznab', label: 'newznab' },
		{ value: 'rss', label: 'rss' }
	];
</script>

<Card.Root aria-labelledby="indexer-form-title">
	<Card.Header>
		<Card.Title id="indexer-form-title">{form.id ? 'Edit indexer' : 'Add indexer'}</Card.Title>
		{#if form.id}
			<Card.Action>
				<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
			</Card.Action>
		{/if}
	</Card.Header>

	<Card.Content>
		<form class="grid gap-4 sm:grid-cols-2" onsubmit={onSave}>
			<div class="space-y-2">
				<Label for="indexer-name">Name</Label>
				<Input id="indexer-name" bind:value={form.name} required maxlength={200} />
			</div>
			<div class="space-y-2">
				<Label>Type</Label>
				<SettingsSelect
					value={form.type}
					options={indexerTypes}
					onValueChange={(value) => (form.type = value as IndexerType)}
				/>
			</div>
			<div class="space-y-2 sm:col-span-2">
				<Label for="indexer-base-url">Base URL</Label>
				<Input
					id="indexer-base-url"
					bind:value={form.baseUrl}
					placeholder="https://indexer.example"
					required
				/>
			</div>
			<div class="space-y-2">
				<Label for="indexer-api-key">API key</Label>
				<Input id="indexer-api-key" bind:value={form.apiKey} autocomplete="off" />
			</div>
			<div class="space-y-2">
				<Label for="indexer-categories">Categories</Label>
				<Input id="indexer-categories" bind:value={form.categoriesText} placeholder="2000, 5000" />
			</div>
			<div class="space-y-2">
				<Label for="indexer-priority">Priority</Label>
				<Input id="indexer-priority" bind:value={form.priority} min="0" max="1000" type="number" />
			</div>
			<div class="flex items-center gap-3 self-end py-2">
				<Switch id="indexer-enabled" bind:checked={form.enabled} />
				<Label for="indexer-enabled">Enabled</Label>
			</div>
			<Button class="w-fit" type="submit" disabled={saving}
				>{saving ? 'Saving' : 'Save indexer'}</Button
			>
		</form>
	</Card.Content>
</Card.Root>
