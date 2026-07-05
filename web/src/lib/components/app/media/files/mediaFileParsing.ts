export function qualityInfo(value: string) {
	const resolution = matchToken(value, ['2160p', '1080p', '720p', '576p', '480p']);
	if (resolution === '-') return '-';
	const source = qualitySource(value);
	return source === '-' ? resolution : `${source}-${resolution}`;
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

function qualitySource(value: string) {
	if (/\bweb[ ._-]?dl\b/i.test(value)) return 'WEBDL';
	if (/\bweb[ ._-]?rip\b/i.test(value)) return 'WEBRip';
	if (/\bblu[ ._-]?ray(?:[ ._-]?rip)?\b|\bbrrip\b/i.test(value)) return 'BluRay';
	if (/\bremux\b/i.test(value)) return 'Remux';
	if (/\bhdtv\b/i.test(value)) return 'HDTV';
	return '-';
}
