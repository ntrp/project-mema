<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MediaSearchResult } from '$lib/settings/types';

	interface Props {
		result: MediaSearchResult;
		adding?: boolean;
		actionLabel: string;
		inLibrary?: boolean;
		onAdd: (_candidate: MediaSearchResult) => void;
		onBlacklist?: (_candidate: MediaSearchResult) => void;
		blacklisting?: boolean;
		showBlacklistAction?: boolean;
	}

	let {
		result,
		adding = false,
		actionLabel,
		inLibrary = false,
		onAdd,
		onBlacklist,
		blacklisting = false,
		showBlacklistAction = false
	}: Props = $props();

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w342${path}`;
	}
</script>

<article class="poster-card">
	<div class="poster-frame">
		{#if posterUrl(result.posterPath)}
			<img src={posterUrl(result.posterPath)} alt="" loading="lazy" />
		{:else}
			<div class="poster-placeholder">{result.type}</div>
		{/if}
		{#if result.externalProvider && result.externalId}
			<a
				class="poster-detail-link"
				href={resolve('/media/[provider]/[type]/[externalId]', {
					provider: result.externalProvider,
					type: result.type,
					externalId: result.externalId
				})}
				aria-label={`Open ${result.title} details`}
			></a>
		{/if}
		<span class="media-badge" class:movie={!inLibrary && result.type === 'movie'}>
			{inLibrary ? 'In library' : result.type}
		</span>
		<div class="poster-hover">
			<span class="poster-year">{result.year ?? 'Unknown'}</span>
			<h3>{result.title}</h3>
			<p>{result.overview ?? 'No overview available.'}</p>
			{#if showBlacklistAction && onBlacklist}
				<button
					type="button"
					class="poster-icon-action blacklist-action"
					disabled={blacklisting}
					aria-label={`Hide ${result.title} from discover`}
					title="Hide from discover"
					onclick={(event) => {
						event.stopPropagation();
						onBlacklist(result);
					}}
				>
					<span class="app-icon" aria-hidden="true">visibility_off</span>
				</button>
			{/if}
			{#if inLibrary}
				<span class="status-pill">In library</span>
			{:else}
				<button
					type="button"
					class="add-action-button"
					disabled={adding}
					onclick={(event) => {
						event.stopPropagation();
						onAdd(result);
					}}
				>
					<span class="app-icon" aria-hidden="true">add</span>
					<span>{adding ? 'Working' : actionLabel}</span>
				</button>
			{/if}
		</div>
	</div>
</article>
