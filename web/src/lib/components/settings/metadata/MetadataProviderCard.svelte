<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import SecretInput from '$lib/components/settings/shared/SecretInput.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import IntegrationTestStatus from '../shared/IntegrationTestStatus.svelte';
	import type { MetadataProviderDefinition } from './metadataProviderDefinitions';
	import type {
		IntegrationTestResponse,
		MetadataProvider,
		MetadataProviderForm
	} from '$lib/settings/types';

	interface Props {
		definition: MetadataProviderDefinition;
		provider?: MetadataProvider;
		onSave: (_form: MetadataProviderForm) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		testingId?: string;
		savingId?: string;
		testResult?: IntegrationTestResponse;
	}

	let { definition, provider, onSave, onTest, testingId, savingId, testResult }: Props = $props();
	let form = $derived(provider ? formFromProvider(provider) : defaultForm(definition));

	function formFromProvider(value: MetadataProvider): MetadataProviderForm {
		const form: MetadataProviderForm = {
			id: value.id,
			name: value.name,
			type: value.type,
			baseUrl: value.baseUrl,
			enabled: value.enabled,
			priority: value.priority
		};
		if (value.type === 'anilist') return form;
		return {
			...form,
			apiKey: value.apiKey ?? '',
			pin: value.pin ?? '',
			accessToken: value.accessToken ?? '',
			apiKeySet: value.apiKeySet,
			pinSet: value.pinSet,
			accessTokenSet: value.accessTokenSet
		};
	}

	function defaultForm(value: MetadataProviderDefinition): MetadataProviderForm {
		const form: MetadataProviderForm = {
			name: value.name,
			type: value.type,
			baseUrl: value.baseUrl,
			enabled: value.type !== 'anidb',
			priority: value.priority
		};
		return value.fields === 'anilist' ? form : { ...form, apiKey: '', pin: '', accessToken: '' };
	}

	function save(event: SubmitEvent) {
		event.preventDefault();
		void onSave(form);
	}

	function updateText(field: 'baseUrl' | 'apiKey' | 'pin' | 'accessToken', value: string) {
		form = { ...form, [field]: value };
	}

	function updateEnabled(value: boolean) {
		form = { ...form, enabled: value };
	}

	function updatePriority(value: number) {
		form = { ...form, priority: Number.isFinite(value) ? value : 0 };
	}
</script>

<div
	class="min-w-0 rounded-md border border-border bg-card p-5"
	aria-labelledby={`${definition.type}-metadata-title`}
>
	<form class="grid gap-3.5" onsubmit={save}>
		<SectionHeading
			title={definition.name}
			titleId={`${definition.type}-metadata-title`}
			kicker="Metadata provider"
		>
			{#snippet actions()}
				<IntegrationTestStatus
					enabled={form.enabled}
					result={testResult}
					testing={form.id ? testingId === form.id : false}
				/>
			{/snippet}
		</SectionHeading>

		<div class="mb-4 grid gap-4 md:grid-cols-2">
			<div class="flex min-h-9 items-center gap-2 self-end">
				<Checkbox
					id={`${definition.type}-metadata-enabled`}
					checked={form.enabled}
					onCheckedChange={(checked) => updateEnabled(checked === true)}
				/>
				<Label for={`${definition.type}-metadata-enabled`}>Enabled</Label>
			</div>
			<label class="grid gap-1.5 md:col-span-2">
				<span class="text-sm font-bold text-muted-foreground">Base URL</span>
				<Input
					value={form.baseUrl}
					required
					maxlength={2000}
					oninput={(event) => updateText('baseUrl', event.currentTarget.value)}
				/>
			</label>
			{#if definition.fields === 'anilist'}
				<p class="text-sm text-muted-foreground md:col-span-2">
					No credentials required. AniList uses its public GraphQL API.
				</p>
			{:else if definition.fields === 'anidb'}
				<label class="grid gap-1.5">
					<span class="text-sm font-bold text-muted-foreground">Registered client name</span>
					<Input
						value={form.apiKey}
						required
						pattern="[a-z0-9]+"
						oninput={(event) => updateText('apiKey', event.currentTarget.value)}
					/>
				</label>
				<label class="grid gap-1.5">
					<span class="text-sm font-bold text-muted-foreground">Client version</span>
					<Input
						value={form.pin}
						required
						min="1"
						type="number"
						oninput={(event) => updateText('pin', event.currentTarget.value)}
					/>
				</label>
			{:else}
				<label class="grid gap-1.5">
					<span class="text-sm font-bold text-muted-foreground">API key</span>
					<SecretInput
						value={form.apiKey}
						autocomplete="off"
						onValueChange={(value) => updateText('apiKey', value)}
					/>
				</label>
				{#if definition.fields === 'tvdb'}
					<label class="grid gap-1.5">
						<span class="text-sm font-bold text-muted-foreground">PIN</span>
						<SecretInput
							value={form.pin}
							autocomplete="off"
							onValueChange={(value) => updateText('pin', value)}
						/>
					</label>
				{/if}
				<label class="grid gap-1.5">
					<span class="text-sm font-bold text-muted-foreground">Access token</span>
					<SecretInput
						value={form.accessToken}
						autocomplete="off"
						onValueChange={(value) => updateText('accessToken', value)}
					/>
				</label>
			{/if}
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">Priority</span>
				<Input
					value={form.priority}
					min="0"
					max="1000"
					type="number"
					oninput={(event) => updatePriority(event.currentTarget.valueAsNumber)}
				/>
			</label>
		</div>

		<div class="flex flex-wrap justify-end gap-2.5">
			<Button type="submit" disabled={savingId === form.id}>
				{savingId === form.id ? 'Saving' : 'Save'}
			</Button>
			<Button
				type="button"
				variant="outline"
				disabled={!form.id || testingId === form.id}
				onclick={() => form.id && onTest(form.id)}
			>
				{testingId === form.id ? 'Testing' : 'Test'}
			</Button>
		</div>
	</form>
</div>
