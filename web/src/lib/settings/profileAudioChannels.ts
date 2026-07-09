const audioChannelAliases: Record<string, string> = {
	mono: '1.0',
	'1': '1.0',
	'10': '1.0',
	'10mono': '1.0',
	'1ch': '1.0',
	'1channel': '1.0',
	'1channels': '1.0',
	stereo: '2.0',
	'2': '2.0',
	'20': '2.0',
	'20stereo': '2.0',
	'2ch': '2.0',
	'2channel': '2.0',
	'2channels': '2.0',
	'3': '3.0',
	'30': '3.0',
	'30ch': '3.0',
	'4': '4.0',
	'40': '4.0',
	'40ch': '4.0',
	'5': '5.0',
	'50': '5.0',
	'50ch': '5.0',
	'51': '5.1',
	'51surround': '5.1',
	'61': '6.1',
	'61surround': '6.1',
	'71': '7.1',
	'71surround': '7.1'
};

export function uniqueAudioChannels(values: string[]) {
	return uniqueTrimmed(values)
		.map(normalizedAudioChannel)
		.filter((value): value is string => !!value)
		.filter((value, index, all) => all.indexOf(value) === index);
}

function normalizedAudioChannel(value: string) {
	return audioChannelAliases[value.trim().toLowerCase().replace(/[^a-z0-9]/g, '')];
}

function uniqueTrimmed(values: string[]) {
	const seen = new Set<string>();
	const result: string[] = [];
	for (const value of values) {
		const text = value.trim();
		if (!text || seen.has(text)) continue;
		seen.add(text);
		result.push(text);
	}
	return result;
}
