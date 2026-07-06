<script lang="ts">
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import IntegrationTestStatus from '$lib/components/settings/shared/IntegrationTestStatus.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import SecretInput from '$lib/components/settings/shared/SecretInput.svelte';
	import { emptySubtitleProviderForm, subtitleProviderFormFromProvider } from '$lib/settings/forms';
	import type {
		IntegrationTestResponse,
		SubtitleProvider,
		SubtitleProviderForm
	} from '$lib/settings/types';

	interface Props {
		provider?: SubtitleProvider;
		onSave: (_form: SubtitleProviderForm) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
		onTest: (_id: string) => void | Promise<void>;
		testingId?: string;
		savingId?: string;
		testResult?: IntegrationTestResponse;
	}

	let { provider, onSave, onDelete, onTest, testingId, savingId, testResult }: Props = $props();
	let form = $derived(
		provider ? subtitleProviderFormFromProvider(provider) : emptySubtitleProviderForm()
	);

	function save(event: SubmitEvent) {
		event.preventDefault();
		void onSave(form);
	}

	function updateText(
		field: 'name' | 'baseUrl' | 'username' | 'password' | 'apiKey',
		value: string
	) {
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
	aria-labelledby="opensubtitles-title"
>
	<form class="grid gap-3.5" onsubmit={save}>
		<SectionHeading title="OpenSubtitles" titleId="opensubtitles-title" kicker="Subtitle provider">
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
					id="opensubtitles-enabled"
					checked={form.enabled}
					onCheckedChange={(checked) => updateEnabled(checked === true)}
				/>
				<Label for="opensubtitles-enabled">Enabled</Label>
			</div>
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
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">Name</span>
				<Input
					value={form.name}
					required
					maxlength={120}
					oninput={(event) => updateText('name', event.currentTarget.value)}
				/>
			</label>
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">Username</span>
				<Input
					value={form.username ?? ''}
					autocomplete="username"
					oninput={(event) => updateText('username', event.currentTarget.value)}
				/>
			</label>
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
				<SecretInput
					value={form.apiKey ?? ''}
					autocomplete="off"
					onValueChange={(value) => updateText('apiKey', value)}
				/>
			</label>
			<label class="grid gap-1.5">
				<span class="text-sm font-bold text-muted-foreground">Password</span>
				<SecretInput
					value={form.password ?? ''}
					autocomplete="current-password"
					onValueChange={(value) => updateText('password', value)}
				/>
			</label>
		</div>

		<div class="flex flex-wrap justify-end gap-2.5">
			{#if form.id}
				<ConfirmActionButton
					label={`Delete ${form.name}`}
					title="Delete subtitle provider"
					description={`Delete subtitle provider "${form.name}"?`}
					confirmLabel="Delete provider"
					onConfirm={() => {
						if (form.id) {
							return onDelete(form.id);
						}
					}}
				>
					Delete
				</ConfirmActionButton>
			{/if}
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
