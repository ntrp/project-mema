import type { ReleaseCandidate } from '$lib/settings/types';
import { qualityMatch } from '$lib/components/app/media/release-display/releaseCandidateDisplay';
import type {
	ReleaseSort,
	ReleaseSortDirection,
	ReleaseSortKey
} from '$lib/components/app/media/release-display/releaseSearchResults';

export const defaultQualityOptions = [
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

const secondarySorts: Required<ReleaseSort>[] = [
	{ key: 'source', direction: 'asc' },
	{ key: 'quality', direction: 'desc' },
	{ key: 'score', direction: 'desc' },
	{ key: 'age', direction: 'asc' }
];

export function compareReleaseCandidates(
	left: ReleaseCandidate,
	right: ReleaseCandidate,
	sort: ReleaseSort
) {
	const severityResult =
		matchSeverityRank(right.match.severity) - matchSeverityRank(left.match.severity);
	if (severityResult !== 0) return severityResult;
	if (!sort.key) return 0;
	const sorters = [
		{ key: sort.key, direction: sort.direction },
		...secondarySorts.filter((candidate) => candidate.key !== sort.key)
	];
	for (const candidate of sorters) {
		const result = compareBySort(left, right, candidate);
		if (result !== 0) return result;
	}
	return 0;
}

function compareBySort(
	left: ReleaseCandidate,
	right: ReleaseCandidate,
	sort: Required<ReleaseSort>
) {
	if (sort.key === 'age') return compareReleaseAge(left, right, sort.direction);
	const result = compareValues(sortValue(left, sort.key), sortValue(right, sort.key));
	return sort.direction === 'asc' ? result : -result;
}

function sortValue(release: ReleaseCandidate, key?: ReleaseSortKey) {
	switch (key) {
		case 'source':
			return protocolRank(release);
		case 'indexer':
			return release.indexerName;
		case 'age':
			return publishedAtTime(release) ?? Number.POSITIVE_INFINITY;
		case 'title':
			return release.title;
		case 'size':
			return release.sizeBytes;
		case 'peers':
			return release.peers ?? release.seeders ?? -1;
		case 'quality':
			return qualityRank(release);
		case 'score':
			return qualityMatch(release).score;
		default:
			return '';
	}
}

function compareReleaseAge(
	left: ReleaseCandidate,
	right: ReleaseCandidate,
	direction: ReleaseSortDirection
) {
	const leftTime = publishedAtTime(left);
	const rightTime = publishedAtTime(right);
	if (leftTime === undefined && rightTime === undefined) return 0;
	if (leftTime === undefined) return 1;
	if (rightTime === undefined) return -1;
	const result = rightTime - leftTime;
	return direction === 'asc' ? result : -result;
}

function protocolRank(release: ReleaseCandidate) {
	if (release.indexerProtocol === 'usenet') return 0;
	if (release.indexerProtocol === 'torrent') return 1;
	return 2;
}

function qualityRank(release: ReleaseCandidate) {
	const rank = defaultQualityOptions.indexOf(qualityMatch(release).label);
	return rank === -1 ? 0 : rank;
}

function publishedAtTime(release: ReleaseCandidate) {
	if (!release.publishedAt) return undefined;
	const time = new Date(release.publishedAt).getTime();
	return Number.isFinite(time) ? time : undefined;
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

function matchSeverityRank(severity: string) {
	if (severity === 'error') return 0;
	if (severity === 'warning') return 1;
	return 2;
}
