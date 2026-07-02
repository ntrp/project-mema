<script lang="ts">
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';
	import type { MediaSearchResult } from '$lib/settings/types';

	interface Props {
		result: MediaSearchResult;
		adding?: boolean;
		label: string;
		class?: string;
		size?: 'default' | 'sm' | 'lg' | 'icon' | 'icon-sm' | 'icon-lg';
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { result, adding = false, label, class: className, size, onAdd }: Props = $props();
</script>

<Button
	type="button"
	class={className}
	{size}
	disabled={adding}
	aria-label={adding ? 'Working' : label}
	onclick={(event) => {
		event.stopPropagation();
		onAdd(result);
	}}
>
	{#if adding}
		<LoaderCircleIcon class="animate-spin" aria-hidden="true" />
	{:else}
		<PlusIcon aria-hidden="true" />
	{/if}
	<span>{adding ? 'Working' : label}</span>
</Button>
