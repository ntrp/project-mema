import type { IndexerCatalogEntry } from '$lib/settings/types';

type IndexerCategory = IndexerCatalogEntry['capabilities']['categories'][number];

export function unique(values: string[]) {
	return [...new Set(values.filter(Boolean))];
}

export function matches(filter: string, value: string) {
	return filter === 'all' || filter === value;
}

export function matchesAny(filters: string[], value: string) {
	return filters.length === 0 || filters.includes(value);
}

export function privacyGroup(value: string) {
	return value === 'semiPrivate' ? 'private' : value;
}

export function privacyLabel(value: string) {
	return value === 'semiPrivate' || value === 'semi-private' ? 'semi-private' : value;
}

export function flattenCategories(categories: IndexerCategory[]): IndexerCategory[] {
	return categories.flatMap((category) => [category, ...flattenCategories(category.children)]);
}

export function uniqueCategories(entries: IndexerCatalogEntry[]) {
	const categoriesByID: Record<number, { id: number; name: string }> = {};
	for (const entry of entries) {
		for (const category of flattenCategories(entry.capabilities.categories)) {
			categoriesByID[category.id] = { id: category.id, name: category.name };
		}
	}
	return Object.values(categoriesByID).sort((left, right) => left.id - right.id);
}

export function entryHasCategory(entry: IndexerCatalogEntry, id: number) {
	return flattenCategories(entry.capabilities.categories).some((category) => category.id === id);
}

export function fuzzyMatch(query: string, values: string[]) {
	const needle = query.trim().toLowerCase();
	if (!needle) {
		return true;
	}
	const haystack = values.filter(Boolean).join(' ').toLowerCase();
	let offset = 0;
	for (const char of needle) {
		const next = haystack.indexOf(char, offset);
		if (next < 0) {
			return false;
		}
		offset = next + 1;
	}
	return true;
}

export function textMatchRank(query: string, values: string[]) {
	const needle = query.trim().toLowerCase();
	if (!needle) {
		return 0;
	}
	const haystack = values.filter(Boolean).join(' ').toLowerCase();
	if (haystack.includes(needle)) {
		return 0;
	}
	return fuzzyMatch(query, values) ? 1 : -1;
}
