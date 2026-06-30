<script lang="ts">
	import type {
		CustomFormatForm as CustomFormatFormValue,
		CustomFormatSpec,
		CustomFormatSpecType
	} from '$lib/settings/types';

	interface Props {
		form: CustomFormatFormValue;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	const specTypes: { value: CustomFormatSpecType; label: string }[] = [
		{ value: 'releaseTitle', label: 'Release title regex' },
		{ value: 'source', label: 'Source' },
		{ value: 'resolution', label: 'Resolution' },
		{ value: 'quality', label: 'Quality' },
		{ value: 'videoCodec', label: 'Video codec' },
		{ value: 'audioCodec', label: 'Audio codec' },
		{ value: 'releaseGroup', label: 'Release group regex' },
		{ value: 'edition', label: 'Edition regex' },
		{ value: 'indexerFlag', label: 'Indexer flag' },
		{ value: 'language', label: 'Language' }
	];

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

<form class="settings-form custom-format-form" onsubmit={onSave}>
	<label>
		<span>Name</span>
		<input bind:value={form.name} type="text" maxlength="200" required />
	</label>

	<div class="custom-format-spec-columns">
		<div class="custom-format-spec-list">
			<div class="profile-quality-header">
				<strong>Required</strong>
				<button type="button" class="add-action-button" onclick={() => addSpec('includeSpecs')}>
					<span class="app-icon" aria-hidden="true">add</span>
					<span>Add condition</span>
				</button>
			</div>
			{#each form.includeSpecs as spec, index (spec.id)}
				<div class="custom-format-spec-row include">
					<input
						value={spec.name}
						type="text"
						placeholder="Label"
						aria-label="Required condition label"
						oninput={(event) =>
							updateSpec('includeSpecs', index, { name: event.currentTarget.value })}
					/>
					<select
						value={spec.type}
						aria-label="Required condition type"
						onchange={(event) =>
							updateSpec('includeSpecs', index, {
								type: event.currentTarget.value as CustomFormatSpecType
							})}
					>
						{#each specTypes as type (type.value)}
							<option value={type.value}>{type.label}</option>
						{/each}
					</select>
					<input
						value={spec.value}
						type="text"
						placeholder="Value or regex"
						aria-label="Required condition value"
						oninput={(event) =>
							updateSpec('includeSpecs', index, { value: event.currentTarget.value })}
					/>
					<label class="toggle custom-format-required-toggle">
						<input
							type="checkbox"
							checked={spec.required}
							onchange={(event) =>
								updateSpec('includeSpecs', index, { required: event.currentTarget.checked })}
						/>
						<span>Required</span>
					</label>
					<button type="button" class="danger" onclick={() => removeSpec('includeSpecs', index)}>
						Remove
					</button>
				</div>
			{/each}
		</div>

		<div class="custom-format-spec-list">
			<div class="profile-quality-header">
				<strong>Rejected</strong>
				<button type="button" class="add-action-button" onclick={() => addSpec('excludeSpecs')}>
					<span class="app-icon" aria-hidden="true">add</span>
					<span>Add condition</span>
				</button>
			</div>
			{#each form.excludeSpecs as spec, index (spec.id)}
				<div class="custom-format-spec-row exclude">
					<input
						value={spec.name}
						type="text"
						placeholder="Label"
						aria-label="Rejected condition label"
						oninput={(event) =>
							updateSpec('excludeSpecs', index, { name: event.currentTarget.value })}
					/>
					<select
						value={spec.type}
						aria-label="Rejected condition type"
						onchange={(event) =>
							updateSpec('excludeSpecs', index, {
								type: event.currentTarget.value as CustomFormatSpecType
							})}
					>
						{#each specTypes as type (type.value)}
							<option value={type.value}>{type.label}</option>
						{/each}
					</select>
					<input
						value={spec.value}
						type="text"
						placeholder="Value or regex"
						aria-label="Rejected condition value"
						oninput={(event) =>
							updateSpec('excludeSpecs', index, { value: event.currentTarget.value })}
					/>
					<label class="toggle custom-format-required-toggle">
						<input
							type="checkbox"
							checked={spec.required}
							onchange={(event) =>
								updateSpec('excludeSpecs', index, { required: event.currentTarget.checked })}
						/>
						<span>Required</span>
					</label>
					<button type="button" class="danger" onclick={() => removeSpec('excludeSpecs', index)}>
						Remove
					</button>
				</div>
			{/each}
		</div>
	</div>

	<div class="form-actions">
		<button type="button" class="secondary" onclick={onCancel}>Cancel</button>
		<button type="submit" disabled={saving || !canSave}>
			{saving ? 'Saving' : form.id ? 'Update format' : 'Create format'}
		</button>
	</div>
</form>
