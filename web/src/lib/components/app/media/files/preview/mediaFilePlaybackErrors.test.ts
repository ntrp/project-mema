import { describe, expect, it } from 'vitest';

import { mediaPlaybackErrorMessage } from '$lib/components/app/media/files/preview/mediaFilePlaybackErrors';

describe('media file playback errors', () => {
	it('explains browser decode failures with a fallback hint', () => {
		const message = mediaPlaybackErrorMessage({ error: { code: 3 } });

		expect(message).toContain('could not decode');
		expect(message).toContain('Chrome, Brave, or Edge');
		expect(message).toContain('VLC');
	});

	it('explains unsupported preview streams', () => {
		const message = mediaPlaybackErrorMessage({ error: { code: 4 } });

		expect(message).toContain('does not support this preview stream');
	});
});
