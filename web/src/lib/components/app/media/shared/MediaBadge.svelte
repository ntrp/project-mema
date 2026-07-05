<script lang="ts">
	import { cn } from '$lib/utils';
	import type { MediaType } from '$lib/settings/types';
	import { mediaBadgeToneClass } from './mediaBadge';

	interface Props {
		type: MediaType;
		inLibrary?: boolean;
	}

	let { type, inLibrary = false }: Props = $props();

	const label = $derived(inLibrary ? 'In library' : type);
	const toneClass = $derived(mediaBadgeToneClass(type));
</script>

<span
	class={cn(
		'absolute top-2 left-2 z-[3] rounded-[3px] border px-[7px] py-1 text-[11px] font-black uppercase',
		toneClass
	)}
>
	{label}
</span>
