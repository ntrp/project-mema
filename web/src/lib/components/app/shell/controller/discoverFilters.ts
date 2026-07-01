import type {
	DiscoverBlacklistItem,
	MediaDiscoverSection,
	MediaMetadataDetails,
	MediaSearchResult
} from '$lib/settings/types';
import type { RelatedSectionKind } from './types';

export function discoverResultKey(candidate: MediaSearchResult) {
	return `${candidate.type}:${candidate.externalProvider ?? ''}:${candidate.externalId ?? ''}:${candidate.title}:${candidate.year ?? ''}`;
}

export function filterDiscoverSections(
	sections: MediaDiscoverSection[],
	blacklist: DiscoverBlacklistItem[]
) {
	return sections.map((section) => filterDiscoverSection(section, blacklist));
}

export function filterDiscoverSection(
	section: MediaDiscoverSection,
	blacklist: DiscoverBlacklistItem[]
): MediaDiscoverSection {
	return {
		...section,
		results: section.results.filter((result) => !isDiscoverBlacklisted(result, blacklist))
	};
}

export function relatedSectionFromDetail(
	detail: MediaMetadataDetails | undefined,
	kind: RelatedSectionKind,
	blacklist: DiscoverBlacklistItem[]
): MediaDiscoverSection | undefined {
	if (!detail) {
		return undefined;
	}
	const results = kind === 'recommendations' ? detail.recommendations : detail.similar;
	return filterDiscoverSection(
		{
			id: kind,
			title:
				kind === 'recommendations'
					? 'Recommendations'
					: detail.type === 'movie'
						? 'Similar Movies'
						: 'Similar Series',
			providerName: detail.externalProvider?.toUpperCase() ?? 'Metadata',
			mediaType: detail.type,
			results: results ?? []
		},
		blacklist
	);
}

export function sameDiscoverBlacklistItem(
	item: DiscoverBlacklistItem,
	result: DiscoverBlacklistItem | MediaSearchResult
) {
	const itemExternalKey = discoverExternalKey(item);
	const resultExternalKey = discoverExternalKey(result);
	if (itemExternalKey && resultExternalKey && itemExternalKey === resultExternalKey) {
		return true;
	}
	return discoverTitleKey(item) === discoverTitleKey(result);
}

function isDiscoverBlacklisted(result: MediaSearchResult, blacklist: DiscoverBlacklistItem[]) {
	return blacklist.some((item) => sameDiscoverBlacklistItem(item, result));
}

function discoverExternalKey(item: DiscoverBlacklistItem | MediaSearchResult) {
	if (!item.externalProvider || !item.externalId) {
		return '';
	}
	return `${item.type}:${item.externalProvider}:${item.externalId}`.trim().toLowerCase();
}

function discoverTitleKey(item: DiscoverBlacklistItem | MediaSearchResult) {
	return `${item.type}:${item.title.trim().toLowerCase()}:${item.year ?? ''}`;
}
