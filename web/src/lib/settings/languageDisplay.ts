const aliases: Record<string, string> = {
	ar: 'ar',
	ara: 'ar',
	arabic: 'ar',
	da: 'da',
	dan: 'da',
	danish: 'da',
	en: 'en',
	eng: 'en',
	english: 'en',
	de: 'de',
	deu: 'de',
	ger: 'de',
	german: 'de',
	fr: 'fr',
	fra: 'fr',
	fre: 'fr',
	french: 'fr',
	fi: 'fi',
	fin: 'fi',
	finnish: 'fi',
	hi: 'hi',
	hin: 'hi',
	hindi: 'hi',
	es: 'es',
	spa: 'es',
	spanish: 'es',
	it: 'it',
	ita: 'it',
	italian: 'it',
	ja: 'ja',
	jpn: 'ja',
	japanese: 'ja',
	ko: 'ko',
	kor: 'ko',
	korean: 'ko',
	nl: 'nl',
	nld: 'nl',
	dut: 'nl',
	dutch: 'nl',
	no: 'no',
	nor: 'no',
	norwegian: 'no',
	pl: 'pl',
	pol: 'pl',
	polish: 'pl',
	pt: 'pt',
	por: 'pt',
	portuguese: 'pt',
	ru: 'ru',
	rus: 'ru',
	russian: 'ru',
	sv: 'sv',
	swe: 'sv',
	swedish: 'sv',
	zh: 'zh',
	zho: 'zh',
	chi: 'zh',
	chinese: 'zh'
};

export function displayLanguage(value?: string) {
	const language = normalizedLanguage(value);
	if (!language) return '-';
	try {
		return (
			new Intl.DisplayNames(undefined, { type: 'language' }).of(language) ?? titleCase(language)
		);
	} catch {
		return titleCase(language);
	}
}

export function languageMatchKey(value?: string) {
	return normalizedLanguage(value);
}

function normalizedLanguage(value?: string) {
	const normalized = value
		?.trim()
		.toLowerCase()
		.replace(/\s+language$/, '');
	if (!normalized || normalized === '-') return '';
	return aliases[normalized] ?? normalized;
}

function titleCase(value: string) {
	return value
		.split(/[-_\s]+/)
		.filter(Boolean)
		.map((part) => `${part.charAt(0).toUpperCase()}${part.slice(1)}`)
		.join(' ');
}
