import type { MediaProfile, MediaProfileForm, MediaProfileRequest } from './types';

export function emptyMediaProfileForm(): MediaProfileForm {
	return {
		name: '',
		isDefault: false,
		finalContainer: 'mkv',
		qualityIds: [],
		upgradesAllowed: true,
		upgradeUntilQualityId: undefined,
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 0,
		minimumCustomFormatScoreIncrement: 1,
		removeUnwantedAudio: false,
		removeUnwantedSubtitles: false,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		videoTarget: defaultVideoTarget(),
		audioTargets: [defaultAudioTarget()],
		subtitleTargets: [],
		customFormatScores: []
	};
}

export function mediaProfileFormFromProfile(profile: MediaProfile): MediaProfileForm {
	return {
		id: profile.id,
		name: profile.name,
		isDefault: profile.isDefault,
		finalContainer: profile.finalContainer,
		qualityIds: [...(profile.qualityIds ?? [])],
		upgradesAllowed: profile.upgradesAllowed,
		upgradeUntilQualityId: profile.upgradeUntilQualityId,
		minimumCustomFormatScore: profile.minimumCustomFormatScore,
		upgradeUntilCustomFormatScore: profile.upgradeUntilCustomFormatScore,
		minimumCustomFormatScoreIncrement: profile.minimumCustomFormatScoreIncrement,
		removeUnwantedAudio: profile.removeUnwantedAudio,
		removeUnwantedSubtitles: profile.removeUnwantedSubtitles,
		preferredProtocol: profile.preferredProtocol,
		seriesPackPreference: profile.seriesPackPreference,
		videoTarget: { ...defaultVideoTarget(), ...(profile.videoTarget ?? {}) },
		audioTargets: (profile.audioTargets ?? []).map((target) => ({ ...target })),
		subtitleTargets: (profile.subtitleTargets ?? []).map((target) => ({ ...target })),
		customFormatScores: (profile.customFormatScores ?? []).map((score) => ({ ...score }))
	};
}

export function normalizeMediaProfileForm(form: MediaProfileForm): MediaProfileRequest {
	const qualityIds = uniqueTrimmed(form.qualityIds);
	const audioTargets = audioTargetsFromForm(form);
	return {
		name: form.name.trim(),
		isDefault: form.isDefault,
		finalContainer: form.finalContainer ?? 'mkv',
		qualityIds,
		upgradesAllowed: form.upgradesAllowed,
		upgradeUntilQualityId:
			form.upgradeUntilQualityId && qualityIds.includes(form.upgradeUntilQualityId)
				? form.upgradeUntilQualityId
				: undefined,
		minimumCustomFormatScore: normalizedInteger(form.minimumCustomFormatScore),
		upgradeUntilCustomFormatScore: normalizedInteger(form.upgradeUntilCustomFormatScore),
		minimumCustomFormatScoreIncrement: Math.max(
			0,
			normalizedInteger(form.minimumCustomFormatScoreIncrement)
		),
		removeUnwantedAudio: form.removeUnwantedAudio,
		removeUnwantedSubtitles: form.removeUnwantedSubtitles,
		preferredProtocol: form.preferredProtocol ?? 'any',
		seriesPackPreference: form.seriesPackPreference ?? 'auto',
		videoTarget: videoTargetFromForm(form),
		audioTargets,
		subtitleTargets: subtitleTargetsFromForm(form),
		customFormatScores: (form.customFormatScores ?? [])
			.filter((score) => score.customFormatId)
			.map((score) => ({
				customFormatId: score.customFormatId,
				score: normalizedInteger(score.score)
			}))
	};
}

export function defaultAudioTarget(): MediaProfileRequest['audioTargets'][number] {
	return {
		languageId: 'EN',
		score: 0,
		required: true,
		lossyTranscodePolicy: 'disabled'
	};
}

export function defaultSubtitleTarget(): MediaProfileRequest['subtitleTargets'][number] {
	return {
		languageId: 'EN',
		score: 0,
		required: true,
		source: 'any',
		formats: ['srt']
	};
}

function defaultVideoTarget(): MediaProfileRequest['videoTarget'] {
	return {
		codecRequired: false,
		codecScore: 0,
		hdrRequired: false,
		hdrScore: 0,
		pixelFormatRequired: false,
		pixelFormatScore: 0
	};
}

function videoTargetFromForm(form: MediaProfileForm): MediaProfileRequest['videoTarget'] {
	const target = { ...defaultVideoTarget(), ...(form.videoTarget ?? {}) };
	return {
		...target,
		codecs: uniqueTrimmed(target.codecs ?? []),
		hdrFormats: uniqueTrimmed(target.hdrFormats ?? []),
		pixelFormats: uniqueTrimmed(target.pixelFormats ?? []),
		codecScore: normalizedInteger(target.codecScore),
		hdrScore: normalizedInteger(target.hdrScore),
		pixelFormatScore: normalizedInteger(target.pixelFormatScore)
	};
}

function audioTargetsFromForm(form: MediaProfileForm): MediaProfileRequest['audioTargets'] {
	const seen = new Set<string>();
	const targets = [];
	for (const value of form.audioTargets ?? []) {
		const languageId = value.languageId.trim();
		if (!languageId || seen.has(languageId)) continue;
		seen.add(languageId);
		targets.push({
			languageId,
			score: normalizedInteger(value.score),
			required: true,
			codecs: uniqueTrimmed(value.codecs ?? []),
			channels: uniqueTrimmed(value.channels ?? []),
			minimumBitrateKbps: positiveInteger(value.minimumBitrateKbps),
			preferredBitrateKbps: positiveInteger(value.preferredBitrateKbps),
			lossyTranscodePolicy: value.lossyTranscodePolicy ?? 'disabled'
		});
	}
	return targets;
}

function subtitleTargetsFromForm(form: MediaProfileForm): MediaProfileRequest['subtitleTargets'] {
	const seen = new Set<string>();
	const targets = [];
	for (const value of form.subtitleTargets ?? []) {
		const languageId = value.languageId.trim();
		if (!languageId || seen.has(languageId)) continue;
		seen.add(languageId);
		targets.push({
			languageId,
			score: normalizedInteger(value.score),
			required: true,
			source: value.source ?? 'any',
			formats: uniqueTrimmed(value.formats ?? [])
		});
	}
	return targets;
}

function uniqueTrimmed(values: string[]) {
	const seen = new Set<string>();
	const result = [];
	for (const value of values) {
		const trimmed = value.trim();
		const key = trimmed.toLowerCase();
		if (!trimmed || seen.has(key)) continue;
		seen.add(key);
		result.push(trimmed);
	}
	return result;
}

function positiveInteger(value: number | undefined) {
	const parsed = normalizedInteger(value);
	return parsed > 0 ? parsed : undefined;
}

function normalizedInteger(value: number | string | undefined) {
	const parsed = Number(value ?? 0);
	return Number.isFinite(parsed) ? Math.trunc(parsed) : 0;
}
