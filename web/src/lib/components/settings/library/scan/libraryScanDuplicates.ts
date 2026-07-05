import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';
import type { MatchDraft } from './libraryScanImport';

export interface DuplicateDraftState {
	duplicate: boolean;
	removalAllowed: boolean;
}

export function duplicateDraftStatesForRows(
	items: LibraryScanItem[],
	drafts: Record<string, MatchDraft>
) {
	const grouped: Record<string, string[]> = {};
	for (const item of items) {
		const key = duplicateValidationKey(item, drafts[item.id]);
		if (!key) continue;
		grouped[key] = [...(grouped[key] ?? []), item.id];
	}
	const states: Record<string, DuplicateDraftState> = {};
	for (const ids of Object.values(grouped)) {
		if (ids.length < 2) continue;
		for (const id of ids) {
			const item = items.find((row) => row.id === id);
			states[id] = {
				duplicate: true,
				removalAllowed: Boolean(item?.duplicateRemovalAllowed)
			};
		}
	}
	return states;
}

export function duplicateSelectionValid(
	items: LibraryScanItem[],
	drafts: Record<string, MatchDraft>
) {
	const selectedByGroup: Record<string, number> = {};
	for (const item of items) {
		const key = duplicateValidationKey(item, drafts[item.id]);
		if (!key || !drafts[item.id]?.selected) continue;
		selectedByGroup[key] = (selectedByGroup[key] ?? 0) + 1;
	}
	return Object.values(selectedByGroup).every((count) => count <= 1);
}

export function normalizeDuplicateDrafts(
	items: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	groupId?: string
) {
	if (!groupId) return;
	const groupItems = items.filter((item) => item.duplicateGroupId === groupId);
	const states = duplicateDraftStatesForRows(groupItems, drafts);
	for (const item of groupItems) {
		const draft = drafts[item.id];
		if (!draft) continue;
		if (!states[item.id]?.duplicate && draft.matched) draft.removeDuplicate = false;
		if (draft.removeDuplicate) draft.selected = false;
	}
}

function duplicateValidationKey(item: LibraryScanItem, draft?: MatchDraft) {
	if (!item.duplicateGroupId || !draft?.matched) return '';
	const matchKey = duplicateMatchKey(item, draft.matched);
	return matchKey ? `${item.duplicateGroupId}:${matchKey}` : '';
}

function duplicateMatchKey(item: LibraryScanItem, match: MediaSearchResult) {
	if (match.id) return `media:${match.id}`;
	if (!match.externalProvider || !match.externalId) return '';
	let key = `external:${match.externalProvider}:${match.externalId}`.toLowerCase();
	if (item.seasonNumber !== undefined && item.episodeNumber !== undefined) {
		key += `:s${item.seasonNumber}e${item.episodeNumber}`;
	}
	return key;
}
