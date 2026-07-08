<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import MusicIcon from '@lucide/svelte/icons/music';
	import VideoIcon from '@lucide/svelte/icons/video';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import { displayLanguage } from '$lib/settings/languageDisplay';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { MediaFileSummaryStatus } from '$lib/components/app/media/files/mediaFileSummaryStatus';

	interface Props {
		type: 'video' | 'audio' | 'subtitle';
		status: MediaFileSummaryStatus;
		row: MediaFileRow;
	}

	interface TooltipSection {
		title: string;
		details: string[];
	}

	let { type, status, row }: Props = $props();
	const label = $derived(`${typeLabel(type)} ${status.label}`);
	const sections = $derived(statusTooltipSections(type, status, row));
	const detail = $derived(
		sections.flatMap((section) => [section.title, ...section.details]).join('. ')
	);
	const iconClass = $derived(
		cn(
			'size-4',
			status.state === 'satisfied' && 'text-emerald-600',
			(status.state === 'partial' || status.state === 'pending') && 'text-orange-500',
			status.state === 'missing' && 'text-destructive',
			status.state === 'ignored' && 'text-muted-foreground'
		)
	);

	function typeLabel(value: Props['type']) {
		if (value === 'video') return 'Video';
		if (value === 'audio') return 'Audio';
		return 'Subtitles';
	}

	function statusTooltipSections(
		value: Props['type'],
		summary: MediaFileSummaryStatus,
		fileRow: MediaFileRow
	): TooltipSection[] {
		const sections: TooltipSection[] = [];
		const trackType = value === 'subtitle' ? 'subtitle' : value;
		for (const track of fileRow.tracks.filter((track) => track.type === trackType)) {
			if (summary.state !== 'satisfied' && track.state?.visualState === 'matching') continue;
			if (!track.state?.details?.length) continue;
			sections.push({
				title: trackTitle(track),
				details: track.state.details
			});
		}
		if (value === 'subtitle') {
			for (const file of fileRow.otherFiles.filter((file) => file.type === 'subtitle')) {
				if (summary.state !== 'satisfied' && file.state?.visualState === 'matching') continue;
				if (!file.state?.details?.length) continue;
				sections.push({
					title: sidecarTitle(file),
					details: file.state.details
				});
			}
		}
		for (const missing of fileRow.missingTracks.filter((missing) => missing.type === trackType)) {
			if (!missing.state.details.length) continue;
			sections.push({
				title: missing.description,
				details: missing.state.details
			});
		}
		if (sections.length > 0) return sections;
		return [{ title: label, details: summary.details }];
	}

	function trackTitle(track: MediaFileRow['tracks'][number]) {
		const number = track.index === undefined ? '-' : String(track.index);
		return `${typeLabel(track.type === 'subtitle' ? 'subtitle' : track.type)} track ${number}${languageSuffix(track.language)}`;
	}

	function sidecarTitle(file: MediaFileRow['otherFiles'][number]) {
		return `Subtitle sidecar${languageSuffix(file.language)}`;
	}

	function languageSuffix(language?: string) {
		const label = displayLanguage(language);
		return label === '-' ? '' : ` - ${label}`;
	}
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<span
				{...props}
				class="inline-flex h-6 w-6 items-center justify-center rounded-sm"
				aria-label={`${label}. ${detail}`}
			>
				{#if type === 'video'}
					<VideoIcon class={iconClass} aria-hidden="true" />
				{:else if type === 'audio'}
					<MusicIcon class={iconClass} aria-hidden="true" />
				{:else}
					<CaptionsIcon class={iconClass} aria-hidden="true" />
				{/if}
			</span>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content class="max-w-96">
		<span class="grid gap-2">
			<strong>{label}</strong>
			{#each sections as section (`${section.title}:${section.details.join('|')}`)}
				<span class="grid gap-1">
					<span class="font-medium">{section.title}</span>
					<ul class="list-disc space-y-1 pl-4">
						{#each section.details as item (item)}
							<li>{item}</li>
						{/each}
					</ul>
				</span>
			{/each}
		</span>
	</Tooltip.Content>
</Tooltip.Root>
