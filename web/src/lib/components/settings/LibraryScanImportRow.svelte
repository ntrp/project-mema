<script lang="ts">
	import LibraryScanMatchCell from '$lib/components/settings/LibraryScanMatchCell.svelte';
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import {
		minimumAvailabilityOptions,
		monitorModeOptions,
		type LibraryScanSortMode,
		type MatchDraft
	} from '$lib/components/settings/libraryScanImport';
	import type {
		LibraryScanItem,
		MediaSearchResult,
		MediaMonitorMode,
		MinimumAvailability,
		QualityProfileOption
	} from '$lib/settings/types';

	interface Props {
		item: LibraryScanItem;
		draft: MatchDraft;
		sortMode: LibraryScanSortMode;
		qualityProfiles: QualityProfileOption[];
		folderLabel: string;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
	}

	let { item, draft, sortMode, qualityProfiles, folderLabel, onSearch, onSelect }: Props = $props();
	let qualityProfileOptions = $derived([
		{ value: '', label: 'Select profile' },
		...qualityProfiles.map((profile) => ({ value: profile.id, label: profile.name }))
	]);
</script>

<tr>
	<td class="align-top">
		<Checkbox
			bind:checked={draft.selected}
			disabled={!draft.matched || item.status !== 'pending'}
		/>
	</td>
	<td class="align-top">
		<strong class="block max-w-[280px] truncate"
			>{sortMode === 'folders' ? folderLabel : item.fileName}</strong
		>
		<span class="block max-w-[280px] truncate">{item.path}</span>
	</td>
	<td class="align-top"><LibraryScanMatchCell {item} {draft} {onSearch} {onSelect} /></td>
	<td class="align-top">
		<SettingsSelect
			value={draft.qualityProfileId}
			options={qualityProfileOptions}
			disabled={!draft.selected || !draft.matched}
			onValueChange={(value) => (draft.qualityProfileId = value)}
		/>
	</td>
	<td class="align-top">
		<SettingsSelect
			value={draft.monitorMode}
			options={monitorModeOptions}
			disabled={!draft.selected || !draft.matched}
			onValueChange={(value) => (draft.monitorMode = value as MediaMonitorMode)}
		/>
	</td>
	<td class="align-top">
		<SettingsSelect
			value={draft.minimumAvailability}
			options={minimumAvailabilityOptions}
			disabled={!draft.selected || !draft.matched}
			onValueChange={(value) => (draft.minimumAvailability = value as MinimumAvailability)}
		/>
	</td>
</tr>
