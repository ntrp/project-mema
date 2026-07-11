<script lang="ts">
	import BugIcon from '@lucide/svelte/icons/bug';
	import CopyIcon from '@lucide/svelte/icons/copy';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import RotateCcwIcon from '@lucide/svelte/icons/rotate-ccw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import UploadIcon from '@lucide/svelte/icons/upload';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import type { DLNARendererProfile } from '$lib/settings/types';

	interface Props {
		profiles: DLNARendererProfile[];
		search: string;
		selectedId?: string;
		onSearch: (_value: string) => void;
		onEdit: (_profile: DLNARendererProfile) => void;
		onClone: (_profile: DLNARendererProfile) => void;
		onReset: (_profile: DLNARendererProfile) => void;
		onExport: (_profile: DLNARendererProfile) => void;
		onDelete: (_profile: DLNARendererProfile) => void;
		onCreate?: () => void;
		onImport?: () => void;
		onTrace?: () => void;
		onRestoreOriginals?: () => void;
	}

	let {
		profiles,
		search,
		selectedId,
		onSearch,
		onEdit,
		onClone,
		onReset,
		onExport,
		onDelete,
		onCreate,
		onImport,
		onTrace,
		onRestoreOriginals
	}: Props = $props();
</script>

<Tooltip.Provider>
	<section class="grid gap-3" aria-label="DLNA renderer profiles">
		<div class="grid gap-2">
			<div class="grid gap-1">
				<h3 class="m-0 text-sm font-semibold">Device profiles</h3>
				<p class="m-0 text-xs text-muted-foreground">Edit seeded profiles in place, clone custom profiles, reset seeded profiles, or export and delete user profiles.</p>
			</div>
			<div class="flex flex-wrap items-center justify-between gap-3">
				<Input class="w-full sm:w-72" aria-label="Search profiles" placeholder="Search profiles" value={search} oninput={(event) => onSearch(event.currentTarget.value)} />
				<div class="flex flex-wrap items-center gap-2">
					{#if onTrace}
						<Tooltip.Root>
							<Tooltip.Trigger>
								{#snippet child({ props })}
									<Button {...props} type="button" variant="outline" size="icon-sm" aria-label="Open decision trace" onclick={onTrace}>
										<BugIcon aria-hidden="true" />
									</Button>
								{/snippet}
							</Tooltip.Trigger>
							<Tooltip.Content>Decision trace</Tooltip.Content>
						</Tooltip.Root>
					{/if}
					{#if onRestoreOriginals}
						<ConfirmActionButton label="Restore original profiles" title="Restore original profiles" description="Restore all seeded DLNA profiles to their original defaults?" confirmLabel="Restore" confirmingLabel="Restoring" variant="outline" size="icon-sm" tooltip="Restore original profiles" onConfirm={onRestoreOriginals}>
							<RotateCcwIcon aria-hidden="true" />
						</ConfirmActionButton>
					{/if}
					{#if onImport}
						<Tooltip.Root><Tooltip.Trigger>{#snippet child({ props })}<Button {...props} type="button" variant="outline" size="icon-sm" aria-label="Import" onclick={onImport}><UploadIcon aria-hidden="true" /></Button>{/snippet}</Tooltip.Trigger><Tooltip.Content>Import profile</Tooltip.Content></Tooltip.Root>
					{/if}
					{#if onCreate}
						<Tooltip.Root><Tooltip.Trigger>{#snippet child({ props })}<Button {...props} type="button" size="icon-sm" aria-label="New profile" onclick={onCreate}><PlusIcon aria-hidden="true" /></Button>{/snippet}</Tooltip.Trigger><Tooltip.Content>New profile</Tooltip.Content></Tooltip.Root>
					{/if}
				</div>
			</div>
		</div>
		<div class="max-h-[34rem] overflow-y-auto rounded-md border border-border">
			<Table.Root class="w-full">
				<Table.Header class="sticky top-0 z-10 bg-card">
					<Table.Row>
						<Table.Head>Name</Table.Head>
						<Table.Head>Family</Table.Head>
						<Table.Head>Class</Table.Head>
						<Table.Head>Enabled</Table.Head>
						<Table.Head>Priority</Table.Head>
						<Table.Head>Version</Table.Head>
						<Table.Head>Customized</Table.Head>
						<Table.Head class="text-right">Actions</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each profiles as profile (profile.id)}
						<Table.Row class={selectedId === profile.id ? 'bg-muted/50' : ''}>
							<Table.Cell><span class="font-medium">{profile.name}</span></Table.Cell>
							<Table.Cell>{profile.vendor || 'Generic'}</Table.Cell>
							<Table.Cell>{profile.deviceClass}</Table.Cell>
							<Table.Cell><Badge variant={profile.enabled ? 'default' : 'secondary'}>{profile.enabled ? 'Enabled' : 'Disabled'}</Badge></Table.Cell>
							<Table.Cell>{profile.priority}</Table.Cell>
							<Table.Cell>{profile.sourceVersion}</Table.Cell>
							<Table.Cell>{profile.customized ? 'Yes' : 'No'}</Table.Cell>
							<Table.Cell>
								<div class="flex justify-end gap-1">
									<Tooltip.Root><Tooltip.Trigger>{#snippet child({ props })}<Button {...props} type="button" variant="outline" size="icon-sm" aria-label={`Edit ${profile.name}`} onclick={() => onEdit(profile)}><PencilIcon aria-hidden="true" /></Button>{/snippet}</Tooltip.Trigger><Tooltip.Content>Edit profile</Tooltip.Content></Tooltip.Root>
									<Tooltip.Root><Tooltip.Trigger>{#snippet child({ props })}<Button {...props} type="button" variant="outline" size="icon-sm" aria-label={`Clone ${profile.name}`} onclick={() => onClone(profile)}><CopyIcon aria-hidden="true" /></Button>{/snippet}</Tooltip.Trigger><Tooltip.Content>Clone profile</Tooltip.Content></Tooltip.Root>
									<ConfirmActionButton label={`Reset ${profile.name}`} title="Reset profile" description={`Reset ${profile.name} to seeded defaults?`} confirmLabel="Reset" confirmingLabel="Resetting" variant="outline" size="icon-sm" disabled={profile.source !== 'mema_seed'} tooltip="Reset to seeded defaults" onConfirm={() => onReset(profile)}><RefreshCwIcon aria-hidden="true" /></ConfirmActionButton>
									<Tooltip.Root><Tooltip.Trigger>{#snippet child({ props })}<Button {...props} type="button" variant="outline" size="icon-sm" aria-label={`Export ${profile.name}`} onclick={() => onExport(profile)}><DownloadIcon aria-hidden="true" /></Button>{/snippet}</Tooltip.Trigger><Tooltip.Content>Export profile</Tooltip.Content></Tooltip.Root>
									{#if profile.source !== 'mema_seed'}
										<ConfirmActionButton label={`Delete ${profile.name}`} title="Delete profile" description={`Delete ${profile.name}? This cannot be undone.`} confirmLabel="Delete" confirmingLabel="Deleting" variant="destructive" size="icon-sm" tooltip="Delete profile" onConfirm={() => onDelete(profile)}><TrashIcon aria-hidden="true" /></ConfirmActionButton>
									{/if}
								</div>
							</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row><Table.Cell colspan={8} class="py-8 text-center text-muted-foreground">No profiles match search</Table.Cell></Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
	</section>
</Tooltip.Provider>
