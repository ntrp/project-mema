<script lang="ts">
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import ArrowLeftRightIcon from '@lucide/svelte/icons/arrow-left-right';
	import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
	import BadgeCheckIcon from '@lucide/svelte/icons/badge-check';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import HardDriveIcon from '@lucide/svelte/icons/hard-drive';
	import ListChecksIcon from '@lucide/svelte/icons/list-checks';
	import ServerIcon from '@lucide/svelte/icons/server';
	import StarIcon from '@lucide/svelte/icons/star';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaItem, ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';
	import ReleaseSearchResultRow from '$lib/components/app/media/release-search/ReleaseSearchResultRow.svelte';
	import type {
		ReleaseSort,
		ReleaseSortKey
	} from '$lib/components/app/media/release-display/releaseSearchResults';
	import { containWheelBoundary } from '$lib/components/app/media/shared/scrollBoundary';

	interface Props {
		item: MediaItem;
		releases: ReleaseCandidate[];
		searching?: boolean;
		sort: ReleaseSort;
		grabbingKey?: string;
		canManage: boolean;
		onSort: (_key: ReleaseSortKey) => void;
		onGrab: (
			_item: MediaItem,
			_release: ReleaseCandidate,
			_overrideMatch?: boolean,
			_details?: ReleaseOverrideDetails
		) => void;
	}

	let {
		item,
		releases,
		searching = false,
		sort,
		grabbingKey,
		canManage,
		onSort,
		onGrab
	}: Props = $props();

	let copiedReleaseId = $state<string | undefined>();
	let renderLimit = $state(80);
	const releaseSignature = $derived(releases.map((release) => release.id).join('\0'));
	const renderedReleases = $derived(releases.slice(0, renderLimit));
	const hiddenReleaseCount = $derived(Math.max(0, releases.length - renderedReleases.length));

	const columns: { key: ReleaseSortKey; label: string; icon?: string }[] = [
		{ key: 'source', label: 'Protocol', icon: 'transfer' },
		{ key: 'indexer', label: 'Indexer', icon: 'server' },
		{ key: 'age', label: 'Age', icon: 'time' },
		{ key: 'title', label: 'Title' },
		{ key: 'size', label: 'Size', icon: 'size' },
		{ key: 'quality', label: 'Quality', icon: 'quality' },
		{ key: 'score', label: 'Score', icon: 'score' },
		{ key: 'match', label: 'Match', icon: 'match' }
	];

	async function copyTitle(release: ReleaseCandidate) {
		await globalThis.navigator.clipboard.writeText(release.title);
		copiedReleaseId = release.id;
		globalThis.setTimeout(() => {
			if (copiedReleaseId === release.id) {
				copiedReleaseId = undefined;
			}
		}, 1200);
	}

	$effect(() => {
		releaseSignature;
		renderLimit = 80;
	});
</script>

<Card
	class="min-h-64 max-h-[min(520px,calc(100vh-340px))] overflow-auto overscroll-contain p-0"
	onwheel={containWheelBoundary}
>
	{#if searching && releases.length === 0}
		<div class="grid min-h-64 place-items-center">
			<InlineSpinner label="Searching releases" />
		</div>
	{:else}
		<Table.Root class="table-auto">
			<colgroup>
				<col class="w-[1%]" />
				<col class="w-[1%]" />
				<col class="w-[1%]" />
				<col class="w-full" />
				<col class="w-[1%]" />
				<col class="w-[1%]" />
				<col class="w-[1%]" />
				<col class="w-[1%]" />
				<col class="w-[1%]" />
			</colgroup>
			<Table.Header class="sticky top-0 z-10 bg-card">
				<Table.Row>
					{#each columns as column (column.key)}
						<Table.Head class={column.key === 'title' ? 'w-full min-w-0 max-w-0' : 'w-[1%]'}>
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="ghost"
											class="-ml-2 h-8 px-2"
											aria-label={`Sort by ${column.label}`}
											onclick={() => onSort(column.key)}
										>
											{#if column.icon === 'transfer'}
												<ArrowLeftRightIcon aria-hidden="true" />
												<span class="sr-only">{column.label}</span>
											{:else if column.icon === 'time'}
												<ClockIcon aria-hidden="true" />
												<span class="sr-only">{column.label}</span>
											{:else if column.icon === 'size'}
												<HardDriveIcon aria-hidden="true" />
												<span class="sr-only">{column.label}</span>
											{:else if column.icon === 'server'}
												<ServerIcon aria-hidden="true" />
												<span class="sr-only">{column.label}</span>
											{:else if column.icon === 'quality'}
												<BadgeCheckIcon aria-hidden="true" />
												<span class="sr-only">{column.label}</span>
											{:else if column.icon === 'score'}
												<StarIcon aria-hidden="true" />
												<span class="sr-only">{column.label}</span>
											{:else if column.icon === 'match'}
												<ListChecksIcon aria-hidden="true" />
												<span class="sr-only">{column.label}</span>
											{:else}
												<span>{column.label}</span>
											{/if}
											{#if sort.key === column.key}
												{#if sort.direction === 'asc'}
													<ArrowUpIcon aria-hidden="true" />
												{:else}
													<ArrowDownIcon aria-hidden="true" />
												{/if}
											{/if}
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>
									{column.label}
								</Tooltip.Content>
							</Tooltip.Root>
						</Table.Head>
					{/each}
					<Table.Head />
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each renderedReleases as release (release.id)}
					<ReleaseSearchResultRow
						{item}
						{release}
						{copiedReleaseId}
						{grabbingKey}
						{canManage}
						onCopy={(value) => void copyTitle(value)}
						{onGrab}
					/>
				{/each}
				{#if hiddenReleaseCount > 0}
					<Table.Row>
						<Table.Cell colspan={9} class="py-3 text-center">
							<Button type="button" variant="outline" onclick={() => (renderLimit += 80)}>
								Show {Math.min(80, hiddenReleaseCount)} more
								<span class="text-muted-foreground">({hiddenReleaseCount} hidden)</span>
							</Button>
						</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
	{/if}
</Card>
