export function containWheelBoundary(event: WheelEvent) {
	const element = event.currentTarget;
	if (!(element instanceof HTMLElement)) return;
	event.stopPropagation();
	if (!wouldOverscroll(element, event.deltaY)) return;
	event.preventDefault();
}

function wouldOverscroll(element: HTMLElement, deltaY: number) {
	if (deltaY === 0 || element.scrollHeight <= element.clientHeight) return false;
	const atTop = element.scrollTop <= 0;
	const atBottom = element.scrollTop + element.clientHeight >= element.scrollHeight - 1;
	return (deltaY < 0 && atTop) || (deltaY > 0 && atBottom);
}
