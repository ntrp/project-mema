<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
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
		patchScores(
			form.customFormatScores.map((item) =>
				item.customFormatId === customFormatId ? { ...item, score: value } : item
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
	<Card.Content>
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

		<div class="mt-4 grid gap-2">
			{#each form.customFormatScores as score (score.customFormatId)}
				<div
					class="grid gap-2 rounded-md border bg-muted/20 p-3 sm:grid-cols-[1fr_120px_auto] sm:items-center"
				>
					<span class="text-sm font-medium">
						{customFormatNames.get(score.customFormatId) ?? score.customFormatId}
					</span>
					<Input
						aria-label={`Score for ${customFormatNames.get(score.customFormatId) ?? score.customFormatId}`}
						type="number"
						value={score.score}
						inputmode="numeric"
						oninput={(event) =>
							updateScore(score.customFormatId, event.currentTarget.valueAsNumber)}
					/>
					<Button
						type="button"
						variant="outline"
						aria-label="Remove custom format"
						onclick={() => removeCustomFormat(score.customFormatId)}
					>
						Remove
					</Button>
				</div>
			{:else}
				<p class="m-0 text-sm text-muted-foreground">No custom formats scored for this profile</p>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
