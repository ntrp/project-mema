import type {
	LibraryScanItem,
	MediaMonitorMode,
	MinimumAvailability,
	SeriesType
} from '$lib/settings/types';
import { matchedDraftKind, type MatchDraft } from './libraryScanImport';

export function matchedRowsByKind(
	rows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	kind: 'movie' | 'series'
) {
	return rows.filter((item) => matchedDraftKind(drafts[item.id]) === kind);
}

export function applyQualityProfile(
	rows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	qualityProfileId: string
) {
	for (const item of rows) {
		const draft = drafts[item.id];
		if (draft) draft.qualityProfileId = qualityProfileId;
	}
}

export function setRowsSelected(
	rows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	selected: boolean
) {
	for (const item of rows) {
		const draft = drafts[item.id];
		if (draft) draft.selected = selected;
	}
}

export function applyMovieOptions(
	rows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	monitorMode: MediaMonitorMode,
	minimumAvailability: MinimumAvailability
) {
	for (const item of rows) {
		const draft = drafts[item.id];
		if (!draft) continue;
		draft.monitorMode = monitorMode;
		draft.minimumAvailability = minimumAvailability;
	}
}

export function applySeriesOptions(
	rows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	monitorMode: MediaMonitorMode,
	seriesType: SeriesType
) {
	for (const item of rows) {
		const draft = drafts[item.id];
		if (!draft) continue;
		draft.monitorMode = monitorMode;
		draft.seriesType = seriesType;
	}
}
