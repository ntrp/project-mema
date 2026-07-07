<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import ReleaseSearchQueryInput from '$lib/components/app/media/release-search/ReleaseSearchQueryInput.svelte';
	import ReleaseSearchStatusLog from '$lib/components/app/media/release-search/ReleaseSearchStatusLog.svelte';
	import {
		createLogEntry,
		placeholderLogEntry,
		type ReleaseSearchLogEntry
	} from '$lib/components/app/media/release-search/releaseSearchLog';
	import { Button } from '$lib/components/ui/button';
	import SearchIcon from '@lucide/svelte/icons/search';
	import SubtitleSearchResultsTable from './SubtitleSearchResultsTable.svelte';
	import {
		subtitleSearchQuery,
		subtitleSearchQueryVariants
	} from '$lib/components/app/media/subtitle-search/subtitleSearchQuery';
	import { searchMediaSubtitles } from '$lib/settings/api';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { GrabSubtitleRequest, MediaItem, SubtitleCandidate } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		row: MediaFileRow;
		languageId: string;
		canManage: boolean;
		onGrab: (_item: MediaItem, _request: GrabSubtitleRequest) => void | Promise<void>;
		onClose: () => void;
	}

	let { item, row, languageId, canManage, onGrab, onClose }: Props = $props();

	let overrideQuery = $state(false);
	let customQuery = $state('');
	let searching = $state(false);
	let grabbingId = $state<string | undefined>();
	let candidates = $state<SubtitleCandidate[]>([]);
	let statusMessages = $state<ReleaseSearchLogEntry[]>([placeholderLogEntry()]);
	const systemQuery = $derived(subtitleSearchQuery(item));
	const queryVariants = $derived(subtitleSearchQueryVariants(item, row));
	const searchQuery = $derived(overrideQuery ? customQuery.trim() : systemQuery);

	async function submitSearch() {
		if (!row.path || !languageId) return;
		searching = true;
		candidates = [];
		statusMessages = [createLogEntry('Search started')];
		try {
			const result = await searchMediaSubtitles(item.id, {
				query: searchQuery,
				languageId,
				filePath: row.path
			});
			candidates = result.candidates;
			statusMessages = [
				...statusMessages,
				...result.logs.map(createLogEntry),
				createLogEntry(`Search finished: ${result.candidates.length} subtitles`)
			].slice(-100);
		} catch (error) {
			statusMessages = [
				...statusMessages,
				createLogEntry(error instanceof Error ? error.message : 'Subtitle search failed')
			].slice(-100);
		} finally {
			searching = false;
		}
	}

	async function grab(candidate: SubtitleCandidate) {
		if (!row.path) return;
		grabbingId = candidate.id;
		try {
			const request: GrabSubtitleRequest = {
				providerId: candidate.providerId,
				title: candidate.title,
				languageId: candidate.languageId,
				format: candidate.format,
				filePath: row.path,
				fileId: candidate.fileId,
				sourceUrl: candidate.sourceUrl,
				sourceReference: candidate.sourceReference
			};
			await onGrab(item, request);
		} finally {
			grabbingId = undefined;
		}
	}
</script>

<SettingsFormModal
	title="Manual subtitle search"
	modalClass="max-h-[calc(100vh-32px)] w-[min(1080px,calc(100vw-32px))]"
	preventScroll={false}
	{onClose}
>
	<div class="grid gap-5">
		<div class="grid gap-3 md:grid-cols-2 md:items-end">
			<ReleaseSearchQueryInput
				bind:overrideQuery
				bind:customQuery
				{queryVariants}
				disabled={!canManage || searching}
			/>
			<div class="flex justify-end">
				<Button
					type="button"
					disabled={!canManage || searching || !searchQuery || !row.path}
					onclick={submitSearch}
				>
					<SearchIcon aria-hidden="true" />
					{searching ? 'Searching' : 'Search'}
				</Button>
			</div>
		</div>
		<ReleaseSearchStatusLog messages={statusMessages} />
		<SubtitleSearchResultsTable {candidates} {searching} {grabbingId} {canManage} onGrab={grab} />
	</div>
</SettingsFormModal>
