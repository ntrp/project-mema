<script lang="ts">
	import { cn } from '$lib/utils';
	import MediaFileDetailStateBadge from '$lib/components/app/media/files/details/MediaFileDetailStateBadge.svelte';
	import MediaFileOtherFileActions from './MediaFileOtherFileActions.svelte';
	import { unwantedMediaRowClass } from '$lib/components/app/media/files/details/mediaFileVisualClasses';
	import {
		otherFileDisplayPath,
		otherFileLanguageLabel,
		otherFileSubtypeLabel,
		otherFileTypeLabel
	} from '$lib/components/app/media/files/other-files/mediaFileOtherFiles';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type {
		MediaFulfillmentActionRequest,
		MediaItemSubtitle,
		MediaItemSubtitleSelectionRequest
	} from '$lib/settings/types';

	interface Props {
		row: MediaFileRow;
		canManage: boolean;
		pendingFulfillmentActionKeys?: string[];
		onSearch?: (_languageId?: string) => void | Promise<void>;
		onManualSearch?: (_languageId?: string) => void;
		onDeleteSubtitle?: (_subtitle: MediaItemSubtitle) => void | Promise<void>;
		onUpdateSubtitle?: (
			_subtitle: MediaItemSubtitle,
			_request: MediaItemSubtitleSelectionRequest
		) => void | Promise<void>;
		onFulfillmentAction?: (
			_row: MediaFileRow,
			_request: MediaFulfillmentActionRequest
		) => void | Promise<void>;
		onDelete: (_row: MediaFileRow) => void;
	}

	let {
		row,
		canManage,
		pendingFulfillmentActionKeys = [],
		onSearch = async () => {},
		onManualSearch = () => {},
		onDeleteSubtitle = async () => {},
		onUpdateSubtitle = async () => {},
		onFulfillmentAction = async () => {},
		onDelete
	}: Props = $props();
	const files = $derived(row.otherFiles ?? []);
	const subtitleMode = $derived(row.subtitleSatisfaction?.mode ?? 'mixed');

	function deleteFile(file: MediaFileRow['otherFiles'][number]) {
		onDelete({ ...row, path: file.path, relativePath: otherFileDisplayPath(row, file) });
	}

	function managedSubtitle(file: MediaFileRow['otherFiles'][number]) {
		return row.externalSubtitles?.find((subtitle) => subtitle.filePath === file.path);
	}

	function subtitleLanguage(file: MediaFileRow['otherFiles'][number]) {
		return file.language ?? managedSubtitle(file)?.languageId;
	}

	function subtitleVisualRow(
		file: MediaFileRow['otherFiles'][number],
		subtitle?: MediaItemSubtitle
	): MediaFileDetailRow | undefined {
		if (file.type !== 'subtitle' || !file.state) return undefined;
		const languageId = subtitleLanguage(file);
		if (!languageId) return undefined;
		return {
			key: file.path,
			filePath: row.path,
			otherFileId: file.id,
			trackNumber: '-',
			type: 'subtitle',
			language: languageId,
			description: otherFileDisplayPath(row, file),
			missing: file.state.visualState === 'missing_placeholder',
			unwanted: file.state.visualState === 'unwanted',
			...file.state
		};
	}
</script>

<div class="border-t border-border bg-card" aria-label="Other files">
	<div
		class="grid items-start gap-3 border-b border-border px-4 pt-3 pb-2 lg:grid-cols-[minmax(220px,1fr)_96px_96px_116px_176px]"
	>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Other files</strong>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Type</strong>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Subtype</strong>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Language</strong>
		<strong class="text-right text-xs font-medium text-muted-foreground uppercase">Actions</strong>
	</div>
	{#if files.length > 0}
		{#each files as file (`${file.status}:${file.type}:${file.path}`)}
			{@const subtitle = managedSubtitle(file)}
			{@const languageId = subtitleLanguage(file)}
			{@const visualRow = subtitleVisualRow(file, subtitle)}
			<div
				class={cn(
					'grid items-start gap-3 border-b border-border p-4 last:border-b-0 lg:grid-cols-[minmax(220px,1fr)_96px_96px_116px_176px]',
					file.status === 'missing' && 'text-muted-foreground italic',
					visualRow?.missing && 'bg-destructive/10 text-destructive',
					visualRow?.unwanted && unwantedMediaRowClass
				)}
			>
				<span class="break-anywhere flex min-h-8 min-w-0 items-center text-sm font-semibold">
					<span class="inline-flex items-center gap-2">
						{otherFileDisplayPath(row, file)}
						{#if visualRow}
							<MediaFileDetailStateBadge row={visualRow} />
						{/if}
					</span>
				</span>
				<span class="flex min-h-8 items-center text-sm">{otherFileTypeLabel(file.type)}</span>
				<span class="flex min-h-8 items-center text-sm">{otherFileSubtypeLabel(file)}</span>
				<span class="flex min-h-8 items-center text-sm">{otherFileLanguageLabel(file)}</span>
				<MediaFileOtherFileActions
					{file}
					{subtitle}
					{languageId}
					{subtitleMode}
					{canManage}
					{pendingFulfillmentActionKeys}
					canSearch={Boolean(row.path)}
					{onSearch}
					{onManualSearch}
					onDelete={() => (subtitle ? onDeleteSubtitle(subtitle) : deleteFile(file))}
					{onUpdateSubtitle}
					onFulfillmentAction={(request) =>
						onFulfillmentAction(row, {
							...request,
							filePath: row.path ?? request.filePath,
							otherFileId: file.id ?? request.otherFileId,
							externalSubtitleId: subtitle?.id
						})}
				/>
			</div>
		{/each}
	{:else}
		<div class="grid items-start gap-3 p-4 lg:grid-cols-[minmax(220px,1fr)_96px_96px_116px_176px]">
			<span class="flex min-h-8 min-w-0 items-center text-sm text-muted-foreground">
				No other files present.
			</span>
			<span class="flex min-h-8 items-center text-sm text-muted-foreground">-</span>
			<span class="flex min-h-8 items-center text-sm text-muted-foreground">-</span>
			<span class="flex min-h-8 items-center text-sm text-muted-foreground">-</span>
			<span class="flex min-h-8 items-center justify-end text-sm text-muted-foreground">-</span>
		</div>
	{/if}
</div>
