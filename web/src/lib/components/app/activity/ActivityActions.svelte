<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import UploadIcon from '@lucide/svelte/icons/upload';
	import XIcon from '@lucide/svelte/icons/x';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cancellable, deletable, manualImportable } from './activityDisplay';
	import type { DownloadActivity } from '$lib/settings/types';

	interface Props {
		activity: DownloadActivity;
		canManage: boolean;
		cancellingId?: string;
		deletingId?: string;
		onManualImport: (_activity: DownloadActivity) => void;
		onCancel: (_activity: DownloadActivity) => void;
		onDelete: (_activity: DownloadActivity) => void;
	}

	let { activity, canManage, cancellingId, deletingId, onManualImport, onCancel, onDelete }: Props =
		$props();
</script>

{#if canManage && manualImportable(activity)}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="outline"
					size="icon-sm"
					aria-label={`Manual import ${activity.releaseTitle}`}
					onclick={() => onManualImport(activity)}
				>
					<UploadIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Manual import</Tooltip.Content>
	</Tooltip.Root>
{/if}
{#if canManage && cancellable(activity)}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="destructive"
					size="icon-sm"
					aria-label={`Cancel ${activity.releaseTitle}`}
					disabled={cancellingId === activity.id}
					onclick={() => onCancel(activity)}
				>
					{#if cancellingId === activity.id}
						<RefreshCwIcon aria-hidden="true" />
					{:else}
						<XIcon aria-hidden="true" />
					{/if}
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Cancel</Tooltip.Content>
	</Tooltip.Root>
{/if}
{#if canManage && deletable(activity)}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="destructive"
					size="icon-sm"
					aria-label={`Delete ${activity.releaseTitle}`}
					disabled={deletingId === activity.id}
					onclick={() => onDelete(activity)}
				>
					{#if deletingId === activity.id}
						<RefreshCwIcon aria-hidden="true" />
					{:else}
						<TrashIcon aria-hidden="true" />
					{/if}
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Delete</Tooltip.Content>
	</Tooltip.Root>
{/if}
