import type { MediaSearchResult, PersonAppearance } from '$lib/settings/types';

export type PersonAppearanceFilter = 'all' | 'movie' | 'series';

export function filteredAppearances(
	appearances: PersonAppearance[],
	filter: PersonAppearanceFilter
) {
	return filter === 'all'
		? appearances
		: appearances.filter((appearance) => appearance.type === filter);
}

export function sortedAppearances(appearances: PersonAppearance[]) {
	return [...appearances].sort((left, right) => {
		const leftDate = appearanceDateValue(left);
		const rightDate = appearanceDateValue(right);
		if (leftDate !== rightDate) return leftDate - rightDate;
		return left.title.localeCompare(right.title);
	});
}

export function appearanceYear(appearance: PersonAppearance) {
	return appearance.year ?? yearFromDate(appearance.releaseDate) ?? 'Unknown';
}

export function appearancesByYear(appearances: PersonAppearance[]) {
	const groups: { year: number | 'Unknown'; appearances: PersonAppearance[] }[] = [];
	for (const appearance of sortedAppearances(appearances)) {
		const year = appearanceYear(appearance);
		const group = groups.find((item) => item.year === year);
		if (group) {
			group.appearances.push(appearance);
		} else {
			groups.push({ year, appearances: [appearance] });
		}
	}
	return groups;
}

export function appearanceResult(appearance: PersonAppearance): MediaSearchResult {
	return {
		title: appearance.title,
		type: appearance.type,
		year: appearance.year,
		externalProvider: appearance.externalProvider,
		externalId: appearance.externalId,
		overview: appearance.overview,
		posterPath: appearance.posterPath,
		backdropPath: appearance.backdropPath,
		releaseDate: appearance.releaseDate
	};
}

export function resultKey(result: MediaSearchResult) {
	return `${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}`;
}

function appearanceDateValue(appearance: PersonAppearance) {
	const date = appearance.releaseDate ? Date.parse(appearance.releaseDate) : NaN;
	if (Number.isFinite(date)) return date;
	return appearance.year ? Date.UTC(appearance.year, 0, 1) : Number.MAX_SAFE_INTEGER;
}

function yearFromDate(value?: string) {
	if (!value) return undefined;
	const year = new Date(value).getUTCFullYear();
	return Number.isFinite(year) ? year : undefined;
}
