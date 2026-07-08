import type { MediaFileRow } from '$lib/components/app/media/file-data/mediaFileRows';

export function missingRow(
	key: string,
	title: string,
	seasonNumber?: number,
	episodeNumber?: number
): MediaFileRow {
	return {
		key,
		relativePath: '-',
		exists: false,
		seasonNumber,
		episodeNumber,
		episodeTitle: title,
		videoCodec: '-',
		audioInfo: '-',
		size: '-',
		sizeBytes: undefined,
		languages: '-',
		quality: '-',
		formats: [],
		tracks: [],
		chapters: [],
		otherFiles: [],
		externalSubtitles: [],
		subtitleSatisfaction: {
			state: 'missing',
			mode: 'mixed',
			wantedLanguages: [],
			matchedLanguages: [],
			missingLanguages: []
		},
		rollup: {
			state: 'missing',
			targetCounts: {
				missing: 0,
				partial: 0,
				pending: 0,
				satisfied: 0,
				upgradeable: 0,
				blocked: 0,
				failed: 0
			},
			reasons: ['File is missing']
		},
		upgrade: { state: 'missing', label: 'Missing', reasons: ['File is missing'] },
		expectedAudioTargets: [],
		expectedLanguages: [],
		expectedRequiredLanguages: [],
		expectedSubtitleLanguages: [],
		removeNonEnabledLanguages: false,
		removeNonEnabledSubtitleLanguages: false,
		score: 0
	};
}

export function episodeKey(season?: number, episode?: number) {
	return `${season ?? 0}:${episode ?? 0}`;
}

export function seasonNumberFromName(name: string) {
	if (name.trim().toLowerCase() === 'specials') {
		return 0;
	}
	const match = /(\d+)/.exec(name);
	return match ? Number(match[1]) : undefined;
}
