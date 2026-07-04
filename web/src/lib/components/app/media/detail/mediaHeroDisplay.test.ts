import { describe, expect, it } from 'vitest';

import {
	imageUrl,
	mediaHeroTopInfo,
	monitorHint,
	monitorStatus,
	runtimeText,
	statusBadgeClass,
	statusLabel
} from '$lib/components/app/media/detail/mediaHeroDisplay';
import type { MediaMetadataDetails } from '$lib/settings/types';

describe('media hero display helpers (SCN-MEDIA-003)', () => {
	it('builds image URLs and runtime labels', () => {
		expect(imageUrl('/poster.jpg')).toBe('https://image.tmdb.org/t/p/w780/poster.jpg');
		expect(imageUrl('/poster.jpg', 'w342')).toBe('https://image.tmdb.org/t/p/w342/poster.jpg');
		expect(imageUrl('https://cdn.test/poster.jpg')).toBe('https://cdn.test/poster.jpg');
		expect(imageUrl()).toBeUndefined();

		expect(runtimeText(125)).toBe('2h 5m');
		expect(runtimeText(120)).toBe('2h');
		expect(runtimeText(45)).toBe('45m');
		expect(runtimeText(0)).toBeUndefined();
	});

	it('summarizes series counts and media status', () => {
		expect(mediaHeroTopInfo({ seasonCount: 2, episodeCount: 18 } as MediaMetadataDetails)).toEqual([
			['Seasons', '2'],
			['Episodes', '18']
		]);
		expect(statusLabel('downloaded')).toBe('Downloaded');
		expect(statusLabel('downloading')).toBe('Downloading');
		expect(statusLabel('missing')).toBe('Missing');
		expect(statusBadgeClass('downloaded')).toContain('emerald');
		expect(statusBadgeClass('downloading')).toContain('primary');
		expect(statusBadgeClass('missing')).toContain('destructive');
	});

	it('builds monitor labels and hints by media type', () => {
		expect(monitorStatus({ monitored: true } as MediaMetadataDetails)).toBe('Monitored');
		expect(monitorStatus({ monitored: false } as MediaMetadataDetails)).toBe('Not monitored');
		expect(monitorHint({ type: 'serie', monitored: true } as MediaMetadataDetails)).toContain(
			'future episodes'
		);
		expect(monitorHint({ type: 'movie', monitored: false } as MediaMetadataDetails)).toContain(
			'this movie'
		);
	});
});
