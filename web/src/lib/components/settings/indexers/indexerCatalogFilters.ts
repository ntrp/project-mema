import type { IndexerCatalogEntry } from '$lib/settings/types';

type IndexerCategory = IndexerCatalogEntry['capabilities']['categories'][number];

export function unique(values: string[]) {
	return [...new Set(values.filter(Boolean))];
}

export function matches(filter: string, value: string) {
	return filter === 'all' || filter === value;
}

export function privacyGroup(value: string) {
	return value === 'semiPrivate' ? 'private' : value;
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
