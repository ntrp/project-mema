import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaRequestArea from './MediaRequestArea.svelte';
import type { LibraryFolder, MediaRequest, QualityProfileOption } from '$lib/settings/types';

const now = '2026-07-03T00:00:00Z';

function mediaRequest(overrides: Partial<MediaRequest> = {}): MediaRequest {
	return {
		id: 'request-1',
		title: 'Edge of Tomorrow',
		type: 'movie',
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		year: 2014,
		externalProvider: 'tmdb',
		externalId: '137113',
		overview: 'A soldier relives the same battle.',
		posterPath: '/poster.jpg',
		tags: ['action', 'sci-fi'],
		status: 'pending',
		requestedByUserId: 'user-1',
		requestedByUsername: 'casey',
		createdAt: now,
		updatedAt: now,
		...overrides
	};
}

const libraryFolders: LibraryFolder[] = [
	{ id: 'folder-1', path: '/media/movies', createdAt: now, updatedAt: now }
];

const qualityProfiles: QualityProfileOption[] = [{ id: 'profile-1', name: 'HD-1080p' }];

describe('media request area (SCN-MEDIA-012)', () => {
	it('renders request cards and empty list state', () => {
		const listed = render(MediaRequestArea, {
			props: {
				requests: [
					mediaRequest(),
					mediaRequest({ id: 'request-2', title: 'Frieren', type: 'serie' })
				],
				libraryFolders,
				qualityProfiles,
				canManage: true,
				onApprove: vi.fn()
			}
		});

		expect(listed.body).toContain('Media requests');
		expect(listed.body).toContain('Edge of Tomorrow');
		expect(listed.body).toContain('Frieren');
		expect(listed.body).toContain('Requested by casey');
		expect(listed.body).toContain('href="/requests/request-1"');

		const empty = render(MediaRequestArea, {
			props: {
				requests: [],
				libraryFolders,
				qualityProfiles,
				canManage: true,
				onApprove: vi.fn()
			}
		});

		expect(empty.body).toContain('No requests');
		expect(empty.body).toContain('Requested media will appear here.');
	});

	it('renders selected request facts, tags, and approval state', () => {
		const selected = render(MediaRequestArea, {
			props: {
				requests: [mediaRequest({ qualityProfileId: 'profile-1', libraryFolderId: 'folder-1' })],
				selectedRequestId: 'request-1',
				libraryFolders,
				qualityProfiles,
				canManage: true,
				approvingRequestId: 'request-1',
				onApprove: vi.fn()
			}
		});

		expect(selected.body).toContain('Back to requests');
		expect(selected.body).toContain('Edge of Tomorrow');
		expect(selected.body).toContain('Requested by');
		expect(selected.body).toContain('casey');
		expect(selected.body).toContain('HD-1080p');
		expect(selected.body).toContain('/media/movies');
		expect(selected.body).toContain('action');
		expect(selected.body).toContain('sci-fi');
		expect(selected.body).toContain('Approving');

		const missing = render(MediaRequestArea, {
			props: {
				requests: [],
				selectedRequestId: 'missing',
				libraryFolders,
				qualityProfiles,
				canManage: true,
				onApprove: vi.fn()
			}
		});

		expect(missing.body).toContain('Request not found');
		expect(missing.body).toContain('The request is not visible to your account.');
	});
});
