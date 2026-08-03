import type { DownloadActivity } from '$lib/settings/types';

const qualityTokens = ['2160p', '1080p', '720p', '576p', '480p'];
const formatTokens = [
	'Remux',
	'BluRay',
	'WEB-DL',
	'WEBDL',
	'WEBRip',
	'HDTV',
	'HDR',
	'DV',
	'Atmos',
	'TrueHD',
	'DTS'
];
const languageTokens = ['Multi', 'Dual', 'English', 'German', 'Japanese', 'Spanish', 'French'];

export interface ActivityDisplay {
	year: string;
	seriesEpisode: string;
	episodeTitle: string;
	languages: string[];
	quality: string;
	formats: string[];
	timeLeft: string;
	progressValue: number;
	progressLabel: string;
}

export function activityDisplay(activity: DownloadActivity): ActivityDisplay {
	const progressValue = activity.status === 'completed' ? 100 : (activity.progressPercent ?? 0);
	const episode =
		activity.mediaType === 'serie' ? episodeFromTitle(activity.releaseTitle) : undefined;
	return {
		year: activity.mediaYear ? String(activity.mediaYear) : yearFromTitle(activity.releaseTitle),
		seriesEpisode: episode ? `${Number(episode[1])}x${Number(episode[2])}` : '',
		episodeTitle: episode
			? episodeTitleFromRelease({ value: activity.releaseTitle, match: episode })
			: '',
		languages: matchedTokens(activity.releaseTitle, languageTokens),
		quality: matchedTokens(activity.releaseTitle, qualityTokens)[0] ?? '-',
		formats: matchedTokens(activity.releaseTitle, formatTokens),
		timeLeft: '-',
		progressValue,
		progressLabel: `${progressValue}%`
	};
}

export function cancellable(activity: DownloadActivity) {
	return ['queued', 'grabbed', 'downloading'].includes(activity.status);
}

export function manualImportable(activity: DownloadActivity) {
	return activity.status === 'failed' && activity.failureType === 'import';
}

export function deletable(activity: DownloadActivity) {
	return activity.status === 'failed' || activity.status === 'cancelled';
}

export function retryable(activity: DownloadActivity) {
	return activity.status === 'failed' || activity.status === 'cancelled';
}

export function statusTone(status: DownloadActivity['status']) {
	return status === 'grabbed' || status === 'completed'
		? 'success'
		: status === 'queued' || status === 'downloading'
			? 'pending'
			: status === 'failed' || status === 'cancelled'
				? 'error'
				: 'muted';
}

export function releaseGroupFromTitle(value: string) {
	const match = /-([A-Za-z0-9][A-Za-z0-9._]{1,24})$/.exec(value.trim());
	return match?.[1] ?? '';
}

function yearFromTitle(value: string) {
	return /\b(19|20)\d{2}\b/.exec(value)?.[0] ?? '';
}

function episodeFromTitle(value: string) {
	return /\bS(\d{1,2})E(\d{1,3})\b/i.exec(value);
}

function episodeTitleFromRelease(input: { value: string; match: RegExpExecArray }) {
	const suffix = input.value
		.slice((input.match.index ?? 0) + input.match[0].length)
		.replace(/^[._ -]+/, '');
	const stopTokens = [...qualityTokens, ...formatTokens, ...languageTokens]
		.map((token) => token.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'))
		.join('|');
	const stop = new RegExp(`(?:^|[._ -])(?:${stopTokens})(?:[._ -]|$)`, 'i').exec(suffix);
	return suffix
		.slice(0, stop?.index ?? suffix.length)
		.replace(/[._]+/g, ' ')
		.trim();
}

function matchedTokens(value: string, tokens: string[]) {
	return tokens.filter((token) => new RegExp(token.replace('+', '\\+'), 'i').test(value));
}
