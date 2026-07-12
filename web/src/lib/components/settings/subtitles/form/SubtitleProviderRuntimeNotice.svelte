<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import type { components } from '$lib/api/generated/schema';
	import { dependencyLabels, runtimeLabel } from '../catalog/subtitleProviderCatalogFilters';
	type SubtitleProviderCatalogEntry = components['schemas']['SubtitleProviderCatalogEntry'];

	interface Props {
		entry?: SubtitleProviderCatalogEntry;
	}

	let { entry }: Props = $props();
	const dependencies = $derived(entry ? dependencyLabels(entry) : []);
	const status = $derived(entry?.runtimeStatus ?? 'supported');
	const message = $derived(entry?.runtimeMessage ?? 'Runtime support is available.');
	const warning = $derived(entry?.warning ?? '');

	function statusClass() {
		if (status === 'supported') return 'border-emerald-500/50 bg-emerald-500/10 text-emerald-300';
		if (status === 'catalog_only') return 'border-amber-500/50 bg-amber-500/10 text-amber-300';
		return 'border-destructive/50 bg-destructive/10 text-destructive';
	}
</script>

<div class="rounded-md border border-border bg-muted/30 p-3 text-sm">
	<div class="mb-2 flex flex-wrap items-center gap-2">
		<Badge variant="outline" class={statusClass()}>{runtimeLabel(status)}</Badge>
		{#each dependencies as dependency (dependency)}
			<Badge variant="outline" class="border-sky-500/50 bg-sky-500/10 text-sky-300">
				{dependency}
			</Badge>
		{/each}
	</div>
	<p class="m-0 text-muted-foreground">{message}</p>
	{#if warning}
		<p class="m-0 mt-2 text-amber-300">{warning}</p>
	{/if}
	{#if status !== 'supported'}
		<p class="m-0 mt-2 text-muted-foreground">
			This entry can be saved only while disabled, and testing is unavailable until runtime support is added.
		</p>
	{/if}
</div>
