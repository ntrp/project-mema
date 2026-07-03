export type PosterRowEdges = {
	canScrollLeft: boolean;
	canScrollRight: boolean;
};

type RowsByKey = Record<string, HTMLDivElement | undefined>;
type EdgesByKey = Record<string, PosterRowEdges | undefined>;

const dragThreshold = 4;
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
		let dragPointerId: number | undefined;
		let dragStartX = 0;
		let dragScrollLeft = 0;
		let didDrag = false;
		let suppressClick = false;

		const finishDrag = () => {
			if (dragPointerId !== undefined && node.hasPointerCapture(dragPointerId)) {
				node.releasePointerCapture(dragPointerId);
			}
			dragPointerId = undefined;
			node.classList.remove('cursor-grabbing', 'select-none');
			if (didDrag) {
				suppressClick = true;
				globalThis.setTimeout(() => {
					suppressClick = false;
				}, 0);
			}
			didDrag = false;
		};

		const onPointerDown = (event: PointerEvent) => {
			if (event.button !== 0 || node.scrollWidth <= node.clientWidth) {
				return;
			}
			dragPointerId = event.pointerId;
			dragStartX = event.clientX;
			dragScrollLeft = node.scrollLeft;
			didDrag = false;
			node.setPointerCapture(event.pointerId);
		};

		const onPointerMove = (event: PointerEvent) => {
			if (event.pointerId !== dragPointerId) {
				return;
			}
			const deltaX = event.clientX - dragStartX;
			if (!didDrag && Math.abs(deltaX) < dragThreshold) {
				return;
			}
			didDrag = true;
			node.classList.add('cursor-grabbing', 'select-none');
			node.scrollLeft = dragScrollLeft - deltaX;
			updateEdges(key);
			event.preventDefault();
		};

		const onPointerUp = (event: PointerEvent) => {
			if (event.pointerId === dragPointerId) {
				finishDrag();
			}
		};

		const onClick = (event: MouseEvent) => {
			if (!suppressClick) {
				return;
			}
			event.preventDefault();
			event.stopPropagation();
		};

		node.addEventListener('scroll', onScroll, { passive: true });
		node.addEventListener('pointerdown', onPointerDown);
		node.addEventListener('pointermove', onPointerMove);
		node.addEventListener('pointerup', onPointerUp);
		node.addEventListener('pointercancel', onPointerUp);
		node.addEventListener('click', onClick, true);
		resizeObserver.observe(node);

		return {
			destroy() {
				globalThis.cancelAnimationFrame(frame);
				node.removeEventListener('scroll', onScroll);
				node.removeEventListener('pointerdown', onPointerDown);
				node.removeEventListener('pointermove', onPointerMove);
				node.removeEventListener('pointerup', onPointerUp);
				node.removeEventListener('pointercancel', onPointerUp);
				node.removeEventListener('click', onClick, true);
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
