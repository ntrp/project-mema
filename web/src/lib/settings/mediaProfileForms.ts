import type {
	MediaProfile,
	MediaProfileComponentFallback,
	MediaProfileComponentSource,
	MediaProfileComponentType,
	MediaProfileForm,
	MediaProfileRequest
} from './types';

export function emptyMediaProfileForm(): MediaProfileForm {
	return {
		name: '',
		isDefault: false,
		qualityIds: [],
		upgradesAllowed: true,
		upgradeUntilQualityId: undefined,
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 0,
		minimumCustomFormatScoreIncrement: 1,
		removeNonEnabledLanguages: false,
		removeNonEnabledSubtitleLanguages: false,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		targetLanguages: ['EN'],
		targetLanguageScores: [{ languageId: 'EN', score: 0, required: false }],
		subtitleLanguages: [{ languageId: 'EN', score: 0, required: false, subtitleType: 'embedded' }],
		componentTargets: [],
		customFormatScores: []
	};
}

export function mediaProfileFormFromProfile(profile: MediaProfile): MediaProfileForm {
	return {
		id: profile.id,
		name: profile.name,
		isDefault: profile.isDefault,
		qualityIds: [...(profile.qualityIds ?? [])],
		upgradesAllowed: profile.upgradesAllowed,
		upgradeUntilQualityId: profile.upgradeUntilQualityId,
		minimumCustomFormatScore: profile.minimumCustomFormatScore,
		upgradeUntilCustomFormatScore: profile.upgradeUntilCustomFormatScore,
		minimumCustomFormatScoreIncrement: profile.minimumCustomFormatScoreIncrement,
		removeNonEnabledLanguages: profile.removeNonEnabledLanguages,
		removeNonEnabledSubtitleLanguages: profile.removeNonEnabledSubtitleLanguages,
		preferredProtocol: profile.preferredProtocol,
		seriesPackPreference: profile.seriesPackPreference,
		targetLanguages: [...(profile.targetLanguages ?? [])],
		targetLanguageScores: languageScoresFromProfile(profile),
		subtitleLanguages: (profile.subtitleLanguages ?? []).map((language) => ({ ...language })),
		componentTargets: (profile.componentTargets ?? []).map((target) => ({ ...target })),
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
	const subtitleLanguages = subtitleLanguagesFromForm(form);
	const componentTargets = componentTargetsFromForm(form);
	return {
		name: form.name.trim(),
		isDefault: form.isDefault,
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
		removeNonEnabledSubtitleLanguages: form.removeNonEnabledSubtitleLanguages,
		preferredProtocol: form.preferredProtocol ?? 'any',
		seriesPackPreference: form.seriesPackPreference ?? 'auto',
		targetLanguages: targetLanguageScores.map((score) => score.languageId),
		targetLanguageScores,
		subtitleLanguages,
		componentTargets,
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

function subtitleLanguagesFromForm(form: MediaProfileForm) {
	const seen = new Set<string>();
	const languages = [];
	for (const value of form.subtitleLanguages ?? []) {
		const languageId = value.languageId.trim();
		if (!languageId || seen.has(languageId)) {
			continue;
		}
		seen.add(languageId);
		languages.push({
			languageId,
			score: normalizedInteger(value.score),
			required: value.required,
			subtitleType: value.subtitleType ?? 'any'
		});
	}
	return languages;
}

function componentTargetsFromForm(form: MediaProfileForm): MediaProfileRequest['componentTargets'] {
	const targets = [];
	for (const value of form.componentTargets ?? []) {
		const componentType = componentTypeFrom(value.componentType);
		const languageId = trimmedValue(value.languageId);
		const codec = trimmedValue(value.codec);
		const channels = componentType === 'audio' ? trimmedValue(value.channels) : undefined;
		if (componentType === 'audio' && !languageId && !codec && !channels) continue;
		if (componentType === 'subtitle' && !languageId && !codec) continue;
		targets.push({
			id: value.id,
			componentType,
			required: value.required,
			languageId: componentType === 'video' ? undefined : languageId,
			codec,
			channels,
			source: componentSourceFrom(value.source, componentType),
			fallbackBehavior: componentFallbackFrom(value.fallbackBehavior)
		});
	}
	return targets;
}

function componentTypeFrom(value: string | undefined): MediaProfileComponentType {
	return value === 'audio' || value === 'subtitle' ? value : 'video';
}

function componentSourceFrom(
	value: string | undefined,
	componentType: MediaProfileComponentType
): MediaProfileComponentSource {
	if (value === 'existing') return 'existing';
	if (componentType === 'subtitle') return 'subtitleProvider';
	return 'release';
}

function componentFallbackFrom(value: string | undefined): MediaProfileComponentFallback {
	if (value === 'preferExisting' || value === 'allowMissing') return value;
	return 'strict';
}

function trimmedValue(value: string | undefined) {
	const trimmed = value?.trim();
	return trimmed ? trimmed : undefined;
}

function normalizedInteger(value: number | string | undefined) {
	const parsed = Number(value ?? 0);
	if (!Number.isFinite(parsed)) {
		return 0;
	}
	return Math.trunc(parsed);
}
