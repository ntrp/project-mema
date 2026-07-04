import type { MediaItem, ReleaseCandidate } from '$lib/settings/types';
import {
	compareReleaseCandidates,
	defaultQualityOptions
} from '$lib/components/app/media/release-display/releaseSortComparators';
import {
	qualityMatch,
	releaseSource
} from '$lib/components/app/media/release-display/releaseCandidateDisplay';

export type ReleaseSortKey =
	| 'source'
	| 'indexer'
	| 'age'
	| 'title'
	| 'size'
	| 'peers'
	| 'quality'
	| 'score';
export type ReleaseSortDirection = 'asc' | 'desc';
export type ReleaseSourceFilter = 'all' | 'usenet' | 'torrent';

export interface ReleaseFilters {
	source: ReleaseSourceFilter;
	minSize: string;
	maxSize: string;
	minScore: string;
	maxScore: string;
	quality: string;
}

export interface ReleaseSort {
	key?: ReleaseSortKey;
	direction: ReleaseSortDirection;
}

export function defaultReleaseFilters(): ReleaseFilters {
	return { source: 'all', minSize: '', maxSize: '', minScore: '', maxScore: '', quality: 'all' };
}

export function activeFilterCount(filters: ReleaseFilters) {
	const defaults = defaultReleaseFilters();
	return (Object.keys(defaults) as (keyof ReleaseFilters)[]).filter(
		(key) => filters[key] !== defaults[key]
	).length;
}

export function releaseQualityOptions(releases: ReleaseCandidate[]) {
	const options = new Set(defaultQualityOptions);
	for (const release of releases) {
		options.add(qualityMatch(release).label);
	}
	return [...options];
}

export function filteredSortedReleases(
	_item: MediaItem,
	releases: ReleaseCandidate[],
	filters: ReleaseFilters,
	sort: ReleaseSort
) {
	const filtered = releases.filter((release) => matchesFilters(release, filters));
	return [...filtered].sort((left, right) => compareReleaseCandidates(left, right, sort));
}

function matchesFilters(release: ReleaseCandidate, filters: ReleaseFilters) {
	const score = qualityMatch(release).score;
	const sizeGiB = release.sizeBytes / 1024 / 1024 / 1024;
	if (filters.source !== 'all' && releaseSource(release) !== filters.source) return false;
	if (filters.quality !== 'all' && qualityMatch(release).label !== filters.quality) return false;
	if (!numberAtLeast(sizeGiB, filters.minSize)) return false;
	if (!numberAtMost(sizeGiB, filters.maxSize)) return false;
	if (!numberAtLeast(score, filters.minScore)) return false;
	if (!numberAtMost(score, filters.maxScore)) return false;
	return true;
}

function numberAtLeast(value: number, minimum: string) {
	const parsed = Number(minimum);
	return minimum.trim() === '' || !Number.isFinite(parsed) || value >= parsed;
}

function numberAtMost(value: number, maximum: string) {
	const parsed = Number(maximum);
	return maximum.trim() === '' || !Number.isFinite(parsed) || value <= parsed;
}
