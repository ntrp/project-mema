<script lang="ts">
	import LibraryScanMatchCell from '$lib/components/settings/LibraryScanMatchCell.svelte';
	import type { MatchDraft } from '$lib/components/settings/libraryScanImport';
	import type {
		LibraryScanItem,
		MediaSearchResult,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		item: LibraryScanItem;
		draft: MatchDraft;
		sortMode: 'folders' | 'mixed';
		qualityProfiles: QualityProfileOption[];
		folderLabel: string;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
	}

	let { item, draft, sortMode, qualityProfiles, folderLabel, onSearch, onSelect }: Props = $props();
</script>

<tr>
	<td>
		<input
			type="checkbox"
			bind:checked={draft.selected}
			disabled={!draft.matched || item.status !== 'pending'}
		/>
	</td>
	<td>
		<strong>{sortMode === 'folders' ? folderLabel : item.fileName}</strong>
		<span>{item.path}</span>
	</td>
	<td><LibraryScanMatchCell {item} {draft} {onSearch} {onSelect} /></td>
	<td>
		<select bind:value={draft.qualityProfileId} disabled={!draft.selected || !draft.matched}>
			<option value="">Select profile</option>
			{#each qualityProfiles as profile (profile.id)}
				<option value={profile.id}>{profile.name}</option>
			{/each}
		</select>
	</td>
	<td>
		<select bind:value={draft.monitorMode} disabled={!draft.selected || !draft.matched}>
			<option value="only_media">Only this media</option>
			<option value="collection">Entire collection</option>
		</select>
	</td>
	<td>
		<select bind:value={draft.minimumAvailability} disabled={!draft.selected || !draft.matched}>
			<option value="released">Released</option>
			<option value="in_cinema">In cinema</option>
			<option value="announced">Announced</option>
		</select>
	</td>
</tr>
