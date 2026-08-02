<script lang="ts">
	import LibraryScanResetImportButton from './LibraryScanResetImportButton.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import {
		isDraftSelectableForImport,
		relativeScanPath,
		type MatchDraft
	} from './libraryScanImport';
	import type { DuplicateDraftState } from './libraryScanDuplicates';
	import type { LibraryScanItem } from '$lib/settings/types';

	interface Props {
		item: LibraryScanItem;
		folderPath: string;
		draft: MatchDraft;
		duplicateState?: DuplicateDraftState;
		importing?: boolean;
		resetting?: boolean;
		onSelectionStart: (_item: LibraryScanItem, _event: PointerEvent) => void;
		onSelectionEnter: (_item: LibraryScanItem) => void;
		onSelectionClick: (_item: LibraryScanItem, _event: MouseEvent) => void;
		onSeasonChange: (_item: LibraryScanItem, _value: number | undefined) => void;
		onEpisodeChange: (_item: LibraryScanItem, _value: number | undefined) => void;
		onResetImport: (_item: LibraryScanItem) => void | Promise<void>;
	}

	let {
		item,
		folderPath,
		draft = $bindable(),
		duplicateState,
		importing = false,
		resetting = false,
		onSelectionStart,
		onSelectionEnter,
		onSelectionClick,
		onSeasonChange,
		onEpisodeChange,
		onResetImport
	}: Props = $props();
	const locked = $derived(item.imported || importing || resetting);
	const importable = $derived(isDraftSelectableForImport(item, draft));
	const canRemoveFile = $derived(
		!item.imported && item.status === 'pending' && (!draft.matched || duplicateState?.duplicate)
	);
	const displayPath = $derived(relativeScanPath(item.path, item.seriesFolderPath ?? folderPath));

	function numberValue(event: Event) {
		const value = (event.currentTarget as HTMLInputElement).valueAsNumber;
		return Number.isInteger(value) && value >= 0 ? value : undefined;
	}
	function setRemoveFile(checked: boolean) {
		draft.removeDuplicate = checked;
		if (checked) draft.selected = false;
	}
</script>

<Table.Row
	class={['h-8 bg-muted/20 text-xs [&>td]:py-1', duplicateState?.duplicate && 'bg-amber-500/5']}
>
	<Table.Cell class="w-px pl-8 align-middle">
		<Checkbox
			checked={draft.selected}
			disabled={!importable || locked}
			aria-label={`Select ${displayPath}`}
			onpointerdown={(event) => onSelectionStart(item, event)}
			onpointerenter={() => onSelectionEnter(item)}
			onclick={(event) => onSelectionClick(item, event)}
		/>
	</Table.Cell>
	<Table.Cell class="min-w-120 pl-4 align-middle">
		<div class="flex min-w-0 items-center gap-3">
			<span class="min-w-0 flex-1 truncate text-muted-foreground">{displayPath}</span>
			{#if canRemoveFile}
				<label class="flex shrink-0 items-center gap-2 text-amber-500">
					<Checkbox
						checked={draft.removeDuplicate}
						disabled={Boolean(
							importing ||
							(draft.matched && duplicateState?.duplicate && !duplicateState.removalAllowed)
						)}
						onCheckedChange={(checked) => setRemoveFile(checked === true)}
					/>
					<span>Remove file</span>
				</label>
			{/if}
			{#if item.imported}
				<Badge variant="outline" class="border-emerald-500/40 bg-emerald-500/10 text-emerald-400"
					>Imported</Badge
				>
			{:else if importing}
				<InlineSpinner label="Importing" />
			{/if}
		</div>
	</Table.Cell>
	<Table.Cell colspan={5} class="align-middle">
		<div class="flex items-center justify-end gap-2">
			<label class="flex items-center gap-1 text-muted-foreground">
				<span>Season</span>
				<Input
					type="number"
					min="0"
					step="1"
					class="h-7 w-16"
					value={draft.seasonNumber}
					disabled={!draft.matched || locked}
					onchange={(event) => onSeasonChange(item, numberValue(event))}
				/>
			</label>
			<label class="flex items-center gap-1 text-muted-foreground">
				<span>Episode</span>
				<Input
					type="number"
					min="0"
					step="1"
					class="h-7 w-16"
					value={draft.episodeNumber}
					disabled={!draft.matched || locked}
					onchange={(event) => onEpisodeChange(item, numberValue(event))}
				/>
			</label>
			{#if item.imported}
				<LibraryScanResetImportButton {item} {resetting} {onResetImport} />
			{/if}
		</div>
	</Table.Cell>
</Table.Row>
