<script lang="ts">
	import ExternalLinkIcon from '@lucide/svelte/icons/external-link';
	import RssIcon from '@lucide/svelte/icons/rss';
	import SearchIcon from '@lucide/svelte/icons/search';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { IndexerCatalogEntry } from '$lib/settings/types';

	interface Props {
		entry: IndexerCatalogEntry;
	}

	let { entry }: Props = $props();
</script>

<div class="flex items-center gap-1.5">
	{#if entry.supportsSearch}
		<Tooltip.Root>
			<Tooltip.Trigger>
				<SearchIcon aria-label="Search" class="size-4 text-muted-foreground" />
			</Tooltip.Trigger>
			<Tooltip.Content>Search</Tooltip.Content>
		</Tooltip.Root>
	{/if}
	{#if entry.supportsRss}
		<Tooltip.Root>
			<Tooltip.Trigger>
				<RssIcon aria-label="RSS" class="size-4 text-muted-foreground" />
			</Tooltip.Trigger>
			<Tooltip.Content>RSS</Tooltip.Content>
		</Tooltip.Root>
	{/if}
	{#if entry.supportsRedirect}
		<Tooltip.Root>
			<Tooltip.Trigger>
				<ExternalLinkIcon aria-label="Redirect" class="size-4 text-muted-foreground" />
			</Tooltip.Trigger>
			<Tooltip.Content>Redirect</Tooltip.Content>
		</Tooltip.Root>
	{/if}
</div>
