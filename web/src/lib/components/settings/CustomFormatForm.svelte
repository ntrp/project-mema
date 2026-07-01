<script lang="ts">
	import CustomFormatSpecEditor from '$lib/components/settings/CustomFormatSpecEditor.svelte';
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import type {
		CustomFormatForm as CustomFormatFormValue,
		CustomFormatSpec
	} from '$lib/settings/types';

	interface Props {
		form: CustomFormatFormValue;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), saving, onSave, onCancel }: Props = $props();
	let canSave = $derived(
		form.name.trim() !== '' && (form.includeSpecs.length > 0 || form.excludeSpecs.length > 0)
	);

	function addSpec(kind: 'includeSpecs' | 'excludeSpecs') {
		const spec: CustomFormatSpec = {
			id: `spec-${Date.now()}-${Math.random().toString(16).slice(2)}`,
			name: '',
			type: 'releaseTitle',
			value: '',
			required: true
		};
		form = { ...form, [kind]: [...form[kind], spec] };
	}

	function updateSpec(
		kind: 'includeSpecs' | 'excludeSpecs',
		index: number,
		patch: Partial<CustomFormatSpec>
	) {
		form = {
			...form,
			[kind]: form[kind].map((spec, specIndex) =>
				specIndex === index ? { ...spec, ...patch } : spec
			)
		};
	}

	function removeSpec(kind: 'includeSpecs' | 'excludeSpecs', index: number) {
		form = { ...form, [kind]: form[kind].filter((_, specIndex) => specIndex !== index) };
	}
</script>

<form class="grid gap-5" onsubmit={onSave}>
	<label class="grid gap-2 text-sm">
		<span class="font-medium">Name</span>
		<Input bind:value={form.name} type="text" maxlength={200} required />
	</label>

	<label class="flex items-center gap-2 text-sm text-muted-foreground">
		<Checkbox
			checked={form.includeInRenameTemplate}
			onCheckedChange={(checked) => (form = { ...form, includeInRenameTemplate: checked === true })}
		/>
		<span>Include in rename template custom_formats</span>
	</label>

	<div class="grid gap-6">
		<section class="grid gap-3">
			<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
				<strong>Conditions</strong>
				<SettingsAddButton label="Add condition" onclick={() => addSpec('includeSpecs')} />
			</div>
			<div class="grid gap-3 min-[860px]:grid-cols-2 min-[1240px]:grid-cols-3">
				{#each form.includeSpecs as spec, index (spec.id)}
					<CustomFormatSpecEditor
						{spec}
						labelPrefix="Required"
						tone="include"
						onChange={(patch) => updateSpec('includeSpecs', index, patch)}
						onRemove={() => removeSpec('includeSpecs', index)}
					/>
				{/each}
			</div>
		</section>

		<section class="grid gap-3">
			<div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
				<strong>Negated conditions</strong>
				<SettingsAddButton label="Add condition" onclick={() => addSpec('excludeSpecs')} />
			</div>
			<div class="grid gap-3 min-[860px]:grid-cols-2 min-[1240px]:grid-cols-3">
				{#each form.excludeSpecs as spec, index (spec.id)}
					<CustomFormatSpecEditor
						{spec}
						labelPrefix="Rejected"
						tone="exclude"
						onChange={(patch) => updateSpec('excludeSpecs', index, patch)}
						onRemove={() => removeSpec('excludeSpecs', index)}
					/>
				{/each}
			</div>
		</section>
	</div>

	<div class="flex justify-end gap-2">
		<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
		<Button type="submit" disabled={saving || !canSave}>
			{saving ? 'Saving' : form.id ? 'Update format' : 'Create format'}
		</Button>
	</div>
</form>
