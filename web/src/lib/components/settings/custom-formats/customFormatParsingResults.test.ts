import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';

import CustomFormatParsingResults from './CustomFormatParsingResults.svelte';
import type { CustomFormatParsingResponse } from '$lib/settings/types';

function parsingResult(
	overrides: Partial<CustomFormatParsingResponse> = {}
): CustomFormatParsingResponse {
	return {
		fileName: 'Movie.Title.2026.2160p.WEB-DL.DV.HDR10.Atmos-GROUP.mkv',
		release: {
			releaseTitle: 'Movie.Title.2026.2160p.WEB-DL.DV.HDR10.Atmos-GROUP',
			movieTitle: 'Movie Title',
			seriesTitle: '',
			year: '2026',
			seasonNumber: null,
			episodeNumber: null,
			seasonPack: false,
			edition: 'Director Cut',
			releaseGroup: 'GROUP',
			releaseHash: 'ABC123'
		},
		quality: {
			qualityId: 'web-2160p',
			quality: 'WEB-2160p',
			source: 'WEB-DL',
			resolution: '2160p',
			videoCodec: 'HEVC',
			audioCodec: 'Atmos',
			audioChannels: '7.1',
			version: 'v2',
			proper: true,
			repack: false,
			real: true
		},
		languages: ['English', 'German'],
		details: {
			releaseType: 'movie',
			customFormatNames: ['Dolby Vision'],
			matchedSpecCount: 2
		},
		matchedProfile: { id: 'profile-1', name: 'Remux UHD' },
		calculatedScore: 125,
		matchedCustomFormats: [
			{
				id: 'format-1',
				name: 'Dolby Vision',
				score: 100,
				matchedSpecs: [
					{
						id: 'spec-1',
						name: 'DV',
						type: 'releaseTitle',
						value: '\\bDV\\b',
						required: true
					},
					{
						id: 'spec-2',
						name: 'Atmos',
						type: 'audioCodec',
						value: 'Atmos',
						required: false
					}
				]
			}
		],
		...overrides
	};
}

describe('custom format parsing results (SCN-SETTINGS-017)', () => {
	it('renders parsed release, quality, language, and custom format matches', () => {
		const { body } = render(CustomFormatParsingResults, {
			props: { result: parsingResult() }
		});

		expect(body).toContain('Release');
		expect(body).toContain('Movie.Title.2026.2160p.WEB-DL.DV.HDR10.Atmos-GROUP');
		expect(body).toContain('Movie Title');
		expect(body).toContain('Director Cut');
		expect(body).toContain('WEB-2160p');
		expect(body).toContain('HEVC');
		expect(body).toContain('Proper');
		expect(body).toContain('Yes');
		expect(body).toContain('English');
		expect(body).toContain('German');
		expect(body).toContain('Remux UHD');
		expect(body).toContain('125');
		expect(body).toContain('Dolby Vision');
		expect(body).toContain('DV');
		expect(body).toContain('Atmos');
	});

	it('renders empty match and language fallbacks', () => {
		const { body } = render(CustomFormatParsingResults, {
			props: {
				result: parsingResult({
					languages: [],
					matchedProfile: undefined,
					calculatedScore: 0,
					matchedCustomFormats: []
				})
			}
		});

		expect(body).toContain('No custom formats matched');
		expect(body).toContain('Matched profile');
		expect(body).toContain('Calculated score');
		expect(body).toContain('0');
	});
});
