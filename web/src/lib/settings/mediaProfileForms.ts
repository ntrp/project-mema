import type { MediaProfile, MediaProfileForm, MediaProfileRequest } from './types';

export function emptyMediaProfileForm(): MediaProfileForm {
	return {
		name: '',
		qualityIds: [],
		upgradesAllowed: true,
		upgradeUntilQualityId: undefined,
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 0,
		minimumCustomFormatScoreIncrement: 1,
		removeNonEnabledLanguages: false,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		targetLanguages: ['english'],
		targetLanguageScores: [{ languageId: 'english', score: 0, required: false }],
		customFormatScores: []
	};
}

export function mediaProfileFormFromProfile(profile: MediaProfile): MediaProfileForm {
	return {
		id: profile.id,
		name: profile.name,
		qualityIds: [...(profile.qualityIds ?? [])],
		upgradesAllowed: profile.upgradesAllowed,
		upgradeUntilQualityId: profile.upgradeUntilQualityId,
		minimumCustomFormatScore: profile.minimumCustomFormatScore,
		upgradeUntilCustomFormatScore: profile.upgradeUntilCustomFormatScore,
		minimumCustomFormatScoreIncrement: profile.minimumCustomFormatScoreIncrement,
		removeNonEnabledLanguages: profile.removeNonEnabledLanguages,
		preferredProtocol: profile.preferredProtocol,
		seriesPackPreference: profile.seriesPackPreference,
		targetLanguages: [...(profile.targetLanguages ?? [])],
		targetLanguageScores: languageScoresFromProfile(profile),
		customFormatScores: (profile.customFormatScores ?? []).map((score) => ({ ...score }))
	};
}

export function normalizeMediaProfileForm(form: MediaProfileForm): MediaProfileRequest {
	const qualityIds = [...new Set(form.qualityIds.map((id) => id.trim()).filter(Boolean))];
	const customFormatScores = form.customFormatScores
		.filter((score) => score.customFormatId)
		.map((score) => ({
			customFormatId: score.customFormatId,
			score: normalizedInteger(score.score)
		}));
	const targetLanguageScores = languageScoresFromForm(form);
	return {
		name: form.name.trim(),
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
		removeNonEnabledLanguages: form.removeNonEnabledLanguages,
		preferredProtocol: form.preferredProtocol ?? 'any',
		seriesPackPreference: form.seriesPackPreference ?? 'auto',
		targetLanguages: targetLanguageScores.map((score) => score.languageId),
		targetLanguageScores,
		customFormatScores
	};
}

function languageScoresFromProfile(profile: MediaProfile) {
	if (profile.targetLanguageScores?.length) {
		return profile.targetLanguageScores.map((score) => ({ ...score }));
	}
	return (profile.targetLanguages ?? []).map((languageId) => ({
		languageId,
		score: 0,
		required: false
	}));
}

function languageScoresFromForm(form: MediaProfileForm) {
	const seen = new Set<string>();
	const source = form.targetLanguageScores?.length
		? form.targetLanguageScores
		: form.targetLanguages.map((languageId) => ({ languageId, score: 0, required: false }));
	const scores = [];
	for (const value of source) {
		const languageId = value.languageId.trim();
		if (!languageId || seen.has(languageId)) {
			continue;
		}
		seen.add(languageId);
		scores.push({ languageId, score: normalizedInteger(value.score), required: value.required });
	}
	return scores;
}

function normalizedInteger(value: number | string | undefined) {
	const parsed = Number(value ?? 0);
	if (!Number.isFinite(parsed)) {
		return 0;
	}
	return Math.trunc(parsed);
}
