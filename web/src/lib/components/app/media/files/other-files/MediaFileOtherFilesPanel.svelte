<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils';
	import MediaFileOtherFileActions from './MediaFileOtherFileActions.svelte';
	import {
		otherFileDisplayPath,
		otherFileLanguageLabel,
		otherFileStatusLabel,
		otherFileTypeLabel
	} from '$lib/components/app/media/files/other-files/mediaFileOtherFiles';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { MediaItemSubtitle, MediaItemSubtitleSelectionRequest } from '$lib/settings/types';
	import { languageMatchKey } from '$lib/settings/languageDisplay';

	interface Props {
		row: MediaFileRow;
		canManage: boolean;
		onSearch?: (_languageId?: string) => void | Promise<void>;
		onManualSearch?: (_languageId?: string) => void;
		onDeleteSubtitle?: (_subtitle: MediaItemSubtitle) => void | Promise<void>;
		onUpdateSubtitle?: (
			_subtitle: MediaItemSubtitle,
			_request: MediaItemSubtitleSelectionRequest
		) => void | Promise<void>;
		onDelete: (_row: MediaFileRow) => void;
	}

	let {
		row,
		canManage,
		onSearch = async () => {},
		onManualSearch = () => {},
		onDeleteSubtitle = async () => {},
		onUpdateSubtitle = async () => {},
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

	function isUnwantedSubtitle(file: MediaFileRow['otherFiles'][number]) {
		const language = languageMatchKey(subtitleLanguage(file));
		if (file.type !== 'subtitle' || language === '' || row.expectedSubtitleLanguages.length === 0) {
			return false;
		}
		return !row.expectedSubtitleLanguages.some((value) => languageMatchKey(value) === language);
	}
</script>

<div class="border-t border-border bg-card" aria-label="Other files">
	<div
		class="grid items-start gap-3 border-b border-border px-4 pt-3 pb-2 lg:grid-cols-[minmax(220px,1fr)_96px_116px_96px_176px]"
	>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Other files</strong>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Type</strong>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Language</strong>
		<strong class="text-xs font-medium text-muted-foreground uppercase">Status</strong>
		<strong class="text-right text-xs font-medium text-muted-foreground uppercase">Actions</strong>
	</div>
	{#if files.length > 0}
		{#each files as file (`${file.status}:${file.type}:${file.path}`)}
			{@const subtitle = managedSubtitle(file)}
			{@const languageId = subtitleLanguage(file)}
			<div
				class={cn(
					'grid items-start gap-3 border-b border-border p-4 last:border-b-0 lg:grid-cols-[minmax(220px,1fr)_96px_116px_96px_176px]',
					file.status === 'missing' && 'text-muted-foreground italic',
					isUnwantedSubtitle(file) && 'bg-secondary/40'
				)}
			>
				<span class="break-anywhere flex min-h-8 min-w-0 items-center text-sm font-semibold">
					{otherFileDisplayPath(row, file)}
				</span>
				<span class="flex min-h-8 items-center text-sm">{otherFileTypeLabel(file.type)}</span>
				<span class="flex min-h-8 items-center text-sm">{otherFileLanguageLabel(file)}</span>
				<span class="flex min-h-8 items-center">
					<Badge
						variant={file.status === 'missing' ? 'destructive' : 'secondary'}
						class="justify-self-start"
					>
						{otherFileStatusLabel(file.status)}
					</Badge>
				</span>
				<MediaFileOtherFileActions
					{file}
					{subtitle}
					{languageId}
					{subtitleMode}
					{canManage}
					canSearch={Boolean(row.path)}
					{onSearch}
					{onManualSearch}
					onDelete={() => (subtitle ? onDeleteSubtitle(subtitle) : deleteFile(file))}
					{onUpdateSubtitle}
				/>
			</div>
		{/each}
	{:else}
		<div class="grid items-start gap-3 p-4 lg:grid-cols-[minmax(220px,1fr)_96px_116px_96px_176px]">
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
