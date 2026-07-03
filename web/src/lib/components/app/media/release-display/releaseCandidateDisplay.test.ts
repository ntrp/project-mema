import { describe, expect, it, vi } from 'vitest';

import {
	ageLabel,
	languageLabels,
	peerLabel,
	qualityMatch,
	releaseSource,
	releaseSourceBadgeClass,
	signedScore,
	sizeLabel
} from '$lib/components/app/media/release-display/releaseCandidateDisplay';
import type { ReleaseCandidate } from '$lib/settings/types';

function release(overrides: Partial<ReleaseCandidate> = {}): ReleaseCandidate {
	return {
		indexerProtocol: 'torrent',
		sizeBytes: 8 * 1024 ** 3,
		seeders: 12,
		peers: 20,
		publishedAt: '2026-07-03T04:00:00Z',
		match: {
			quality: 'WEBDL-1080p',
			score: 150,
			languages: ['German', 'English']
		},
		...overrides
	} as ReleaseCandidate;
}

describe('release candidate display labels (SCN-MEDIA-002)', () => {
	it('labels release source and source badge by indexer type', () => {
		expect(releaseSource(release({ indexerProtocol: 'torrent' }))).toBe('torrent');
		expect(releaseSource(release({ indexerProtocol: 'usenet' }))).toBe('nzb');
		expect(releaseSourceBadgeClass(release({ indexerProtocol: 'torrent' }))).toContain('emerald');
		expect(releaseSourceBadgeClass(release({ indexerProtocol: 'usenet' }))).toContain('sky');
	});

	it('formats age, size, peer, language, score, and quality labels', () => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date('2026-07-04T04:30:00Z'));

		expect(ageLabel(release())).toBe('24h');
		expect(sizeLabel(8 * 1024 ** 3)).toBe('8.0 GiB');
		expect(sizeLabel(15 * 1024 ** 3)).toBe('15 GiB');
		expect(peerLabel(release())).toBe('20 / 12');
		expect(languageLabels(release())).toEqual(['German', 'English']);
		expect(signedScore(10)).toBe('+10');
		expect(signedScore(0)).toBe('0');
		expect(qualityMatch(release())).toEqual({ label: 'WEBDL-1080p', score: 150 });

		vi.useRealTimers();
	});

	it('uses placeholders for missing or invalid release metadata', () => {
		expect(ageLabel(release({ publishedAt: undefined }))).toBe('-');
		expect(ageLabel(release({ publishedAt: 'not-a-date' }))).toBe('-');
		expect(sizeLabel(0)).toBe('-');
		expect(peerLabel(release({ seeders: undefined, peers: undefined }))).toBe('-');
		expect(peerLabel(release({ seeders: 4, peers: undefined }))).toBe('- / 4');
		expect(
			qualityMatch(release({ match: { quality: '', score: -10 } } as Partial<ReleaseCandidate>))
		).toEqual({
			label: 'Unknown',
			score: -10
		});
	});
});
