<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import {
		otherFileDisplayPath,
		otherFileLanguageLabel,
		otherFileStatusLabel,
		otherFileTypeLabel
	} from '$lib/components/app/media/files/other-files/mediaFileOtherFiles';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

	interface Props {
		row: MediaFileRow;
		canManage: boolean;
		onDelete: (_row: MediaFileRow) => void;
	}

	let { row, canManage, onDelete }: Props = $props();
	const files = $derived(row.otherFiles ?? []);

	function deleteFile(file: MediaFileRow['otherFiles'][number]) {
		onDelete({ ...row, path: file.path, relativePath: otherFileDisplayPath(row, file) });
	}
</script>

<div class="border-t border-border bg-card" aria-label="Other files">
	{#if files.length > 0}
		{#each files as file (`${file.status}:${file.type}:${file.path}`)}
			<div
				class={cn(
					'grid items-start gap-3 border-t border-border p-4 first:border-t-0 lg:grid-cols-[minmax(260px,1fr)_118px_118px_118px]',
					file.status === 'missing' && 'text-muted-foreground italic'
				)}
			>
				<div class="grid min-w-0 content-start gap-1">
					<strong class="text-xs font-medium text-muted-foreground uppercase">Other files</strong>
					<span class="break-anywhere flex min-h-8 items-center text-sm font-semibold">
						{otherFileDisplayPath(row, file)}
					</span>
				</div>
				<span class="grid content-start gap-1">
					<strong class="text-xs font-medium text-muted-foreground uppercase">Type</strong>
					<span class="flex min-h-8 items-center text-sm">{otherFileTypeLabel(file.type)}</span>
				</span>
				<span class="grid content-start gap-1">
					<strong class="text-xs font-medium text-muted-foreground uppercase">Language</strong>
					<span class="flex min-h-8 items-center text-sm">{otherFileLanguageLabel(file)}</span>
				</span>
				<span class="grid content-start gap-1">
					<strong class="text-xs font-medium text-muted-foreground uppercase">Status</strong>
					<span class="flex min-h-8 items-center gap-2">
						<Badge
							variant={file.status === 'missing' ? 'destructive' : 'secondary'}
							class="justify-self-start"
						>
							{otherFileStatusLabel(file.status)}
						</Badge>
						{#if file.status === 'available'}
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="destructive"
											size="icon-sm"
											aria-label="Delete other file"
											disabled={!canManage}
											onclick={() => deleteFile(file)}
										>
											<TrashIcon aria-hidden="true" />
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>Delete other file</Tooltip.Content>
							</Tooltip.Root>
						{/if}
					</span>
				</span>
			</div>
		{/each}
	{:else}
		<div class="grid items-start gap-3 p-4 lg:grid-cols-[minmax(260px,1fr)_118px_118px_118px]">
			<div class="grid min-w-0 content-start gap-1">
				<strong class="text-xs font-medium text-muted-foreground uppercase">Other files</strong>
				<span class="flex min-h-8 items-center text-sm text-muted-foreground">
					No other files present.
				</span>
			</div>
			<span class="grid content-start gap-1">
				<strong class="text-xs font-medium text-muted-foreground uppercase">Type</strong>
				<span class="flex min-h-8 items-center text-sm text-muted-foreground">-</span>
			</span>
			<span class="grid content-start gap-1">
				<strong class="text-xs font-medium text-muted-foreground uppercase">Language</strong>
				<span class="flex min-h-8 items-center text-sm text-muted-foreground">-</span>
			</span>
			<span class="grid content-start gap-1">
				<strong class="text-xs font-medium text-muted-foreground uppercase">Status</strong>
				<span class="flex min-h-8 items-center text-sm text-muted-foreground">-</span>
			</span>
		</div>
	{/if}
</div>
