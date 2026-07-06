<script lang="ts">
	import DownloadIcon from '@lucide/svelte/icons/download';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import {
		embeddedSubtitleRows,
		externalSubtitlesForRow,
		subtitleFileLabel,
		subtitleSourceLabel,
		subtitleStateRows,
		type SubtitleStateRow
	} from './mediaSubtitles';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { MediaItemSubtitle } from '$lib/settings/types';

	interface Props {
		row: MediaFileRow;
		canManage: boolean;
		searching?: boolean;
		onSearch: (_languageId?: string) => void | Promise<void>;
		onDelete: (_subtitle: MediaItemSubtitle) => void | Promise<void>;
	}

	let { row, canManage, searching = false, onSearch, onDelete }: Props = $props();

	const wantedRows = $derived(subtitleStateRows(row, searching));
	const embeddedRows = $derived(embeddedSubtitleRows(row));
	const externalRows = $derived(externalSubtitlesForRow(row));
	const hasRows = $derived(
		wantedRows.length > 0 || embeddedRows.length > 0 || externalRows.length > 0
	);

	function statusVariant(state: SubtitleStateRow['state']) {
		return state === 'missing' ? 'destructive' : 'secondary';
	}
</script>

{#if hasRows}
	<div class="grid gap-3 border-t border-border bg-muted/20 px-4 py-3 text-sm">
		<div class="flex flex-wrap items-center justify-between gap-2">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Subtitle state</strong>
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label="Refresh subtitles"
							disabled={!canManage || searching || !row.path}
							onclick={() => onSearch()}
						>
							<RefreshCwIcon aria-hidden="true" />
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>Refresh subtitles</Tooltip.Content>
			</Tooltip.Root>
		</div>

		{#if wantedRows.length > 0}
			<div class="grid gap-1.5">
				{#each wantedRows as item (item.key)}
					<div class="grid grid-cols-[minmax(0,1fr)_auto_auto] items-center gap-2">
						<span>{item.language}</span>
						<Badge variant={statusVariant(item.state)}>{item.label}</Badge>
						<Tooltip.Root>
							<Tooltip.Trigger>
								{#snippet child({ props })}
									<Button
										{...props}
										type="button"
										variant="outline"
										size="icon-sm"
										aria-label={`Download ${item.language} subtitles`}
										disabled={!canManage || searching || !row.path}
										onclick={() => onSearch(item.languageId)}
									>
										<DownloadIcon aria-hidden="true" />
									</Button>
								{/snippet}
							</Tooltip.Trigger>
							<Tooltip.Content>Download {item.language}</Tooltip.Content>
						</Tooltip.Root>
					</div>
				{/each}
			</div>
		{/if}

		{#if embeddedRows.length > 0}
			<div class="grid gap-1.5">
				<strong class="text-xs font-medium text-muted-foreground uppercase">Embedded tracks</strong>
				{#each embeddedRows as track (track.key)}
					<div class="grid grid-cols-[minmax(0,8rem)_1fr] gap-2">
						<span>{track.language}</span>
						<span class="truncate text-muted-foreground">{track.description}</span>
					</div>
				{/each}
			</div>
		{/if}

		{#if externalRows.length > 0}
			<div class="grid gap-1.5">
				<strong class="text-xs font-medium text-muted-foreground uppercase">External files</strong>
				{#each externalRows as subtitle (subtitle.id)}
					<div class="grid grid-cols-[minmax(0,1fr)_auto] items-center gap-2">
						<span class="min-w-0">
							<span>{subtitle.languageId}</span>
							<span class="text-muted-foreground">
								{subtitle.selected ? 'Active' : 'Inactive'} · {subtitle.retentionMode} · {subtitle.format}
								· {subtitleSourceLabel(subtitle)} · {subtitleFileLabel(row, subtitle)}
							</span>
						</span>
						<ConfirmActionButton
							label={`Delete ${subtitle.languageId} subtitle`}
							title="Delete subtitle"
							description={`Delete external subtitle "${subtitle.filePath}"?`}
							confirmLabel="Delete subtitle"
							confirmingLabel="Deleting"
							size="icon-sm"
							tooltip="Delete subtitle"
							disabled={!canManage}
							onConfirm={() => onDelete(subtitle)}
						>
							<TrashIcon aria-hidden="true" />
						</ConfirmActionButton>
					</div>
				{/each}
			</div>
		{/if}
	</div>
{/if}
