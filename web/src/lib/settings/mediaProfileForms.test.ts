import { describe, expect, it } from 'vitest';

import {
	emptyMediaProfileForm,
	mediaProfileFormFromProfile,
	normalizeMediaProfileForm
} from './mediaProfileForms';
import type { MediaProfile, MediaProfileForm } from './types';

describe('media profile forms (SCN-SETTINGS-012)', () => {
	it('starts with an english target language score', () => {
		expect(emptyMediaProfileForm()).toMatchObject({
			name: '',
			qualityIds: [],
			targetLanguages: ['english'],
			targetLanguageScores: [{ languageId: 'english', score: 0, required: false }],
			subtitleLanguages: [{ languageId: 'english', required: true, subtitleType: 'any' }]
		});
	});

	it('copies arrays from an existing profile', () => {
		const profile = {
			id: 'profile-1',
			name: 'Main',
			qualityIds: ['webdl-1080p'],
			upgradesAllowed: true,
			upgradeUntilQualityId: 'webdl-1080p',
			minimumCustomFormatScore: 0,
			upgradeUntilCustomFormatScore: 100,
			minimumCustomFormatScoreIncrement: 1,
			removeNonEnabledLanguages: false,
			preferredProtocol: 'any',
			seriesPackPreference: 'auto',
			targetLanguages: ['english'],
			targetLanguageScores: [{ languageId: 'english', score: 10, required: true }],
			subtitleLanguages: [{ languageId: 'english', required: true, subtitleType: 'embedded' }],
			customFormatScores: [{ customFormatId: 'cf-1', score: 50 }]
		} as MediaProfile;

		const form = mediaProfileFormFromProfile(profile);
		form.qualityIds.push('bluray-2160p');
		form.targetLanguageScores[0].score = 20;
		form.subtitleLanguages[0].subtitleType = 'external';

		expect(profile.qualityIds).toEqual(['webdl-1080p']);
		expect(profile.targetLanguageScores?.[0].score).toBe(10);
		expect(profile.subtitleLanguages?.[0].subtitleType).toBe('embedded');
	});

	it('normalizes profile request payloads', () => {
		const form = {
			...emptyMediaProfileForm(),
			name: '  Main  ',
			qualityIds: [' webdl-1080p ', 'webdl-1080p', 'bluray-2160p'],
			upgradeUntilQualityId: 'raw-hd',
			minimumCustomFormatScore: '10',
			upgradeUntilCustomFormatScore: '20.9',
			minimumCustomFormatScoreIncrement: '-5',
			preferredProtocol: undefined,
			seriesPackPreference: undefined,
			targetLanguageScores: [
				{ languageId: 'english', score: '100', required: true },
				{ languageId: 'english', score: 50, required: false },
				{ languageId: ' german ', score: Number.NaN, required: false }
			],
			subtitleLanguages: [
				{ languageId: 'english', required: true, subtitleType: 'embedded' },
				{ languageId: 'english', required: false, subtitleType: 'external' }
			],
			customFormatScores: [
				{ customFormatId: 'cf-1', score: '25.9' },
				{ customFormatId: '', score: 100 }
			]
		} as unknown as MediaProfileForm;

		expect(normalizeMediaProfileForm(form)).toMatchObject({
			name: 'Main',
			qualityIds: ['webdl-1080p', 'bluray-2160p'],
			upgradeUntilQualityId: undefined,
			minimumCustomFormatScore: 10,
			upgradeUntilCustomFormatScore: 20,
			minimumCustomFormatScoreIncrement: 0,
			preferredProtocol: 'any',
			seriesPackPreference: 'auto',
			targetLanguages: ['english', 'german'],
			targetLanguageScores: [
				{ languageId: 'english', score: 100, required: true },
				{ languageId: 'german', score: 0, required: false }
			],
			subtitleLanguages: [{ languageId: 'english', required: true, subtitleType: 'embedded' }],
			customFormatScores: [{ customFormatId: 'cf-1', score: 25 }]
		});
	});
});
