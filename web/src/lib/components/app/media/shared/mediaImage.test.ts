import { describe, expect, it } from 'vitest';

import { mediaImageUrl } from './mediaImage';

describe('mediaImageUrl', () => {
	it('keeps local artwork local and expands TMDB paths', () => {
		expect(mediaImageUrl('/api/media/artwork/anidb/456.jpg', 'w342')).toBe(
			'/api/media/artwork/anidb/456.jpg'
		);
		expect(mediaImageUrl('/poster.jpg', 'w342')).toBe('https://image.tmdb.org/t/p/w342/poster.jpg');
	});
});
