<script lang="ts">
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';

	interface Props {
		label: string;
	}

	let { label }: Props = $props();
</script>

<span class="inline-flex items-center gap-2 text-[13px] font-extrabold text-muted-foreground">
	<LoaderCircleIcon class="size-3.5 animate-spin text-primary" aria-hidden="true" />
	{label}
</span>
