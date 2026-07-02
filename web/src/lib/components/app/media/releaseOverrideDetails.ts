import { languageCodeFromValue, languageLabelFromCatalog } from '$lib/settings/languageCatalog';
import type {
	Language,
	MediaItem,
	ReleaseCandidate,
	ReleaseOverrideDetails
} from '$lib/settings/types';

export interface ReleaseOverrideDraft {
	movieTitle: string;
	seriesTitle: string;
	seasonNumber: string;
	episodeNumbers: string;
	releaseGroup: string;
	quality: string;
	languages: string[];
}

export function overrideDraftFromRelease(
	item: MediaItem,
	release: ReleaseCandidate,
	languages: Language[]
): ReleaseOverrideDraft {
	const parsed = parseSeriesReleaseTitle(release.title);
	return {
		movieTitle: release.match.matchedMedia || item.title,
		seriesTitle: release.match.matchedMedia || parsed.seriesTitle || item.title,
		seasonNumber: numberText(parsed.seasonNumber),
		episodeNumbers: parsed.episodeNumbers.map(String).join(', '),
		releaseGroup: parsed.releaseGroup,
		quality: release.match.quality || '',
		languages: release.match.languages
			.map((value) => languageCodeFromValue(value, languages))
			.filter(Boolean)
	};
}

export function detailsFromOverrideDraft(
	draft: ReleaseOverrideDraft,
	languages: Language[]
): ReleaseOverrideDetails {
	return {
		movieTitle: trimmedOrUndefined(draft.movieTitle),
		seriesTitle: trimmedOrUndefined(draft.seriesTitle),
		seasonNumber: integerOrUndefined(draft.seasonNumber),
		episodeNumbers: integerList(draft.episodeNumbers),
		releaseGroup: trimmedOrUndefined(draft.releaseGroup),
		quality: trimmedOrUndefined(draft.quality),
		languages: draft.languages.map((value) => languageLabelFromCatalog(value, languages))
	};
}

function parseSeriesReleaseTitle(title: string) {
	const releaseGroup = title.includes('-') ? title.slice(title.lastIndexOf('-') + 1).trim() : '';
	const seasonEpisode = /s(\d{1,2})((?:e\d{1,3})+)/i.exec(title);
	if (!seasonEpisode) {
		return { seriesTitle: '', seasonNumber: undefined, episodeNumbers: [], releaseGroup };
	}
	const seriesTitle = title
		.slice(0, seasonEpisode.index)
		.replace(/[._-]+/g, ' ')
		.trim();
	return {
		seriesTitle,
		seasonNumber: Number(seasonEpisode[1]),
		episodeNumbers: [...seasonEpisode[2].matchAll(/e(\d{1,3})/gi)].map((match) => Number(match[1])),
		releaseGroup
	};
}

function numberText(value: number | undefined) {
	return value === undefined || !Number.isFinite(value) ? '' : String(value);
}

function trimmedOrUndefined(value: string) {
	const trimmed = value.trim();
	return trimmed === '' ? undefined : trimmed;
}

function integerOrUndefined(value: string) {
	const parsed = Number(value);
	return Number.isInteger(parsed) ? parsed : undefined;
}

function integerList(value: string) {
	return value
		.split(/[\s,]+/)
		.map((part) => Number(part))
		.filter((part) => Number.isInteger(part) && part > 0);
}
