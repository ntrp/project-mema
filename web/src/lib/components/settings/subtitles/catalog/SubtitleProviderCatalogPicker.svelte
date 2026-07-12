<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import IndexerCatalogFilterMultiSelect from '$lib/components/settings/indexers/IndexerCatalogFilterMultiSelect.svelte';
	import type { components } from '$lib/api/generated/schema';
	import {
		dependencyLabels,
		matchesAny,
		runtimeLabel,
		textMatchRank,
		unique
	} from './subtitleProviderCatalogFilters';
	import SubtitleProviderCatalogTable from './SubtitleProviderCatalogTable.svelte';
	type SubtitleProviderCatalogEntry = components['schemas']['SubtitleProviderCatalogEntry'];

	interface Props {
		catalog: SubtitleProviderCatalogEntry[];
		onSelect: (_entry: SubtitleProviderCatalogEntry) => void;
	}

	let { catalog, onSelect }: Props = $props();
	let runtimeFilter = $state<string[]>([]);
	let dependencyFilter = $state<string[]>([]);
	let textFilter = $state('');

	const runtimeOptions = $derived(
		unique(catalog.map((entry) => runtimeLabel(entry.runtimeStatus))).map((value) => ({
			value,
			label: value,
			class: 'border-border bg-muted text-muted-foreground'
		}))
	);
	const dependencyOptions = $derived(
		unique(catalog.flatMap(dependencyLabels)).map((value) => ({
			value,
			label: value,
			class: 'border-sky-500/50 bg-sky-500/10 text-sky-300'
		}))
	);
	const filteredCatalog = $derived(filterCatalog());

	function filterCatalog() {
		const query = textFilter.trim();
		const facetMatches = catalog.filter((entry) => {
			const deps = dependencyLabels(entry);
			return (
				matchesAny(runtimeFilter, runtimeLabel(entry.runtimeStatus)) &&
				(dependencyFilter.length === 0 || dependencyFilter.some((item) => deps.includes(item)))
			);
		});
		if (!query) return facetMatches;
		return facetMatches
			.map((entry) => ({
				entry,
				rank: textMatchRank(query, [entry.displayName, entry.key, entry.runtimeMessage])
			}))
			.filter(({ rank }) => rank >= 0)
			.sort(
				(left, right) => left.rank - right.rank || left.entry.key.localeCompare(right.entry.key)
			)
			.map(({ entry }) => entry);
	}
</script>

<div class="grid gap-4">
	<div class="grid gap-3 sm:grid-cols-2">
		<IndexerCatalogFilterMultiSelect
			id="subtitle-provider-runtime-filter"
			label="Support status"
			values={runtimeFilter}
			options={runtimeOptions}
			placeholder="All support states"
			onChange={(values) => (runtimeFilter = values)}
		/>
		<IndexerCatalogFilterMultiSelect
			id="subtitle-provider-dependency-filter"
			label="Requirements"
			values={dependencyFilter}
			options={dependencyOptions}
			placeholder="All requirements"
			onChange={(values) => (dependencyFilter = values)}
		/>
	</div>
	<p class="m-0 text-xs text-muted-foreground">
		Support status shows whether Mema can use the provider now. Requirements are extra tools or
		services the provider needs, such as FFmpeg, an API account, or captcha support.
	</p>
	<label class="grid gap-1 text-sm font-bold text-muted-foreground">
		Filter
		<Input bind:value={textFilter} placeholder="Search subtitle providers by name" />
	</label>
	<SubtitleProviderCatalogTable entries={filteredCatalog} {onSelect} />
	<p class="m-0 text-xs text-muted-foreground">
		Showing {filteredCatalog.length} of {catalog.length} catalog providers.
	</p>
</div>
