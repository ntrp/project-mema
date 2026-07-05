<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { profileLanguageOptions } from '$lib/settings/languageCatalog';
	import type {
		Language,
		MediaProfileForm,
		MediaProfileSubtitleLanguage
	} from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		languages: Language[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, languages, onChange }: Props = $props();
	let languageFilter = $state('');
	let fallbackSubtitleType = $state<MediaProfileSubtitleLanguage['subtitleType']>('embedded');
	let selectedLanguages = $derived(
		new Map(form.subtitleLanguages.map((item) => [item.languageId, item]))
	);
	let sharedSubtitleType = $derived(form.subtitleLanguages[0]?.subtitleType ?? fallbackSubtitleType);
	let options = $derived(
		profileLanguageOptions(
			languages,
			form.subtitleLanguages.map((item) => item.languageId)
		)
	);
	let filteredOptions = $derived(
		options
			.filter((option) =>
				option.displayLabel.toLowerCase().includes(languageFilter.trim().toLowerCase())
			)
			.toSorted((left, right) => {
				const selected =
					Number(selectedLanguages.has(right.id)) - Number(selectedLanguages.has(left.id));
				return selected || left.displayLabel.localeCompare(right.displayLabel);
			})
	);

	function patch(languages: MediaProfileSubtitleLanguage[]) {
		onChange({ ...form, subtitleLanguages: languages });
	}
	function toggleLanguage(language: string) {
		const next = [...form.subtitleLanguages];
		const index = next.findIndex((item) => item.languageId === language);
		if (index >= 0) {
			if (next.length === 1) {
				fallbackSubtitleType = next[index].subtitleType;
			}
			next.splice(index, 1);
		} else {
			next.push({ languageId: language, score: 0, required: true, subtitleType: sharedSubtitleType });
		}
		patch(next);
	}
	function updateRequired(language: string, required: boolean) {
		patch(
			form.subtitleLanguages.map((item) =>
				item.languageId === language ? { ...item, required } : item
			)
		);
	}
	function updateScore(language: string, score: number) {
		patch(
			form.subtitleLanguages.map((item) =>
				item.languageId === language ? { ...item, score: Number.isFinite(score) ? score : 0 } : item
			)
		);
	}
	function updateType(
		subtitleType: MediaProfileSubtitleLanguage['subtitleType']
	) {
		fallbackSubtitleType = subtitleType;
		patch(form.subtitleLanguages.map((item) => ({ ...item, subtitleType })));
	}
	function typeLabel(value?: string) {
		if (value === 'embedded') return 'Embedded';
		if (value === 'external') return 'External';
		return 'Any';
	}
</script>

<Card.Root class="md:col-span-2">
	<Card.Header><Card.Title>Subtitle languages</Card.Title></Card.Header>
	<Card.Content class="mt-2 grid gap-4">
		<div class="grid gap-2 text-sm">
			<span class="flex items-center gap-2 text-muted-foreground">
				<Checkbox
					checked={form.removeNonEnabledSubtitleLanguages}
					onCheckedChange={(checked) =>
						onChange({ ...form, removeNonEnabledSubtitleLanguages: checked === true })}
				/>
				Remove subtitle tracks that are not wanted
			</span>
		</div>
		<div class="grid gap-2 text-sm sm:max-w-60">
			<Label>Subtitle source</Label>
			<Select.Root
				type="single"
				value={sharedSubtitleType}
				onValueChange={(value: string) =>
					updateType(value as MediaProfileSubtitleLanguage['subtitleType'])}
			>
				<Select.Trigger class="w-full">{typeLabel(sharedSubtitleType)}</Select.Trigger>
				<Select.Content>
					<Select.Item value="embedded" label="Embedded" />
					<Select.Item value="external" label="External" />
					<Select.Item value="any" label="Any" />
				</Select.Content>
			</Select.Root>
		</div>
		<div class="grid gap-2 text-sm">
			<Label>Quick filter</Label>
			<Input bind:value={languageFilter} type="search" placeholder="Filter subtitle languages" />
		</div>
		<div class="mt-4 grid max-h-80 gap-2 overflow-auto rounded-md bg-muted/30 p-2">
			{#each filteredOptions as option (option.id)}
				{@const selected = selectedLanguages.get(option.id)}
				<div
					class="grid gap-2 rounded-md bg-muted/20 p-2 sm:grid-cols-[1fr_120px_80px] sm:items-center"
				>
					<Label class="flex items-center gap-2 text-sm">
						<Checkbox
							checked={selected !== undefined}
							onCheckedChange={() => toggleLanguage(option.id)}
						/>
						<span>{option.displayLabel}</span>
					</Label>
					<Label class="flex items-center gap-2 text-sm">
						<Checkbox
							checked={selected?.required === true}
							disabled={!selected}
							onCheckedChange={(checked) => updateRequired(option.id, checked === true)}
						/>
						<span>Required</span>
					</Label>
					<Input
						class="w-20"
						type="number"
						aria-label={`${option.displayLabel} subtitle score`}
						value={selected?.score ?? 0}
						disabled={!selected}
						inputmode="numeric"
						oninput={(event) => updateScore(option.id, event.currentTarget.valueAsNumber)}
					/>
				</div>
			{:else}
				<p class="m-0 p-4 text-center text-sm text-muted-foreground">
					No languages match the filter.
				</p>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
