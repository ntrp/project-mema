import { describe, expect, it } from 'vitest';

import { groupQualitiesByResolution } from './qualityGroups';

describe('quality resolution grouping (SCN-SETTINGS-010)', () => {
	it('groups qualities by resolution in display order', () => {
		const groups = groupQualitiesByResolution([
			{ qualityId: 'webdl-1080p', name: 'WEBDL-1080p' },
			{ qualityId: 'cam', name: 'CAM' },
			{ qualityId: 'bluray-2160p', name: 'Bluray-2160p' },
			{ qualityId: 'br-disk', name: 'BR-DISK' }
		]);

		expect(groups.map((group) => group.id)).toEqual(['sd', '1080p', '4k', 'native']);
		expect(groups.find((group) => group.id === 'native')?.qualities[0].qualityId).toBe('br-disk');
	});

	it('places unrecognized qualities in the other group', () => {
		const groups = groupQualitiesByResolution([{ qualityId: 'unknown', name: 'Mystery' }]);

		expect(groups).toEqual([
			{
				id: 'other',
				label: 'Other',
				qualities: [{ qualityId: 'unknown', name: 'Mystery' }]
			}
		]);
	});
});
