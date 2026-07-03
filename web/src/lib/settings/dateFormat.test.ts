import { afterEach, describe, expect, it, vi } from 'vitest';

import {
	formatCompactDateTime,
	formatDate,
	formatDateTime,
	formatDateTimeWithSeconds,
	formatLongDateTime,
	formatShortDate,
	formatShortDateTime,
	formatTimeWithSeconds
} from './dateFormat';

describe('date formatting helpers (SCN-SETTINGS-012)', () => {
	afterEach(() => {
		vi.unstubAllGlobals();
	});

	it('uses browser locales for date-only values without shifting the day', () => {
		vi.stubGlobal('navigator', { languages: ['en-US'] });

		expect(formatDate('2026-01-15')).toBe('Jan 15, 2026');
		expect(formatShortDate('2026-01-15')).toBe('1/15/26');
	});

	it('formats date-time variants and returns invalid input unchanged', () => {
		vi.stubGlobal('navigator', { language: 'en-US' });
		const value = '2026-01-15T12:34:56Z';

		expect(formatDateTime(value)).toContain('Jan 15, 2026');
		expect(formatShortDateTime(value)).toContain('1/15/26');
		expect(formatCompactDateTime(value)).toContain('Jan');
		expect(formatDateTimeWithSeconds(value)).toMatch(/56/);
		expect(formatLongDateTime(value)).toContain('2026');
		expect(formatTimeWithSeconds(value)).toMatch(/56/);
		expect(formatDate('not a date')).toBe('not a date');
	});
});
