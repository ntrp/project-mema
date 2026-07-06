<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import SecretInput from '$lib/components/settings/shared/SecretInput.svelte';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import IntegrationTestStatus from '../shared/IntegrationTestStatus.svelte';
	import type {
		DownloadClientForm,
		DownloadClientType,
		IntegrationTestResponse
	} from '$lib/settings/types';
	import { downloadClientProtocolForType } from '$lib/settings/forms';

	interface Props {
		form: DownloadClientForm;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onTest: () => boolean | void | Promise<boolean | void>;
		showTypeSelect?: boolean;
		testing?: boolean;
		testResult?: IntegrationTestResponse;
	}

	let {
		form = $bindable(),
		saving,
		onSave,
		onCancel,
		onTest,
		showTypeSelect = true,
		testing = false,
		testResult
	}: Props = $props();
	const downloadClientTypes: { value: DownloadClientType; label: string }[] = [
		{ value: 'transmission', label: 'transmission' },
		{ value: 'sabnzbd', label: 'sabnzbd' }
	];
	const protocolLabel = $derived(downloadClientProtocolForType(form.type).toUpperCase());

	function setDownloadClientType(value: string) {
		form.type = value as DownloadClientType;
		form.protocol = downloadClientProtocolForType(form.type);
	}
</script>

<Card.Root aria-labelledby="download-client-form-title">
	<Card.Header>
		<Card.Title id="download-client-form-title">
			{form.id ? 'Edit download client' : 'Add download client'}
		</Card.Title>
		{#if form.id}
			<Card.Action>
				<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
			</Card.Action>
		{/if}
	</Card.Header>

	<Card.Content>
		<form class="grid gap-4 sm:grid-cols-2" onsubmit={onSave}>
			<div class="space-y-2">
				<Label for="download-client-name">Name</Label>
				<Input id="download-client-name" bind:value={form.name} required maxlength={200} />
			</div>
			{#if showTypeSelect}
				<div class="space-y-2">
					<Label>Type</Label>
					<SettingsSelect
						value={form.type}
						options={downloadClientTypes}
						onValueChange={setDownloadClientType}
					/>
				</div>
			{/if}
			<div class="space-y-2">
				<Label>Protocol</Label>
				<div class="flex h-9 items-center">
					<Badge variant="outline" class="uppercase">{protocolLabel}</Badge>
				</div>
			</div>
			<div class="space-y-2 sm:col-span-2">
				<Label for="download-client-base-url">Base URL</Label>
				<Input
					id="download-client-base-url"
					bind:value={form.baseUrl}
					placeholder="http://host:port"
					required
				/>
			</div>
			{#if form.type === 'transmission'}
				<div class="space-y-2">
					<Label for="download-client-username">Username</Label>
					<Input id="download-client-username" bind:value={form.username} autocomplete="off" />
				</div>
				<div class="space-y-2">
					<Label for="download-client-password">Password</Label>
					<SecretInput
						id="download-client-password"
						bind:value={form.password}
						autocomplete="off"
					/>
				</div>
				<div class="space-y-2">
					<Label for="download-client-category">Category</Label>
					<Input id="download-client-category" bind:value={form.category} placeholder="movies" />
				</div>
			{:else}
				<div class="space-y-2 sm:col-span-2">
					<Label for="download-client-api-key">API key</Label>
					<SecretInput id="download-client-api-key" bind:value={form.apiKey} autocomplete="off" />
				</div>
				<div class="space-y-2">
					<Label for="download-client-category">Category</Label>
					<Input id="download-client-category" bind:value={form.category} placeholder="movies" />
				</div>
			{/if}
			<div class="space-y-2">
				<Label for="download-client-priority">Priority</Label>
				<Input
					id="download-client-priority"
					bind:value={form.priority}
					min="0"
					max="1000"
					type="number"
				/>
			</div>
			<div class="flex items-center gap-3 self-end py-2">
				<Switch id="download-client-enabled" bind:checked={form.enabled} />
				<Label for="download-client-enabled">Enabled</Label>
			</div>
			<div class="sm:col-span-2">
				<IntegrationTestStatus enabled={form.enabled} result={testResult} {testing} />
			</div>
			<div class="flex flex-wrap justify-end gap-2 sm:col-span-2">
				<Button type="button" variant="outline" disabled={saving || testing} onclick={onTest}>
					{testing ? 'Testing' : 'Test connection'}
				</Button>
				<Button type="submit" disabled={saving || testing}>
					{testing ? 'Testing' : saving ? 'Saving' : 'Save client'}
				</Button>
			</div>
		</form>
	</Card.Content>
</Card.Root>
