import { describe, expect, it } from 'vitest';

import {
	languageCodeFromValue,
	languageLabelFromCatalog,
	languageOptionsFromCatalog,
	profileLanguageOptions
} from './languageCatalog';
import type { Language } from './types';

const catalog: Language[] = [
	{ code: 'JA', displayName: 'Japanese', aliases: ['Nihongo'] },
	{ code: 'DE', displayName: 'German', aliases: ['Deutsch'] },
	{ code: 'EN', displayName: 'English', aliases: ['ENG'] },
	{ code: 'EO', displayName: 'Esperanto', aliases: [] }
] as Language[];

describe('language catalog selectors (SCN-SETTINGS-010)', () => {
	it('orders common languages by the preferred settings order', () => {
		expect(languageOptionsFromCatalog(catalog).map((option) => option.code)).toEqual([
			'EN',
			'DE',
			'JA',
			'EO'
		]);
	});

	it('resolves labels and codes from ids, names, and aliases', () => {
		expect(languageLabelFromCatalog('deutsch', catalog)).toBe('German');
		expect(languageCodeFromValue('Nihongo', catalog)).toBe('JA');
		expect(languageCodeFromValue(' esperanto ', catalog)).toBe('EO');
	});

	it('keeps unknown selected profile languages visible', () => {
		const options = profileLanguageOptions(catalog, ['custom-language']);

		expect(options.at(-1)).toEqual({
			id: 'custom-language',
			code: 'custom-language',
			label: 'custom-language',
			displayLabel: 'custom-language'
		});
	});
});
