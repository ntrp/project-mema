import type { LibraryFolder, MediaSearchResult, QualityProfileOption } from '$lib/settings/types';

type SmartMediaCandidate = MediaSearchResult & {
	genres?: string[];
	originalLanguage?: string;
};

export function preselectQualityProfileId(
	candidate: MediaSearchResult,
	qualityProfiles: QualityProfileOption[]
) {
	return bestScored(qualityProfiles, (profile) => profileScore(candidate, profile))?.id ?? '';
}

export function preselectLibraryFolderId(
	candidate: MediaSearchResult,
	libraryFolders: LibraryFolder[]
) {
	return bestScored(libraryFolders, (folder) => folderScore(candidate, folder))?.id ?? '';
}

export function mediaPosterUrl(path?: string) {
	if (!path) {
		return undefined;
	}
	if (path.startsWith('http://') || path.startsWith('https://')) {
		return path;
	}
	return `https://image.tmdb.org/t/p/w780${path}`;
}

function profileScore(candidate: MediaSearchResult, profile: QualityProfileOption) {
	const text = normalizedText(`${profile.id} ${profile.name}`);
	let score = 0;
	if (isAnimeCandidate(candidate)) {
		score += hasAny(text, ['anime']) ? 50 : 0;
	} else if (hasAny(text, ['anime'])) {
		score -= 20;
	}
	score += hasAny(text, ['1080', '1080p']) ? 40 : 0;
	score += hasAny(text, ['2160', '2160p', '4k', 'uhd']) ? 25 : 0;
	score += hasAny(text, ['default', 'best']) ? 15 : 0;
	score -= hasAny(text, ['any acceptable', 'any']) ? 10 : 0;
	return score;
}

function folderScore(candidate: MediaSearchResult, folder: LibraryFolder) {
	const path = normalizedText(folder.path);
	const hasAnime = hasAny(path, ['anime']);
	const hasMovie = hasAny(path, ['movie', 'movies', 'film', 'films']);
	const hasSeries = hasAny(path, ['series', 'tv', 'show', 'shows']);
	if (candidate.type === 'series') {
		return (
			(hasSeries ? 100 : 0) +
			(hasAnime && isAnimeCandidate(candidate) ? 20 : 0) -
			(hasMovie ? 25 : 0)
		);
	}
	if (isAnimeCandidate(candidate)) {
		return (
			(hasAnime && hasMovie ? 120 : 0) +
			(hasAnime && !hasMovie ? 90 : 0) +
			(!hasAnime && hasMovie ? 60 : 0) -
			(hasSeries ? 25 : 0)
		);
	}
	return (hasMovie ? 100 : 0) - (hasAnime ? 15 : 0) - (hasSeries ? 25 : 0);
}

function isAnimeCandidate(candidate: MediaSearchResult) {
	const smartCandidate = candidate as SmartMediaCandidate;
	const genres = smartCandidate.genres ?? [];
	const text = normalizedText(
		[candidate.title, candidate.overview, genres.join(' '), smartCandidate.originalLanguage].join(
			' '
		)
	);
	const isJapaneseAnimation =
		smartCandidate.originalLanguage?.toLowerCase() === 'ja' &&
		genres.some((genre) => normalizedText(genre).includes('animation'));
	return text.includes('anime') || isJapaneseAnimation;
}

function bestScored<T>(items: T[], scoreItem: (item: T) => number) {
	return items
		.map((item, index) => ({ item, score: scoreItem(item), index }))
		.sort((left, right) => right.score - left.score || left.index - right.index)[0]?.item;
}

function normalizedText(value: string) {
	return value.toLowerCase();
}

function hasAny(text: string, needles: string[]) {
	return needles.some((needle) => text.includes(needle));
}
