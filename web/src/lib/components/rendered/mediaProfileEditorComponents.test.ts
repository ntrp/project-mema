import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaProfileCustomFormatScores from '$lib/components/settings/profiles/MediaProfileCustomFormatScores.svelte';
import MediaProfileQualitySelector from '$lib/components/settings/profiles/MediaProfileQualitySelector.svelte';
import MediaProfileRules from '$lib/components/settings/profiles/MediaProfileRules.svelte';
import { emptyMediaProfileForm } from '$lib/settings/forms';
import type { CustomFormat, QualitySizeSetting } from '$lib/settings/types';

describe('rendered media profile editor controls (SCN-SETTINGS-023)', () => {
	it('renders grouped quality choices and loading errors', () => {
		const { body } = render(MediaProfileQualitySelector, {
			props: {
				form: { ...emptyMediaProfileForm(), qualityIds: ['webdl-1080p'] },
				qualities: qualities(),
				loading: false,
				error: 'Could not load quality sizes',
				onChange: vi.fn()
			}
		});

		expect(body).toContain('Qualities');
		expect(body).toContain('Select all');
		expect(body).toContain('Could not load quality sizes');
		expect(body).toContain('WEB-DL 1080p');
		expect(body).toContain('Remux 2160p');

		const loading = render(MediaProfileQualitySelector, {
			props: {
				form: emptyMediaProfileForm(),
				qualities: [],
				loading: true,
				error: '',
				onChange: vi.fn()
			}
		});
		expect(loading.body).toContain('Loading qualities');
	});

	it('renders profile rules, language scores, and selected upgrade labels', () => {
		const form = {
			...emptyMediaProfileForm(),
			qualityIds: ['webdl-1080p'],
			upgradeUntilQualityId: 'webdl-1080p',
			preferredProtocol: 'usenet' as const,
			seriesPackPreference: 'preferPacks' as const,
			removeUnwantedAudio: true,
			removeUnwantedSubtitles: true,
			audioTargets: [
				{
					languageId: 'japanese',
					score: 100,
					required: true,
					lossyTranscodePolicy: 'disabled' as const
				}
			],
			subtitleTargets: [
				{
					languageId: 'english',
					score: 25,
					required: false,
					source: 'any' as const
				}
			]
		};
		const { body } = render(MediaProfileRules, {
			props: {
				form,
				qualities: qualities(),
				onChange: vi.fn()
			}
		});

		expect(body).toContain('General');
		expect(body).toContain('WEB-DL 1080p');
		expect(body).toContain('Prefer Usenet');
		expect(body).toContain('Prefer season packs');
	});

	it('renders custom format scores and empty state', () => {
		const formatId = '00000000-0000-4000-8000-000000000301';
		const customFormats = [{ id: formatId, name: 'Preferred Group' }] as CustomFormat[];
		const scored = render(MediaProfileCustomFormatScores, {
			props: {
				form: {
					...emptyMediaProfileForm(),
					customFormatScores: [{ customFormatId: formatId, score: 25 }]
				},
				customFormats,
				onChange: vi.fn()
			}
		});
		expect(scored.body).toContain('Custom formats');
		expect(scored.body).toContain('Preferred Group');
		expect(scored.body).toContain('Remove');
		expect(scored.body).toContain('25');

		const empty = render(MediaProfileCustomFormatScores, {
			props: { form: emptyMediaProfileForm(), customFormats, onChange: vi.fn() }
		});
		expect(empty.body).toContain('No custom formats scored for this profile');
	});
});

function qualities(): QualitySizeSetting[] {
	return [
		{
			qualityId: 'webdl-1080p',
			name: 'WEB-DL 1080p',
			minimumSizeMbPerMinute: 10,
			preferredSizeMbPerMinute: 20,
			maximumSizeMbPerMinute: 30
		},
		{
			qualityId: 'remux-2160p',
			name: 'Remux 2160p',
			minimumSizeMbPerMinute: 30,
			preferredSizeMbPerMinute: 60,
			maximumSizeMbPerMinute: 90
		}
	] as QualitySizeSetting[];
}
