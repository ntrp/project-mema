<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as Table from '$lib/components/ui/table';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import type {
		DLNAClientDiagnostic,
		DLNARendererDeviceOverride,
		DLNARendererDeviceOverrideRequest,
		DLNARendererProfile
	} from '$lib/settings/types';
	interface Props {
		devices: DLNAClientDiagnostic[];
		overrides: DLNARendererDeviceOverride[];
		profiles: DLNARendererProfile[];
		overrideForm: DLNARendererDeviceOverrideRequest;
		overrideJsonText: string;
		saving?: boolean;
		onOverrideJson: (_value: string) => void;
		onSave: () => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
		onQuickAssign: (_device: DLNAClientDiagnostic, _profileId: string) => void | Promise<void>;
	}
	let {
		devices,
		overrides,
		profiles,
		overrideForm = $bindable(),
		overrideJsonText,
		saving = false,
		onOverrideJson,
		onSave,
		onDelete,
		onQuickAssign
	}: Props = $props();
	function profileName(id: string) {
		return profiles.find((profile) => profile.id === id)?.name ?? id;
	}
	function overrideForIp(ip: string) {
		return overrides.find((override) => override.ipAddress === ip);
	}
</script>

<section class="grid gap-5" aria-label="DLNA device overrides">
	<div class="grid gap-3">
		<h3 class="m-0 text-sm font-semibold">Recent devices</h3>
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>IP</Table.Head>
					<Table.Head>UUID</Table.Head>
					<Table.Head>Friendly name</Table.Head>
					<Table.Head>Last seen</Table.Head>
					<Table.Head>Selected profile</Table.Head>
					<Table.Head>Override</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each devices as device (device.ip)}
					{@const override = overrideForIp(device.ip)}
					<Table.Row>
						<Table.Cell class="font-medium">{device.ip}</Table.Cell>
						<Table.Cell class="text-muted-foreground">Unknown</Table.Cell>
						<Table.Cell>{device.userAgent || 'Unknown'}</Table.Cell>
						<Table.Cell>{new Date(device.lastSeen).toLocaleString()}</Table.Cell>
						<Table.Cell>{profileName(device.profileId)}</Table.Cell>
						<Table.Cell>
							<select
								class="h-9 rounded-md border border-input bg-background px-2 text-sm"
								aria-label={`Override profile for ${device.ip}`}
								value={override?.profileId ?? ''}
								onchange={(event) => void onQuickAssign(device, event.currentTarget.value)}
							>
								<option value="">Automatic</option>
								{#each profiles as profile (profile.id)}
									<option value={profile.id}>{profile.name}</option>
								{/each}
							</select>
						</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={6} class="py-8 text-center text-muted-foreground">
							No recent DLNA devices
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
	<form
		class="grid gap-4 rounded-md border border-border p-4"
		onsubmit={(event) => {
			event.preventDefault();
			void onSave();
		}}
	>
		<h3 class="m-0 text-sm font-semibold">Manual override</h3>
		<div class="grid gap-4 sm:grid-cols-3">
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
				<Input id="dlna-override-uuid" bind:value={overrideForm.rendererUuid} placeholder="uuid" />
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
			<label class="flex items-center gap-3 pt-7">
				<Switch bind:checked={overrideForm.allowed} />
				<span class="text-sm font-medium">Allowed</span>
			</label>
			<div class="space-y-2">
				<Label for="dlna-override-notes">Notes</Label>
				<Input id="dlna-override-notes" bind:value={overrideForm.notes} />
			</div>
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
		<div class="flex justify-end">
			<Button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save override'}</Button>
		</div>
	</form>
	<div class="grid gap-3">
		<h3 class="m-0 text-sm font-semibold">Saved overrides</h3>
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Target</Table.Head>
					<Table.Head>Name</Table.Head>
					<Table.Head>Profile</Table.Head>
					<Table.Head>Allowed</Table.Head>
					<Table.Head class="text-right">Actions</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each overrides as override (override.id)}
					<Table.Row>
						<Table.Cell>{override.rendererUuid || override.ipAddress || 'Any'}</Table.Cell>
						<Table.Cell>{override.displayName || 'Unnamed'}</Table.Cell>
						<Table.Cell>{profileName(override.profileId)}</Table.Cell>
						<Table.Cell>{override.allowed ? 'Yes' : 'No'}</Table.Cell>
						<Table.Cell class="text-right">
							<ConfirmActionButton
								label={`Delete ${override.displayName || override.id}`}
								title="Delete override"
								description="Delete this DLNA renderer override?"
								confirmLabel="Delete"
								confirmingLabel="Deleting"
								size="sm"
								onConfirm={() => onDelete(override.id)}
							>
								Delete
							</ConfirmActionButton>
						</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={5} class="py-8 text-center text-muted-foreground">
							No device overrides
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
</section>
