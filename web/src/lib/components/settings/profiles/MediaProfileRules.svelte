<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { targetLanguageOptions } from '$lib/settings/languageOptions';
	import type { MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';
	import MediaProfileLanguageRow from './MediaProfileLanguageRow.svelte';

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
			: form.targetLanguages.map((languageId) => ({ languageId, score: 0, required: false }))
	);
	let selectedLanguageScores = $derived(
		new Map(activeLanguageScores.map((score) => [score.languageId, score]))
	);
	let filteredLanguageOptions = $derived(
		targetLanguageOptions
			.filter((option) =>
				option.displayLabel.toLowerCase().includes(languageFilter.trim().toLowerCase())
			)
			.toSorted((left, right) => {
				const selected =
					Number(selectedLanguageScores.has(right.id)) -
					Number(selectedLanguageScores.has(left.id));
				return selected || left.displayLabel.localeCompare(right.displayLabel);
			})
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
		patchLanguages([...nextScores, { languageId: language, score: 0, required: false }]);
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
	function updateLanguageRequired(language: string, required: boolean) {
		patchLanguages(
			activeLanguageScores.map((value) =>
				value.languageId === language ? { ...value, required } : value
			)
		);
	}

	function patchLanguages(scores: MediaProfileForm['targetLanguageScores']) {
		patch({
			targetLanguages: scores.map((score) => score.languageId),
			targetLanguageScores: scores
		});
	}

	function upgradeUntilLabel() {
		return (
			selectedQualities.find((quality) => quality.id === form.upgradeUntilQualityId)?.name ??
			'No quality cutoff'
		);
	}
</script>

<Card.Root class="md:col-span-2">
	<Card.Header>
		<Card.Title>General</Card.Title>
	</Card.Header>
	<Card.Content class="grid gap-4">
		<div class="grid gap-2 text-sm mt-2">
			<span class="flex items-center gap-2 text-muted-foreground">
				<Checkbox
					checked={form.upgradesAllowed}
					onCheckedChange={(checked) => patch({ upgradesAllowed: checked === true })}
				/>
				Allow replacing existing releases when the candidate is better
			</span>
		</div>

		<div class="grid gap-2 text-sm">
			<Label>Upgrade until</Label>
			<Select.Root
				type="single"
				value={form.upgradeUntilQualityId ?? ''}
				disabled={!form.upgradesAllowed || selectedQualities.length === 0}
				onValueChange={(value: string) => patch({ upgradeUntilQualityId: value || undefined })}
			>
				<Select.Trigger class="w-full">{upgradeUntilLabel()}</Select.Trigger>
				<Select.Content>
					<Select.Item value="" label="No quality cutoff" />
					{#each selectedQualities as quality (quality.id)}
						<Select.Item value={quality.id} label={quality.name} />
					{/each}
				</Select.Content>
			</Select.Root>
		</div>

		<div class="grid gap-2 text-sm">
			<Label>Minimum custom format score</Label>
			<Input
				type="number"
				value={form.minimumCustomFormatScore}
				inputmode="numeric"
				oninput={(event) => patch({ minimumCustomFormatScore: event.currentTarget.valueAsNumber })}
			/>
		</div>

		<div class="grid gap-2 text-sm">
			<Label>Upgrade until custom format score</Label>
			<Input
				type="number"
				value={form.upgradeUntilCustomFormatScore}
				inputmode="numeric"
				oninput={(event) =>
					patch({ upgradeUntilCustomFormatScore: event.currentTarget.valueAsNumber })}
			/>
		</div>

		<div class="grid gap-2 text-sm">
			<Label>Minimum score increment</Label>
			<Input
				type="number"
				min="0"
				value={form.minimumCustomFormatScoreIncrement}
				inputmode="numeric"
				oninput={(event) =>
					patch({ minimumCustomFormatScoreIncrement: event.currentTarget.valueAsNumber })}
			/>
		</div>
	</Card.Content>
</Card.Root>

<Card.Root class="md:col-span-2">
	<Card.Header>
		<Card.Title>Target languages</Card.Title>
	</Card.Header>
	<Card.Content class="grid gap-4 mt-2">
		<div class="grid gap-2 text-sm">
			<span class="flex items-center gap-2 text-muted-foreground">
				<Checkbox
					checked={form.removeNonEnabledLanguages}
					onCheckedChange={(checked) => patch({ removeNonEnabledLanguages: checked === true })}
				/>
				Remove audio tracks that are not wanted
			</span>
		</div>
		<div class="grid gap-2 text-sm">
			<Label>Quick filter</Label>
			<Input bind:value={languageFilter} type="search" placeholder="Filter languages" />
		</div>
		<div class="mt-4 grid max-h-80 gap-2 overflow-auto rounded-md bg-muted/30 p-2">
			{#each filteredLanguageOptions as option (option.id)}
				<MediaProfileLanguageRow
					id={option.id}
					label={option.displayLabel}
					score={selectedLanguageScores.get(option.id)}
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
