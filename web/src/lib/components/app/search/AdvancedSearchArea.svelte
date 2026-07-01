<script lang="ts">
	import { onMount } from 'svelte';
	import type {
		MediaAdvancedSearchRequest,
		MediaSearchGroup,
		MediaSearchResult,
		MediaType,
		MetadataProvider
	} from '$lib/settings/types';
	import AdvancedSearchResults from './AdvancedSearchResults.svelte';

	interface Props {
		initialQuery: string;
		metadataProviders: MetadataProvider[];
		groups: MediaSearchGroup[];
		searching: boolean;
		addingKey?: string;
		actionLabel: string;
		onSearch: (_request: MediaAdvancedSearchRequest) => void | Promise<void>;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let {
		initialQuery,
		metadataProviders,
		groups,
		searching,
		addingKey,
		actionLabel,
		onSearch,
		onAdd
	}: Props = $props();

	let query = $state('');
	let type = $state<MediaType | 'any'>('any');
	let year = $state('');
	let selectedProviderIds = $state<string[]>([]);

	const enabledProviders = $derived(metadataProviders.filter((provider) => provider.enabled));
	const resultCount = $derived(groups.reduce((count, group) => count + group.results.length, 0));

	onMount(() => {
		query = initialQuery;
		selectedProviderIds = enabledProviders.map((provider) => provider.id);
		if (query.trim().length > 0) {
			void submitSearch();
		}
	});

	function toggleProvider(id: string) {
		selectedProviderIds = selectedProviderIds.includes(id)
			? selectedProviderIds.filter((providerId) => providerId !== id)
			: [...selectedProviderIds, id];
	}

	function submit(event: SubmitEvent) {
		event.preventDefault();
		void submitSearch();
	}

	async function submitSearch() {
		const parsedYear = Number.parseInt(year, 10);
		const request: MediaAdvancedSearchRequest = {
			query: query.trim(),
			type: type === 'any' ? undefined : type,
			year: Number.isFinite(parsedYear) ? parsedYear : undefined,
			providerIds: selectedProviderIds.length > 0 ? selectedProviderIds : undefined,
			limit: 30
		};
		await onSearch(request);
	}
</script>

<section class="workspace-main advanced-search" aria-labelledby="advanced-search-title">
	<div class="page-heading">
		<p>Search</p>
		<h1 id="advanced-search-title">Advanced media search</h1>
	</div>

	<form class="advanced-search-form panel" onsubmit={submit}>
		<label class="wide">
			<span>Title</span>
			<input bind:value={query} placeholder="Movie or series title" autocomplete="off" />
		</label>
		<label>
			<span>Type</span>
			<select bind:value={type}>
				<option value="any">Any</option>
				<option value="movie">Movie</option>
				<option value="series">Series</option>
			</select>
		</label>
		<label>
			<span>Year</span>
			<input bind:value={year} inputmode="numeric" placeholder="Optional" />
		</label>
		<fieldset class="provider-picker wide">
			<legend>Metadata providers</legend>
			{#if enabledProviders.length > 0}
				<div class="provider-options">
					{#each enabledProviders as provider (provider.id)}
						<label>
							<input
								type="checkbox"
								checked={selectedProviderIds.includes(provider.id)}
								onchange={() => toggleProvider(provider.id)}
							/>
							<span>{provider.name}</span>
						</label>
					{/each}
				</div>
			{:else}
				<p>No enabled metadata providers are configured.</p>
			{/if}
		</fieldset>
		<div class="form-actions wide">
			<button type="submit" disabled={searching || query.trim().length === 0}>
				{searching ? 'Searching' : 'Search'}
			</button>
			<span class="muted">{resultCount} results</span>
		</div>
	</form>

	<AdvancedSearchResults {groups} {addingKey} {actionLabel} {onAdd} />
</section>
