<script lang="ts">
	import { targetLanguageOptions } from '$lib/settings/languageOptions';
	import type { MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		qualities: QualitySizeSetting[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, qualities, onChange }: Props = $props();

	let qualityNames = $derived(
		new Map(qualities.map((quality) => [quality.qualityId, quality.name]))
	);
	let selectedQualities = $derived(
		form.qualityIds.map((qualityId) => ({
			id: qualityId,
			name: qualityNames.get(qualityId) ?? qualityId
		}))
	);
	let languageFilter = $state('');
	let activeLanguageScores = $derived(
		form.targetLanguageScores?.length
			? form.targetLanguageScores
			: form.targetLanguages.map((languageId) => ({ languageId, score: 0 }))
	);
	let selectedLanguageScores = $derived(
		new Map(activeLanguageScores.map((score) => [score.languageId, score.score]))
	);
	let filteredLanguageOptions = $derived(
		targetLanguageOptions.filter((option) =>
			option.label.toLowerCase().includes(languageFilter.trim().toLowerCase())
		)
	);

	function patch(patchValue: Partial<MediaProfileForm>) {
		onChange({ ...form, ...patchValue });
	}

	function toggleLanguage(language: string) {
		const nextScores = [...activeLanguageScores];
		const index = nextScores.findIndex((score) => score.languageId === language);
		if (index >= 0) {
			nextScores.splice(index, 1);
			patchLanguages(nextScores);
			return;
		}
		patchLanguages([...nextScores, { languageId: language, score: 0 }]);
	}

	function updateLanguageScore(language: string, score: number) {
		patchLanguages(
			activeLanguageScores.map((value) =>
				value.languageId === language
					? { ...value, score: Number.isFinite(score) ? score : 0 }
					: value
			)
		);
	}

	function patchLanguages(scores: { languageId: string; score: number }[]) {
		patch({
			targetLanguages: scores.map((score) => score.languageId),
			targetLanguageScores: scores
		});
	}
</script>

<fieldset class="profile-fieldset wide">
	<legend>General</legend>
	<label>
		<span>Upgrades allowed</span>
		<span class="toggle">
			<input
				type="checkbox"
				checked={form.upgradesAllowed}
				onchange={(event) => patch({ upgradesAllowed: event.currentTarget.checked })}
			/>
			<span>Allow replacing existing releases when the candidate is better</span>
		</span>
	</label>

	<label>
		<span>Upgrade until</span>
		<select
			value={form.upgradeUntilQualityId ?? ''}
			disabled={!form.upgradesAllowed || selectedQualities.length === 0}
			onchange={(event) => patch({ upgradeUntilQualityId: event.currentTarget.value || undefined })}
		>
			<option value="">No quality cutoff</option>
			{#each selectedQualities as quality (quality.id)}
				<option value={quality.id}>{quality.name}</option>
			{/each}
		</select>
	</label>

	<label>
		<span>Minimum custom format score</span>
		<input
			type="number"
			value={form.minimumCustomFormatScore}
			inputmode="numeric"
			oninput={(event) => patch({ minimumCustomFormatScore: event.currentTarget.valueAsNumber })}
		/>
	</label>

	<label>
		<span>Upgrade until custom format score</span>
		<input
			type="number"
			value={form.upgradeUntilCustomFormatScore}
			inputmode="numeric"
			oninput={(event) =>
				patch({ upgradeUntilCustomFormatScore: event.currentTarget.valueAsNumber })}
		/>
	</label>

	<label>
		<span>Minimum score increment</span>
		<input
			type="number"
			min="0"
			value={form.minimumCustomFormatScoreIncrement}
			inputmode="numeric"
			oninput={(event) =>
				patch({ minimumCustomFormatScoreIncrement: event.currentTarget.valueAsNumber })}
		/>
	</label>
</fieldset>

<fieldset class="profile-fieldset wide">
	<legend>Target languages</legend>
	<label class="profile-language-filter">
		<span>Quick filter</span>
		<input bind:value={languageFilter} type="search" placeholder="Filter languages" />
	</label>
	<div class="profile-language-list">
		{#each filteredLanguageOptions as option (option.id)}
			<div class="profile-language-row">
				<label class="quality-checkbox">
					<input
						type="checkbox"
						checked={selectedLanguageScores.has(option.id)}
						onchange={() => toggleLanguage(option.id)}
					/>
					<span>{option.label}</span>
				</label>
				<input
					type="number"
					aria-label={`${option.label} score`}
					value={selectedLanguageScores.get(option.id) ?? 0}
					disabled={!selectedLanguageScores.has(option.id)}
					inputmode="numeric"
					oninput={(event) => updateLanguageScore(option.id, event.currentTarget.valueAsNumber)}
				/>
			</div>
		{:else}
			<p class="empty">No languages match the filter.</p>
		{/each}
	</div>
</fieldset>
