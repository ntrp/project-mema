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
	languages: string[];
	quality: string;
	formats: string[];
	timeLeft: string;
	progressValue?: number;
	progressLabel: string;
}

export function activityDisplay(activity: DownloadActivity): ActivityDisplay {
	const progressValue =
		activity.status === 'completed' ? 100 : (activity.progressPercent ?? undefined);
	return {
		year: activity.mediaYear ? String(activity.mediaYear) : yearFromTitle(activity.releaseTitle),
		languages: matchedTokens(activity.releaseTitle, languageTokens),
		quality: matchedTokens(activity.releaseTitle, qualityTokens)[0] ?? '-',
		formats: matchedTokens(activity.releaseTitle, formatTokens),
		timeLeft: '-',
		progressValue,
		progressLabel: typeof progressValue === 'number' ? `${progressValue}%` : 'Waiting'
	};
}

export function cancellable(activity: DownloadActivity) {
	return ['queued', 'grabbed', 'downloading'].includes(activity.status);
}

export function manualImportable(activity: DownloadActivity) {
	return activity.status === 'failed';
}

export function createdLabel(value: string) {
	return new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short'
	}).format(new Date(value));
}

export function releaseGroupFromTitle(value: string) {
	const match = /-([A-Za-z0-9][A-Za-z0-9._-]{1,24})$/.exec(value.trim());
	return match?.[1] ?? '';
}

function yearFromTitle(value: string) {
	return /\b(19|20)\d{2}\b/.exec(value)?.[0] ?? '-';
}

function matchedTokens(value: string, tokens: string[]) {
	return tokens.filter((token) => new RegExp(token.replace('+', '\\+'), 'i').test(value));
}
