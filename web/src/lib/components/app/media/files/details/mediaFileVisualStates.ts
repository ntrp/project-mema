import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import {
	audioTargetMismatchDetails,
	audioTrackMatchesTarget
} from '$lib/components/app/media/files/mediaFileAudioTargetMatching';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';

type MediaFileTrack = MediaFileRow['tracks'][number];
type VisualState =
	| 'matching'
	| 'partial'
	| 'unwanted'
	| 'pending_operation'
	| 'missing_placeholder';

export interface DetailVisualState {
	visualState?: VisualState;
	statusLabel?: string;
	details?: string[];
	operationLabel?: string;
}

export function trackVisualState(row: MediaFileRow, track: MediaFileTrack): DetailVisualState {
	if (track.type === 'audio') return audioVisualState(row, track);
	if (track.type === 'subtitle') return embeddedSubtitleVisualState(row, track.language);
	if (track.type === 'video') return matchingState('Video track is available.');
	return {};
}

export function externalSubtitleVisualState(
	row: MediaFileRow,
	languageId: string
): DetailVisualState {
	if (unwantedSubtitleLanguage(row, languageId)) {
		return unwantedState(`${displayLanguage(languageId)} is outside enabled subtitle targets.`);
	}
	if (!subtitleTargetMatches(row, languageId)) return {};
	if (row.subtitleSatisfaction?.mode === 'embedded') {
		return pendingState(
			'Embed subtitle',
			'External subtitle can satisfy the target after embedding.'
		);
	}
	return matchingState('External subtitle satisfies the subtitle target.');
}

export function missingPlaceholderState(
	kind: 'audio' | 'subtitle',
	languageId: string
): DetailVisualState {
	return {
		visualState: 'missing_placeholder',
		statusLabel: 'Missing',
		details: [`Missing expected ${kind}: ${displayLanguage(languageId)}`]
	};
}

function audioVisualState(row: MediaFileRow, track: MediaFileTrack): DetailVisualState {
	if (unwantedAudioLanguage(row, track.language)) {
		return unwantedState(`${displayLanguage(track.language)} is outside enabled audio targets.`);
	}
	const targets = audioTargets(row);
	if (targets.length === 0) return matchingState('Audio track is available.');
	const sameLanguage = targets.filter((target) =>
		languageMatches(track.language, target.languageId)
	);
	if (sameLanguage.some((target) => audioTrackMatchesTarget(track, target))) {
		return matchingState('Audio track satisfies a profile target.');
	}
	if (sameLanguage.length > 0) {
		return {
			visualState: 'partial',
			statusLabel: 'Partial',
			details: audioTargetMismatchDetails(track, sameLanguage[0]).map(
				(detail) => `${displayLanguage(sameLanguage[0].languageId)} audio ${detail}`
			)
		};
	}
	return {};
}

function embeddedSubtitleVisualState(
	row: MediaFileRow,
	language: string | undefined
): DetailVisualState {
	if (row.subtitleSatisfaction?.mode === 'external') {
		if (subtitleTargetMatches(row, language)) {
			return pendingState(
				'Extract subtitle',
				'Embedded subtitle can satisfy the target after extraction.'
			);
		}
		return unwantedState('Embedded subtitles conflict with external subtitle mode.');
	}
	if (unwantedSubtitleLanguage(row, language)) {
		return unwantedState(`${displayLanguage(language)} is outside enabled subtitle targets.`);
	}
	if (subtitleTargetMatches(row, language) || (row.expectedSubtitleLanguages ?? []).length === 0) {
		return matchingState('Embedded subtitle satisfies the subtitle target.');
	}
	return {};
}

function audioTargets(row: MediaFileRow) {
	const expectedAudioTargets = row.expectedAudioTargets ?? [];
	if (expectedAudioTargets.length > 0) return expectedAudioTargets;
	return (row.expectedRequiredLanguages ?? []).map((languageId) => ({ languageId }));
}

function unwantedAudioLanguage(row: MediaFileRow, language: string | undefined) {
	return unwantedLanguage(row.removeNonEnabledLanguages, row.expectedLanguages, language);
}

function unwantedSubtitleLanguage(row: MediaFileRow, language: string | undefined) {
	return unwantedLanguage(
		row.removeNonEnabledSubtitleLanguages,
		row.expectedSubtitleLanguages,
		language
	);
}

function unwantedLanguage(
	enabled: boolean,
	expected: string[] | undefined,
	language: string | undefined
) {
	const expectedLanguages = expected ?? [];
	if (!enabled || expectedLanguages.length === 0) return false;
	const wanted = new Set(expectedLanguages.map(languageMatchKey).filter(Boolean));
	const key = languageMatchKey(language);
	return key !== '' && !wanted.has(key);
}

function subtitleTargetMatches(row: MediaFileRow, language: string | undefined) {
	return (row.expectedSubtitleLanguages ?? []).some((target) => languageMatches(language, target));
}

function languageMatches(value: string | undefined, target: string) {
	return languageMatchKey(value) === languageMatchKey(target);
}

function matchingState(detail: string): DetailVisualState {
	return { visualState: 'matching', statusLabel: 'Matching', details: [detail] };
}

function pendingState(operationLabel: string, detail: string): DetailVisualState {
	return {
		visualState: 'pending_operation',
		statusLabel: 'Pending',
		operationLabel,
		details: [detail]
	};
}

function unwantedState(detail: string): DetailVisualState {
	return { visualState: 'unwanted', statusLabel: 'Unwanted', details: [detail] };
}
