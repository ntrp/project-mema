<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Textarea } from '$lib/components/ui/textarea';
	import type { DLNARendererDeviceOverrideRequest, DLNARendererProfile } from '$lib/settings/types';

	interface Props {
		open: boolean;
		profiles: DLNARendererProfile[];
		overrideForm: DLNARendererDeviceOverrideRequest;
		overrideJsonText: string;
		saving?: boolean;
		onOverrideJson: (_value: string) => void;
		onSave: () => unknown | Promise<unknown>;
		onClose: () => void;
	}

	let {
		open,
		profiles,
		overrideForm = $bindable(),
		overrideJsonText,
		saving = false,
		onOverrideJson,
		onSave,
		onClose
	}: Props = $props();

	async function saveOverride() {
		const saved = await onSave();
		if (saved !== false) onClose();
	}
</script>

{#if open}
	<SettingsFormModal
		title="Manual device override"
		{onClose}
		modalClass="w-[min(760px,calc(100vw-32px))]"
	>
		<form
			class="grid gap-4"
			onsubmit={(event) => {
				event.preventDefault();
				void saveOverride();
			}}
		>
			<div class="grid gap-4 sm:grid-cols-2">
				<div class="space-y-2">
					<Label for="dlna-override-ip">IP</Label>
					<Input
						id="dlna-override-ip"
						bind:value={overrideForm.ipAddress}
						placeholder="192.168.1.40"
					/>
				</div>
				<div class="space-y-2">
					<Label for="dlna-override-uuid">UUID</Label>
					<Input
						id="dlna-override-uuid"
						bind:value={overrideForm.rendererUuid}
						placeholder="uuid"
					/>
				</div>
				<div class="space-y-2">
					<Label for="dlna-override-name">Friendly name</Label>
					<Input id="dlna-override-name" bind:value={overrideForm.displayName} />
				</div>
				<div class="space-y-2">
					<Label for="dlna-override-profile">Profile</Label>
					<select
						id="dlna-override-profile"
						class="h-9 w-full rounded-md border border-input bg-background px-2 text-sm"
						bind:value={overrideForm.profileId}
						required
					>
						<option value="">Choose profile</option>
						{#each profiles as profile (profile.id)}
							<option value={profile.id}>{profile.name}</option>
						{/each}
					</select>
				</div>
				<div class="space-y-2">
					<Label for="dlna-override-notes">Notes</Label>
					<Input id="dlna-override-notes" bind:value={overrideForm.notes} />
				</div>
				<label class="flex items-center gap-3 pt-7">
					<Switch bind:checked={overrideForm.allowed} />
					<span class="text-sm font-medium">Allowed</span>
				</label>
			</div>
			<div class="space-y-2">
				<Label for="dlna-override-policy">Delivery policy overrides</Label>
				<Textarea
					id="dlna-override-policy"
					class="min-h-24 font-mono text-xs"
					value={overrideJsonText}
					oninput={(event) => onOverrideJson(event.currentTarget.value)}
					spellcheck={false}
				/>
			</div>
			<div class="flex justify-end gap-2">
				<Button type="button" variant="secondary" onclick={onClose}>Cancel</Button>
				<Button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save override'}</Button>
			</div>
		</form>
	</SettingsFormModal>
{/if}
