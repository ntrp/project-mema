import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaRequestCard from '$lib/components/app/requests/MediaRequestCard.svelte';
import ReleaseMatchInfo from '$lib/components/app/media/release-display/ReleaseMatchInfo.svelte';
import CustomFormatCard from '$lib/components/settings/custom-formats/CustomFormatCard.svelte';
import MediaProfileTable from '$lib/components/settings/profiles/MediaProfileTable.svelte';
import type {
	CustomFormat,
	MediaProfile,
	MediaRequest,
	QualitySizeSetting
} from '$lib/settings/types';
import { renderWithTooltip } from './renderHelpers';

describe('rendered settings cards (SCN-SETTINGS-009)', () => {
	it('renders custom format specs in required, excluded, optional order', () => {
		const { body } = render(CustomFormatCard, {
			props: {
				format: customFormat(),
				deleting: true,
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});

		expect(body).toContain('Scenario Format');
		expect(body).toContain('Required Source');
		expect(body).toContain('Rejected Language');
		expect(body).toContain('Optional Codec');
		expect(body).toContain('Deleting Scenario Format');
	});

	it('renders profile quality, upgrade, language, and score summaries', () => {
		const { body } = render(MediaProfileTable, {
			props: {
				profiles: [mediaProfile()],
				qualities: qualities(),
				deletingId: 'profile-1',
				onDelete: vi.fn()
			}
		});

		expect(body).toContain('Scenario Profile');
		expect(body).toContain('HD-1080p');
		expect(body).toContain('UHD-2160p');
		expect(body).toContain('eng, deu');
		expect(body).toContain('10 min / 50 cutoff / 1 scored');
		expect(body).toContain('Deleting Scenario Profile');
	});

	it('renders an empty profile table state', () => {
		const { body } = render(MediaProfileTable, {
			props: { profiles: [], qualities: [], onDelete: vi.fn() }
		});

		expect(body).toContain('No profiles configured');
	});
});

describe('rendered media cards (SCN-MEDIA-003)', () => {
	it('renders media request metadata, tags, and approved status', () => {
		const { body } = render(MediaRequestCard, {
			props: {
				request: {
					id: 'request-1',
					title: 'Scenario Movie',
					type: 'movie',
					year: 2026,
					status: 'approved',
					requestedByUsername: 'admin',
					overview: 'A request visible to reviewers.',
					posterPath: '/poster.jpg',
					tags: ['family', 'uhd', 'priority', 'hidden']
				} as MediaRequest
			}
		});

		expect(body).toContain('Scenario Movie');
		expect(body).toContain('movie · 2026 · Requested by admin');
		expect(body).toContain('family');
		expect(body).toContain('uhd');
		expect(body).not.toContain('>hidden<');
		expect(body).toContain('approved');
		expect(body).toContain('https://image.tmdb.org/t/p/w185/poster.jpg');
	});

	it('renders release match severity labels', () => {
		const { body } = renderWithTooltip(ReleaseMatchInfo, {
			info: {
				severity: 'warning',
				details: ['Quality is below cutoff'],
				customFormatContributors: [{ label: 'HDR', score: 10 }]
			} as never,
			mediaType: 'movie' as const
		});

		expect(body).toContain('Release warning');
		expect(body).toContain('text-amber-500');
	});
});

function customFormat(): CustomFormat {
	return {
		id: 'format-1',
		name: 'Scenario Format',
		includeInRenameTemplate: false,
		includeSpecs: [
			{
				id: 'required-source',
				name: 'Required Source',
				type: 'source',
				value: 'WEB-DL',
				required: true
			},
			{
				id: 'optional-codec',
				name: 'Optional Codec',
				type: 'videoCodec',
				value: 'x265',
				required: false
			}
		],
		excludeSpecs: [
			{
				id: 'rejected-language',
				name: 'Rejected Language',
				type: 'language',
				value: 'German',
				required: true
			}
		]
	} as CustomFormat;
}

function mediaProfile(): MediaProfile {
	return {
		id: 'profile-1',
		name: 'Scenario Profile',
		qualityIds: ['q-1080p', 'q-2160p'],
		upgradesAllowed: true,
		upgradeUntilQualityId: 'q-2160p',
		targetLanguages: ['eng', 'deu'],
		minimumCustomFormatScore: 10,
		upgradeUntilCustomFormatScore: 50,
		customFormatScores: [{ customFormatId: 'format-1', score: 10 }],
		updatedAt: '2026-07-03T01:02:03Z'
	} as MediaProfile;
}

function qualities(): QualitySizeSetting[] {
	return [
		{ qualityId: 'q-1080p', name: 'HD-1080p' },
		{ qualityId: 'q-2160p', name: 'UHD-2160p' }
	] as QualitySizeSetting[];
}
