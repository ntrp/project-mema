import type { AppearanceTimelineItem, AppearanceTimelineYear } from './personTimeline';

type ScrollMetrics = Record<'scrollLeft' | 'clientWidth' | 'scrollWidth', number>;

export function nextTimelineYearScrollLeft(
	years: AppearanceTimelineYear[],
	metrics: ScrollMetrics,
	direction: -1 | 1
) {
	const center = metrics.scrollLeft + metrics.clientWidth / 2;
	const currentIndex = centeredYearIndex(years, center);
	const targetIndex = Math.max(0, Math.min(years.length - 1, currentIndex + direction));
	const target = years[targetIndex]?.x ?? center;
	return clampScrollLeft(target - metrics.clientWidth / 2, metrics);
}

export function nextTimelineCardScrollLeft(
	items: AppearanceTimelineItem[],
	metrics: ScrollMetrics,
	cardWidth: number,
	direction: -1 | 1
) {
	const edgeTolerance = 2;
	const right = metrics.scrollLeft + metrics.clientWidth;
	const cards = items
		.map((item) => ({ left: item.cardX, right: item.cardX + cardWidth }))
		.sort((a, b) => a.left - b.left);
	if (direction > 0) {
		return (
			cards.find((card) => card.right > right + edgeTolerance)?.left ??
			metrics.scrollWidth - metrics.clientWidth
		);
	}
	const card = [...cards].reverse().find((item) => item.left < metrics.scrollLeft - edgeTolerance);
	return card ? clampScrollLeft(card.right - metrics.clientWidth, metrics) : 0;
}

function centeredYearIndex(years: AppearanceTimelineYear[], center: number) {
	let index = 0;
	let distance = Number.POSITIVE_INFINITY;
	for (const [candidateIndex, year] of years.entries()) {
		const candidateDistance = Math.abs(year.x - center);
		if (candidateDistance < distance) {
			index = candidateIndex;
			distance = candidateDistance;
		}
	}
	return index;
}

function clampScrollLeft(value: number, metrics: ScrollMetrics) {
	return Math.max(0, Math.min(value, Math.max(0, metrics.scrollWidth - metrics.clientWidth)));
}
