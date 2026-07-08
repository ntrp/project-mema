import { describe, expect, it } from 'vitest';

import { wantedDisplayRows } from '$lib/components/app/home/wanted/wantedRows';
import type { MediaItem } from '$lib/settings/types';

describe('wanted display rows', () => {
	it('builds media, target, and custom-format rows', () => {
		const rows = wantedDisplayRows([
			mediaItem(),
			mediaItem({
				id: 'downloaded-1',
				status: 'downloaded',
				targetSatisfaction: {
					targets: [
						{
							id: 'audio:1',
							type: 'audio',
							state: 'pending',
							mediaItemId: 'downloaded-1',
							languageId: 'english',
							requiredOperation: {
								type: 'audio_transcode',
								manual: true,
								automatic: true,
								reason: 'Transcode audio.'
							},
							reasons: ['Transcode audio.']
						}
					],
					candidates: []
				},
				files: [
					{
						path: '/library/Movie.mkv',
						status: 'available',
						rollup: {
							state: 'upgradeable',
							targetCounts: emptyCounts(),
							reasons: ['Custom format score can improve.']
						}
					}
				]
			})
		]);

		expect(rows.map((row) => row.kind)).toEqual(['media', 'target', 'custom_format_upgrade']);
		expect(rows[1].operation).toBe('Transcode audio.');
		expect(rows[2].context).toBe('Movie.mkv');
	});
});

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	return {
		id: 'media-1',
		title: 'Scenario Movie',
		type: 'movie',
		status: 'missing',
		monitored: true,
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		filePaths: [],
		metadataFilePaths: [],
		createdAt: '2026-07-08T00:00:00Z',
		updatedAt: '2026-07-08T00:00:00Z',
		...overrides
	};
}

function emptyCounts() {
	return {
		missing: 0,
		partial: 0,
		pending: 0,
		satisfied: 0,
		upgradeable: 0,
		blocked: 0,
		failed: 0
	};
}
