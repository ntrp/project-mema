import { describe, expect, it } from 'vitest';

import {
	emptyLanguageForm,
	languageFormFromLanguage,
	normalizeLanguageForm,
	normalizeLanguageUpdateForm
} from './languageForms';
import type { Language } from './types';

describe('language forms (SCN-SETTINGS-012)', () => {
	it('creates and populates language forms', () => {
		expect(emptyLanguageForm()).toEqual({ code: '', displayName: '', aliasesText: '' });
		expect(
			languageFormFromLanguage({
				code: 'PT-BR',
				displayName: 'Brazilian Portuguese',
				aliases: ['Portuguese BR', 'Português']
			} as Language)
		).toEqual({
			code: 'PT-BR',
			originalCode: 'PT-BR',
			displayName: 'Brazilian Portuguese',
			aliasesText: 'Portuguese BR, Português'
		});
	});

	it('normalizes create and update requests', () => {
		const form = {
			code: ' pt-br ',
			displayName: ' Brazilian Portuguese ',
			aliasesText: 'Portuguese BR,\n Português,  '
		};

		expect(normalizeLanguageForm(form)).toEqual({
			code: 'PT-BR',
			displayName: 'Brazilian Portuguese',
			aliases: ['Portuguese BR', 'Português']
		});
		expect(normalizeLanguageUpdateForm(form)).toEqual({
			displayName: 'Brazilian Portuguese',
			aliases: ['Portuguese BR', 'Português']
		});
	});
});
