<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type { MediaProfileForm } from '$lib/settings/types';

	type AudioLossyTranscodePolicy = MediaProfileForm['audioLossyTranscodePolicy'];

	interface Props {
		value: AudioLossyTranscodePolicy;
		onChange: (_value: AudioLossyTranscodePolicy) => void;
	}

	let { value = 'disabled', onChange }: Props = $props();
	let labels = $derived(
		new Map<AudioLossyTranscodePolicy, string>([
			['disabled', 'Disabled'],
			['losslessToLossy', 'From lossless'],
			['lossyToLossy', 'From lossy']
		])
	);
	let label = $derived(labels.get(value) ?? 'Disabled');
</script>

<div class="grid gap-2">
	<Label for="audio-lossy-transcode-policy">Conversion</Label>
	<Select.Root
		type="single"
		{value}
		onValueChange={(selected) => onChange(selected as AudioLossyTranscodePolicy)}
	>
		<Select.Trigger id="audio-lossy-transcode-policy">
			{label}
		</Select.Trigger>
		<Select.Content>
			<Select.Item value="disabled" label="Disabled" />
			<Select.Item value="losslessToLossy" label="From lossless" />
			<Select.Item value="lossyToLossy" label="From lossy" />
		</Select.Content>
	</Select.Root>
</div>
