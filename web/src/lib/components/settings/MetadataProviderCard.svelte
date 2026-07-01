<script lang="ts">
	import IntegrationTestStatus from './shared/IntegrationTestStatus.svelte';
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

<div class="panel provider-panel" aria-labelledby={`${definition.type}-metadata-title`}>
	<form class="provider-form" onsubmit={save}>
		<div class="section-heading">
			<div>
				<p class="section-kicker">Metadata provider</p>
				<h2 id={`${definition.type}-metadata-title`}>{definition.name}</h2>
			</div>
			<IntegrationTestStatus
				enabled={form.enabled}
				result={testResult}
				testing={form.id ? testingId === form.id : false}
			/>
		</div>

		<div class="settings-form provider-form-fields">
			<label class="toggle provider-enabled">
				<input
					checked={form.enabled}
					type="checkbox"
					onchange={(event) => updateEnabled(event.currentTarget.checked)}
				/>
				<span>Enabled</span>
			</label>
			<label class="wide">
				<span>Base URL</span>
				<input
					value={form.baseUrl}
					required
					maxlength="2000"
					oninput={(event) => updateText('baseUrl', event.currentTarget.value)}
				/>
			</label>
			<label>
				<span>API key</span>
				<input
					value={form.apiKey}
					autocomplete="off"
					oninput={(event) => updateText('apiKey', event.currentTarget.value)}
				/>
			</label>
			{#if definition.fields === 'tvdb'}
				<label>
					<span>PIN</span>
					<input
						value={form.pin}
						autocomplete="off"
						oninput={(event) => updateText('pin', event.currentTarget.value)}
					/>
				</label>
			{/if}
			<label>
				<span>Access token</span>
				<input
					value={form.accessToken}
					autocomplete="off"
					oninput={(event) => updateText('accessToken', event.currentTarget.value)}
				/>
			</label>
			<label>
				<span>Priority</span>
				<input
					value={form.priority}
					min="0"
					max="1000"
					type="number"
					oninput={(event) => updatePriority(event.currentTarget.valueAsNumber)}
				/>
			</label>
		</div>

		<div class="form-actions provider-actions">
			<button type="submit" disabled={savingId === form.id}>
				{savingId === form.id ? 'Saving' : 'Save'}
			</button>
			<button
				type="button"
				class="secondary"
				disabled={!form.id || testingId === form.id}
				onclick={() => form.id && onTest(form.id)}
			>
				{testingId === form.id ? 'Testing' : 'Test'}
			</button>
		</div>
	</form>
</div>
