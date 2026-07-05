export function fileName(path: string) {
	return path.replaceAll('\\', '/').split('/').filter(Boolean).pop() ?? path;
}

export function relativePath(root: string | undefined, path: string) {
	if (!root) return fileName(path);
	const normalizedRoot = root.replaceAll('\\', '/').replace(/\/+$/, '');
	const normalizedPath = path.replaceAll('\\', '/');
	return normalizedPath.startsWith(`${normalizedRoot}/`)
		? normalizedPath.slice(normalizedRoot.length + 1)
		: fileName(path);
}

export function episodeParts(path: string) {
	const match = /s(\d{1,2})e(\d{1,3})/i.exec(path);
	if (!match) return {};
	return { seasonNumber: Number(match[1]), episodeNumber: Number(match[2]) };
}
