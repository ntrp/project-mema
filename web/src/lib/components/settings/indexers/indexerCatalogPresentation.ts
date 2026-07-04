import type { IndexerCatalogEntry } from '$lib/settings/types';

export function protocolBadgeClass(protocol: string) {
	return protocol === 'usenet'
		? 'uppercase border-sky-500/50 bg-sky-500/10 text-sky-300'
		: 'uppercase border-emerald-500/50 bg-emerald-500/10 text-emerald-300';
}

export function privacyBadgeClass(privacy: string) {
	if (privacy === 'public') {
		return 'uppercase border-emerald-500/50 bg-emerald-500/10 text-emerald-300';
	}
	if (privacy === 'semiPrivate' || privacy === 'semi-private') {
		return 'uppercase border-orange-500/50 bg-orange-500/10 text-orange-300';
	}
	return 'uppercase border-destructive/50 bg-destructive/10 text-destructive';
}

export function categoryBadges(
	entry: IndexerCatalogEntry,
	flatten: (categories: Category[]) => Category[]
) {
	return flatten(entry.capabilities.categories).slice(0, 4);
}

type Category = IndexerCatalogEntry['capabilities']['categories'][number];
