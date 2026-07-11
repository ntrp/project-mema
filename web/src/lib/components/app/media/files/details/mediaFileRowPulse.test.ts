import { afterEach, describe, expect, it, vi } from 'vitest';
import { createMediaFileRowPulse } from './mediaFileRowPulse';
import type { MediaFileDetailRow } from '../mediaFileDetails';

describe('media file row pulse tracking', () => {
	afterEach(() => {
		vi.useRealTimers();
		vi.unstubAllGlobals();
	});

	it('pulses only existing rows whose visible signature changes', () => {
		let expire = () => {};
		vi.stubGlobal('window', {
			setTimeout: (callback: () => void) => {
				expire = callback;
				return 1;
			},
			clearTimeout: vi.fn()
		});
		const track = createMediaFileRowPulse();
		const initial = row('audio', 'Ready');
		expect(track([initial]).has('audio')).toBe(false);
		expect(track([{ ...initial, statusLabel: 'Processing' }]).has('audio')).toBe(true);
		expire();
		expect(track([{ ...initial, statusLabel: 'Processing' }]).has('audio')).toBe(false);
	});
});

function row(key: string, statusLabel: string): MediaFileDetailRow {
	return {
		key,
		description: 'English audio',
		statusLabel,
		details: []
	} as unknown as MediaFileDetailRow;
}
