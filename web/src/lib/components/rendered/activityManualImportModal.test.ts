import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import ActivityManualImportForm from '$lib/components/app/activity/ActivityManualImportForm.svelte';
import {
	initialManualImportForm,
	manualImportRequestFromForm
} from '$lib/components/app/activity/activityManualImportForm';
import type { DownloadActivity } from '$lib/settings/types';

describe('rendered activity manual import form (SCN-ACTIVITY-002)', () => {
	it('renders movie import fields, parsed release hints, and an error state', () => {
		const { body } = render(ActivityManualImportForm, {
			props: {
				activity: downloadActivity({
					releaseTitle: 'Scenario.Movie.2026.German.2160p.WEB-DL.Atmos-GROUP'
				}),
				importing: false,
				error: 'Manual import failed',
				onImport: vi.fn(),
				onClose: vi.fn()
			}
		});

		expect(body).toContain('Scenario Movie');
		expect(body).toContain('Manual import');
		expect(body).toContain('Scenario.Movie.2026.German.2160p.WEB-DL.Atmos-GROUP');
		expect(body).toContain('Source path');
		expect(body).toContain('/downloads/release/file.mkv');
		expect(body).toContain('Target filename override');
		expect(body).toContain('Movie');
		expect(body).toContain('Year');
		expect(body).toContain('Release group');
		expect(body).toContain('Edition');
		expect(body).toContain('Quality');
		expect(body).toContain('Languages');
		expect(body).toContain('Manual import failed');
		expect(body).toContain('Import');
	});

	it('renders series season and episode fields while importing', () => {
		const { body } = render(ActivityManualImportForm, {
			props: {
				activity: downloadActivity({
					mediaTitle: 'Scenario Series',
					mediaType: 'series',
					releaseTitle: 'Scenario.Series.S02E05.1080p.WEB-DL-GROUP'
				}),
				importing: true,
				onImport: vi.fn(),
				onClose: vi.fn()
			}
		});

		expect(body).toContain('Scenario Series');
		expect(body).toContain('Series');
		expect(body).toContain('Season');
		expect(body).toContain('Episode');
		expect(body).toContain('Episode title');
		expect(body).toContain('Importing');
		expect(body).not.toContain('Manual import failed');
	});

	it('builds request payloads from defaults and trimmed optional fields', () => {
		const form = initialManualImportForm(
			downloadActivity({
				mediaType: 'series',
				mediaTitle: 'Scenario Series',
				releaseTitle: 'Scenario.Series.S02E05.German.1080p.WEB-DL-GROUP'
			})
		);

		expect(form).toMatchObject({
			movieTitle: 'Scenario Series',
			seasonNumber: 1,
			episodeNumber: 1,
			releaseGroup: 'GROUP',
			quality: '1080p',
			languagesText: 'German'
		});

		expect(
			manualImportRequestFromForm({
				...form,
				sourcePath: '/downloads/scenario.mkv',
				targetFileName: ' Scenario Series - S02E05.mkv ',
				episodeTitle: ' The Case ',
				edition: ' ',
				languagesText: ' English, German, '
			})
		).toEqual({
			sourcePath: '/downloads/scenario.mkv',
			targetFileName: 'Scenario Series - S02E05.mkv',
			movieTitle: 'Scenario Series',
			year: 2026,
			seasonNumber: 1,
			episodeNumber: 1,
			episodeTitle: 'The Case',
			releaseGroup: 'GROUP',
			edition: undefined,
			quality: '1080p',
			languages: ['English', 'German']
		});
	});
});

function downloadActivity(overrides: Partial<DownloadActivity> = {}): DownloadActivity {
	return {
		id: 'activity-1',
		mediaItemId: 'media-1',
		mediaTitle: 'Scenario Movie',
		mediaType: 'movie',
		mediaYear: 2026,
		releaseTitle: 'Scenario.Movie.2026.German.1080p.WEB-DL.Atmos-GROUP',
		indexerName: 'Scenario Indexer',
		downloadClientName: 'Scenario Client',
		downloadUrl: 'https://example.test/download',
		status: 'failed',
		progressPercent: 100,
		failureType: 'import',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as DownloadActivity;
}
