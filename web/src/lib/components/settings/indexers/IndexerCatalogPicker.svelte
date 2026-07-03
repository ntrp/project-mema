<script lang="ts">
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import {
		entryHasCategory,
		flattenCategories,
		matches,
		privacyGroup,
		unique,
		uniqueCategories
	} from '$lib/components/settings/indexers/indexerCatalogFilters';
	import type { IndexerCatalogEntry } from '$lib/settings/types';

	interface Props {
		catalog: IndexerCatalogEntry[];
		onSelect: (_entry: IndexerCatalogEntry) => void;
	}

	let { catalog, onSelect }: Props = $props();
	let protocolFilter = $state('all');
	let languageFilter = $state('all');
	let privacyFilter = $state('all');
	let categoryFilter = $state('all');

	const protocols = $derived(['all', ...unique(catalog.map((entry) => entry.protocol))]);
	const languages = $derived(['all', ...unique(catalog.map((entry) => entry.language))]);
	const privacy = $derived(['all', ...unique(catalog.map((entry) => privacyGroup(entry.privacy)))]);
	const categories = $derived([
		{ value: 'all', label: 'All categories' },
		...uniqueCategories(catalog).map((category) => ({
			value: String(category.id),
			label: `${category.id} ${category.name}`
		}))
	]);
	const filteredCatalog = $derived(
		catalog.filter(
			(entry) =>
				matches(protocolFilter, entry.protocol) &&
				matches(languageFilter, entry.language) &&
				matches(privacyFilter, privacyGroup(entry.privacy)) &&
				(categoryFilter === 'all' || entryHasCategory(entry, Number(categoryFilter)))
		)
	);

	function categorySummary(entry: IndexerCatalogEntry) {
		return flattenCategories(entry.capabilities.categories)
			.map((category) => category.name)
			.slice(0, 4)
			.join(', ');
	}

	function supportSummary(entry: IndexerCatalogEntry) {
		return [
			entry.supportsSearch ? 'Search' : '',
			entry.supportsRss ? 'RSS' : '',
			entry.supportsRedirect ? 'Redirect' : ''
		]
			.filter(Boolean)
			.join(' · ');
	}

	function handleRowKeydown(event: KeyboardEvent, entry: IndexerCatalogEntry) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			onSelect(entry);
		}
	}
</script>

<div class="grid gap-4">
	<div class="grid gap-3 sm:grid-cols-4">
		<label class="grid gap-1 text-sm font-bold text-muted-foreground">
			Protocol
			<SettingsSelect
				value={protocolFilter}
				options={protocols.map((value) => ({ value, label: value }))}
				onValueChange={(value) => (protocolFilter = value)}
			/>
		</label>
		<label class="grid gap-1 text-sm font-bold text-muted-foreground">
			Language
			<SettingsSelect
				value={languageFilter}
				options={languages.map((value) => ({ value, label: value }))}
				onValueChange={(value) => (languageFilter = value)}
			/>
		</label>
		<label class="grid gap-1 text-sm font-bold text-muted-foreground">
			Privacy
			<SettingsSelect
				value={privacyFilter}
				options={privacy.map((value) => ({ value, label: value }))}
				onValueChange={(value) => (privacyFilter = value)}
			/>
		</label>
		<label class="grid gap-1 text-sm font-bold text-muted-foreground">
			Category
			<SettingsSelect
				value={categoryFilter}
				options={categories}
				onValueChange={(value) => (categoryFilter = value)}
			/>
		</label>
	</div>

	<div class="max-h-[min(560px,calc(100vh-300px))] overflow-auto rounded-md border border-border">
		<table class="w-full min-w-[760px] border-collapse text-sm">
			<thead class="sticky top-0 bg-card text-left text-xs font-extrabold text-muted-foreground">
				<tr class="border-b border-border">
					<th class="px-3 py-2">Name</th>
					<th class="px-3 py-2">Protocol</th>
					<th class="px-3 py-2">Privacy</th>
					<th class="px-3 py-2">Language</th>
					<th class="px-3 py-2">Supports</th>
					<th class="px-3 py-2">Categories</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredCatalog as entry (entry.definitionId)}
					<tr
						class="cursor-pointer border-b border-border last:border-0 hover:bg-muted/60 focus-visible:bg-muted/60 focus-visible:outline-none"
						tabindex="0"
						role="button"
						onclick={() => onSelect(entry)}
						onkeydown={(event) => handleRowKeydown(event, entry)}
					>
						<td class="px-3 py-2">
							<div class="font-bold text-foreground">{entry.name}</div>
							<div class="line-clamp-1 text-xs text-muted-foreground">{entry.description}</div>
						</td>
						<td class="px-3 py-2">{entry.protocol}</td>
						<td class="px-3 py-2">{privacyGroup(entry.privacy)}</td>
						<td class="px-3 py-2">{entry.language}</td>
						<td class="px-3 py-2">{supportSummary(entry)}</td>
						<td class="px-3 py-2">{categorySummary(entry)}</td>
					</tr>
				{:else}
					<tr>
						<td class="px-3 py-8 text-center text-muted-foreground" colspan="6">
							No catalog indexers match the selected filters.
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
