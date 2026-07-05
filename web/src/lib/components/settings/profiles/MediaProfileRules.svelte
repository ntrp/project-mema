<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type { Language, MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';
	import MediaProfileDecisionSettings from './MediaProfileDecisionSettings.svelte';
	import MediaProfileLanguageSelector from './MediaProfileLanguageSelector.svelte';
	import MediaProfileSubtitleSelector from './MediaProfileSubtitleSelector.svelte';

	interface Props {
		form: MediaProfileForm;
		qualities: QualitySizeSetting[];
		languages: Language[];
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, qualities, languages, onChange }: Props = $props();

	let qualityNames = $derived(
		new Map(qualities.map((quality) => [quality.qualityId, quality.name]))
	);
	let selectedQualities = $derived(
		form.qualityIds.map((qualityId) => ({
			id: qualityId,
			name: qualityNames.get(qualityId) ?? qualityId
		}))
	);
	function patch(patchValue: Partial<MediaProfileForm>) {
		onChange({ ...form, ...patchValue });
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

		<MediaProfileDecisionSettings {form} onChange={(value) => onChange(value)} />
	</Card.Content>
</Card.Root>

<MediaProfileLanguageSelector {form} {languages} {onChange} />
<MediaProfileSubtitleSelector {form} {languages} {onChange} />
