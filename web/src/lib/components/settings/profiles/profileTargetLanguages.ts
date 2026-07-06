import {
	languageCodeFromValue,
	languageLabelFromCatalog,
	languageOptionsFromCatalog
} from '$lib/settings/languageCatalog';
import type { Language } from '$lib/settings/types';

export function targetLanguageKey(languageId: string, languages: Language[]) {
	return targetLanguageValue(languageId, languages).toLowerCase();
}

export function targetLanguageValue(languageId: string, languages: Language[]) {
	return languageCodeFromValue(languageId, languages) || languageId;
}

export function targetLanguageDisplayLabel(languageId: string, languages: Language[]) {
	const code = languageCodeFromValue(languageId, languages);
	if (!code) return languageId;
	return `${languageLabelFromCatalog(languageId, languages)} (${code})`;
}

export function nextTargetLanguageId(languages: Language[], selected: Set<string>) {
	return languageOptionsFromCatalog(languages).find(
		(option) => !selected.has(targetLanguageKey(option.id, languages))
	)?.id;
}

export function targetLanguageChoices(
	languages: Language[],
	currentLanguageId: string,
	selected: Set<string>
) {
	const currentKey = targetLanguageKey(currentLanguageId, languages);
	const choices = languageOptionsFromCatalog(languages).filter(
		(option) =>
			targetLanguageKey(option.id, languages) !== currentKey &&
			!selected.has(targetLanguageKey(option.id, languages))
	);
	const currentId = targetLanguageValue(currentLanguageId, languages);
	if (!choices.some((option) => option.id === currentId)) {
		choices.unshift({
			id: currentId,
			code: currentId,
			label: languageLabelFromCatalog(currentLanguageId, languages),
			displayLabel: targetLanguageDisplayLabel(currentLanguageId, languages)
		});
	}
	return choices;
}
