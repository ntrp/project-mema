<script lang="ts">
	import InfoIcon from '@lucide/svelte/icons/info';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils';
	import {
		mediaFileInfoSections,
		mediaTrackActionClass,
		type MediaFilePlaybackStats,
		type MediaFilePreviewInfo
	} from '$lib/components/app/media/files/preview/mediaFilePreviewInfo';

	interface Props {
		info?: MediaFilePreviewInfo;
		playbackStats?: MediaFilePlaybackStats;
		loading: boolean;
		error?: string;
	}

	let { info, playbackStats, loading, error }: Props = $props();
	const sections = $derived(mediaFileInfoSections(info, playbackStats));
</script>

<section
	class="pointer-events-auto grid w-[min(18rem,calc(100vw-3rem))] gap-2 rounded-md border border-white/15 bg-black/75 p-2 font-mono text-[10px] leading-tight text-white shadow-2xl backdrop-blur-sm"
	aria-live="polite"
>
	<div class="flex flex-wrap items-center justify-between gap-2">
		<h3 class="m-0 inline-flex items-center gap-1.5 text-[10px] font-semibold text-white uppercase">
			<InfoIcon aria-hidden="true" />
			Media info
		</h3>
		{#if loading}
			<span class="inline-flex items-center gap-1 text-[9px] text-white/70">
				<LoaderCircleIcon class="size-2.5 animate-spin" aria-hidden="true" />
				Refreshing
			</span>
		{/if}
	</div>
	{#if error}
		<p
			class="m-0 rounded border border-red-300/30 bg-red-500/20 px-2 py-1.5 text-[10px] text-red-100"
		>
			{error}
		</p>
	{/if}
	<div class="grid gap-2">
		{#each sections as section (section.key)}
			<div class="grid content-start gap-1">
				<h4
					class="m-0 inline-flex items-center gap-1.5 text-[10px] font-semibold tracking-normal text-white/80 uppercase"
				>
					{section.title}
					{#if section.action}
						<Badge
							variant="outline"
							class={cn(
								'h-4 rounded px-1 font-mono text-[9px] uppercase',
								mediaTrackActionClass(section.action)
							)}
						>
							{section.action}
						</Badge>
					{/if}
				</h4>
				<dl class="m-0 grid gap-1">
					{#each section.rows as row (row.label)}
						<div class="grid grid-cols-[4.6rem_minmax(0,1fr)] gap-1.5">
							<dt class="text-white/55">{row.label}</dt>
							<dd class="m-0 break-anywhere font-medium text-white">{row.value}</dd>
						</div>
					{/each}
				</dl>
			</div>
		{/each}
	</div>
</section>
