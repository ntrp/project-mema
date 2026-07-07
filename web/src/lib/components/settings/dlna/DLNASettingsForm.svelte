<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import type { DLNASettingsRequest } from '$lib/settings/types';

	interface Props {
		form: DLNASettingsRequest;
		interfacesText: string;
		allowedText: string;
		saving?: boolean;
		onSave: (event: SubmitEvent) => void | Promise<void>;
	}

	let {
		form = $bindable(),
		interfacesText = $bindable(),
		allowedText = $bindable(),
		saving = false,
		onSave
	}: Props = $props();
</script>

<form class="grid gap-4" onsubmit={onSave}>
	<div class="grid gap-4 sm:grid-cols-2">
		<label class="flex items-center gap-3">
			<Switch bind:checked={form.enabled} />
			<span class="text-sm font-medium">Enable DLNA</span>
		</label>
		<div class="space-y-2">
			<Label for="dlna-friendly-name">Friendly name</Label>
			<Input id="dlna-friendly-name" bind:value={form.friendlyName} />
		</div>
		<div class="space-y-2">
			<Label for="dlna-profile">Default renderer profile</Label>
			<Input id="dlna-profile" bind:value={form.defaultRendererProfile} />
		</div>
		<div class="space-y-2">
			<Label for="dlna-announce">Announce interval seconds</Label>
			<Input id="dlna-announce" type="number" min="60" bind:value={form.announceIntervalSeconds} />
		</div>
	</div>
	<div class="grid gap-4 sm:grid-cols-3">
		<label class="flex items-center gap-3">
			<Switch bind:checked={form.transcodeEnabled} />
			<span class="text-sm font-medium">Transcoding</span>
		</label>
		<label class="flex items-center gap-3">
			<Switch bind:checked={form.thumbnailsEnabled} />
			<span class="text-sm font-medium">Thumbnails</span>
		</label>
		<label class="flex items-center gap-3">
			<Switch bind:checked={form.subtitlesEnabled} />
			<span class="text-sm font-medium">Subtitles</span>
		</label>
	</div>
	<div class="grid gap-4 sm:grid-cols-2">
		<div class="space-y-2">
			<Label for="dlna-interfaces">Interfaces</Label>
			<textarea
				class="min-h-24 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
				id="dlna-interfaces"
				bind:value={interfacesText}></textarea>
		</div>
		<div class="space-y-2">
			<Label for="dlna-cidrs">Allowed CIDRs</Label>
			<textarea
				class="min-h-24 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
				id="dlna-cidrs"
				bind:value={allowedText}></textarea>
		</div>
	</div>
	<div class="flex justify-end">
		<Button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save settings'}</Button>
	</div>
</form>
