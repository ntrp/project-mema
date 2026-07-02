<script lang="ts">
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import type { MediaProfileForm } from '$lib/settings/types';

	interface Props {
		form: MediaProfileForm;
		onChange: (_form: MediaProfileForm) => void;
	}

	let { form, onChange }: Props = $props();

	function patch(patchValue: Partial<MediaProfileForm>) {
		onChange({ ...form, ...patchValue });
	}

	function protocolLabel() {
		if (form.preferredProtocol === 'torrent') return 'Prefer torrents';
		if (form.preferredProtocol === 'usenet') return 'Prefer Usenet';
		return 'No protocol preference';
	}

	function packPreferenceLabel() {
		if (form.seriesPackPreference === 'preferPacks') return 'Prefer season packs';
		if (form.seriesPackPreference === 'preferEpisodes') return 'Prefer episodes';
		return 'Automatic';
	}
</script>

<div class="grid gap-2 text-sm">
	<Label>Preferred protocol</Label>
	<Select.Root
		type="single"
		value={form.preferredProtocol}
		onValueChange={(value: string) =>
			patch({ preferredProtocol: value as MediaProfileForm['preferredProtocol'] })}
	>
		<Select.Trigger class="w-full">{protocolLabel()}</Select.Trigger>
		<Select.Content>
			<Select.Item value="any" label="No protocol preference" />
			<Select.Item value="torrent" label="Prefer torrents" />
			<Select.Item value="usenet" label="Prefer Usenet" />
		</Select.Content>
	</Select.Root>
</div>

<div class="grid gap-2 text-sm">
	<Label>Series pack preference</Label>
	<Select.Root
		type="single"
		value={form.seriesPackPreference}
		onValueChange={(value: string) =>
			patch({ seriesPackPreference: value as MediaProfileForm['seriesPackPreference'] })}
	>
		<Select.Trigger class="w-full">{packPreferenceLabel()}</Select.Trigger>
		<Select.Content>
			<Select.Item value="auto" label="Automatic" />
			<Select.Item value="preferPacks" label="Prefer season packs" />
			<Select.Item value="preferEpisodes" label="Prefer episodes" />
		</Select.Content>
	</Select.Root>
</div>
