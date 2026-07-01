<script lang="ts">
	import { cn } from '$lib/utils';
	import type { MediaType } from '$lib/settings/types';

	interface Props {
		type: MediaType;
		inLibrary?: boolean;
	}

	let { type, inLibrary = false }: Props = $props();

	const label = $derived(inLibrary ? 'In library' : type);
	const toneClass = $derived(
		type === 'movie'
			? 'border-yellow-500 bg-yellow-400 text-yellow-950'
			: 'border-blue-600 bg-blue-600 text-white'
	);
</script>

<span
	class={cn(
		'absolute top-2 left-2 z-[3] rounded-md border px-[7px] py-1 text-[11px] font-black uppercase',
		toneClass
	)}
>
	{label}
</span>
