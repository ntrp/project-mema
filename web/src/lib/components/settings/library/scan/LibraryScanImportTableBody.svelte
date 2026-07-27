<script lang="ts">
	import LibraryScanEpisodeRow from '$lib/components/settings/library/scan/LibraryScanEpisodeRow.svelte';
	import LibraryScanImportRow from '$lib/components/settings/library/scan/LibraryScanImportRow.svelte';
	import { isScanItemGroupRepresentative, type MatchDraft } from './libraryScanImport';
	import type { DuplicateDraftState } from './libraryScanDuplicates';
	import * as Table from '$lib/components/ui/table';
	import type {
		LibraryScanItem,
		MediaSearchResult,
		MetadataProvider,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		rows: LibraryScanItem[];
		folderPath: string;
		drafts: Record<string, MatchDraft>;
		duplicateStates: Record<string, DuplicateDraftState>;
		qualityProfiles: QualityProfileOption[];
		metadataProviders: MetadataProvider[];
		importingItemId: string;
		resettingItemId: string;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
		onProviderChange: (_item: LibraryScanItem, _providerId: string) => void;
		onSeasonChange: (_item: LibraryScanItem, _value: number | undefined) => void;
		onEpisodeChange: (_item: LibraryScanItem, _value: number | undefined) => void;
		onResetImport: (_item: LibraryScanItem) => void | Promise<void>;
	}

	let {
		rows,
		folderPath,
		drafts = $bindable(),
		duplicateStates,
		qualityProfiles,
		metadataProviders,
		importingItemId,
		resettingItemId,
		onSearch,
		onSelect,
		onProviderChange,
		onSeasonChange,
		onEpisodeChange,
		onResetImport
	}: Props = $props();
</script>

<Table.Body>
	{#each rows as item (item.id)}
		{#if drafts[item.id]}
			{#if isScanItemGroupRepresentative(item, rows)}
				<LibraryScanImportRow
					{item}
					{folderPath}
					bind:draft={drafts[item.id]}
					{qualityProfiles}
					{metadataProviders}
					duplicateState={duplicateStates[item.id]}
					importing={importingItemId === item.id}
					resetting={resettingItemId === item.id}
					{onSearch}
					{onSelect}
					{onProviderChange}
					{onResetImport}
				/>
			{/if}
			{#if item.seriesFolderPath}
				<LibraryScanEpisodeRow
					{item}
					{folderPath}
					bind:draft={drafts[item.id]}
					duplicateState={duplicateStates[item.id]}
					importing={importingItemId === item.id}
					resetting={resettingItemId === item.id}
					{onSeasonChange}
					{onEpisodeChange}
					{onResetImport}
				/>
			{/if}
		{/if}
	{:else}
		<Table.Row>
			<Table.Cell colspan={7} class="text-muted-foreground">No files to import.</Table.Cell>
		</Table.Row>
	{/each}
</Table.Body>
