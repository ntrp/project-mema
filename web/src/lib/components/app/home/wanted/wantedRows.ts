import type { MediaItem } from '$lib/settings/types';

type MediaTarget = NonNullable<MediaItem['targetSatisfaction']>['targets'][number];
type MediaFile = NonNullable<MediaItem['files']>[number];

export type WantedRowKind = 'media' | 'target' | 'custom_format_upgrade';

export interface WantedDisplayRow {
	key: string;
	kind: WantedRowKind;
	item: MediaItem;
	title: string;
	context: string;
	state: string;
	operation?: string;
	filePath?: string;
}

const wantedTargetStates = new Set(['missing', 'partial', 'pending', 'blocked', 'failed']);

export function wantedDisplayRows(items: MediaItem[]): WantedDisplayRow[] {
	return items.flatMap((item) => [
		...mediaWantedRows(item),
		...targetWantedRows(item),
		...customFormatWantedRows(item)
	]);
}

function mediaWantedRows(item: MediaItem): WantedDisplayRow[] {
	if (item.status !== 'missing' && item.rollup?.state !== 'missing') return [];
	return [
		{
			key: `media:${item.id}`,
			kind: 'media',
			item,
			title: item.title,
			context: item.type === 'serie' ? 'Series media' : 'Movie media',
			state: item.rollup?.state ?? item.status
		}
	];
}

function targetWantedRows(item: MediaItem): WantedDisplayRow[] {
	return (item.targetSatisfaction?.targets ?? [])
		.filter((target) => wantedTargetStates.has(target.state))
		.map((target) => targetWantedRow(item, target));
}

function targetWantedRow(item: MediaItem, target: MediaTarget): WantedDisplayRow {
	const language = target.languageId ? ` / ${target.languageId}` : '';
	const file = fileLabel(mediaFilePath(item, target.mediaFileId));
	return {
		key: `target:${target.id}`,
		kind: 'target',
		item,
		title: item.title,
		context: `${target.type}${language}${file ? ` / ${file}` : ''}`,
		state: target.state,
		operation: target.requiredOperation?.reason,
		filePath: mediaFilePath(item, target.mediaFileId)
	};
}

function customFormatWantedRows(item: MediaItem): WantedDisplayRow[] {
	return (item.files ?? [])
		.filter((file) => file.rollup?.state === 'upgradeable')
		.map((file) => customFormatWantedRow(item, file));
}

function customFormatWantedRow(item: MediaItem, file: MediaFile): WantedDisplayRow {
	return {
		key: `custom-format:${item.id}:${file.path}`,
		kind: 'custom_format_upgrade',
		item,
		title: item.title,
		context: fileLabel(file.path),
		state: 'upgradeable',
		operation: 'Search for a higher-scoring release',
		filePath: file.path
	};
}

function mediaFilePath(item: MediaItem, mediaFileId?: string) {
	if (!mediaFileId) return undefined;
	return item.files?.find((file) => file.path.includes(mediaFileId))?.path;
}

function fileLabel(path?: string) {
	if (!path) return '';
	return path.split('/').filter(Boolean).at(-1) ?? path;
}
