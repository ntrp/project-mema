import type {
	LibraryScanImportRequest,
	LibraryScanItem,
	LibraryScanItemMatchRequest,
	MediaMonitorMode,
	MediaSearchResult,
	MinimumAvailability,
	SeriesType
} from '$lib/settings/types';
import { duplicateDraftStatesForRows } from './libraryScanDuplicates';
import { effectiveDraftValue, isSeriesKind, type MatchDraft } from './libraryScanImport';

export interface ImportBulkOptions {
	qualityProfileId: string;
	monitorMode: MediaMonitorMode;
	minimumAvailability: MinimumAvailability;
	seriesType: SeriesType;
}

export function importRequestForDraft(
	draft: MatchDraft,
	match: MediaSearchResult,
	bulk: ImportBulkOptions
): LibraryScanItemMatchRequest {
	const series = isSeriesKind(draft.mediaKind);
	const monitorMode = effectiveDraftValue(draft.monitorMode, bulk.monitorMode);
	return {
		mediaKind: draft.mediaKind,
		title: match.title,
		year: match.year,
		monitored: monitorMode !== 'none',
		qualityProfileId: effectiveDraftValue(draft.qualityProfileId, bulk.qualityProfileId),
		monitorMode,
		minimumAvailability: series
			? 'released'
			: effectiveDraftValue(draft.minimumAvailability, bulk.minimumAvailability),
		seriesType: series ? effectiveDraftValue(draft.seriesType, bulk.seriesType) : undefined,
		metadataProviderId: draft.metadataProviderId || undefined,
		mediaItemId: match.id,
		externalProvider: match.externalProvider,
		externalId: match.externalId,
		overview: match.overview,
		posterPath: match.posterPath
	};
}

export function importPayloadForRows(
	importRows: LibraryScanItem[],
	allRows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	bulk: ImportBulkOptions
): LibraryScanImportRequest {
	return {
		items: importRows.map((item) => ({
			itemId: item.id,
			match: importRequestForDraft(drafts[item.id], drafts[item.id].matched!, bulk)
		})),
		removeDuplicatePaths: duplicateRemovalPathsForRows(allRows, drafts)
	};
}

export function importPayloadForSingleRow(
	item: LibraryScanItem,
	drafts: Record<string, MatchDraft>,
	bulk: ImportBulkOptions,
	removeDuplicatePaths: string[] = []
): LibraryScanImportRequest {
	return {
		items: [
			{
				itemId: item.id,
				match: importRequestForDraft(drafts[item.id], drafts[item.id].matched!, bulk)
			}
		],
		removeDuplicatePaths
	};
}

export function duplicateRemovalPathsForRows(
	allRows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>
) {
	const duplicateStates = duplicateDraftStatesForRows(allRows, drafts);
	return Object.entries(drafts)
		.filter(
			([id, draft]) => draft.removeDuplicate && (!draft.matched || duplicateStates[id]?.duplicate)
		)
		.map(([id]) => allRows.find((row) => row.id === id)?.path)
		.filter((path): path is string => Boolean(path));
}
