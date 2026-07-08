<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import MusicIcon from '@lucide/svelte/icons/music';
	import VideoIcon from '@lucide/svelte/icons/video';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import type { MediaFileSummaryStatus } from '$lib/components/app/media/files/mediaFileSummaryStatus';

	interface Props {
		type: 'video' | 'audio' | 'subtitle';
		status: MediaFileSummaryStatus;
	}

	let { type, status }: Props = $props();
	const label = $derived(`${typeLabel(type)} ${status.label}`);
	const detail = $derived(status.details.join('. '));
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
	<Tooltip.Content>
		<span class="grid gap-1">
			<strong>{label}</strong>
			{#each status.details as item (item)}
				<span>{item}</span>
			{/each}
		</span>
	</Tooltip.Content>
</Tooltip.Root>
