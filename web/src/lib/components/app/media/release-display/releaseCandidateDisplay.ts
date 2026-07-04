import type { ReleaseCandidate } from '$lib/settings/types';

export type MatchSeverity = 'info' | 'warning' | 'error';
export type MatchInfo = ReleaseCandidate['match'];

export function releaseSource(release: ReleaseCandidate) {
	return release.indexerProtocol === 'torrent' ? 'torrent' : 'usenet';
}

export function releaseSourceBadgeClass(release: ReleaseCandidate) {
	return releaseSource(release) === 'torrent'
		? 'border-emerald-300 bg-emerald-100 text-emerald-800'
		: 'border-sky-300 bg-sky-100 text-sky-800';
}

export function ageLabel(release: ReleaseCandidate) {
	if (!release.publishedAt) return '-';
	const published = new Date(release.publishedAt).getTime();
	if (!Number.isFinite(published)) return '-';
	const minutes = Math.max(1, Math.floor((Date.now() - published) / 60000));
	if (minutes < 60) return `${minutes}m`;
	const hours = Math.floor(minutes / 60);
	if (hours < 48) return `${hours}h`;
	const days = Math.floor(hours / 24);
	if (days < 365) return `${days}d`;
	return `${Math.floor(days / 365)}y`;
}

export function sizeLabel(sizeBytes: number) {
	if (!sizeBytes) return '-';
	const gib = sizeBytes / 1024 / 1024 / 1024;
	return `${gib.toFixed(gib >= 10 ? 0 : 1)} GiB`;
}

export function peerLabel(release: ReleaseCandidate) {
	if (release.seeders === undefined && release.peers === undefined) return '-';
	return `${release.peers ?? '-'} / ${release.seeders ?? '-'}`;
}

export function peerBadgeClass(release: ReleaseCandidate) {
	if ((release.seeders ?? 0) > 0) {
		return 'border-emerald-300 bg-emerald-100 text-emerald-800';
	}
	if ((release.peers ?? 0) > 5) {
		return 'border-yellow-300 bg-yellow-100 text-yellow-800';
	}
	return 'border-red-300 bg-red-100 text-red-800';
}

export function languageLabels(release: ReleaseCandidate) {
	return release.match.languages;
}

export function signedScore(score: number) {
	return score > 0 ? `+${score}` : String(score);
}

export function qualityMatch(release: ReleaseCandidate) {
	return { label: release.match.quality || 'Unknown', score: release.match.score };
}
