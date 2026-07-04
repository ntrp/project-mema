import { describe, expect, it } from 'vitest';

import {
	formatBitRate,
	mediaFileInfoSections
} from '$lib/components/app/media/files/preview/mediaFilePreviewInfo';

describe('media file preview info helpers', () => {
	it('formats bitrates for stream stats', () => {
		expect(formatBitRate('4640000')).toBe('4.64 Mbps');
		expect(formatBitRate('64000')).toBe('64.0 Kbps');
		expect(formatBitRate(undefined)).toBe('-');
	});

	it('builds video and audio sections with codec actions', () => {
		const sections = mediaFileInfoSections(
			{
				streamingMode: 'transcode',
				outputVideoCodec: 'libx264',
				outputAudioCodec: 'aac',
				liveBitRate: '4640000',
				videoTrack: {
					type: 'video',
					codec: 'hevc',
					width: 1920,
					height: 1080,
					frameRate: '24000/1001'
				},
				audioTrack: {
					type: 'audio',
					codec: 'dts',
					language: 'eng',
					channels: 2
				}
			},
			{
				playing: true,
				variableBitRate: true,
				liveBitRate: '5500000',
				activeSubtitleLabel: 'English / eng'
			}
		);

		expect(sections).toHaveLength(3);
		expect(sections[0]).toMatchObject({ key: 'video', action: 'transcode' });
		expect(sections[0].rows).toContainEqual({ label: 'Codec', value: 'hevc -> h264' });
		expect(sections[0].rows).toContainEqual({ label: 'Resolution', value: '1920x1080' });
		expect(sections[0].rows).toContainEqual({ label: 'Frame rate', value: '23.976 fps' });
		expect(sections[0].rows).toContainEqual({ label: 'Bitrate', value: '5.50 Mbps' });
		expect(sections[1]).toMatchObject({ key: 'audio', action: 'transcode' });
		expect(sections[1].rows).toContainEqual({ label: 'Codec', value: 'dts -> aac' });
		expect(sections[1].rows).toContainEqual({ label: 'Channels', value: '2' });
		expect(sections[2].rows).toContainEqual({ label: 'Track', value: 'English / eng' });
	});
});
