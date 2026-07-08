export function uniqueSubtitleFormats(values: string[]) {
	const seen = new Set<string>();
	const result = [];
	for (const value of values) {
		const normalized = normalizeSubtitleFormat(value);
		if (!normalized || seen.has(normalized)) continue;
		seen.add(normalized);
		result.push(normalized);
	}
	return result;
}

export function normalizeSubtitleFormat(value: string) {
	const clean = value.trim().toLowerCase().replace(/^\./, '');
	if (clean === 'srt' || clean === 'subrip') return 'subrip';
	if (clean === 'webvtt') return 'vtt';
	if (clean === 'sup') return 'pgs';
	return clean;
}
