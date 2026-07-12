import type { components } from '$lib/api/generated/schema';

type SubtitleProviderCatalogEntry = components['schemas']['SubtitleProviderCatalogEntry'];

export function unique(values: string[]) {
	return Array.from(new Set(values.filter(Boolean))).sort((left, right) => left.localeCompare(right));
}

export function dependencyLabels(entry: SubtitleProviderCatalogEntry) {
	return Object.entries(entry.dependencies)
		.filter(([, enabled]) => enabled)
		.map(([key]) => key.replaceAll('_', ' '));
}

export function runtimeLabel(status: SubtitleProviderCatalogEntry['runtimeStatus']) {
	if (status === 'supported') return 'Supported';
	if (status === 'catalog_only') return 'Catalog only';
	return 'Unsupported';
}

export function matchesAny(filters: string[], value: string) {
	return filters.length === 0 || filters.includes(value);
}

export function textMatchRank(query: string, values: string[]) {
	const normalizedQuery = query.trim().toLowerCase();
	if (!normalizedQuery) return 0;
	let best = -1;
	for (const value of values) {
		const normalized = value.toLowerCase();
		const index = normalized.indexOf(normalizedQuery);
		if (index >= 0 && (best === -1 || index < best)) best = index;
	}
	return best;
}
