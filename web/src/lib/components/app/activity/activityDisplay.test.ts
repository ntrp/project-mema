import { describe, expect, it } from 'vitest';
import {
	activityDisplay,
	cancellable,
	deletable,
	manualImportable,
	releaseGroupFromTitle
} from './activityDisplay';
import type { DownloadActivity } from '$lib/settings/types';

const activity: DownloadActivity = {
	id: 'activity-1',
	mediaItemId: 'media-1',
	mediaTitle: 'Example Movie',
	mediaType: 'movie',
	mediaYear: 2026,
	releaseTitle: 'Example.Movie.2026.German.1080p.WEB-DL.Atmos-GROUP',
	indexerName: 'Scenario Indexer',
	downloadClientName: 'Scenario Client',
	downloadUrl: '',
	status: 'downloading',
	progressPercent: 42,
	createdAt: '2026-07-03T00:00:00Z',
	updatedAt: '2026-07-03T00:01:00Z'
};

describe('activity display', () => {
	it('SCN-ACTIVITY-001 summarizes release details and available actions', () => {
		expect(activityDisplay(activity)).toMatchObject({
			year: '2026',
			languages: ['German'],
			quality: '1080p',
			formats: ['WEB-DL', 'Atmos'],
			progressValue: 42,
			progressLabel: '42%'
		});
		expect(cancellable(activity)).toBe(true);
		expect(manualImportable(activity)).toBe(false);
		expect(deletable(activity)).toBe(false);
		expect(releaseGroupFromTitle(activity.releaseTitle)).toBe('GROUP');

		expect(manualImportable({ ...activity, status: 'failed', failureType: 'import' })).toBe(true);
		expect(deletable({ ...activity, status: 'cancelled' })).toBe(true);
		expect(
			activityDisplay({ ...activity, status: 'completed', progressPercent: undefined })
				.progressLabel
		).toBe('100%');
	});
});
