import { describe, expect, it } from 'vitest';

import {
	emptyMediaProfileForm,
	mediaProfileFormFromProfile,
	normalizeMediaProfileForm
} from './mediaProfileForms';
import type { MediaProfile, MediaProfileForm } from './types';

describe('media profile forms (SCN-SETTINGS-012)', () => {
	it('starts with required video, audio, and container targets', () => {
		expect(emptyMediaProfileForm()).toMatchObject({
			name: '',
			isDefault: false,
			finalContainer: 'mkv',
			qualityIds: [],
			audioLossyTranscodePolicy: 'disabled',
			subtitleMode: 'mixed',
			allowSubtitleReleaseFallback: false,
			audioTargets: [
				{
					languageId: 'EN',
					score: 0
				}
			],
			subtitleTargets: []
		});
	});

	it('copies arrays from an existing profile', () => {
		const profile = mediaProfile();

		const form = mediaProfileFormFromProfile(profile);
		form.qualityIds.push('bluray-2160p');
		form.audioTargets[0].score = 20;
		form.subtitleTargets[0].score = 30;

		expect(profile.qualityIds).toEqual(['webdl-1080p']);
		expect(profile.audioTargets?.[0].score).toBe(10);
		expect(profile.subtitleTargets?.[0].score).toBe(25);
		expect(form.isDefault).toBe(true);
		expect(form.removeUnwantedSubtitles).toBe(true);
	});

	it('normalizes profile request payloads', () => {
		const form = {
			...emptyMediaProfileForm(),
			name: '  Main  ',
			isDefault: true,
			finalContainer: 'mp4',
			qualityIds: [' webdl-1080p ', 'webdl-1080p', 'bluray-2160p'],
			upgradeUntilQualityId: 'raw-hd',
			minimumCustomFormatScore: '10',
			upgradeUntilCustomFormatScore: '20.9',
			minimumCustomFormatScoreIncrement: '-5',
			subtitleMode: 'embedded',
			allowSubtitleReleaseFallback: true,
			preferredProtocol: undefined,
			seriesPackPreference: undefined,
			videoTarget: {
				codecs: [' h265 ', 'h265', 'av1'],
				codecScore: '15',
				hdrFormats: [' HDR10 ']
			},
			audioTargets: [
				{
					languageId: 'english',
					score: '100',
					targetCodec: ' AAC ',
					targetChannels: [' Stereo ', '2.0', 'Atmos', '7.1'],
					minimumBitrateKbps: '384'
				},
				{
					languageId: 'english',
					score: 50
				}
			],
			subtitleTargets: [
				{
					languageId: 'english',
					score: '25',
					formats: [' SRT ', 'subrip', 'ass']
				},
				{
					languageId: 'english',
					score: 10
				}
			],
			customFormatScores: [
				{ customFormatId: 'cf-1', score: '25.9' },
				{ customFormatId: '', score: 100 }
			]
		} as unknown as MediaProfileForm;

		expect(normalizeMediaProfileForm(form)).toMatchObject({
			name: 'Main',
			isDefault: true,
			finalContainer: 'mp4',
			qualityIds: ['webdl-1080p', 'bluray-2160p'],
			upgradeUntilQualityId: undefined,
			minimumCustomFormatScore: 10,
			upgradeUntilCustomFormatScore: 20,
			minimumCustomFormatScoreIncrement: 0,
			preferredProtocol: 'any',
			seriesPackPreference: 'auto',
			audioLossyTranscodePolicy: 'disabled',
			subtitleMode: 'embedded',
			allowSubtitleReleaseFallback: true,
			videoTarget: { codecs: ['h265', 'av1'], codecScore: 15, hdrFormats: ['HDR10'] },
			audioTargets: [
				{
					languageId: 'english',
					score: 100,
					targetCodec: 'AAC',
					targetChannels: ['2.0', '7.1'],
					minimumBitrateKbps: 384
				}
			],
			subtitleTargets: [
				{
					languageId: 'english',
					score: 25,
					formats: ['subrip', 'ass']
				}
			],
			customFormatScores: [{ customFormatId: 'cf-1', score: 25 }]
		});
	});
});

function mediaProfile(): MediaProfile {
	return {
		id: 'profile-1',
		name: 'Main',
		isDefault: true,
		finalContainer: 'mkv',
		qualityIds: ['webdl-1080p'],
		upgradesAllowed: true,
		upgradeUntilQualityId: 'webdl-1080p',
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 100,
		minimumCustomFormatScoreIncrement: 1,
		removeUnwantedAudio: false,
		audioLossyTranscodePolicy: 'disabled',
		removeUnwantedSubtitles: true,
		subtitleMode: 'embedded',
		allowSubtitleReleaseFallback: true,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		videoTarget: {},
		audioTargets: [
			{
				languageId: 'english',
				score: 10
			}
		],
		subtitleTargets: [
			{
				languageId: 'english',
				score: 25
			}
		],
		customFormatScores: [{ customFormatId: 'cf-1', score: 50 }],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z'
	};
}
