<script lang="ts">
	import type { MediaItem, ReleaseCandidate, ReleaseSearchState } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		releaseResults?: ReleaseSearchState;
		searching?: boolean;
		grabbingKey?: string;
		canManage: boolean;
		onSearch: (_item: MediaItem) => void;
		onGrab: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onClose: () => void;
	}

	let {
		item,
		releaseResults,
		searching = false,
		grabbingKey,
		canManage,
		onSearch,
		onGrab,
		onClose
	}: Props = $props();

	function releaseKey(release: ReleaseCandidate) {
		return `${item.id}:${release.id}`;
	}

	function sizeLabel(sizeBytes: number) {
		if (!sizeBytes) return '-';
		const gib = sizeBytes / 1024 / 1024 / 1024;
		return `${gib.toFixed(gib >= 10 ? 0 : 1)} GiB`;
	}
</script>

<div class="modal-backdrop" role="presentation" onclick={onClose}>
	<div
		class="modal-shell settings-modal media-file-modal"
		role="dialog"
		aria-modal="true"
		aria-labelledby="manual-search-title"
		tabindex="-1"
		onclick={(event) => event.stopPropagation()}
		onkeydown={(event) => event.stopPropagation()}
	>
		<div class="modal-heading">
			<h2 id="manual-search-title">Manual search</h2>
			<button type="button" class="icon-button" aria-label="Close" onclick={onClose}>
				<span class="app-icon" aria-hidden="true">close</span>
			</button>
		</div>
		<div class="settings-toolbar">
			<button type="button" disabled={!canManage || searching} onclick={() => onSearch(item)}>
				{searching ? 'Searching' : 'Search releases'}
			</button>
		</div>
		{#if releaseResults?.errors.length}
			<div class="inline-errors">
				{#each releaseResults.errors as error (error)}
					<p>{error}</p>
				{/each}
			</div>
		{/if}
		<div class="table-wrap">
			<table>
				<thead>
					<tr>
						<th>Release</th>
						<th>Indexer</th>
						<th>Size</th>
						<th>Seeders</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{#each releaseResults?.releases ?? [] as release (release.id)}
						<tr>
							<td>{release.title}</td>
							<td>{release.indexerName}</td>
							<td>{sizeLabel(release.sizeBytes)}</td>
							<td>{release.seeders ?? '-'}</td>
							<td class="row-actions">
								{#if canManage}
									<button
										type="button"
										disabled={grabbingKey === releaseKey(release)}
										onclick={() => onGrab(item, release)}
									>
										{grabbingKey === releaseKey(release) ? 'Queueing' : 'Grab'}
									</button>
								{/if}
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="5" class="empty">
								{releaseResults?.loaded
									? 'No release candidates found.'
									: 'No search results loaded.'}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
