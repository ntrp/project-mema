export type PosterRowEdges = {
	canScrollLeft: boolean;
	canScrollRight: boolean;
};

type RowsByKey = Record<string, HTMLDivElement | undefined>;
type EdgesByKey = Record<string, PosterRowEdges | undefined>;

const defaultEdges: PosterRowEdges = {
	canScrollLeft: false,
	canScrollRight: false
};

export function createPosterRowScroller() {
	const rows = $state<RowsByKey>({});
	const edges = $state<EdgesByKey>({});

	function edgeState(key: string) {
		return edges[key] ?? defaultEdges;
	}

	function trackRow(node: HTMLDivElement, key: string) {
		rows[key] = node;
		let frame = globalThis.requestAnimationFrame(() => updateEdges(key));
		const resizeObserver = new globalThis.ResizeObserver(() => updateEdges(key));
		const onScroll = () => updateEdges(key);

		node.addEventListener('scroll', onScroll, { passive: true });
		resizeObserver.observe(node);

		return {
			destroy() {
				globalThis.cancelAnimationFrame(frame);
				node.removeEventListener('scroll', onScroll);
				resizeObserver.disconnect();
				delete rows[key];
				delete edges[key];
			}
		};
	}

	function scrollRow(key: string, direction: -1 | 1, itemOffset = 140, minimum = 220) {
		const row = rows[key];
		if (!row) {
			return;
		}
		row.scrollBy({
			left: direction * Math.max(row.clientWidth - itemOffset, minimum),
			behavior: 'smooth'
		});
	}

	function updateEdges(key: string) {
		const row = rows[key];
		if (!row) {
			edges[key] = defaultEdges;
			return;
		}
		const maxScrollLeft = Math.max(row.scrollWidth - row.clientWidth, 0);
		edges[key] = {
			canScrollLeft: row.scrollLeft > 1,
			canScrollRight: row.scrollLeft < maxScrollLeft - 1
		};
	}

	return {
		edgeState,
		scrollRow,
		trackRow
	};
}
