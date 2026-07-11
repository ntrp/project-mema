import { describe, expect, it } from 'vitest';
import { probeFromMediaPreview } from './dlnaMediaProbe';
import type { MediaFilePreviewInfo } from '$lib/settings/types';

describe('DLNA media probe helper', () => {
	it('maps preview info to trace probe fields', () => {
		const previewInfo: MediaFilePreviewInfo = {
			streamingMode: 'direct',
			deliveryProtocol: 'file',
			outputVideoCodec: 'copy',
			outputAudioCodec: 'copy',
			containerFormatName: 'mov,mp4,m4a,3gp,3g2,mj2',
			videoTrack: {
				type: 'video',
				codec: 'h264',
				height: 1080
			},
			audioTrack: {
				type: 'audio',
				codec: 'aac'
			}
		};

		expect(probeFromMediaPreview(previewInfo)).toEqual({
			container: 'mov,mp4,m4a,3gp,3g2,mj2',
			videoCodec: 'h264',
			audioCodec: 'aac',
			height: 1080
		});
	});

	it('returns undefined for empty preview info', () => {
		expect(probeFromMediaPreview(undefined)).toBeUndefined();
	});
});
