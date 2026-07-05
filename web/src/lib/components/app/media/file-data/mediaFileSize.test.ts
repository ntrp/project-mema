import { describe, expect, it } from 'vitest';

import {
	formatBytes,
	mediaFileInfo,
	mediaFileSize
} from '$lib/components/app/media/file-data/mediaFileSize';
import type { MediaItem } from '$lib/settings/types';

describe('media file size helpers (SCN-MEDIA-001)', () => {
	it('formats byte values for display', () => {
		expect(formatBytes(-1)).toBe('-');
		expect(formatBytes(Number.NaN)).toBe('-');
		expect(formatBytes(512)).toBe('512 B');
		expect(formatBytes(1536)).toBe('1.50 KiB');
		expect(formatBytes(12 * 1024 ** 2)).toBe('12.0 MiB');
	});

	it('looks up file info by path', () => {
		const item = {
			files: [{ path: '/movie.mkv', status: 'available', sizeBytes: 1024 }]
		} as MediaItem;

		expect(mediaFileInfo(item, '/movie.mkv')?.sizeBytes).toBe(1024);
		expect(mediaFileInfo(item, '/missing.mkv')).toBeUndefined();
		expect(mediaFileSize(item, '/movie.mkv')).toBe('1.00 KiB');
		expect(mediaFileSize(item, '/missing.mkv')).toBe('-');
	});
});
