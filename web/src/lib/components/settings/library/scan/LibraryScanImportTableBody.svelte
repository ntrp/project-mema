<script lang="ts">
	import LibraryScanImportRow from '$lib/components/settings/library/scan/LibraryScanImportRow.svelte';
	import { type MatchDraft } from './libraryScanImport';
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
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
		onProviderChange: (_item: LibraryScanItem, _providerId: string) => void;
	}

	let {
		rows,
		folderPath,
		drafts = $bindable(),
		duplicateStates,
		qualityProfiles,
		metadataProviders,
		importingItemId,
		onSearch,
		onSelect,
		onProviderChange
	}: Props = $props();
</script>

<Table.Body>
	{#each rows as item (item.id)}
		{#if drafts[item.id]}
			<LibraryScanImportRow
				{item}
				{folderPath}
				bind:draft={drafts[item.id]}
				{qualityProfiles}
				{metadataProviders}
				duplicateState={duplicateStates[item.id]}
				importing={importingItemId === item.id}
				{onSearch}
				{onSelect}
				{onProviderChange}
			/>
		{/if}
	{:else}
		<Table.Row>
			<Table.Cell colspan={7} class="text-muted-foreground">No files to import.</Table.Cell>
		</Table.Row>
	{/each}
</Table.Body>
