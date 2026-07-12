<script lang="ts">
	import IntegrationTestStatus from '$lib/components/settings/shared/IntegrationTestStatus.svelte';
	import MockSubtitleRows from '$lib/components/settings/subtitles/MockSubtitleRows.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { components } from '$lib/api/generated/schema';
	import type {
		IntegrationTestResponse,
		SubtitleProviderForm as FormValue
	} from '$lib/settings/types';
	import SubtitleProviderDynamicField from './SubtitleProviderDynamicField.svelte';
	import SubtitleProviderRuntimeNotice from './SubtitleProviderRuntimeNotice.svelte';
	type CatalogEntry = components['schemas']['SubtitleProviderCatalogEntry'];

	interface Props {
		form: FormValue;
		entry?: CatalogEntry;
		saving?: boolean;
		testing?: boolean;
		testResult?: IntegrationTestResponse;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onTest: () => void | Promise<void>;
	}

	let {
		form = $bindable(),
		entry,
		saving = false,
		testing = false,
		testResult,
		onSave,
		onCancel,
		onTest
	}: Props = $props();
	const supported = $derived(
		(entry?.runtimeStatus ?? form.runtimeStatus ?? 'supported') === 'supported'
	);

	function update(patch: Partial<FormValue>) {
		form = { ...form, ...patch };
	}
</script>

<form class="grid gap-4" onsubmit={onSave}>
	<div class="flex flex-wrap items-center justify-between gap-3">
		<IntegrationTestStatus enabled={form.enabled} result={testResult} {testing} />
		<SubtitleProviderRuntimeNotice {entry} />
	</div>
	<div class="grid gap-4 md:grid-cols-2">
		<div class="flex min-h-9 items-center gap-2 self-end">
			<Checkbox
				id="subtitle-provider-enabled"
				checked={form.enabled}
				disabled={!supported}
				onCheckedChange={(checked) => update({ enabled: supported && checked === true })}
			/>
			<Label for="subtitle-provider-enabled">Enabled</Label>
		</div>
		<label class="grid gap-1.5">
			<span class="text-sm font-bold text-muted-foreground">Priority</span>
			<Input
				value={form.priority}
				min="0"
				max="1000"
				type="number"
				oninput={(event) => update({ priority: event.currentTarget.valueAsNumber || 0 })}
			/>
		</label>
		<label class="grid gap-1.5 md:col-span-2">
			<span class="text-sm font-bold text-muted-foreground">Name</span>
			<Input
				value={form.name}
				required
				maxlength={120}
				oninput={(event) => update({ name: event.currentTarget.value })}
			/>
		</label>
		{#each entry?.fields ?? [] as field (field.key)}
			<SubtitleProviderDynamicField {field} {form} onChange={(next) => (form = next)} />
		{/each}
		{#if form.type === 'mock'}
			<div class="md:col-span-2">
				<MockSubtitleRows
					rows={form.mockSubtitles ?? []}
					onChange={(rows) => update({ mockSubtitles: rows })}
				/>
			</div>
		{/if}
	</div>
	<div class="flex flex-wrap justify-end gap-2.5">
		<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
		<Button type="button" variant="outline" disabled={!supported || testing} onclick={onTest}>
			{testing ? 'Testing' : 'Test configuration'}
		</Button>
		<Button type="submit" disabled={saving || (!supported && form.enabled)}>
			{saving ? 'Saving' : 'Save'}
		</Button>
	</div>
</form>
