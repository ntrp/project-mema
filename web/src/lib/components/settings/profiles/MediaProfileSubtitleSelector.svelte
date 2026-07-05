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
	let selectedLanguages = $derived(
		new Map(form.subtitleLanguages.map((item) => [item.languageId, item]))
	);
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
		if (index >= 0) next.splice(index, 1);
		else next.push({ languageId: language, required: true, subtitleType: 'any' });
		patch(next);
	}
	function updateRequired(language: string, required: boolean) {
		patch(
			form.subtitleLanguages.map((item) =>
				item.languageId === language ? { ...item, required } : item
			)
		);
	}
	function updateType(
		language: string,
		subtitleType: MediaProfileSubtitleLanguage['subtitleType']
	) {
		patch(
			form.subtitleLanguages.map((item) =>
				item.languageId === language ? { ...item, subtitleType } : item
			)
		);
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
			<Label>Quick filter</Label>
			<Input bind:value={languageFilter} type="search" placeholder="Filter subtitle languages" />
		</div>
		<div class="mt-4 grid max-h-80 gap-2 overflow-auto rounded-md bg-muted/30 p-2">
			{#each filteredOptions as option (option.id)}
				{@const selected = selectedLanguages.get(option.id)}
				<div
					class="grid gap-2 rounded-md bg-muted/20 p-2 sm:grid-cols-[1fr_120px_140px] sm:items-center"
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
					<Select.Root
						type="single"
						value={selected?.subtitleType ?? 'any'}
						disabled={!selected}
						onValueChange={(value: string) =>
							updateType(option.id, value as MediaProfileSubtitleLanguage['subtitleType'])}
					>
						<Select.Trigger class="w-full">{typeLabel(selected?.subtitleType)}</Select.Trigger>
						<Select.Content>
							<Select.Item value="any" label="Any" />
							<Select.Item value="embedded" label="Embedded" />
							<Select.Item value="external" label="External" />
						</Select.Content>
					</Select.Root>
				</div>
			{:else}
				<p class="m-0 p-4 text-center text-sm text-muted-foreground">
					No languages match the filter.
				</p>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
