<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import type {
		DLNAInterfaceDiagnostic,
		DLNARendererProfile,
		DLNASettingsRequest
	} from '$lib/settings/types';

	interface Props {
		form: DLNASettingsRequest;
		availableInterfaces?: DLNAInterfaceDiagnostic[];
		profiles?: DLNARendererProfile[];
		allowedText: string;
		saving?: boolean;
		onSave: (event: SubmitEvent) => void | Promise<void>;
	}

	let {
		form = $bindable(),
		availableInterfaces = [],
		profiles = [],
		allowedText = $bindable(),
		saving = false,
		onSave
	}: Props = $props();

	const selectedInterfaces = $derived(new Set(form.interfaces));
	const allInterfaces = $derived(form.interfaces.length === 0);

	function interfaceChecked(name: string) {
		return allInterfaces || selectedInterfaces.has(name);
	}

	function setAllInterfaces() {
		form.interfaces = [];
	}

	function toggleInterface(name: string, checked: boolean) {
		const availableNames = availableInterfaces.map((item) => item.name);
		const current = allInterfaces ? availableNames : form.interfaces;
		form.interfaces = checked
			? unique([...current, name])
			: current.filter((item) => item !== name);
	}

	function unique(items: string[]) {
		return Array.from(new Set(items));
	}
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
			<select
				id="dlna-profile"
				class="h-9 w-full rounded-md border border-input bg-background px-3 text-sm"
				bind:value={form.defaultRendererProfile}
			>
				{#each profiles as profile (profile.id)}
					<option value={profile.id}>{profile.name}</option>
				{/each}
			</select>
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
			<div class="flex items-center justify-between gap-3">
				<Label>Interfaces</Label>
				<Button type="button" variant="outline" size="sm" onclick={setAllInterfaces}>All</Button>
			</div>
			<div
				class="grid max-h-56 gap-2 overflow-y-auto rounded-md border border-input bg-background p-2"
			>
				{#each availableInterfaces as item (item.name)}
					<label class="grid grid-cols-[18px_minmax(0,1fr)] gap-2 rounded-md p-2 hover:bg-muted">
						<Checkbox
							checked={interfaceChecked(item.name)}
							onCheckedChange={(checked) => toggleInterface(item.name, checked === true)}
						/>
						<span class="grid min-w-0 gap-0.5">
							<span class="truncate text-sm font-medium">{item.name}</span>
							<span class="truncate text-xs text-muted-foreground">{item.address}</span>
						</span>
					</label>
				{:else}
					<p class="m-0 p-2 text-sm text-muted-foreground">No interfaces available</p>
				{/each}
			</div>
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
