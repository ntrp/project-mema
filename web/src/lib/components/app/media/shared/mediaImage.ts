export function mediaImageUrl(path?: string, size = 'w780') {
	if (!path) return undefined;
	if (path.startsWith('http://') || path.startsWith('https://') || path.startsWith('/api/')) {
		return path;
	}
	return `https://image.tmdb.org/t/p/${size}${path}`;
}
