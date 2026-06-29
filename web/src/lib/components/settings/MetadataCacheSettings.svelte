<script lang="ts">
	import type { MetadataCacheResponse } from '$lib/settings/types';

	interface Props {
		cache: MetadataCacheResponse;
		pattern: string;
		loading: boolean;
		clearing: boolean;
		onRefresh: () => void | Promise<void>;
		onClearAll: () => void | Promise<void>;
		onClearPattern: (_event: SubmitEvent) => void | Promise<void>;
	}

	let {
		cache,
		pattern = $bindable(),
		loading,
		clearing,
		onRefresh,
		onClearAll,
		onClearPattern
	}: Props = $props();

	function formatDate(value: string) {
		return new Intl.DateTimeFormat(undefined, {
			dateStyle: 'medium',
			timeStyle: 'short'
		}).format(new Date(value));
	}
</script>

<section class="panel cache-panel" aria-labelledby="metadata-cache-title">
	<div class="section-heading">
		<div>
			<p class="section-kicker">Cache</p>
			<h2 id="metadata-cache-title">Metadata provider cache</h2>
		</div>
		<button type="button" class="secondary" disabled={loading} onclick={() => void onRefresh()}>
			{loading ? 'Refreshing' : 'Refresh'}
		</button>
	</div>

	<div class="status-grid cache-stats" aria-label="Metadata cache stats">
		<div>
			<span>Total entries</span>
			<strong>{cache.stats.totalEntries}</strong>
		</div>
		<div>
			<span>Active</span>
			<strong>{cache.stats.activeEntries}</strong>
		</div>
		<div>
			<span>Expired</span>
			<strong>{cache.stats.expiredEntries}</strong>
		</div>
		<div>
			<span>Providers</span>
			<strong>{cache.stats.providerCount}</strong>
		</div>
	</div>

	<form class="cache-actions" onsubmit={onClearPattern}>
		<label>
			<span>Reset by regex</span>
			<input bind:value={pattern} placeholder="discover:|details:123|matrix" autocomplete="off" />
		</label>
		<div class="form-actions">
			<button type="submit" class="danger" disabled={clearing || pattern.trim().length === 0}>
				{clearing ? 'Resetting' : 'Reset matching'}
			</button>
			<button type="button" class="danger" disabled={clearing} onclick={() => void onClearAll()}>
				Reset all
			</button>
		</div>
	</form>

	{#if cache.entries.length > 0}
		<div class="table-wrap">
			<table>
				<thead>
					<tr>
						<th>Provider</th>
						<th>Kind</th>
						<th>Media</th>
						<th>Key</th>
						<th>Items</th>
						<th>Expires</th>
					</tr>
				</thead>
				<tbody>
					{#each cache.entries as entry (`${entry.providerName}:${entry.mediaType}:${entry.query}:${entry.year}`)}
						<tr>
							<td>
								<strong>{entry.providerName}</strong>
								<span>{entry.providerType}</span>
							</td>
							<td>{entry.cacheKind}</td>
							<td>{entry.mediaType}{entry.year ? ` · ${entry.year}` : ''}</td>
							<td><code>{entry.query}</code></td>
							<td>{entry.itemCount}</td>
							<td>
								<span class:status-disabled={entry.expired} class:status-enabled={!entry.expired}>
									{entry.expired ? 'Expired' : formatDate(entry.expiresAt)}
								</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<p class="empty">No metadata cache entries yet.</p>
	{/if}
</section>
