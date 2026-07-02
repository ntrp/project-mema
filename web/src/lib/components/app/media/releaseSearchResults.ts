import type { MediaItem, ReleaseCandidate } from '$lib/settings/types';
import { qualityMatch, releaseSource } from './releaseCandidateDisplay';

export type ReleaseSortKey =
	| 'source'
	| 'indexer'
	| 'age'
	| 'title'
	| 'size'
	| 'peers'
	| 'quality'
	| 'score'
	| 'match';
export type ReleaseSortDirection = 'asc' | 'desc';
export type ReleaseSourceFilter = 'all' | 'nzb' | 'torrent';

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

const defaultQualityOptions = [
	'Unknown',
	'WORKPRINT',
	'CAM',
	'TELESYNC',
	'TELECINE',
	'REGIONAL',
	'DVDSCR',
	'SDTV',
	'DVD',
	'DVD-R',
	'WEBDL-480p',
	'WEBRip-480p',
	'Bluray-480p',
	'Bluray-576p',
	'HDTV-720p',
	'WEBDL-720p',
	'WEBRip-720p',
	'Bluray-720p',
	'HDTV-1080p',
	'WEBDL-1080p',
	'WEBRip-1080p',
	'Bluray-1080p',
	'Remux-1080p',
	'HDTV-2160p',
	'WEBDL-2160p',
	'WEBRip-2160p',
	'Bluray-2160p',
	'Remux-2160p',
	'BR-DISK',
	'Raw-HD'
];

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
	item: MediaItem,
	releases: ReleaseCandidate[],
	filters: ReleaseFilters,
	sort: ReleaseSort
) {
	const filtered = releases.filter((release) => matchesFilters(release, filters));
	return [...filtered].sort((left, right) => compareReleases(item, left, right, sort));
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

function compareReleases(
	item: MediaItem,
	left: ReleaseCandidate,
	right: ReleaseCandidate,
	sort: ReleaseSort
) {
	const severityResult =
		matchSeverityRank(right.match.severity) - matchSeverityRank(left.match.severity);
	if (severityResult !== 0) return severityResult;
	if (!sort.key) return 0;
	const result = compareValues(sortValue(item, left, sort.key), sortValue(item, right, sort.key));
	return sort.direction === 'asc' ? result : -result;
}

function sortValue(item: MediaItem, release: ReleaseCandidate, key?: ReleaseSortKey) {
	switch (key) {
		case 'source':
			return releaseSource(release);
		case 'indexer':
			return release.indexerName;
		case 'age':
			return release.publishedAt ? new Date(release.publishedAt).getTime() : 0;
		case 'title':
			return release.title;
		case 'size':
			return release.sizeBytes;
		case 'peers':
			return release.peers ?? release.seeders ?? -1;
		case 'quality':
			return qualityMatch(release).label;
		case 'score':
			return qualityMatch(release).score;
		case 'match':
			return matchSeverityRank(release.match.severity);
		default:
			return '';
	}
}

function compareValues(left: number | string, right: number | string) {
	if (typeof left === 'number' && typeof right === 'number') {
		return left - right;
	}
	return String(left).localeCompare(String(right), undefined, {
		numeric: true,
		sensitivity: 'base'
	});
}

function numberAtLeast(value: number, minimum: string) {
	const parsed = Number(minimum);
	return minimum.trim() === '' || !Number.isFinite(parsed) || value >= parsed;
}

function numberAtMost(value: number, maximum: string) {
	const parsed = Number(maximum);
	return maximum.trim() === '' || !Number.isFinite(parsed) || value <= parsed;
}

function matchSeverityRank(severity: string) {
	if (severity === 'error') return 0;
	if (severity === 'warning') return 1;
	return 2;
}
