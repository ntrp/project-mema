<script lang="ts">
	import type {
		LibraryMediaKind,
		LibraryScan,
		LibraryScanItem,
		LibraryScanItemMatchRequest,
		MediaSearchResult
	} from '$lib/settings/types';

	interface MatchDraft {
		mediaKind: LibraryMediaKind;
		query: string;
		title: string;
		year?: number;
		externalProvider?: string;
		externalId?: string;
		overview?: string;
		posterPath?: string;
		results: MediaSearchResult[];
		searching: boolean;
		matching: boolean;
	}

	interface Props {
		scan?: LibraryScan;
		loading: boolean;
		onSearchMatch: (_kind: LibraryMediaKind, _query: string) => Promise<MediaSearchResult[]>;
		onMatch: (_item: LibraryScanItem, _request: LibraryScanItemMatchRequest) => Promise<void>;
	}

	let { scan, loading, onSearchMatch, onMatch }: Props = $props();

	const mediaKinds: { value: LibraryMediaKind; label: string }[] = [
		{ value: 'movie', label: 'Movie' },
		{ value: 'series', label: 'Series' },
		{ value: 'anime_movie', label: 'Anime movie' },
		{ value: 'anime_series', label: 'Anime series' }
	];

	let drafts = $state<Record<string, MatchDraft>>({});
	let pendingItems = $derived(scan?.items.filter((item) => item.status === 'pending') ?? []);
	let addedItems = $derived(scan?.items.filter((item) => item.status !== 'pending') ?? []);

	$effect(() => {
		for (const item of scan?.items ?? []) {
			if (item.status !== 'pending' || drafts[item.id]) {
				continue;
			}
			const detectedKind = item.detectedMediaKind === 'unknown' ? 'movie' : item.detectedMediaKind;
			drafts[item.id] = {
				mediaKind: detectedKind,
				query: item.detectedTitle,
				title: item.detectedTitle,
				year: item.detectedYear,
				results: [],
				searching: false,
				matching: false
			};
		}
	});

	async function search(item: LibraryScanItem) {
		const draft = drafts[item.id];
		if (!draft || draft.query.trim() === '') {
			return;
		}
		draft.searching = true;
		try {
			draft.results = await onSearchMatch(draft.mediaKind, draft.query);
		} finally {
			draft.searching = false;
		}
	}

	function selectResult(item: LibraryScanItem, result: MediaSearchResult) {
		const draft = drafts[item.id];
		if (!draft) {
			return;
		}
		draft.title = result.title;
		draft.query = result.title;
		draft.year = result.year;
		draft.externalProvider = result.externalProvider;
		draft.externalId = result.externalId;
		draft.overview = result.overview;
		draft.posterPath = result.posterPath;
		draft.results = [];
	}

	async function match(item: LibraryScanItem) {
		const draft = drafts[item.id];
		if (!draft || draft.title.trim() === '') {
			return;
		}
		draft.matching = true;
		try {
			await onMatch(item, {
				mediaKind: draft.mediaKind,
				title: draft.title.trim(),
				year: draft.year,
				monitored: true,
				externalProvider: draft.externalProvider,
				externalId: draft.externalId,
				overview: draft.overview,
				posterPath: draft.posterPath
			});
		} finally {
			draft.matching = false;
		}
	}
</script>

<div class="panel" aria-labelledby="library-scan-title">
	<div class="section-heading">
		<div>
			<p class="section-kicker">Library review</p>
			<h2 id="library-scan-title">
				{#if scan}
					{scan.folderPath}
				{:else}
					Latest scan
				{/if}
			</h2>
		</div>
		{#if scan}
			<div class="summary">
				<span>{scan.totalFiles} files</span>
				<span>{scan.autoMatchedCount} auto-added</span>
				<span>{scan.manualCount} pending</span>
			</div>
		{/if}
	</div>

	{#if loading}
		<p class="muted">Loading library scan</p>
	{:else if !scan}
		<p class="empty">Add a library folder to review discovered media.</p>
	{:else}
		<div class="review-stack">
			<section class="review-section" aria-labelledby="pending-title">
				<h3 id="pending-title">Manual matches</h3>
				<div class="review-list">
					{#each pendingItems as item (item.id)}
						{@const draft = drafts[item.id]}
						<div class="review-card">
							<div class="review-file">
								<strong>{item.detectedTitle || item.fileName}</strong>
								<span>{item.path}</span>
							</div>

							{#if draft}
								<div class="match-grid">
									<label>
										<span>Type</span>
										<select bind:value={draft.mediaKind}>
											{#each mediaKinds as kind (kind.value)}
												<option value={kind.value}>{kind.label}</option>
											{/each}
										</select>
									</label>
									<label>
										<span>Search</span>
										<input bind:value={draft.query} />
									</label>
									<label>
										<span>Year</span>
										<input bind:value={draft.year} min="1800" max="3000" type="number" />
									</label>
									<div class="form-actions">
										<button
											type="button"
											class="secondary"
											disabled={draft.searching}
											onclick={() => search(item)}
										>
											{draft.searching ? 'Searching' : 'Search'}
										</button>
										<button
											type="button"
											disabled={draft.matching || draft.title.trim() === ''}
											onclick={() => match(item)}
										>
											{draft.matching ? 'Adding' : 'Add match'}
										</button>
									</div>
								</div>

								{#if draft.results.length > 0}
									<div class="autocomplete-list">
										{#each draft.results as result (`${result.type}:${result.title}:${result.year ?? ''}`)}
											<button type="button" onclick={() => selectResult(item, result)}>
												<strong>{result.title}</strong>
												<span>{result.type}{result.year ? ` · ${result.year}` : ''}</span>
											</button>
										{/each}
									</div>
								{/if}
							{/if}
						</div>
					{:else}
						<p class="empty">No manual matches needed.</p>
					{/each}
				</div>
			</section>

			<section class="review-section" aria-labelledby="matched-title">
				<h3 id="matched-title">Added automatically</h3>
				<div class="review-list compact">
					{#each addedItems as item (item.id)}
						<div class="review-row">
							<div>
								<strong>{item.matchedTitle ?? item.detectedTitle}</strong>
								<span>{item.path}</span>
							</div>
							<span class="status-enabled">{item.status === 'auto_added' ? 'Auto' : 'Manual'}</span>
						</div>
					{:else}
						<p class="empty">No files were auto-added.</p>
					{/each}
				</div>
			</section>
		</div>
	{/if}
</div>
