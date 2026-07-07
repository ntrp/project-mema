import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import {
	audioTargetMismatchDetails,
	audioTrackMatchesTarget
} from '$lib/components/app/media/files/mediaFileAudioTargetMatching';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';

type MediaFileAudioTrack = MediaFileRow['tracks'][number];
type MediaFileAudioTarget = MediaFileRow['expectedAudioTargets'][number];

interface AudioTargetResult {
	matched: boolean;
	missingLanguage: boolean;
	details: string[];
}

export function audioSatisfaction(row: MediaFileRow) {
	if (!row.exists) return missingStatus(['File is missing']);
	const audioTracks = row.tracks.filter((track) => track.type === 'audio');
	if (audioTracks.length === 0) return missingStatus(['Audio track is missing']);
	const targetResults = audioTargets(row).map((target) => audioTargetResult(target, audioTracks));
	const targetIssues = targetResults.flatMap((result) => result.details);
	if (targetIssues.length > 0) {
		const anyMatched = targetResults.some((result) => result.matched);
		const anyWrongDetails = targetResults.some((result) => !result.missingLanguage);
		return anyMatched || anyWrongDetails
			? partialStatus(targetIssues)
			: missingStatus(targetIssues);
	}
	const unwanted = unwantedAudio(row);
	if (unwanted.length > 0) {
		return partialStatus([`Unwanted audio tracks: ${languageList(unwanted)}`]);
	}
	return okStatus([okAudioDetail(row)]);
}

function audioTargets(row: MediaFileRow): MediaFileAudioTarget[] {
	if (row.expectedAudioTargets.length > 0) return row.expectedAudioTargets;
	return row.expectedRequiredLanguages.map((languageId) => ({ languageId }));
}

function audioTargetResult(
	target: MediaFileAudioTarget,
	audioTracks: MediaFileAudioTrack[]
): AudioTargetResult {
	const candidates = audioTracks.filter((track) =>
		languageMatches(track.language, target.languageId)
	);
	if (candidates.length === 0) {
		return {
			matched: false,
			missingLanguage: true,
			details: [`Missing required audio: ${displayLanguage(target.languageId)}`]
		};
	}
	if (candidates.some((track) => audioTrackMatchesTarget(track, target))) {
		return { matched: true, missingLanguage: false, details: [] };
	}
	return {
		matched: false,
		missingLanguage: false,
		details: audioTargetMismatchDetails(candidates[0], target).map(
			(detail) => `${displayLanguage(target.languageId)} audio ${detail}`
		)
	};
}

function unwantedAudio(row: MediaFileRow) {
	if (!row.removeNonEnabledLanguages || row.expectedLanguages.length === 0) return [];
	const expected = new Set(row.expectedLanguages.map(languageMatchKey).filter(Boolean));
	return uniqueLanguages(
		row.tracks
			.filter((track) => track.type === 'audio')
			.map((track) => track.language)
			.filter((language): language is string => Boolean(language))
			.filter((language) => {
				const key = languageMatchKey(language);
				return key !== '' && !expected.has(key);
			})
	);
}

function okAudioDetail(row: MediaFileRow) {
	if (row.expectedAudioTargets.length > 0) {
		return `Audio requirements met: ${languageList(row.expectedAudioTargets.map((target) => target.languageId))}`;
	}
	if (row.expectedRequiredLanguages.length > 0) {
		return `Required audio present: ${languageList(row.expectedRequiredLanguages)}`;
	}
	return 'Audio is available';
}

function languageMatches(value: string | undefined, target: string) {
	return languageMatchKey(value) === languageMatchKey(target);
}

function okStatus(details: string[]) {
	return { state: 'satisfied' as const, label: 'Ok', details };
}

function partialStatus(details: string[]) {
	return { state: 'partial' as const, label: 'Partial', details };
}

function missingStatus(details: string[]) {
	return { state: 'missing' as const, label: 'Missing', details };
}

function languageList(values: string[]) {
	return values.map(displayLanguage).join(', ') || '-';
}

function uniqueLanguages(values: string[]) {
	return [...new Map(values.map((value) => [languageMatchKey(value) || value, value])).values()];
}
