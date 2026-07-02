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
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaItem, ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';
	import ReleaseGrabActions from './ReleaseGrabActions.svelte';
	import ReleaseMatchInfo from './ReleaseMatchInfo.svelte';
	import ReleaseScoreCell from './ReleaseScoreCell.svelte';
	import ReleaseTitleCell from './ReleaseTitleCell.svelte';
	import {
		ageLabel,
		peerLabel,
		qualityMatch,
		releaseSource,
		releaseSourceBadgeClass,
		sizeLabel
	} from './releaseCandidateDisplay';
	import type { ReleaseSort, ReleaseSortKey } from './releaseSearchResults';

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
</script>

<Card class="min-h-64 max-h-[min(520px,calc(100vh-340px))] overflow-auto p-0">
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
				{#each releases as release (release.id)}
					<Table.Row>
						<Table.Cell class="whitespace-nowrap">
							<Badge
								variant="outline"
								class={`relative overflow-visible uppercase ${releaseSourceBadgeClass(release)}`}
							>
								{releaseSource(release)}
								{#if releaseSource(release) === 'torrent' && peerLabel(release) !== '-'}
									<span
										class="absolute -right-2 -bottom-2 rounded-[3px] border border-background bg-background px-1 text-[9px] leading-3 font-black text-foreground shadow-sm"
									>
										{peerLabel(release)}
									</span>
								{/if}
							</Badge>
						</Table.Cell>
						<Table.Cell class="max-w-[160px] truncate whitespace-nowrap"
							>{release.indexerName}</Table.Cell
						>
						<Table.Cell class="whitespace-nowrap">{ageLabel(release)}</Table.Cell>
						<Table.Cell class="w-full min-w-0 max-w-0">
							<ReleaseTitleCell
								{release}
								{copiedReleaseId}
								onCopy={(value) => void copyTitle(value)}
							/>
						</Table.Cell>
						<Table.Cell class="whitespace-nowrap">{sizeLabel(release.sizeBytes)}</Table.Cell>
						<Table.Cell class="whitespace-nowrap">
							<Badge variant="secondary" class="bg-muted text-muted-foreground">
								{qualityMatch(release).label}
							</Badge>
						</Table.Cell>
						<Table.Cell class="whitespace-nowrap"
							><ReleaseScoreCell match={release.match} /></Table.Cell
						>
						<Table.Cell class="whitespace-nowrap">
							<ReleaseMatchInfo info={release.match} />
						</Table.Cell>
						<Table.Cell class="text-right">
							{#if canManage}
								<ReleaseGrabActions {item} {release} {grabbingKey} {onGrab} />
							{/if}
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
</Card>
