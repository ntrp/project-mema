import { describe, expect, it } from 'vitest';
import { selectedFirst } from './multiSelectOrdering';

describe('selectedFirst', () => {
	it('SCN-SETTINGS-001 keeps selected options first without changing relative order', () => {
		const options = [
			{ code: 'de', label: 'German' },
			{ code: 'en', label: 'English' },
			{ code: 'ja', label: 'Japanese' }
		];

		const ordered = selectedFirst(options, new Set(['en', 'ja']), (option) => option.code);

		expect(ordered.map((option) => option.code)).toEqual(['en', 'ja', 'de']);
	});
});
