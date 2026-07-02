import type { Language } from './types';

export interface LanguageOption {
	id: string;
	code: string;
	label: string;
	displayLabel: string;
}

const commonCodes = ['EN', 'ES', 'FR', 'DE', 'JA'];

export function languageOptionsFromCatalog(languages: Language[]): LanguageOption[] {
	return languages
		.map((language) => ({
			id: language.code,
			code: language.code,
			label: language.displayName,
			displayLabel: `${language.displayName} (${language.code})`
		}))
		.toSorted((left, right) => {
			const leftCommon = commonCodes.indexOf(left.code);
			const rightCommon = commonCodes.indexOf(right.code);
			if (leftCommon !== rightCommon) {
				if (leftCommon === -1) return 1;
				if (rightCommon === -1) return -1;
				return leftCommon - rightCommon;
			}
			return left.label.localeCompare(right.label);
		});
}

export function languageLabelFromCatalog(id: string, languages: Language[]) {
	const language = findLanguage(id, languages);
	return language?.displayName ?? id;
}

export function languageCodeFromValue(value: string, languages: Language[]) {
	return findLanguage(value, languages)?.code ?? '';
}

export function profileLanguageOptions(
	languages: Language[],
	selectedIds: string[]
): LanguageOption[] {
	const options = languageOptionsFromCatalog(languages);
	const known = new Set(options.map((option) => option.id.toLowerCase()));
	for (const id of selectedIds) {
		if (known.has(id.toLowerCase())) {
			continue;
		}
		const language = findLanguage(id, languages);
		options.push({
			id,
			code: language?.code ?? id,
			label: language?.displayName ?? id,
			displayLabel: language ? `${language.displayName} (${id})` : id
		});
		known.add(id.toLowerCase());
	}
	return options;
}

function findLanguage(value: string, languages: Language[]) {
	const normalized = normalize(value);
	if (!normalized) {
		return undefined;
	}
	return languages.find((language) => {
		if (normalize(language.code) === normalized || normalize(language.displayName) === normalized) {
			return true;
		}
		return language.aliases.some((alias) => normalize(alias) === normalized);
	});
}

function normalize(value: string) {
	return value
		.trim()
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '');
}
