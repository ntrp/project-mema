import { displayLanguage } from '$lib/settings/languageDisplay';

const fileLanguagePattern =
	/\b(en|eng|english|de|ger|deu|german|fr|fre|fra|french|es|spa|spanish|ja|jpn|japanese|ko|kor|korean|zh|zho|chi|chinese)\b/i;

export function mediaFileLanguageInfo(value: string) {
	if (/\bmulti\b/i.test(value)) return 'Multi';
	if (/\bdual\b/i.test(value)) return 'Dual';
	const match = fileLanguagePattern.exec(value);
	return match ? displayLanguage(match[1]) : '-';
}
