<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type { MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';
	import MediaProfileDecisionSettings from './MediaProfileDecisionSettings.svelte';

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
		<div class="grid gap-4 md:grid-cols-2">
			<div class="grid gap-2 text-sm">
				<Label>Name</Label>
				<Input
					value={form.name}
					type="text"
					maxlength={200}
					required
					oninput={(event) => patch({ name: event.currentTarget.value })}
				/>
			</div>
			<Label class="flex h-10 items-center gap-2 self-end text-muted-foreground">
				<Checkbox
					checked={form.isDefault}
					onCheckedChange={(checked) => patch({ isDefault: checked === true })}
				/>
				<span>Default profile</span>
			</Label>
		</div>

		<div class="grid gap-4 md:grid-cols-2">
			<div class="grid gap-2 text-sm">
				<Label>Target container</Label>
				<Select.Root
					type="single"
					value={form.finalContainer}
					onValueChange={(value: string) => patch({ finalContainer: value as 'mkv' | 'mp4' })}
				>
					<Select.Trigger class="w-full">{form.finalContainer.toUpperCase()}</Select.Trigger>
					<Select.Content>
						<Select.Item value="mkv" label="MKV" />
						<Select.Item value="mp4" label="MP4" />
					</Select.Content>
				</Select.Root>
			</div>
		</div>

		<Label class="flex items-center gap-2 text-muted-foreground">
			<Checkbox
				checked={form.upgradesAllowed}
				onCheckedChange={(checked) => patch({ upgradesAllowed: checked === true })}
			/>
			<span>Allow replacing existing releases when the candidate is better</span>
		</Label>

		<div class="grid gap-4 md:grid-cols-2">
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
			<MediaProfileDecisionSettings {form} onChange={(value) => onChange(value)} />

			<div class="grid gap-2 text-sm">
				<Label>Minimum custom format score</Label>
				<Input
					type="number"
					value={form.minimumCustomFormatScore}
					inputmode="numeric"
					oninput={(event) =>
						patch({ minimumCustomFormatScore: event.currentTarget.valueAsNumber })}
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
		</div>
	</Card.Content>
</Card.Root>
