<script lang="ts">
	import type { CustomFormat, MediaProfileForm } from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		customFormats: CustomFormat[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, customFormats, onChange }: Props = $props();
	let selectedCustomFormatId = $state('');
	let customFormatNames = $derived(
		new Map(customFormats.map((format) => [format.id, format.name]))
	);
	let availableFormats = $derived(
		customFormats.filter(
			(format) => !form.customFormatScores.some((score) => score.customFormatId === format.id)
		)
	);

	function patchScores(scores: MediaProfileForm['customFormatScores']) {
		onChange({ ...form, customFormatScores: scores });
	}

	function addCustomFormat() {
		if (!selectedCustomFormatId) {
			return;
		}
		patchScores([...form.customFormatScores, { customFormatId: selectedCustomFormatId, score: 0 }]);
		selectedCustomFormatId = '';
	}

	function updateScore(customFormatId: string, value: number) {
		patchScores(
			form.customFormatScores.map((item) =>
				item.customFormatId === customFormatId ? { ...item, score: value } : item
			)
		);
	}

	function removeCustomFormat(customFormatId: string) {
		patchScores(form.customFormatScores.filter((score) => score.customFormatId !== customFormatId));
	}
</script>

<fieldset class="profile-fieldset wide">
	<legend>Custom formats</legend>
	<div class="profile-add-row">
		<select bind:value={selectedCustomFormatId}>
			<option value="">Select custom format</option>
			{#each availableFormats as format (format.id)}
				<option value={format.id}>{format.name}</option>
			{/each}
		</select>
		<button
			type="button"
			class="secondary"
			disabled={!selectedCustomFormatId}
			onclick={addCustomFormat}
		>
			Add
		</button>
	</div>

	<div class="profile-score-list">
		{#each form.customFormatScores as score (score.customFormatId)}
			<div class="profile-score-row">
				<span>{customFormatNames.get(score.customFormatId) ?? score.customFormatId}</span>
				<input
					aria-label={`Score for ${customFormatNames.get(score.customFormatId) ?? score.customFormatId}`}
					type="number"
					value={score.score}
					inputmode="numeric"
					oninput={(event) => updateScore(score.customFormatId, event.currentTarget.valueAsNumber)}
				/>
				<button
					type="button"
					class="secondary"
					aria-label="Remove custom format"
					onclick={() => removeCustomFormat(score.customFormatId)}
				>
					Remove
				</button>
			</div>
		{:else}
			<p class="muted">No custom formats scored for this profile</p>
		{/each}
	</div>
</fieldset>
