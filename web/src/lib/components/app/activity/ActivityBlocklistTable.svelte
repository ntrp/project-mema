<script lang="ts">
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import type { ReleaseBlocklistItem } from '$lib/settings/types';

	interface Props {
		items: ReleaseBlocklistItem[];
	}

	let { items }: Props = $props();
</script>

<Card class="overflow-hidden p-0">
	<Table.Root class="[&_td]:whitespace-nowrap [&_th]:whitespace-nowrap">
		<Table.Header>
			<Table.Row>
				<Table.Head class="min-w-70 whitespace-normal">Release</Table.Head>
				<Table.Head>Media</Table.Head>
				<Table.Head>Indexer</Table.Head>
				<Table.Head>Reason</Table.Head>
				<Table.Head>Expires</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each items as block (block.id)}
				<Table.Row>
					<Table.Cell class="max-w-120 whitespace-normal">
						<strong>{block.releaseTitle}</strong>
						<small class="block text-xs text-muted-foreground">{block.source}</small>
					</Table.Cell>
					<Table.Cell>{block.mediaTitle}</Table.Cell>
					<Table.Cell>{block.indexerName}</Table.Cell>
					<Table.Cell class="max-w-80 whitespace-normal">{block.reason}</Table.Cell>
					<Table.Cell
						>{block.expiresAt ? new Date(block.expiresAt).toLocaleString() : '-'}</Table.Cell
					>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</Card>
