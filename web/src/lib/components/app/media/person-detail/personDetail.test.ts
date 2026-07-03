import { describe, expect, it } from 'vitest';

import { appearancesByYear, appearanceYear, sortedAppearances } from './personDetail';
import { appearanceTimelineData } from './personTimeline';
import { nextTimelineYearScrollLeft } from './personTimelineScroll';
import type { PersonAppearance } from '$lib/settings/types';

describe('person detail appearance helpers', () => {
	it('sorts appearances chronologically by release date, then year, then title', () => {
		const appearances = [
			appearance({ title: 'Later', year: 2022 }),
			appearance({ title: 'Exact date', year: 2020, releaseDate: '2020-05-01' }),
			appearance({ title: 'Same date B', year: 2020, releaseDate: '2020-05-01' }),
			appearance({ title: 'Same date A', year: 2020, releaseDate: '2020-05-01' }),
			appearance({ title: 'Earlier', year: 2019 })
		];

		expect(sortedAppearances(appearances).map((item) => item.title)).toEqual([
			'Earlier',
			'Exact date',
			'Same date A',
			'Same date B',
			'Later'
		]);
	});

	it('derives timeline years from year or release date', () => {
		expect(appearanceYear(appearance({ year: 2024 }))).toBe(2024);
		expect(appearanceYear(appearance({ year: undefined, releaseDate: '2025-03-01' }))).toBe(2025);
		expect(appearanceYear(appearance({ year: undefined, releaseDate: undefined }))).toBe('Unknown');
	});

	it('groups timeline appearances by year once', () => {
		const groups = appearancesByYear([
			appearance({ title: 'Second 2020', year: 2020 }),
			appearance({ title: 'First 2020', year: 2020 }),
			appearance({ title: 'Later', year: 2021 })
		]);

		expect(groups.map((group) => group.year)).toEqual([2020, 2021]);
		expect(groups[0].appearances.map((item) => item.title)).toEqual(['First 2020', 'Second 2020']);
	});

	it('builds a full-year timeline with alternating card lanes', () => {
		const timeline = appearanceTimelineData(
			[
				appearance({ title: 'First', year: 2020 }),
				appearance({ title: 'Second', year: 2022 }),
				appearance({ title: 'Third', year: 2023 })
			],
			{
				cardWidth: 116,
				cardGap: 24,
				cardTopY: 24,
				cardBottomY: 270,
				emptyYearWidth: 96,
				yearEntryWidth: 84,
				paddingX: 80
			}
		);

		expect(timeline.years.map((year) => year.year)).toEqual([2020, 2021, 2022, 2023]);
		expect(timeline.items.map((item) => item.top)).toEqual([true, false, true]);
		expect(timeline.items.map((item) => item.cardY)).toEqual([24, 270, 24]);
	});

	it('spaces years by appearance density', () => {
		const timeline = appearanceTimelineData(
			[
				appearance({ title: 'First 2020', year: 2020 }),
				appearance({ title: 'Second 2020', year: 2020 }),
				appearance({ title: 'Only 2022', year: 2022 })
			],
			{
				cardWidth: 116,
				cardGap: 24,
				cardTopY: 24,
				cardBottomY: 270,
				emptyYearWidth: 96,
				yearEntryWidth: 84,
				paddingX: 80
			}
		);

		expect(timeline.years).toEqual([
			{ year: 2020, x: 80 },
			{ year: 2021, x: 344 },
			{ year: 2022, x: 440 }
		]);
	});

	it('marks future or unknown release cards as unreleased', () => {
		const timeline = appearanceTimelineData(
			[
				appearance({ title: 'Released', year: 2020, releaseDate: '2020-01-01' }),
				appearance({ title: 'Future', year: 3000, releaseDate: '3000-01-01' }),
				appearance({ title: 'Unknown', year: 2022 })
			],
			timelineOptions()
		);

		expect(timeline.items.map((item) => [item.appearance.title, item.unreleased])).toEqual([
			['Released', false],
			['Unknown', true],
			['Future', true]
		]);
	});

	it('increments or decrements the centered timeline year', () => {
		const years = [
			{ year: 2020, x: 80 },
			{ year: 2021, x: 344 },
			{ year: 2022, x: 440 }
		];
		const metrics = { scrollLeft: 244, clientWidth: 200, scrollWidth: 700 };

		expect(nextTimelineYearScrollLeft(years, metrics, 1)).toBe(340);
		expect(nextTimelineYearScrollLeft(years, metrics, -1)).toBe(0);
	});
});

function appearance(overrides: Partial<PersonAppearance>): PersonAppearance {
	return {
		title: 'Appearance',
		type: 'movie',
		externalProvider: 'tmdb',
		externalId: overrides.title ?? '1',
		...overrides
	};
}

function timelineOptions() {
	return {
		cardWidth: 116,
		cardGap: 24,
		cardTopY: 24,
		cardBottomY: 270,
		emptyYearWidth: 96,
		yearEntryWidth: 84,
		paddingX: 80
	};
}
