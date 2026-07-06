<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
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
	let customFormatById = $derived(new Map(customFormats.map((format) => [format.id, format])));
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
		const score = Number.isFinite(value) ? value : 0;
		patchScores(
			form.customFormatScores.map((item) =>
				item.customFormatId === customFormatId ? { ...item, score } : item
			)
		);
	}

	function removeCustomFormat(customFormatId: string) {
		patchScores(form.customFormatScores.filter((score) => score.customFormatId !== customFormatId));
	}

	function selectedFormatLabel() {
		return customFormatNames.get(selectedCustomFormatId) ?? 'Select custom format';
	}
</script>

<Card.Root class="md:col-span-2">
	<Card.Header>
		<Card.Title>Custom formats</Card.Title>
	</Card.Header>
	<Card.Content class="grid gap-4">
		<div class="flex flex-col gap-2 sm:flex-row">
			<Select.Root type="single" bind:value={selectedCustomFormatId}>
				<Select.Trigger class="w-full sm:flex-1">{selectedFormatLabel()}</Select.Trigger>
				<Select.Content>
					{#each availableFormats as format (format.id)}
						<Select.Item value={format.id} label={format.name} />
					{/each}
				</Select.Content>
			</Select.Root>
			<Button type="button" disabled={!selectedCustomFormatId} onclick={addCustomFormat}>
				<PlusIcon aria-hidden="true" />
				<span>Add</span>
			</Button>
		</div>

		<div class="grid gap-3 [grid-template-columns:repeat(auto-fit,minmax(min(100%,340px),1fr))]">
			{#each form.customFormatScores as score (score.customFormatId)}
				{@const format = customFormatById.get(score.customFormatId)}
				<article class="grid min-h-34 gap-4 rounded-md border border-border bg-muted p-5">
					<div class="flex items-start justify-between gap-3">
						<h3 class="m-0 min-w-0 break-words text-xl font-bold text-muted-foreground">
							{format?.name ?? score.customFormatId}
						</h3>
						<Button
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label="Remove custom format"
							onclick={() => removeCustomFormat(score.customFormatId)}
						>
							<TrashIcon aria-hidden="true" />
						</Button>
					</div>

					<div class="grid gap-1 text-sm">
						<span class="font-bold text-muted-foreground">Score</span>
						<Input
							aria-label={`Score for ${format?.name ?? score.customFormatId}`}
							type="number"
							value={score.score}
							step="1"
							oninput={(event) =>
								updateScore(score.customFormatId, event.currentTarget.valueAsNumber)}
						/>
					</div>
				</article>
			{:else}
				<p class="m-0 text-sm text-muted-foreground">No custom formats scored for this profile</p>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
