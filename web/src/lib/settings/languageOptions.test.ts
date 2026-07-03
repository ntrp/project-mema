import { describe, expect, it } from 'vitest';

import { languageLabel, targetLanguageOptions } from './languageOptions';

describe('language option catalog (SCN-SETTINGS-010)', () => {
	it('puts common target languages first with display codes', () => {
		expect(targetLanguageOptions.slice(0, 5)).toEqual([
			expect.objectContaining({ id: 'english', label: 'English', code: 'EN' }),
			expect.objectContaining({ id: 'spanish', label: 'Spanish', code: 'ES' }),
			expect.objectContaining({ id: 'french', label: 'French', code: 'FR' }),
			expect.objectContaining({ id: 'german', label: 'German', code: 'DE' }),
			expect.objectContaining({ id: 'japanese', label: 'Japanese', code: 'JA' })
		]);
		expect(targetLanguageOptions[0].displayLabel).toBe('English (EN)');
	});

	it('keeps less common languages searchable by stable ids', () => {
		const option = targetLanguageOptions.find((item) => item.id === 'norwegian-bokmal');

		expect(option).toEqual(
			expect.objectContaining({
				label: 'Norwegian Bokmal',
				code: 'NB',
				displayLabel: 'Norwegian Bokmal (NB)'
			})
		);
		expect(languageLabel('norwegian-bokmal')).toBe('Norwegian Bokmal');
	});

	it('falls back to the original value for unknown labels', () => {
		expect(languageLabel('invented-language')).toBe('invented-language');
	});
});
