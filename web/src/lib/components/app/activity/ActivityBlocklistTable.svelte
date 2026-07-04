<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { ReleaseBlocklistItem } from '$lib/settings/types';

	interface Props {
		items: ReleaseBlocklistItem[];
		canManage: boolean;
		deletingId?: string;
		clearing?: boolean;
		onDelete: (_item: ReleaseBlocklistItem) => void;
		onClear: () => void;
	}

	let { items, canManage, deletingId, clearing = false, onDelete, onClear }: Props = $props();
	let deleteTarget = $state<ReleaseBlocklistItem | undefined>();
	let clearOpen = $state(false);

	function protocolLabel(protocol: string) {
		if (protocol === 'usenet') return 'Usenet';
		if (protocol === 'torrent') return 'Torrent';
		return '-';
	}

	function protocolBadgeClass(protocol: string) {
		const base =
			'inline-flex min-w-18 justify-center rounded-md border px-2 py-1 text-xs font-black uppercase tracking-normal';
		if (protocol === 'usenet') return `${base} border-sky-300 bg-sky-50 text-sky-700`;
		if (protocol === 'torrent') return `${base} border-emerald-300 bg-emerald-50 text-emerald-700`;
		return `${base} border-border bg-muted text-muted-foreground`;
	}

	function displayText(value?: string) {
		const trimmed = value?.trim();
		return trimmed ? trimmed : '-';
	}

	function confirmDelete() {
		const item = deleteTarget;
		if (!item) return;
		deleteTarget = undefined;
		onDelete(item);
	}

	function confirmClear() {
		clearOpen = false;
		onClear();
	}
</script>

<Card class="overflow-hidden p-0">
	{#if canManage}
		<div class="flex justify-end border-b border-border px-3 py-2.5">
			<Button
				type="button"
				variant="destructive"
				size="sm"
				disabled={clearing}
				onclick={() => (clearOpen = true)}
			>
				{#if clearing}
					<RefreshCwIcon class="size-4 animate-spin" aria-hidden="true" />
				{:else}
					<TrashIcon class="size-4" aria-hidden="true" />
				{/if}
				<span>{clearing ? 'Clearing' : 'Clear all'}</span>
			</Button>
		</div>
	{/if}
	<Table.Root class="[&_td]:whitespace-nowrap [&_th]:whitespace-nowrap">
		<Table.Header>
			<Table.Row>
				<Table.Head>Protocol</Table.Head>
				<Table.Head class="min-w-70 whitespace-normal">Release</Table.Head>
				<Table.Head>Media</Table.Head>
				<Table.Head>Indexer</Table.Head>
				<Table.Head>Client</Table.Head>
				<Table.Head>Reason</Table.Head>
				<Table.Head>Expires</Table.Head>
				{#if canManage}
					<Table.Head class="text-right">Actions</Table.Head>
				{/if}
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each items as block (block.id)}
				<Table.Row>
					<Table.Cell>
						<span class={protocolBadgeClass(block.indexerProtocol)}>
							{protocolLabel(block.indexerProtocol)}
						</span>
					</Table.Cell>
					<Table.Cell class="max-w-120 whitespace-normal">
						<strong>{block.releaseTitle}</strong>
					</Table.Cell>
					<Table.Cell>{block.mediaTitle}</Table.Cell>
					<Table.Cell>{block.indexerName}</Table.Cell>
					<Table.Cell>{displayText(block.downloadClientName)}</Table.Cell>
					<Table.Cell class="max-w-80 whitespace-normal">{block.reason}</Table.Cell>
					<Table.Cell
						>{block.expiresAt ? new Date(block.expiresAt).toLocaleString() : '-'}</Table.Cell
					>
					{#if canManage}
						<Table.Cell class="text-right">
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="destructive"
											size="icon-sm"
											aria-label={`Remove ${block.releaseTitle} from blocklist`}
											disabled={deletingId === block.id || clearing}
											onclick={() => (deleteTarget = block)}
										>
											{#if deletingId === block.id}
												<RefreshCwIcon class="animate-spin" aria-hidden="true" />
											{:else}
												<TrashIcon aria-hidden="true" />
											{/if}
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>Remove blocklist entry</Tooltip.Content>
							</Tooltip.Root>
						</Table.Cell>
					{/if}
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>

<Dialog.Root open={!!deleteTarget} onOpenChange={(open) => !open && (deleteTarget = undefined)}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Remove blocklist entry</Dialog.Title>
			<Dialog.Description>
				Remove {deleteTarget?.releaseTitle} from the release blocklist?
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={() => (deleteTarget = undefined)}
				>Cancel</Button
			>
			<Button type="button" variant="destructive" onclick={confirmDelete}>Remove</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root open={clearOpen} onOpenChange={(open) => (clearOpen = open)}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Clear blocklist</Dialog.Title>
			<Dialog.Description>Remove every release from the blocklist?</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button
				type="button"
				variant="outline"
				disabled={clearing}
				onclick={() => (clearOpen = false)}
			>
				Cancel
			</Button>
			<Button type="button" variant="destructive" disabled={clearing} onclick={confirmClear}>
				{clearing ? 'Clearing' : 'Clear all'}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
