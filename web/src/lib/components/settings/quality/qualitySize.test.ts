import { describe, expect, it } from 'vitest';
import {
	activeTrackStyle,
	gibPerHourToMbPerMinute,
	mbPerMinuteToGibPerHour,
	mbPerMinuteTitle,
	nextSliderQuality,
	qualityRequest,
	rowError,
	sliderValues
} from './qualitySize';
import type { QualitySizeSetting } from '$lib/settings/types';

const baseQuality: QualitySizeSetting = {
	qualityId: 'bluray-1080p',
	name: 'Bluray 1080p',
	sortOrder: 10,
	minimumSizeMbPerMinute: 10,
	preferredSizeMbPerMinute: 20,
	maximumSizeMbPerMinute: 30,
	createdAt: '2026-07-03T00:00:00Z',
	updatedAt: '2026-07-03T00:00:00Z'
};

describe('quality size controls', () => {
	it('SCN-SETTINGS-008 validates and converts visible size ranges', () => {
		expect(rowError({ ...baseQuality, minimumSizeMbPerMinute: -1 })).toBe(
			'Minimum must be zero or greater'
		);
		expect(rowError({ ...baseQuality, preferredSizeMbPerMinute: 5 })).toBe(
			'Preferred must be at least minimum'
		);
		expect(rowError({ ...baseQuality, maximumSizeMbPerMinute: 5 })).toBe(
			'Maximum must be at least minimum'
		);
		expect(rowError({ ...baseQuality, preferredSizeMbPerMinute: 40 })).toBe(
			'Preferred must be at most maximum'
		);
		expect(rowError(baseQuality)).toBe('');

		expect(mbPerMinuteToGibPerHour(20)).toBe(1.17);
		expect(gibPerHourToMbPerMinute(1.17)).toBe(19.97);
		expect(mbPerMinuteTitle('Min', 1.17)).toBe('Min: 19.97 MiB/min');
		expect(sliderValues(baseQuality)).toEqual({ minimum: 0.59, preferred: 1.17, maximum: 1.76 });
		expect(qualityRequest(baseQuality)).toMatchObject({ qualityId: 'bluray-1080p' });
		expect(activeTrackStyle({ minimum: 12, maximum: 24 })).toBe('left: 10%; width: 10%');

		const next = nextSliderQuality(baseQuality, 'preferred', '1.5');
		expect(next.preferredSizeMbPerMinute).toBe(25.6);
		expect(next.minimumSizeMbPerMinute).toBeLessThanOrEqual(next.preferredSizeMbPerMinute ?? 0);
	});
});
