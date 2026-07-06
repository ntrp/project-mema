import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({
	GET: vi.fn(),
	POST: vi.fn(),
	PUT: vi.fn(),
	DELETE: vi.fn()
}));

vi.mock('$lib/api/client', () => ({ client: clientMock }));

import {
	enqueueMediaComponentExtraction,
	evaluateMediaComponentCompatibility,
	getMediaComponentSource,
	listMediaComponentSources,
	releaseMediaComponentSource,
	reviewMediaComponentCompatibility,
	retainMediaComponentSource
} from './mediaComponentSources';

describe('media component source API helpers', () => {
	beforeEach(() => {
		clientMock.GET.mockReset();
		clientMock.POST.mockReset();
		clientMock.PUT.mockReset();
		clientMock.DELETE.mockReset();
	});

	it('maps source, extraction, and compatibility commands', async () => {
		clientMock.GET.mockResolvedValueOnce({
			data: { sources: [{ id: 'source-1' }] }
		}).mockResolvedValueOnce({ data: { id: 'source-1' } });
		clientMock.POST.mockResolvedValueOnce({ data: { id: 'source-1' } })
			.mockResolvedValueOnce({ data: { id: 'source-1' } })
			.mockResolvedValueOnce({
				data: { jobId: 42, message: 'queued', artifact: { id: 'artifact-1' } }
			})
			.mockResolvedValueOnce({
				data: { id: 'decision-1', confidenceState: 'uncertain' }
			});
		clientMock.PUT.mockResolvedValueOnce({
			data: { id: 'decision-1', reviewState: 'approved' }
		});

		await expect(listMediaComponentSources('media-1')).resolves.toEqual({
			sources: [{ id: 'source-1' }]
		});
		await expect(
			retainMediaComponentSource('media-1', {
				sourceRole: 'baseVideo',
				sourceFilePath: '/library/Movie/Base.mkv'
			})
		).resolves.toEqual({ id: 'source-1' });
		await expect(getMediaComponentSource('media-1', 'source-1')).resolves.toEqual({
			id: 'source-1'
		});
		await expect(releaseMediaComponentSource('media-1', 'source-1')).resolves.toEqual({
			id: 'source-1'
		});
		await expect(
			enqueueMediaComponentExtraction('media-1', 'source-1', {
				streamId: 2,
				streamType: 'audio'
			})
		).resolves.toEqual({ jobId: 42, message: 'queued', artifact: { id: 'artifact-1' } });
		await expect(
			evaluateMediaComponentCompatibility('media-1', 'source-1', {
				baseSourceId: 'base-1'
			})
		).resolves.toEqual({ id: 'decision-1', confidenceState: 'uncertain' });
		await expect(
			reviewMediaComponentCompatibility('media-1', 'source-1', 'decision-1', {
				reviewState: 'approved',
				reason: 'manual sync check passed'
			})
		).resolves.toEqual({ id: 'decision-1', reviewState: 'approved' });
		expect(clientMock.PUT).toHaveBeenLastCalledWith(
			'/media/items/{id}/component-sources/{sourceId}/compatibility/{decisionId}/review',
			{
				params: { path: { id: 'media-1', sourceId: 'source-1', decisionId: 'decision-1' } },
				body: { reviewState: 'approved', reason: 'manual sync check passed' }
			}
		);
	});
});
