<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import IntegrationTestStatus from '../shared/IntegrationTestStatus.svelte';
	import type {
		IntegrationTestResponse,
		MetadataProvider,
		MetadataProviderForm,
		MetadataProviderType
	} from '$lib/settings/types';

	interface ProviderDefinition {
		type: MetadataProviderType;
		name: string;
		baseUrl: string;
		priority: number;
		fields: 'tmdb' | 'tvdb';
	}

	interface Props {
		definition: ProviderDefinition;
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
		return {
			id: value.id,
			name: value.name,
			type: value.type,
			baseUrl: value.baseUrl,
			apiKey: value.apiKey ?? '',
			pin: value.pin ?? '',
			accessToken: value.accessToken ?? '',
			enabled: value.enabled,
			priority: value.priority
		};
	}

	function defaultForm(value: ProviderDefinition): MetadataProviderForm {
		return {
			name: value.name,
			type: value.type,
			baseUrl: value.baseUrl,
			apiKey: '',
			pin: '',
			accessToken: '',
			enabled: true,
			priority: value.priority
		};
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

	function updatePriority(value: number | undefined) {
		form = { ...form, priority: value ?? 0 };
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
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">API key</span>
				<Input
					value={form.apiKey}
					autocomplete="off"
					oninput={(event) => updateText('apiKey', event.currentTarget.value)}
				/>
			</label>
			{#if definition.fields === 'tvdb'}
				<label class="grid gap-1.5">
					<span class="text-sm font-bold text-muted-foreground">PIN</span>
					<Input
						value={form.pin}
						autocomplete="off"
						oninput={(event) => updateText('pin', event.currentTarget.value)}
					/>
				</label>
			{/if}
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">Access token</span>
				<Input
					value={form.accessToken}
					autocomplete="off"
					oninput={(event) => updateText('accessToken', event.currentTarget.value)}
				/>
			</label>
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
