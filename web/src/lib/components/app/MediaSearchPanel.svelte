<script lang="ts">
	import type { MediaSearchRequest, MediaSearchResult, MediaType } from '$lib/settings/types';

	interface Props {
		results: MediaSearchResult[];
		searching: boolean;
		addingKey?: string;
		mediaItemsCount: number;
		onSearch: (_request: MediaSearchRequest) => void;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { results, searching, addingKey, mediaItemsCount, onSearch, onAdd }: Props = $props();

	let title = $state('');
	let type = $state<MediaType>('series');
	let year = $state('');

	const digest = [
		{ title: 'Dune: Part Two', type: 'Movie', year: '2024', status: 'Trending' },
		{ title: "Frieren: Beyond Journey's End", type: 'Series', year: '2023', status: 'Anime' },
		{ title: 'The Last of Us', type: 'Series', year: '2023', status: 'Popular' },
		{ title: 'Suzume', type: 'Movie', year: '2022', status: 'Anime' }
	];

	function submit(event: SubmitEvent) {
		event.preventDefault();
		const parsedYear = Number.parseInt(year, 10);
		onSearch({
			query: title.trim(),
			type,
			year: Number.isFinite(parsedYear) ? parsedYear : undefined
		});
	}

	function resultKey(result: MediaSearchResult) {
		return `${result.type}:${result.title}:${result.year ?? ''}`;
	}
</script>

<div class="page-heading">
	<p>Explore</p>
	<h1 id="home-title">Find media to monitor</h1>
</div>

<form class="media-search-form panel" onsubmit={submit}>
	<label class="wide">
		<span>Title</span>
		<input bind:value={title} placeholder="Search for a movie or series" autocomplete="off" />
	</label>
	<label>
		<span>Type</span>
		<select bind:value={type}>
			<option value="series">Series</option>
			<option value="movie">Movie</option>
		</select>
	</label>
	<label>
		<span>Year</span>
		<input bind:value={year} inputmode="numeric" placeholder="Optional" />
	</label>
	<div class="form-actions wide">
		<button type="submit" disabled={searching}>{searching ? 'Searching' : 'Search'}</button>
		<span class="muted">{mediaItemsCount} monitored items</span>
	</div>
</form>

{#if results.length > 0}
	<div class="data-list">
		{#each results as result (resultKey(result))}
			<div class="data-row media-result-row">
				<div>
					<strong>{result.title}</strong>
					<span>{result.type}{result.year ? ` · ${result.year}` : ''}</span>
				</div>
				<small>Candidate</small>
				<button
					type="button"
					disabled={addingKey === resultKey(result)}
					onclick={() => onAdd(result)}
				>
					{addingKey === resultKey(result) ? 'Adding' : 'Monitor'}
				</button>
			</div>
		{/each}
	</div>
{:else}
	<div class="digest-grid" aria-label="Latest digest">
		{#each digest as item (item.title)}
			<article class="media-tile">
				<div class="poster-placeholder">{item.type}</div>
				<h2>{item.title}</h2>
				<p>{item.year} · {item.status}</p>
			</article>
		{/each}
	</div>
{/if}
