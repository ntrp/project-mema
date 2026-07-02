import type { Language, LanguageForm, LanguageRequest, LanguageUpdateRequest } from './types';

export function emptyLanguageForm(): LanguageForm {
	return { code: '', displayName: '', aliasesText: '' };
}

export function languageFormFromLanguage(language: Language): LanguageForm {
	return {
		code: language.code,
		originalCode: language.code,
		displayName: language.displayName,
		aliasesText: language.aliases.join(', ')
	};
}

export function normalizeLanguageForm(form: LanguageForm): LanguageRequest {
	return {
		code: form.code.trim().toUpperCase(),
		displayName: form.displayName.trim(),
		aliases: languageAliasValues(form.aliasesText)
	};
}

export function normalizeLanguageUpdateForm(form: LanguageForm): LanguageUpdateRequest {
	return {
		displayName: form.displayName.trim(),
		aliases: languageAliasValues(form.aliasesText)
	};
}

function languageAliasValues(value: string) {
	return value
		.split(/[,\n]/)
		.map((alias) => alias.trim())
		.filter(Boolean);
}
