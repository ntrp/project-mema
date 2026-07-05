export function qualityInfo(value: string) {
	return matchToken(value, ['2160p', '1080p', '720p', '576p', '480p']);
}

export function audioInfo(value: string) {
	const tokens = ['TrueHD', 'Atmos', 'DTS-HD', 'DTS', 'DDP', 'DD+', 'EAC3', 'AC3', 'AAC']
		.map((token) => matchToken(value, [token]))
		.filter((token) => token !== '-');
	return tokens.join(' ') || '-';
}

export function matchToken(value: string, tokens: string[]) {
	return tokens.find((token) => new RegExp(token, 'i').test(value)) ?? '-';
}
