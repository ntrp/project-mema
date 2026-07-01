<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import type { MediaItem, ReleaseCandidate, ReleaseSearchState } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		releaseResults?: ReleaseSearchState;
		searching?: boolean;
		grabbingKey?: string;
		canManage: boolean;
		onSearch: (_item: MediaItem) => void;
		onGrab: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onClose: () => void;
	}

	let {
		item,
		releaseResults,
		searching = false,
		grabbingKey,
		canManage,
		onSearch,
		onGrab,
		onClose
	}: Props = $props();

	function releaseKey(release: ReleaseCandidate) {
		return `${item.id}:${release.id}`;
	}

	function sizeLabel(sizeBytes: number) {
		if (!sizeBytes) return '-';
		const gib = sizeBytes / 1024 / 1024 / 1024;
		return `${gib.toFixed(gib >= 10 ? 0 : 1)} GiB`;
	}
</script>

<SettingsFormModal title="Manual search" modalClass="w-[min(1960px,calc(100vw-32px))]" {onClose}>
	<div class="flex justify-end">
		<Button type="button" disabled={!canManage || searching} onclick={() => onSearch(item)}>
			{searching ? 'Searching' : 'Search releases'}
		</Button>
	</div>
	{#if releaseResults?.errors.length}
		<div class="grid gap-1 rounded-md bg-secondary px-3 py-2.5 font-bold text-secondary-foreground">
			{#each releaseResults.errors as error (error)}
				<p class="m-0">{error}</p>
			{/each}
		</div>
	{/if}
	<Card class="overflow-hidden p-0">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Release</Table.Head>
					<Table.Head>Indexer</Table.Head>
					<Table.Head>Size</Table.Head>
					<Table.Head>Seeders</Table.Head>
					<Table.Head></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each releaseResults?.releases ?? [] as release (release.id)}
					<Table.Row>
						<Table.Cell>{release.title}</Table.Cell>
						<Table.Cell>{release.indexerName}</Table.Cell>
						<Table.Cell>{sizeLabel(release.sizeBytes)}</Table.Cell>
						<Table.Cell>{release.seeders ?? '-'}</Table.Cell>
						<Table.Cell class="text-right">
							{#if canManage}
								<Button
									type="button"
									size="sm"
									disabled={grabbingKey === releaseKey(release)}
									onclick={() => onGrab(item, release)}
								>
									{grabbingKey === releaseKey(release) ? 'Queueing' : 'Grab'}
								</Button>
							{/if}
						</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={5} class="text-muted-foreground">
							{releaseResults?.loaded
								? 'No release candidates found.'
								: 'No search results loaded.'}
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Card>
</SettingsFormModal>
