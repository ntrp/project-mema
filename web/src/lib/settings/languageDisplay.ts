const aliases: Record<string, string> = {
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
	es: 'es',
	spa: 'es',
	spanish: 'es',
	ja: 'ja',
	jpn: 'ja',
	japanese: 'ja',
	ko: 'ko',
	kor: 'ko',
	korean: 'ko',
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
