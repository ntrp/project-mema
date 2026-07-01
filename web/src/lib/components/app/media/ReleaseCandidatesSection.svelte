<script lang="ts">
	import { Button } from '$lib/components/ui/button';
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
		<h2 id="release-candidates-title" class="m-0 text-3xl font-semibold text-foreground">
			Release Candidates
		</h2>

		{#if releaseResults?.errors.length}
			<div
				class="grid gap-1 rounded-md bg-secondary px-3 py-2.5 font-bold text-secondary-foreground"
			>
				{#each releaseResults.errors as searchError (searchError)}
					<p class="m-0">{searchError}</p>
				{/each}
			</div>
		{/if}

		{#if releaseResults?.loaded}
			<div class="grid gap-2">
				{#each releaseResults.releases as release (release.id)}
					<div
						class="flex items-center justify-between gap-4 rounded-md border border-border bg-card p-2.5 max-sm:flex-col max-sm:items-stretch"
					>
						<div class="grid min-w-0 gap-1">
							<strong>{release.title}</strong>
							<span class="text-sm text-muted-foreground">
								{release.indexerName} · {sizeLabel(release.sizeBytes)}
								{release.seeders !== undefined ? ` · ${release.seeders} seeders` : ''}
							</span>
						</div>
						{#if canManage}
							<Button
								type="button"
								disabled={grabbingKey === releaseKey(item, release)}
								onclick={() => onGrabRelease(item, release)}
							>
								{grabbingKey === releaseKey(item, release) ? 'Queueing' : 'Grab'}
							</Button>
						{/if}
					</div>
				{:else}
					<p class="m-0 text-sm leading-6 text-muted-foreground">
						No release candidates stored yet
					</p>
				{/each}
			</div>
		{/if}
	</section>
{/if}
