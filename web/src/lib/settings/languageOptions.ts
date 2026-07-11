import { languageCodes } from './languageOptionCodes';
import { languageNames } from './languageOptionNames';

const commonLanguageOrder = ['english', 'spanish', 'french', 'german', 'japanese'];

export const targetLanguageOptions = languageNames
	.map((label) => ({
		id: languageId(label),
		label,
		code: languageCode(label),
		displayLabel: `${label} (${languageCode(label)})`
	}))
	.sort((left, right) => {
		const leftCommon = commonLanguageOrder.indexOf(left.id);
		const rightCommon = commonLanguageOrder.indexOf(right.id);
		if (leftCommon !== rightCommon) {
			if (leftCommon === -1) return 1;
			if (rightCommon === -1) return -1;
			return leftCommon - rightCommon;
		}
		return left.label.localeCompare(right.label);
	});

export function languageLabel(id: string) {
	return targetLanguageOptions.find((option) => option.id === id)?.label ?? id;
}

function languageId(label: string) {
	return label.toLowerCase().replace(/[^a-z0-9]+/g, '-');
}

function languageCode(label: string) {
	return languageCodes[label] ?? label.slice(0, 2).toUpperCase();
}
