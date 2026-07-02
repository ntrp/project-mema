<script lang="ts">
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
	import CopyIcon from '@lucide/svelte/icons/copy';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaItem, ReleaseCandidate } from '$lib/settings/types';
	import ReleaseMatchInfo from './ReleaseMatchInfo.svelte';
	import {
		ageLabel,
		languageLabels,
		peerLabel,
		qualityMatch,
		releaseSource,
		releaseSourceBadgeClass,
		signedScore,
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
		onGrab: (_item: MediaItem, _release: ReleaseCandidate) => void;
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

	const columns: { key: ReleaseSortKey; label: string }[] = [
		{ key: 'source', label: 'Protocol' },
		{ key: 'indexer', label: 'Indexer' },
		{ key: 'age', label: 'Age' },
		{ key: 'title', label: 'Title' },
		{ key: 'size', label: 'Size' },
		{ key: 'peers', label: 'Peers' },
		{ key: 'languages', label: 'Languages' },
		{ key: 'quality', label: 'Identified quality' },
		{ key: 'score', label: 'Score' },
		{ key: 'match', label: 'Match info' }
	];

	function releaseKey(release: ReleaseCandidate) {
		return `${item.id}:${release.id}`;
	}

	function grabDisabled(release: ReleaseCandidate) {
		return grabbingKey === releaseKey(release) || release.match.severity === 'error';
	}

	function grabTooltip(release: ReleaseCandidate) {
		if (release.match.severity === 'error') {
			return 'Cannot grab a release that does not match this series/movie';
		}
		if (grabbingKey === releaseKey(release)) {
			return 'Queueing release';
		}
		return 'Grab release';
	}

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
		<Table.Root class="min-w-[1180px]">
			<Table.Header class="sticky top-0 z-10 bg-card">
				<Table.Row>
					{#each columns as column (column.key)}
						<Table.Head>
							<Button
								type="button"
								variant="ghost"
								class="-ml-2 h-8 px-2"
								onclick={() => onSort(column.key)}
							>
								<span>{column.label}</span>
								{#if sort.key === column.key}
									{#if sort.direction === 'asc'}
										<ArrowUpIcon aria-hidden="true" />
									{:else}
										<ArrowDownIcon aria-hidden="true" />
									{/if}
								{/if}
							</Button>
						</Table.Head>
					{/each}
					<Table.Head>Actions</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each releases as release (release.id)}
					<Table.Row>
						<Table.Cell>
							<Badge variant="outline" class={`uppercase ${releaseSourceBadgeClass(release)}`}>
								{releaseSource(release)}
							</Badge>
						</Table.Cell>
						<Table.Cell>{release.indexerName}</Table.Cell>
						<Table.Cell>{ageLabel(release)}</Table.Cell>
						<Table.Cell class="max-w-96">
							<div class="flex max-w-96 items-center gap-1.5">
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<span {...props} class="block min-w-0 flex-1 truncate">{release.title}</span>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content class="max-w-160">{release.title}</Tooltip.Content>
								</Tooltip.Root>
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<Button
												{...props}
												type="button"
												variant="ghost"
												size="icon-sm"
												aria-label="Copy release title"
												onclick={() => void copyTitle(release)}
											>
												<CopyIcon aria-hidden="true" />
											</Button>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content>
										{copiedReleaseId === release.id ? 'Copied title' : 'Copy title'}
									</Tooltip.Content>
								</Tooltip.Root>
							</div>
						</Table.Cell>
						<Table.Cell>{sizeLabel(release.sizeBytes)}</Table.Cell>
						<Table.Cell>{releaseSource(release) === 'torrent' ? peerLabel(release) : ''}</Table.Cell
						>
						<Table.Cell>
							<div class="flex max-w-56 flex-wrap gap-1">
								{#each languageLabels(release) as language (language)}
									<Badge variant="secondary" class="bg-muted text-muted-foreground"
										>{language}</Badge
									>
								{/each}
							</div>
						</Table.Cell>
						<Table.Cell>
							<Badge variant="secondary" class="bg-muted text-muted-foreground">
								{qualityMatch(release).label}
							</Badge>
						</Table.Cell>
						<Table.Cell>{signedScore(qualityMatch(release).score)}</Table.Cell>
						<Table.Cell>
							<ReleaseMatchInfo info={release.match} />
						</Table.Cell>
						<Table.Cell class="text-right">
							{#if canManage}
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<Button
												{...props}
												type="button"
												size="icon-sm"
												aria-label="Grab release"
												disabled={grabDisabled(release)}
												onclick={() => onGrab(item, release)}
											>
												<DownloadIcon aria-hidden="true" />
											</Button>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content>
										{grabTooltip(release)}
									</Tooltip.Content>
								</Tooltip.Root>
							{/if}
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
</Card>
