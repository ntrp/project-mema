import { render } from 'svelte/server';
import type { Component, ComponentProps } from 'svelte';
import { describe, expect, it, vi } from 'vitest';

import MediaDetail from '$lib/components/app/media/detail/MediaDetail.svelte';
import RenderWithTooltip from '$lib/components/rendered/RenderWithTooltip.svelte';
import { assemblyRun, componentSources } from './mediaComponentAssemblyTestData';
import type {
	Language,
	LibraryFolder,
	MediaItem,
	MediaSearchResult,
	QualityProfileOption
} from '$lib/settings/types';

describe('rendered media detail area (SCN-MEDIA-004)', () => {
	it('renders loading and not-found states', () => {
		expect(renderDetail({ loading: true }).body).toContain('Loading media details');
		expect(renderDetail({ item: undefined, requestedItemId: 'missing-1' }).body).toContain(
			'Media item not found'
		);
		expect(renderDetail({ item: undefined, requestedItemId: 'missing-1' }).body).toContain(
			'No monitored movie matches missing-1.'
		);
	});

	it('renders movie metadata, management actions, file rows, and related media', () => {
		const { body } = renderDetail();

		expect(body).toContain('Scenario Movie');
		expect(body).toContain('Downloaded');
		expect(body).toContain('2026');
		expect(body).toContain('1h 42m');
		expect(body).toContain('R');
		expect(body).toContain('Science Fiction');
		expect(body).toContain('Overview for the scenario movie.');
		expect(body).toContain('Director');
		expect(body).toContain('Ada Director');
		expect(body).toContain('Cast');
		expect(body).toContain('Lead Actor');
		expect(body).toContain('/movies/media-1/cast');
		expect(body).toContain('Quality Profile');
		expect(body).toContain('Min. Availability');
		expect(body).toContain('Refresh metadata');
		expect(body).toContain('Delete media');
		expect(body).toContain('Files');
		expect(body).toContain('Media root');
		expect(body).toContain('/library/Scenario Movie');
		expect(body).toContain('Scenario.Movie.2026.1080p.WEB-DL.DDP5.1.EN.mkv');
		expect(body).toContain('5.00 GiB');
		expect(body).toContain('WEB-DL');
		expect(body).toContain('Automatic search');
		expect(body).toContain('Manual search');
		expect(body).toContain('Delete file');
		expect(body).toContain('Recommendations');
		expect(body).toContain('Next Movie');
		expect(body).toContain('Add to library');
	});

	it('renders season file size for added series', () => {
		const filePath = '/library/Scenario Series/Season 01/Scenario.Series.S01E01.mkv';
		const { body } = renderDetail({
			mediaType: 'serie',
			item: mediaItem({
				title: 'Scenario Series',
				type: 'serie',
				filePaths: [filePath],
				files: [{ path: filePath, status: 'available', sizeBytes: 5 * 1024 * 1024 * 1024 }],
				seasons: [
					{
						name: 'Season 1',
						episodeCount: 1,
						episodes: [{ episodeNumber: 1, name: 'Pilot', monitored: true }]
					}
				]
			})
		});

		expect(body).toContain('Season 1');
		expect(body).toContain('5.00 GiB');
	});

	it('renders retained, blocked, running, failed, and completed component assembly states', () => {
		const retainedBody = renderDetail({
			item: mediaItem({ componentSources: componentSources('pending') })
		}).body;
		expect(retainedBody).toContain('Components');
		expect(retainedBody).toContain('2 retained');
		expect(retainedBody).toContain('1 blocked');
		expect(retainedBody).toContain('review needed');
		expect(retainedBody).toContain('Approve compatibility');

		const runningBody = renderDetail({
			item: mediaItem({
				componentSources: componentSources('approved'),
				assemblyRuns: [assemblyRun({ status: 'running' })]
			})
		}).body;
		expect(runningBody).toContain('Mux job running');

		const failedBody = renderDetail({
			item: mediaItem({
				componentSources: componentSources('approved'),
				assemblyRuns: [assemblyRun({ status: 'failed', errorMessage: 'mkvmerge failed' })]
			})
		}).body;
		expect(failedBody).toContain('Retry assembly');
		expect(failedBody).toContain('mkvmerge failed');

		const completedBody = renderDetail({
			item: mediaItem({
				componentSources: componentSources('approved'),
				assemblyRuns: [assemblyRun({ status: 'succeeded' })]
			})
		}).body;
		expect(completedBody).toContain('assembled.mkv');
		expect(completedBody).toContain('Provenance');
		expect(completedBody).toContain('stream 1');
	});
});

type DetailProps = ComponentProps<typeof MediaDetail>;

function renderDetail(overrides: Partial<DetailProps> = {}) {
	const props: DetailProps = {
		mediaType: 'movie',
		item: mediaItem(),
		loading: false,
		mediaItems: [],
		libraryFolders: [libraryFolder()],
		languages: [language()],
		qualityProfiles: [qualityProfile()],
		requestedItemId: 'media-1',
		activities: [],
		canManage: true,
		actionLabel: 'Add to library',
		onAutoSearchMedia: vi.fn(),
		onRefreshMediaMetadata: vi.fn(),
		onSaveMediaItemOptions: vi.fn(),
		onDeleteMediaFile: vi.fn(),
		onDeleteMedia: vi.fn(),
		onGrabRelease: vi.fn(),
		onAddMedia: vi.fn(),
		...overrides
	};
	return render(RenderWithTooltip, {
		props: {
			component: MediaDetail as unknown as Component<Record<string, unknown>>,
			componentProps: { ...props }
		}
	});
}

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	const filePath = '/library/Scenario Movie/Scenario.Movie.2026.1080p.WEB-DL.DDP5.1.EN.mkv';
	return {
		id: 'media-1',
		title: 'Scenario Movie',
		type: 'movie',
		year: 2026,
		status: 'downloaded',
		monitored: true,
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		externalProvider: 'tmdb',
		externalId: 'scenario-1',
		overview: 'Overview for the scenario movie.',
		posterPath: '/poster.jpg',
		runtimeMinutes: 102,
		genres: ['Science Fiction'],
		facts: [
			{ label: 'Certification', value: 'R' },
			{ label: 'Director', value: 'Ada Director' }
		],
		cast: [{ name: 'Lead Actor', role: 'Hero', profilePath: '/actor.jpg' }],
		recommendations: [searchResult({ title: 'Next Movie' })],
		qualityProfileId: 'profile-1',
		qualityProfileName: 'HD Profile',
		libraryFolderId: 'folder-1',
		libraryFolderPath: '/library',
		mediaFolderPath: '/library/Scenario Movie',
		filePaths: [filePath],
		files: [{ path: filePath, status: 'available', sizeBytes: 5 * 1024 * 1024 * 1024 }],
		metadataFilePaths: [],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as MediaItem;
}

function searchResult(overrides: Partial<MediaSearchResult> = {}): MediaSearchResult {
	return {
		title: 'Scenario Result',
		type: 'movie',
		year: 2027,
		externalProvider: 'tmdb',
		externalId: 'next-1',
		overview: 'Related result.',
		...overrides
	} as MediaSearchResult;
}

function qualityProfile(): QualityProfileOption {
	return {
		id: 'profile-1',
		name: 'HD Profile',
		audioTargets: [
			{
				languageId: 'english',
				score: 0,
				required: true
			}
		]
	};
}

function libraryFolder(): LibraryFolder {
	return {
		id: 'folder-1',
		path: '/library',
		kind: 'movie',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z'
	};
}

function language(): Language {
	return {
		code: 'eng',
		displayName: 'English',
		aliases: ['english'],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z'
	};
}
