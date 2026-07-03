export function dragScroll(node: HTMLElement) {
	let pointerId = -1;
	let startX = 0;
	let startLeft = 0;

	function endDrag() {
		if (pointerId >= 0 && node.hasPointerCapture(pointerId)) node.releasePointerCapture(pointerId);
		pointerId = -1;
		node.classList.remove('cursor-grabbing', 'select-none');
	}

	function onPointerDown(event: PointerEvent) {
		if (event.button !== 0 || interactiveTarget(event.target)) return;
		pointerId = event.pointerId;
		startX = event.clientX;
		startLeft = node.scrollLeft;
		node.setPointerCapture(pointerId);
		node.classList.add('cursor-grabbing', 'select-none');
	}

	function onPointerMove(event: PointerEvent) {
		if (event.pointerId !== pointerId) return;
		const delta = event.clientX - startX;
		if (Math.abs(delta) > 2) event.preventDefault();
		node.scrollLeft = startLeft - delta;
	}

	node.addEventListener('pointerdown', onPointerDown);
	node.addEventListener('pointermove', onPointerMove);
	node.addEventListener('pointerup', endDrag);
	node.addEventListener('pointercancel', endDrag);

	return {
		destroy() {
			endDrag();
			node.removeEventListener('pointerdown', onPointerDown);
			node.removeEventListener('pointermove', onPointerMove);
			node.removeEventListener('pointerup', endDrag);
			node.removeEventListener('pointercancel', endDrag);
		}
	};
}

function interactiveTarget(target: EventTarget | null) {
	return target instanceof Element && Boolean(target.closest('a,button,input,select,textarea'));
}
