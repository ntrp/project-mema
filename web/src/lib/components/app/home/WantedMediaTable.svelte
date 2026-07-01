<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MediaItem } from '$lib/settings/types';

	interface Props {
		items: MediaItem[];
		searchingItemId?: string;
		canManage: boolean;
		onFindReleases: (_item: MediaItem) => void;
	}

	let { items, searchingItemId, canManage, onFindReleases }: Props = $props();

	function monitorLabel(item: MediaItem) {
		if (!item.monitored || item.monitorMode === 'none') {
			return 'None';
		}
		return item.monitorMode === 'collection' ? 'Entire collection' : 'This media only';
	}
</script>

<div class="page-heading">
	<p>Library</p>
	<h1 id="home-title">Wanted</h1>
</div>

{#if items.length}
	<div class="table-wrap">
		<table class="data-table wanted-table">
			<thead>
				<tr>
					<th>Title</th>
					<th>Type</th>
					<th>Year</th>
					<th>Monitor</th>
					<th>Profile</th>
					<th>Availability</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each items as item (item.id)}
					<tr>
						<td>
							<a
								href={item.type === 'movie'
									? resolve('/movies/[id]', { id: item.id })
									: resolve('/series/[id]', { id: item.id })}
							>
								{item.title}
							</a>
						</td>
						<td>{item.type}</td>
						<td>{item.year ?? '-'}</td>
						<td>{monitorLabel(item)}</td>
						<td>{item.qualityProfileName ?? '-'}</td>
						<td>{item.minimumAvailability}</td>
						<td>
							{#if canManage}
								<button
									type="button"
									class="secondary compact-action"
									disabled={searchingItemId === item.id}
									onclick={() => onFindReleases(item)}
								>
									{searchingItemId === item.id ? 'Searching' : 'Search'}
								</button>
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else}
	<div class="panel">
		<p class="empty">No missing media.</p>
	</div>
{/if}
