import type { CustomFormat } from '$lib/settings/types';

export function filterCustomFormats(formats: CustomFormat[], query: string) {
	const terms = searchTerms(query);
	if (terms.length === 0) return formats;
	return formats.filter((format) =>
		terms.every((term) => format.name.toLowerCase().includes(term))
	);
}

function searchTerms(value: string) {
	return value.trim().toLowerCase().split(/\s+/).filter(Boolean);
}
