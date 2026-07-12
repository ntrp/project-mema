<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import DLNADeviceOverrideModal from './DLNADeviceOverrideModal.svelte';
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
		onSave: () => unknown | Promise<unknown>;
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
	let overrideModalOpen = $state(false);

	function profileName(id: string) {
		return profiles.find((profile) => profile.id === id)?.name ?? id;
	}
	function overrideForIp(ip: string) {
		return overrides.find((override) => override.ipAddress === ip);
	}
</script>

<section class="grid gap-5" aria-label="DLNA device overrides">
	<div class="grid gap-3">
		<h2 class="m-0 font-semibold">Recent devices</h2>
		<Table.Root class="table-auto">
			<Table.Header>
				<Table.Row>
					<Table.Head class="whitespace-nowrap">IP</Table.Head>
					<Table.Head class="whitespace-nowrap">UUID</Table.Head>
					<Table.Head class="w-full">Friendly name</Table.Head>
					<Table.Head class="whitespace-nowrap">Last seen</Table.Head>
					<Table.Head class="whitespace-nowrap">Selected profile</Table.Head>
					<Table.Head class="whitespace-nowrap">Override</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each devices as device (device.ip)}
					{@const override = overrideForIp(device.ip)}
					<Table.Row>
						<Table.Cell class="whitespace-nowrap font-medium">{device.ip}</Table.Cell>
						<Table.Cell class="whitespace-nowrap text-muted-foreground">Unknown</Table.Cell>
						<Table.Cell class="w-full max-w-0">
							<span class="block truncate">{device.userAgent || 'Unknown'}</span>
						</Table.Cell>
						<Table.Cell class="whitespace-nowrap"
							>{new Date(device.lastSeen).toLocaleString()}</Table.Cell
						>
						<Table.Cell class="whitespace-nowrap">{profileName(device.profileId)}</Table.Cell>
						<Table.Cell class="whitespace-nowrap">
							<select
								class="h-9 rounded-md border border-input bg-background px-2 text-sm"
								aria-label={`Override profile for ${device.ip}`}
								value={override?.profileId ?? ''}
								onchange={(event) => void onQuickAssign(device, event.currentTarget.value)}
							>
								<option value="">No Override</option>
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
	<div class="grid gap-3">
		<div class="flex flex-wrap items-center justify-between gap-3">
			<h2 class="m-0 font-semibold">Saved overrides</h2>
			<Button type="button" size="sm" onclick={() => (overrideModalOpen = true)}>
				Add manual override
			</Button>
		</div>
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

<DLNADeviceOverrideModal
	open={overrideModalOpen}
	{profiles}
	bind:overrideForm
	{overrideJsonText}
	{saving}
	{onOverrideJson}
	onSave={async () => (await onSave()) !== false}
	onClose={() => (overrideModalOpen = false)}
/>
