import { displayLanguage } from '$lib/settings/languageDisplay';

const fileLanguagePattern =
	/\b(ar|ara|arabic|da|dan|danish|de|ger|deu|german|dut|dutch|en|eng|english|es|spa|spanish|fi|fin|finnish|fr|fre|fra|french|hi|hin|hindi|it|ita|italian|ja|jpn|japanese|ko|kor|korean|nl|nld|no|nor|norwegian|pl|pol|polish|por|portuguese|pt|ru|rus|russian|sv|swe|swedish|zh|zho|chi|chinese)\b/i;

export function mediaFileLanguageInfo(value: string) {
	if (/\bmulti\b/i.test(value)) return 'Multi';
	if (/\bdual\b/i.test(value)) return 'Dual';
	const match = fileLanguagePattern.exec(value);
	return match ? displayLanguage(match[1]) : '-';
}
