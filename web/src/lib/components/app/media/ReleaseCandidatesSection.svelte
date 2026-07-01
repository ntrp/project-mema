<script lang="ts">
	import type { MediaItem, ReleaseCandidate, ReleaseSearchState } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		releaseResults?: ReleaseSearchState;
		grabbingKey?: string;
		canManage: boolean;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
	}

	let { item, releaseResults, grabbingKey, canManage, onGrabRelease }: Props = $props();
	const shouldShow = $derived(Boolean(releaseResults?.loaded || releaseResults?.errors.length));

	function releaseKey(mediaItem: MediaItem, release: ReleaseCandidate) {
		return `${mediaItem.id}:${release.id}`;
	}

	function sizeLabel(sizeBytes: number) {
		if (!sizeBytes) {
			return 'Unknown size';
		}
		const gib = sizeBytes / 1024 / 1024 / 1024;
		return `${gib.toFixed(gib >= 10 ? 0 : 1)} GiB`;
	}
</script>

{#if shouldShow}
	<section aria-labelledby="release-candidates-title">
		<h2 id="release-candidates-title">Release Candidates</h2>

		{#if releaseResults?.errors.length}
			<div class="inline-errors">
				{#each releaseResults.errors as searchError (searchError)}
					<p>{searchError}</p>
				{/each}
			</div>
		{/if}

		{#if releaseResults?.loaded}
			<div class="release-list standalone">
				{#each releaseResults.releases as release (release.id)}
					<div class="release-row">
						<div>
							<strong>{release.title}</strong>
							<span>
								{release.indexerName} · {sizeLabel(release.sizeBytes)}
								{release.seeders !== undefined ? ` · ${release.seeders} seeders` : ''}
							</span>
						</div>
						{#if canManage}
							<button
								type="button"
								disabled={grabbingKey === releaseKey(item, release)}
								onclick={() => onGrabRelease(item, release)}
							>
								{grabbingKey === releaseKey(item, release) ? 'Queueing' : 'Grab'}
							</button>
						{/if}
					</div>
				{:else}
					<p class="empty">No release candidates stored yet</p>
				{/each}
			</div>
		{/if}
	</section>
{/if}
