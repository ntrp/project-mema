<script lang="ts">
	import { Skeleton } from '$lib/components/ui/skeleton';
	import MediaMetadataShell from '$lib/components/app/media/metadata/MediaMetadataShell.svelte';
</script>

<MediaMetadataShell labelledby="media-detail-loading-title">
	<div class="grid gap-7" aria-busy="true" aria-live="polite">
		<div
			class="grid items-end gap-6.5 min-[641px]:grid-cols-[clamp(190px,18vw,270px)_minmax(0,1fr)]"
		>
			<Skeleton class="aspect-[2/3] rounded-md max-sm:w-[170px]" />
			<div class="grid gap-3">
				<h1 id="media-detail-loading-title" class="sr-only">Loading media details</h1>
				<Skeleton class="h-14 w-[min(720px,82%)] max-sm:h-11" />
				<Skeleton class="h-5 w-36" />
				<div class="flex flex-wrap gap-2">
					{#each Array.from({ length: 5 }) as _, index (index)}
						<Skeleton class="h-6 w-24 rounded-4xl" />
					{/each}
				</div>
				<div class="mt-1 flex flex-wrap gap-2.5">
					<Skeleton class="h-9 w-32" />
					<Skeleton class="h-9 w-32" />
				</div>
			</div>
		</div>

		<main class="grid min-w-0 gap-6">
			<section
				class="mt-4 grid items-start gap-5.5 min-[981px]:grid-cols-[minmax(0,1fr)_minmax(280px,390px)]"
				aria-hidden="true"
			>
				<div class="grid min-w-0 gap-3">
					<Skeleton class="h-8 w-36" />
					<Skeleton class="h-4 w-full" />
					<Skeleton class="h-4 w-[94%]" />
					<Skeleton class="h-4 w-[78%]" />
					<Skeleton class="mt-2 h-7 w-24" />
					<div class="grid gap-2 md:grid-cols-3">
						{#each Array.from({ length: 6 }) as _, index (index)}
							<Skeleton class="h-12 w-full" />
						{/each}
					</div>
				</div>
				<div class="grid gap-2 rounded-md border border-border bg-card p-4">
					{#each Array.from({ length: 6 }) as _, index (index)}
						<Skeleton class="h-5 w-full" />
					{/each}
				</div>
			</section>

			<section class="grid min-w-0 gap-2.5" aria-hidden="true">
				<Skeleton class="h-8 w-28" />
				<div
					class="-mx-3.5 grid auto-cols-[minmax(190px,220px)] grid-flow-col gap-5 overflow-hidden px-3.5 py-2 max-sm:auto-cols-[minmax(128px,150px)] max-sm:gap-3"
				>
					{#each Array.from({ length: 6 }) as _, index (index)}
						<Skeleton class="aspect-[2/3] rounded-md" />
					{/each}
				</div>
			</section>
		</main>
	</div>
</MediaMetadataShell>
