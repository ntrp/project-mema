import type { PersonAppearance } from '$lib/settings/types';
import { appearanceYear, sortedAppearances } from './personDetail';

export interface AppearanceTimelineItem {
	appearance: PersonAppearance;
	cardX: number;
	cardY: number;
	markerX: number;
	top: boolean;
	unreleased: boolean;
}

export interface AppearanceTimelineOptions {
	cardWidth: number;
	cardGap: number;
	cardTopY: number;
	cardBottomY: number;
	emptyYearWidth: number;
	yearEntryWidth: number;
	paddingX: number;
}

export type AppearanceTimelineYear = { year: number; x: number };

export function appearanceTimelineData(
	source: PersonAppearance[],
	options: AppearanceTimelineOptions
) {
	const sorted = sortedAppearances(source);
	const numericYears = uniqueYears(sorted);
	const minYear = numericYears[0] ?? new Date().getFullYear();
	const maxYear = numericYears[numericYears.length - 1] ?? minYear;
	const counts = appearanceCountsByYear(sorted);
	const years = timelineYears(minYear, maxYear, counts, options);
	const lineWidth = options.paddingX + timelineEndX(years, counts, options);
	const items = sorted.map((appearance, index) => {
		const markerX = timelineMarkerX(appearance, years, counts, lineWidth, options);
		const top = index % 2 === 0;
		const unreleased = unreleasedAppearance(appearance);
		return {
			appearance,
			cardX: markerX - options.cardWidth / 2,
			cardY: top ? options.cardTopY : options.cardBottomY,
			markerX,
			top,
			unreleased
		};
	});
	spaceTimelineCards(
		items.filter((item) => item.top),
		options
	);
	spaceTimelineCards(
		items.filter((item) => !item.top),
		options
	);
	const contentWidth = Math.max(
		lineWidth,
		...items.map((item) => item.cardX + options.cardWidth + options.paddingX)
	);
	return { items, years, contentWidth };
}

function uniqueYears(source: PersonAppearance[]) {
	return [
		...new Set(
			source
				.map((appearance) => appearanceYear(appearance))
				.filter((year): year is number => typeof year === 'number')
		)
	];
}

function appearanceCountsByYear(source: PersonAppearance[]) {
	const counts = new Map<number, number>();
	for (const appearance of source) {
		const year = appearanceYear(appearance);
		if (typeof year === 'number') {
			counts.set(year, (counts.get(year) ?? 0) + 1);
		}
	}
	return counts;
}

function timelineYears(
	minYear: number,
	maxYear: number,
	counts: Map<number, number>,
	options: AppearanceTimelineOptions
) {
	const years: AppearanceTimelineYear[] = [];
	let x = options.paddingX;
	for (let year = minYear; year <= maxYear; year += 1) {
		years.push({ year, x });
		x += timelineYearWidth(year, counts, options);
	}
	return years;
}

function timelineEndX(
	years: AppearanceTimelineYear[],
	counts: Map<number, number>,
	options: AppearanceTimelineOptions
) {
	const last = years.at(-1);
	return last
		? last.x + timelineYearWidth(last.year, counts, options)
		: options.paddingX + options.emptyYearWidth;
}

function timelineYearWidth(
	year: number,
	counts: Map<number, number>,
	options: AppearanceTimelineOptions
) {
	return options.emptyYearWidth + (counts.get(year) ?? 0) * options.yearEntryWidth;
}

function timelineMarkerX(
	appearance: PersonAppearance,
	years: AppearanceTimelineYear[],
	counts: Map<number, number>,
	fallback: number,
	options: AppearanceTimelineOptions
) {
	const year = appearanceYear(appearance);
	if (typeof year !== 'number') return fallback - options.paddingX;
	const tick = years.find((item) => item.year === year);
	if (!tick) return fallback - options.paddingX;
	return (
		tick.x + releaseFraction(appearance.releaseDate) * timelineYearWidth(year, counts, options)
	);
}

function unreleasedAppearance(appearance: PersonAppearance) {
	if (!appearance.releaseDate) return true;
	const value = Date.parse(`${appearance.releaseDate}T00:00:00Z`);
	return !Number.isFinite(value) || value > Date.now();
}

function releaseFraction(value?: string) {
	if (!value) return 0.5;
	const date = new Date(`${value}T00:00:00Z`);
	if (Number.isNaN(date.getTime())) return 0.5;
	const start = Date.UTC(date.getUTCFullYear(), 0, 1);
	const end = Date.UTC(date.getUTCFullYear() + 1, 0, 1);
	return (date.getTime() - start) / (end - start);
}

function spaceTimelineCards(items: AppearanceTimelineItem[], options: AppearanceTimelineOptions) {
	let nextX = options.paddingX;
	for (const item of items.sort((left, right) => left.markerX - right.markerX)) {
		item.cardX = Math.max(item.markerX - options.cardWidth / 2, nextX);
		nextX = item.cardX + options.cardWidth + options.cardGap;
	}
}
