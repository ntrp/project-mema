export function selectedFirst<T>(
	options: T[],
	selectedValues: Set<string | number>,
	valueOf: (_option: T) => string | number
) {
	return [
		...options.filter((option) => selectedValues.has(valueOf(option))),
		...options.filter((option) => !selectedValues.has(valueOf(option)))
	];
}
