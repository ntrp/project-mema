<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import type { components } from '$lib/api/generated/schema';
	import { dependencyLabels, runtimeLabel } from './subtitleProviderCatalogFilters';
	type SubtitleProviderCatalogEntry = components['schemas']['SubtitleProviderCatalogEntry'];

	interface Props {
		entries: SubtitleProviderCatalogEntry[];
		onSelect: (_entry: SubtitleProviderCatalogEntry) => void;
	}

	let { entries, onSelect }: Props = $props();

	function badgeClass(status: SubtitleProviderCatalogEntry['runtimeStatus']) {
		if (status === 'supported') return 'border-emerald-500/50 bg-emerald-500/10 text-emerald-300';
		if (status === 'catalog_only') return 'border-amber-500/50 bg-amber-500/10 text-amber-300';
		return 'border-destructive/50 bg-destructive/10 text-destructive';
	}
</script>

<div class="max-h-[min(560px,calc(100vh-300px))] overflow-auto rounded-md border border-border">
	<table class="w-full min-w-190 table-auto border-collapse text-sm">
		<thead class="sticky top-0 z-10 bg-card text-left text-xs font-extrabold text-muted-foreground">
			<tr class="border-b border-border">
				<th class="px-3 py-2">Provider</th>
				<th class="px-3 py-2">Support</th>
				<th class="px-3 py-2">Media</th>
				<th class="px-3 py-2">Requirements</th>
				<th class="w-px px-3 py-2">Action</th>
			</tr>
		</thead>
		<tbody>
			{#each entries as entry (entry.key)}
				<tr class="border-b border-border last:border-0 hover:bg-muted/60">
					<td class="px-3 py-2">
						<div class="font-bold text-foreground">{entry.displayName}</div>
						<div class="line-clamp-2 text-xs text-muted-foreground">{entry.runtimeMessage}</div>
					</td>
					<td class="px-3 py-2">
						<Badge variant="outline" class={badgeClass(entry.runtimeStatus)}>
							{runtimeLabel(entry.runtimeStatus)}
						</Badge>
					</td>
					<td class="px-3 py-2 text-xs uppercase text-muted-foreground">
						{entry.mediaTypes.join(', ')}
					</td>
					<td class="px-3 py-2">
						<div class="flex flex-wrap gap-1">
							{#each dependencyLabels(entry) as label (label)}
								<Badge variant="outline" class="border-sky-500/50 bg-sky-500/10 text-sky-300">
									{label}
								</Badge>
							{:else}
								<span class="text-xs text-muted-foreground">None</span>
							{/each}
						</div>
					</td>
					<td class="px-3 py-2">
						<Button type="button" size="sm" onclick={() => onSelect(entry)}>Configure</Button>
					</td>
				</tr>
			{:else}
				<tr>
					<td class="px-3 py-8 text-center text-muted-foreground" colspan="5">
						No subtitle providers match the selected filters.
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
