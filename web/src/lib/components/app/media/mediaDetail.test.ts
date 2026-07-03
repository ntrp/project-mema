import { describe, expect, it } from 'vitest';

import { mediaMetadataDetail } from './mediaDetail';
import type { MediaItem } from '$lib/settings/types';

describe('media detail projection (SCN-MEDIA-004)', () => {
	it('projects stored media items into metadata details with local fallbacks', () => {
		const detail = mediaMetadataDetail({
			id: 'media-1',
			type: 'movie',
			title: 'Scenario Movie',
			year: 2026,
			monitored: true,
			overview: 'Overview',
			posterPath: '/poster.jpg',
			metadataStatus: 'available',
			genres: ['Drama'],
			recommendations: [{ title: 'Next' }]
		} as MediaItem);

		expect(detail).toMatchObject({
			title: 'Scenario Movie',
			type: 'movie',
			year: 2026,
			monitored: true,
			externalProvider: 'local',
			externalId: 'media-1',
			overview: 'Overview',
			posterPath: '/poster.jpg',
			status: 'available',
			genres: ['Drama'],
			recommendations: [{ title: 'Next' }]
		});
	});
});
