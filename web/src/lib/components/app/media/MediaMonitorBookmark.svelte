<script lang="ts">
	import BookmarkIcon from '@lucide/svelte/icons/bookmark';
	import { cn } from '$lib/utils';

	interface Props {
		monitored?: boolean;
		label?: string;
	}

	let { monitored = false, label }: Props = $props();
	const title = $derived(label ?? (monitored ? 'Monitored' : 'Not monitored'));
</script>

<span
	class={cn(
		'inline-flex shrink-0 items-center justify-center leading-none',
		monitored ? 'text-secondary-foreground' : 'text-muted-foreground'
	)}
	{title}
	aria-label={title}
>
	<BookmarkIcon aria-hidden="true" size="64" class={monitored ? 'fill-current' : undefined} />
</span>
