<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import MusicIcon from '@lucide/svelte/icons/music';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import AudioConversionSelect from './audio-targets/AudioConversionSelect.svelte';
	import ProfileTargetMultiSelect from './ProfileTargetMultiSelect.svelte';
	import { defaultAudioTarget } from '$lib/settings/forms';
	import { audioChannelOptions, audioCodecOptions } from './profileTargetOptions';
	import {
		nextTargetLanguageId,
		targetLanguageChoices,
		targetLanguageDisplayLabel,
		targetLanguageKey,
		targetLanguageValue
	} from './profileTargetLanguages';
	import type { Language, MediaProfileForm } from '$lib/settings/types';

	type AudioTarget = MediaProfileForm['audioTargets'][number];

	interface Props {
		form: MediaProfileForm;
		languages: Language[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, languages, onChange }: Props = $props();
	let targets = $derived(form.audioTargets ?? []);
	let selected = $derived(
		new Set(targets.map((target) => targetLanguageKey(target.languageId, languages)))
	);
	let nextLanguageId = $derived(nextTargetLanguageId(languages, selected));
	let codecLabels = $derived(
		new Map(audioCodecOptions.map((option) => [option.value, option.label]))
	);

	function patch(audioTargets: AudioTarget[]) {
		onChange({ ...form, audioTargets });
	}
	function add() {
		if (!nextLanguageId) return;
		patch([...targets, { ...defaultAudioTarget(), languageId: nextLanguageId }]);
	}
	function update(index: number, value: Partial<AudioTarget>) {
		patch(targets.map((target, row) => (row === index ? { ...target, ...value } : target)));
	}
	function remove(index: number) {
		if (targets.length <= 1) return;
		patch(targets.filter((_, row) => row !== index));
	}
</script>

<Card.Root>
	<Card.Header class="flex flex-row items-center justify-between gap-3">
		<Card.Title class="flex items-center gap-2">
			<MusicIcon aria-hidden="true" />
			<span>Audio targets</span>
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
		<div class="grid gap-3 md:grid-cols-2">
			<Label class="flex items-center gap-2 text-sm text-muted-foreground">
				<Checkbox
					checked={form.removeUnwantedAudio}
					onCheckedChange={(checked) =>
						onChange({ ...form, removeUnwantedAudio: checked === true })}
				/>
				Remove audio tracks that are not wanted
			</Label>
			<AudioConversionSelect
				value={form.audioLossyTranscodePolicy ?? 'disabled'}
				onChange={(audioLossyTranscodePolicy) => onChange({ ...form, audioLossyTranscodePolicy })}
			/>
		</div>
		<Table.Root class="w-full table-fixed">
			<Table.Header>
				<Table.Row>
					<Table.Head class="w-60 text-left">Language</Table.Head>
					<Table.Head class="w-36 text-left">Target codec</Table.Head>
					<Table.Head class="w-36 text-left">Target channels</Table.Head>
					<Table.Head class="w-24 text-left">Minimum Bitrate</Table.Head>
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
							<Select.Root
								type="single"
								value={target.targetCodec ?? '__any'}
								onValueChange={(value) =>
									update(index, { targetCodec: value === '__any' ? undefined : value })}
							>
								<Select.Trigger class="w-full min-w-0 overflow-hidden"
									>{codecLabels.get(target.targetCodec ?? '') ?? 'Any codec'}</Select.Trigger
								>
								<Select.Content>
									<Select.Item value="__any" label="Any codec" />
									{#each audioCodecOptions as option (option.value)}
										<Select.Item value={option.value} label={option.label} />
									{/each}
								</Select.Content>
							</Select.Root>
						</Table.Cell>
						<Table.Cell>
							<ProfileTargetMultiSelect
								id={`audio-target-channels-${target.languageId}`}
								label="Target channels"
								labelClass="sr-only"
								values={target.targetChannels ?? []}
								options={audioChannelOptions}
								placeholder="Any channels"
								onChange={(values) => update(index, { targetChannels: values })}
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
						<Table.Cell class="text-right">
							<Input
								aria-label="Audio score"
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
