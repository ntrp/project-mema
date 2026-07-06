<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import MusicIcon from '@lucide/svelte/icons/music';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import ProfileLanguageAutocomplete from './ProfileLanguageAutocomplete.svelte';
	import ProfileTargetMultiSelect from './ProfileTargetMultiSelect.svelte';
	import { defaultAudioTarget } from '$lib/settings/forms';
	import { languageLabelFromCatalog } from '$lib/settings/languageCatalog';
	import { audioChannelOptions, audioCodecOptions } from './profileTargetOptions';
	import type { Language, MediaProfileAudioTarget, MediaProfileForm } from '$lib/settings/types';
	import type { MediaProfileLossyTranscodePolicy } from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		languages: Language[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, languages, onChange }: Props = $props();
	let targets = $derived(form.audioTargets ?? []);
	let selected = $derived(new Set(targets.map((target) => target.languageId)));

	function patch(audioTargets: MediaProfileAudioTarget[]) {
		onChange({ ...form, audioTargets });
	}
	function add(languageId: string) {
		if (selected.has(languageId)) return;
		patch([...targets, { ...defaultAudioTarget(), languageId }]);
	}
	function update(index: number, value: Partial<MediaProfileAudioTarget>) {
		patch(targets.map((target, row) => (row === index ? { ...target, ...value } : target)));
	}
	function remove(index: number) {
		if (targets.length <= 1) return;
		patch(targets.filter((_, row) => row !== index));
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title class="flex items-center gap-2">
			<MusicIcon aria-hidden="true" />
			<span>Audio targets</span>
		</Card.Title>
	</Card.Header>
	<Card.Content class="mt-2 grid gap-4">
		<div class="grid gap-3">
			<Label class="flex items-center gap-2 text-sm text-muted-foreground">
				<Checkbox
					checked={form.removeUnwantedAudio}
					onCheckedChange={(checked) =>
						onChange({ ...form, removeUnwantedAudio: checked === true })}
				/>
				Remove audio tracks that are not wanted
			</Label>

			<ProfileLanguageAutocomplete
				id="audio-target-add-language"
				label="Add language"
				placeholder="Search languages"
				{languages}
				selectedIds={[...selected]}
				onSelect={add}
			/>
		</div>

		<Table.Root class="min-w-[760px] table-fixed">
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-36">Language</Table.Head>
					<Table.Head class="w-16">Score</Table.Head>
					<Table.Head>Codec</Table.Head>
					<Table.Head>Channels</Table.Head>
					<Table.Head class="w-20">Min kbps</Table.Head>
					<Table.Head class="w-32">Lossy</Table.Head>
					<Table.Head class="w-9"><span class="sr-only">Actions</span></Table.Head>
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
								aria-label="Audio score"
								type="number"
								value={target.score}
								oninput={(event) => update(index, { score: event.currentTarget.valueAsNumber })}
							/>
						</Table.Cell>
						<Table.Cell>
							<ProfileTargetMultiSelect
								id={`audio-target-codecs-${target.languageId}`}
								label="Codec"
								labelClass="sr-only"
								values={target.codecs ?? []}
								options={audioCodecOptions}
								placeholder="Any codec"
								onChange={(values) => update(index, { codecs: values })}
							/>
						</Table.Cell>
						<Table.Cell>
							<ProfileTargetMultiSelect
								id={`audio-target-channels-${target.languageId}`}
								label="Channels"
								labelClass="sr-only"
								values={target.channels ?? []}
								options={audioChannelOptions}
								placeholder="Any channels"
								onChange={(values) => update(index, { channels: values })}
							/>
						</Table.Cell>
						<Table.Cell>
							<Input
								aria-label="Minimum bitrate kbps"
								type="number"
								value={target.minimumBitrateKbps ?? ''}
								oninput={(event) =>
									update(index, { minimumBitrateKbps: event.currentTarget.valueAsNumber })}
							/>
						</Table.Cell>
						<Table.Cell>
							<Select.Root
								type="single"
								value={target.lossyTranscodePolicy}
								onValueChange={(value) =>
									update(index, {
										lossyTranscodePolicy: value as MediaProfileLossyTranscodePolicy
									})}
							>
								<Select.Trigger>{target.lossyTranscodePolicy}</Select.Trigger>
								<Select.Content>
									<Select.Item value="disabled" label="Disabled" />
									<Select.Item value="losslessToLossy" label="Lossless to lossy" />
									<Select.Item value="lossyToLossy" label="Lossy to lossy" />
								</Select.Content>
							</Select.Root>
						</Table.Cell>
						<Table.Cell class="text-right">
							<Button
								type="button"
								variant="destructive"
								size="icon-sm"
								disabled={targets.length <= 1}
								onclick={() => remove(index)}
							>
								<TrashIcon aria-label="Remove audio target" />
							</Button>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card.Content>
</Card.Root>
