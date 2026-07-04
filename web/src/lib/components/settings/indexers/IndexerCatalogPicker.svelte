<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import IndexerCatalogFilterMultiSelect from './IndexerCatalogFilterMultiSelect.svelte';
	import IndexerCatalogTable from './IndexerCatalogTable.svelte';
	import {
		entryHasCategory,
		matchesAny,
		privacyLabel,
		textMatchRank,
		unique,
		uniqueCategories
	} from '$lib/components/settings/indexers/indexerCatalogFilters';
	import {
		privacyBadgeClass,
		protocolBadgeClass
	} from '$lib/components/settings/indexers/indexerCatalogPresentation';
	import type { IndexerCatalogEntry } from '$lib/settings/types';
	interface Props {
		catalog: IndexerCatalogEntry[];
		onSelect: (_entry: IndexerCatalogEntry) => void;
	}
	let { catalog, onSelect }: Props = $props();
	let protocolFilter = $state<string[]>([]);
	let languageFilter = $state<string[]>([]);
	let privacyFilter = $state<string[]>([]);
	let categoryFilter = $state<string[]>([]);
	let textFilter = $state('');
	let visibleCount = $state(25);
	let previousFilterKey = $state('');

	const pageSize = 25;
	const protocols = $derived(
		unique(catalog.map((entry) => entry.protocol)).map((value) => ({
			value,
			label: value,
			class: protocolBadgeClass(value)
		}))
	);
	const languages = $derived(
		unique(catalog.map((entry) => entry.language)).map((value) => ({
			value,
			label: value,
			class: 'uppercase border-border bg-muted text-muted-foreground'
		}))
	);
	const privacy = $derived(
		unique(catalog.map((entry) => privacyLabel(entry.privacy))).map((value) => ({
			value,
			label: value,
			class: privacyBadgeClass(value)
		}))
	);
	const categories = $derived(
		uniqueCategories(catalog).map((category) => ({
			value: String(category.id),
			label: category.name,
			class: 'border-sky-500/50 bg-sky-500/10 text-sky-300'
		}))
	);
	const filteredCatalog = $derived(filterCatalog());
	const visibleCatalog = $derived(filteredCatalog.slice(0, visibleCount));
	const filterKey = $derived(
		JSON.stringify([textFilter, protocolFilter, languageFilter, privacyFilter, categoryFilter])
	);

	$effect(() => {
		if (previousFilterKey && previousFilterKey !== filterKey) {
			visibleCount = pageSize;
		}
		previousFilterKey = filterKey;
	});

	function filterCatalog() {
		const query = textFilter.trim();
		const facetMatches = catalog.filter(matchesFacets);
		if (!query) return facetMatches;
		return facetMatches
			.map((entry) => ({ entry, rank: textMatchRank(query, [entry.name]) }))
			.filter(({ rank }) => rank >= 0)
			.sort((left, right) => left.rank - right.rank || left.entry.name.localeCompare(right.entry.name))
			.map(({ entry }) => entry);
	}

	function matchesFacets(entry: IndexerCatalogEntry) {
		return (
			matchesAny(protocolFilter, entry.protocol) &&
			matchesAny(languageFilter, entry.language) &&
			matchesAny(privacyFilter, privacyLabel(entry.privacy)) &&
			(categoryFilter.length === 0 ||
				categoryFilter.some((category) => entryHasCategory(entry, Number(category))))
		);
	}

	function loadMore() {
		if (visibleCount >= filteredCatalog.length) return;
		visibleCount = Math.min(visibleCount + pageSize, filteredCatalog.length);
	}
</script>

<div class="grid gap-4">
	<label class="grid gap-1 text-sm font-bold text-muted-foreground">
		Filter
		<Input bind:value={textFilter} placeholder="Search catalog by name" />
	</label>
	<div class="grid gap-3 sm:grid-cols-4">
		<IndexerCatalogFilterMultiSelect
			id="indexer-catalog-protocol-filter"
			label="Protocol"
			values={protocolFilter}
			options={protocols}
			placeholder="All protocols"
			onChange={(values) => (protocolFilter = values)}
		/>
		<IndexerCatalogFilterMultiSelect
			id="indexer-catalog-language-filter"
			label="Language"
			values={languageFilter}
			options={languages}
			placeholder="All languages"
			onChange={(values) => (languageFilter = values)}
		/>
		<IndexerCatalogFilterMultiSelect
			id="indexer-catalog-privacy-filter"
			label="Privacy"
			values={privacyFilter}
			options={privacy}
			placeholder="All privacy"
			onChange={(values) => (privacyFilter = values)}
		/>
		<IndexerCatalogFilterMultiSelect
			id="indexer-catalog-category-filter"
			label="Category"
			values={categoryFilter}
			options={categories}
			placeholder="All categories"
			onChange={(values) => (categoryFilter = values)}
		/>
	</div>

	<IndexerCatalogTable
		entries={visibleCatalog}
		hasMore={visibleCatalog.length < filteredCatalog.length}
		onEndReached={loadMore}
		{onSelect}
	/>
	{#if visibleCatalog.length < filteredCatalog.length}
		<p class="m-0 text-xs text-muted-foreground">
			Showing {visibleCatalog.length} of {filteredCatalog.length}
		</p>
	{/if}
</div>
