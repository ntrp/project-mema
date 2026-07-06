<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import ClapperboardIcon from '@lucide/svelte/icons/clapperboard';
	import FileTextIcon from '@lucide/svelte/icons/file-text';
	import MusicIcon from '@lucide/svelte/icons/music';
	import VideoIcon from '@lucide/svelte/icons/video';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';

	interface Props {
		type: MediaFileDetailRow['type'];
	}

	let { type }: Props = $props();
	const label = $derived(trackTypeLabel(type));

	function trackTypeLabel(value: MediaFileDetailRow['type']) {
		switch (value) {
			case 'video':
				return 'Video track';
			case 'audio':
				return 'Audio track';
			case 'subtitle':
				return 'Subtitle track';
			case 'chapter':
				return 'Chapter';
			default:
				return 'Track';
		}
	}
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<span {...props} class="inline-flex items-center" aria-label={label}>
				{#if type === 'video'}
					<VideoIcon aria-hidden="true" />
				{:else if type === 'audio'}
					<MusicIcon aria-hidden="true" />
				{:else if type === 'subtitle'}
					<CaptionsIcon aria-hidden="true" />
				{:else if type === 'chapter'}
					<ClapperboardIcon aria-hidden="true" />
				{:else}
					<FileTextIcon aria-hidden="true" />
				{/if}
			</span>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content>{label}</Tooltip.Content>
</Tooltip.Root>
