<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import { profileLanguageOptions } from '$lib/settings/languageCatalog';
	import type { Language, MediaProfileForm } from '$lib/settings/types';
	import MediaProfileLanguageRow from './MediaProfileLanguageRow.svelte';

	interface Props {
		form: MediaProfileForm;
		languages: Language[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, languages, onChange }: Props = $props();
	let languageFilter = $state('');
	let activeLanguageScores = $derived(
		form.targetLanguageScores?.length
			? form.targetLanguageScores
			: form.targetLanguages.map((languageId) => ({ languageId, score: 0, required: false }))
	);
	let selectedScores = $derived(
		new Map(activeLanguageScores.map((score) => [score.languageId, score]))
	);
	let options = $derived(
		profileLanguageOptions(
			languages,
			activeLanguageScores.map((score) => score.languageId)
		)
	);
	let filteredOptions = $derived(
		options
			.filter((option) =>
				option.displayLabel.toLowerCase().includes(languageFilter.trim().toLowerCase())
			)
			.toSorted((left, right) => {
				const selected = Number(selectedScores.has(right.id)) - Number(selectedScores.has(left.id));
				return selected || left.displayLabel.localeCompare(right.displayLabel);
			})
	);

	function patch(scores: MediaProfileForm['targetLanguageScores']) {
		onChange({
			...form,
			targetLanguages: scores.map((score) => score.languageId),
			targetLanguageScores: scores
		});
	}
	function toggleLanguage(language: string) {
		const nextScores = [...activeLanguageScores];
		const index = nextScores.findIndex((score) => score.languageId === language);
		if (index >= 0) nextScores.splice(index, 1);
		else nextScores.push({ languageId: language, score: 0, required: false });
		patch(nextScores);
	}
	function updateLanguageScore(language: string, score: number) {
		patch(
			activeLanguageScores.map((value) =>
				value.languageId === language
					? { ...value, score: Number.isFinite(score) ? score : 0 }
					: value
			)
		);
	}
	function updateLanguageRequired(language: string, required: boolean) {
		patch(
			activeLanguageScores.map((value) =>
				value.languageId === language ? { ...value, required } : value
			)
		);
	}
</script>

<Card.Root class="md:col-span-2">
	<Card.Header><Card.Title>Target languages</Card.Title></Card.Header>
	<Card.Content class="grid gap-4 mt-2">
		<div class="grid gap-2 text-sm">
			<span class="flex items-center gap-2 text-muted-foreground">
				<Checkbox
					checked={form.removeNonEnabledLanguages}
					onCheckedChange={(checked) =>
						onChange({ ...form, removeNonEnabledLanguages: checked === true })}
				/>
				Remove audio tracks that are not wanted
			</span>
		</div>
		<div class="grid gap-2 text-sm">
			<Label>Quick filter</Label>
			<Input bind:value={languageFilter} type="search" placeholder="Filter languages" />
		</div>
		<div class="mt-4 grid max-h-80 gap-2 overflow-auto rounded-md bg-muted/30 p-2">
			{#each filteredOptions as option (option.id)}
				<MediaProfileLanguageRow
					id={option.id}
					label={option.displayLabel}
					score={selectedScores.get(option.id)}
					onToggle={toggleLanguage}
					onScoreChange={updateLanguageScore}
					onRequiredChange={updateLanguageRequired}
				/>
			{:else}
				<p class="m-0 p-4 text-center text-sm text-muted-foreground">
					No languages match the filter.
				</p>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
