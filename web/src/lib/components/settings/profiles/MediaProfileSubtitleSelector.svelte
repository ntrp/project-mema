<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import ProfileLanguageAutocomplete from './ProfileLanguageAutocomplete.svelte';
	import ProfileTargetMultiSelect from './ProfileTargetMultiSelect.svelte';
	import { defaultSubtitleTarget } from '$lib/settings/forms';
	import { languageLabelFromCatalog } from '$lib/settings/languageCatalog';
	import { subtitleFormatOptions } from './profileTargetOptions';
	import type { Language, MediaProfileForm, MediaProfileSubtitleTarget } from '$lib/settings/types';
	import type { MediaProfileSubtitleSource } from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		languages: Language[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, languages, onChange }: Props = $props();
	let targets = $derived(form.subtitleTargets ?? []);
	let selected = $derived(new Set(targets.map((target) => target.languageId)));

	function patch(subtitleTargets: MediaProfileSubtitleTarget[]) {
		onChange({ ...form, subtitleTargets });
	}
	function add(languageId: string) {
		if (selected.has(languageId)) return;
		patch([...targets, { ...defaultSubtitleTarget(), languageId }]);
	}
	function update(index: number, value: Partial<MediaProfileSubtitleTarget>) {
		patch(targets.map((target, row) => (row === index ? { ...target, ...value } : target)));
	}
	function remove(index: number) {
		patch(targets.filter((_, row) => row !== index));
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex items-center gap-2">
			<CaptionsIcon aria-hidden="true" />
			<span>Subtitle targets</span>
		</Card.Title>
	</Card.Header>
	<Card.Content class="mt-2 grid gap-4">
		<div class="grid gap-3">
			<Label class="flex items-center gap-2 text-sm text-muted-foreground">
				<Checkbox
					checked={form.removeUnwantedSubtitles}
					onCheckedChange={(checked) =>
						onChange({ ...form, removeUnwantedSubtitles: checked === true })}
				/>
				Remove subtitle tracks that are not wanted
			</Label>
			<ProfileLanguageAutocomplete
				id="subtitle-target-add-language"
				label="Add subtitle language"
				placeholder="Search subtitle languages"
				{languages}
				selectedIds={[...selected]}
				onSelect={add}
			/>
		</div>

		<Table.Root class="min-w-[680px] table-fixed">
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-40">Language</Table.Head>
					<Table.Head class="w-18">Score</Table.Head>
					<Table.Head class="w-32">Source</Table.Head>
					<Table.Head>Format</Table.Head>
					<Table.Head class="w-10"><span class="sr-only">Actions</span></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each targets as target, index (target.languageId)}
					<Table.Row>
						<Table.Cell>
							<span class="font-medium"
								>{languageLabelFromCatalog(target.languageId, languages)}</span
							>
						</Table.Cell>
						<Table.Cell>
							<Input
								aria-label="Subtitle score"
								type="number"
								value={target.score}
								oninput={(event) => update(index, { score: event.currentTarget.valueAsNumber })}
							/>
						</Table.Cell>
						<Table.Cell>
							<Select.Root
								type="single"
								value={target.source}
								onValueChange={(value) =>
									update(index, { source: value as MediaProfileSubtitleSource })}
							>
								<Select.Trigger>{target.source}</Select.Trigger>
								<Select.Content>
									<Select.Item value="any" label="Any" />
									<Select.Item value="embedded" label="Embedded" />
									<Select.Item value="external" label="External" />
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
						<Table.Cell colspan={5} class="py-6 text-center text-muted-foreground">
							No subtitle targets configured.
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>
