<script lang="ts">
	import StatusPill from '$lib/components/shared/StatusPill.svelte';

	interface Props {
		clientCount: number;
		enabledClientCount: number;
		indexerCount: number;
		enabledIndexerCount: number;
	}

	let { clientCount, enabledClientCount, indexerCount, enabledIndexerCount }: Props = $props();
</script>

<header class="mb-6 grid gap-6 md:flex md:items-end md:justify-between">
	<div>
		<p class="mb-2 text-xs font-bold text-muted-foreground uppercase">Settings</p>
		<h1 class="m-0 text-[40px] leading-tight font-semibold text-foreground">
			Download clients and indexers
		</h1>
	</div>
	<div class="flex flex-wrap gap-2 md:justify-end">
		<StatusPill>{clientCount} clients</StatusPill>
		<StatusPill>{enabledClientCount} enabled</StatusPill>
		<StatusPill>{indexerCount} indexers</StatusPill>
		<StatusPill>{enabledIndexerCount} enabled</StatusPill>
	</div>
</header>
