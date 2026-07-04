import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import type { MediaFileTextTrack } from '$lib/components/app/media/files/preview/mediaFilePlayback';

interface ChapterOption {
	key: string;
	title: string;
	time: string;
	seconds: number;
}

export function mediaFileChapterTrack(row: MediaFileRow): MediaFileTextTrack | undefined {
	const chapters = mediaFileChapterOptions(row);
	if (chapters.length === 0) return undefined;
	return {
		key: 'chapters',
		kind: 'chapters',
		label: 'Chapters',
		src: webVttDataUrl(chaptersWebVtt(chapters)),
		default: true
	};
}

function mediaFileChapterOptions(row: MediaFileRow): ChapterOption[] {
	return row.chapters
		.map((chapter, index): ChapterOption | undefined => {
			const seconds = chapterSeconds(chapter.startTime);
			if (seconds === undefined) return undefined;
			return {
				key: `chapter-${chapter.index}-${index}`,
				title: chapter.title?.trim() || `Chapter ${chapter.index + 1 || index + 1}`,
				time: formatPlaybackTime(seconds),
				seconds
			};
		})
		.filter((chapter): chapter is ChapterOption => Boolean(chapter));
}

function chapterSeconds(value?: string) {
	const trimmed = value?.trim();
	if (!trimmed) return undefined;
	const numeric = Number(trimmed);
	if (Number.isFinite(numeric)) return numeric;
	const parts = trimmed.split(':').map(Number);
	if (parts.some((part) => !Number.isFinite(part))) return undefined;
	return parts.reduce((total, part) => total * 60 + part, 0);
}

function formatPlaybackTime(seconds: number) {
	if (!Number.isFinite(seconds) || seconds < 0) return '0:00';
	const rounded = Math.floor(seconds);
	const hours = Math.floor(rounded / 3600);
	const minutes = Math.floor((rounded % 3600) / 60);
	const rest = String(rounded % 60).padStart(2, '0');
	return hours > 0 ? `${hours}:${String(minutes).padStart(2, '0')}:${rest}` : `${minutes}:${rest}`;
}

function chaptersWebVtt(chapters: ChapterOption[]) {
	const lines = ['WEBVTT', ''];
	for (const [index, chapter] of chapters.entries()) {
		const end = chapters[index + 1]?.seconds ?? chapter.seconds + 1;
		if (end <= chapter.seconds) continue;
		lines.push(
			`${webVttTime(chapter.seconds)} --> ${webVttTime(end)}`,
			webVttText(chapter.title),
			''
		);
	}
	return lines.join('\n');
}

function webVttTime(seconds: number) {
	const safe = Math.max(0, seconds);
	const hours = Math.floor(safe / 3600);
	const minutes = Math.floor((safe % 3600) / 60);
	const wholeSeconds = Math.floor(safe % 60);
	const millis = Math.floor((safe - Math.floor(safe)) * 1000);
	return [
		String(hours).padStart(2, '0'),
		String(minutes).padStart(2, '0'),
		`${String(wholeSeconds).padStart(2, '0')}.${String(millis).padStart(3, '0')}`
	].join(':');
}

function webVttDataUrl(content: string) {
	return `data:text/vtt;charset=utf-8,${encodeURIComponent(content)}`;
}

function webVttText(value: string) {
	return value.replaceAll('\r', '').replaceAll('-->', '->').trim() || 'Chapter';
}
