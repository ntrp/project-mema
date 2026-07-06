<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import ProfileTargetMultiSelect from './ProfileTargetMultiSelect.svelte';
	import { defaultSubtitleTarget } from '$lib/settings/forms';
	import { subtitleFormatOptions } from './profileTargetOptions';
	import {
		nextTargetLanguageId,
		targetLanguageChoices,
		targetLanguageDisplayLabel,
		targetLanguageKey,
		targetLanguageValue
	} from './profileTargetLanguages';
	import type { Language, MediaProfileForm, MediaProfileSubtitleTarget } from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		languages: Language[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, languages, onChange }: Props = $props();
	let targets = $derived(form.subtitleTargets ?? []);
	let selected = $derived(
		new Set(targets.map((target) => targetLanguageKey(target.languageId, languages)))
	);
	let nextLanguageId = $derived(nextTargetLanguageId(languages, selected));

	function patch(subtitleTargets: MediaProfileSubtitleTarget[]) {
		onChange({ ...form, subtitleTargets });
	}
	function add() {
		if (!nextLanguageId) return;
		patch([...targets, { ...defaultSubtitleTarget(), languageId: nextLanguageId }]);
	}
	function update(index: number, value: Partial<MediaProfileSubtitleTarget>) {
		patch(targets.map((target, row) => (row === index ? { ...target, ...value } : target)));
	}
	function remove(index: number) {
		patch(targets.filter((_, row) => row !== index));
	}
	function patchForm(value: Partial<MediaProfileForm>) {
		onChange({ ...form, ...value });
	}
	function subtitleModeLabel(value: MediaProfileForm['subtitlePreferredMode'] | undefined) {
		switch (value) {
			case 'embedded':
				return 'Embedded';
			case 'external':
				return 'External';
			default:
				return 'Mixed';
		}
	}
</script>

<Card.Root>
	<Card.Header class="flex flex-row items-center justify-between gap-3">
		<Card.Title class="flex items-center gap-2">
			<CaptionsIcon aria-hidden="true" />
			<span>Subtitle targets</span>
		</Card.Title>
		<Button
			type="button"
			size="sm"
			class="bg-accent text-accent-foreground hover:bg-accent/80"
			disabled={!nextLanguageId}
			onclick={add}
		>
			<PlusIcon aria-hidden="true" />
			<span>Add</span>
		</Button>
	</Card.Header>
	<Card.Content class="mt-2 grid gap-4">
		<div class="grid gap-3">
			<div class="grid gap-2 text-sm md:max-w-xs">
				<Label>Preferred Mode</Label>
				<Select.Root
					type="single"
					value={form.subtitlePreferredMode ?? 'mixed'}
					onValueChange={(value: string) =>
						patchForm({
							subtitlePreferredMode: value as MediaProfileForm['subtitlePreferredMode']
						})}
				>
					<Select.Trigger class="w-full">
						{subtitleModeLabel(form.subtitlePreferredMode)}
					</Select.Trigger>
					<Select.Content>
						<Select.Item value="mixed" label="Mixed" />
						<Select.Item value="embedded" label="Embedded" />
						<Select.Item value="external" label="External" />
					</Select.Content>
				</Select.Root>
			</div>
			<Label class="flex items-center gap-2 text-sm text-muted-foreground">
				<Checkbox
					checked={form.removeUnwantedSubtitles}
					onCheckedChange={(checked) => patchForm({ removeUnwantedSubtitles: checked === true })}
				/>
				Remove subtitle tracks that are not wanted
			</Label>
			<Label class="flex items-center gap-2 text-sm text-muted-foreground">
				<Checkbox
					checked={form.allowSubtitleReleaseFallback}
					onCheckedChange={(checked) =>
						patchForm({ allowSubtitleReleaseFallback: checked === true })}
				/>
				Allow searching subtitles in other releases
			</Label>
		</div>

		<Table.Root class="w-full table-fixed">
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-60 text-left">Language</Table.Head>
					<Table.Head class="w-44 text-left">Format</Table.Head>
					<Table.Head class="w-20 text-right">Score</Table.Head>
					<Table.Head class="w-14 text-right"><span class="sr-only">Actions</span></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each targets as target, index (target.languageId)}
					<Table.Row>
						<Table.Cell>
							<Select.Root
								type="single"
								value={targetLanguageValue(target.languageId, languages)}
								onValueChange={(value) => update(index, { languageId: value })}
							>
								<Select.Trigger class="w-full"
									>{targetLanguageDisplayLabel(target.languageId, languages)}</Select.Trigger
								>
								<Select.Content>
									{#each targetLanguageChoices(languages, target.languageId, selected) as option (option.id)}
										<Select.Item value={option.id} label={option.displayLabel} />
									{/each}
								</Select.Content>
							</Select.Root>
						</Table.Cell>
						<Table.Cell>
							<ProfileTargetMultiSelect
								id={`subtitle-target-formats-${target.languageId}`}
								label="Format"
								labelClass="sr-only"
								values={target.formats ?? []}
								options={subtitleFormatOptions}
								placeholder="Any format"
								onChange={(values) => update(index, { formats: values })}
							/>
						</Table.Cell>
						<Table.Cell class="text-right">
							<Input
								aria-label="Subtitle score"
								class="ml-auto w-20 text-right"
								type="number"
								value={target.score}
								oninput={(event) => update(index, { score: event.currentTarget.valueAsNumber })}
							/>
						</Table.Cell>
						<Table.Cell class="pl-3 text-right">
							<Button
								type="button"
								variant="destructive"
								size="icon-sm"
								onclick={() => remove(index)}
							>
								<TrashIcon aria-label="Remove subtitle target" />
							</Button>
						</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={4} class="py-6 text-center text-muted-foreground">
							No subtitle targets configured.
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>
