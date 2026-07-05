<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaItem, ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';
	import ReleaseGrabActions from '$lib/components/app/media/release-display/ReleaseGrabActions.svelte';
	import ReleaseMatchInfo from '$lib/components/app/media/release-display/ReleaseMatchInfo.svelte';
	import ReleaseScoreCell from '$lib/components/app/media/release-display/ReleaseScoreCell.svelte';
	import ReleaseTitleCell from '$lib/components/app/media/release-display/ReleaseTitleCell.svelte';
	import {
		ageLabel,
		peerBadgeClass,
		peerLabel,
		qualityMatch,
		releaseSource,
		releaseSourceBadgeClass,
		sizeLabel
	} from '$lib/components/app/media/release-display/releaseCandidateDisplay';

	interface Props {
		item: MediaItem;
		release: ReleaseCandidate;
		copiedReleaseId?: string;
		grabbingKey?: string;
		canManage: boolean;
		onCopy: (_release: ReleaseCandidate) => void;
		onGrab: (
			_item: MediaItem,
			_release: ReleaseCandidate,
			_overrideMatch?: boolean,
			_details?: ReleaseOverrideDetails
		) => void;
	}

	let { item, release, copiedReleaseId, grabbingKey, canManage, onCopy, onGrab }: Props = $props();

	const source = $derived(releaseSource(release));
	const peers = $derived(peerLabel(release));
	const releaseSources = $derived(release.sources ?? []);
	const sourceCount = $derived(releaseSources.length);
</script>

<Table.Row>
	<Table.Cell class="whitespace-nowrap">
		<Badge
			variant="outline"
			class={`relative overflow-visible uppercase ${releaseSourceBadgeClass(release)}`}
		>
			{source}
			{#if source === 'torrent' && peers !== '-'}
				<span
					class={`absolute -right-2 -bottom-2 rounded-[3px] border px-1 text-[9px] leading-3 font-black shadow-sm ${peerBadgeClass(release)}`}
				>
					{peers}
				</span>
			{/if}
		</Badge>
	</Table.Cell>
	<Table.Cell class="max-w-[180px] whitespace-nowrap">
		{#if sourceCount > 1}
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<span {...props} class="flex min-w-0 items-center gap-1">
							<span class="min-w-0 truncate">{release.indexerName}</span>
							<Badge variant="outline" class="px-1 py-0 text-[10px] leading-4">
								+{sourceCount - 1}
							</Badge>
						</span>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content class="max-w-96">
					<div class="space-y-2">
						{#each releaseSources as releaseSource (releaseSource.downloadUrl)}
							<div class="space-y-0.5">
								<div class="font-semibold">
									{releaseSource.indexerName}
									<span class="text-muted-foreground">({releaseSource.indexerProtocol})</span>
								</div>
								<div class="text-muted-foreground truncate">{releaseSource.downloadUrl}</div>
							</div>
						{/each}
					</div>
				</Tooltip.Content>
			</Tooltip.Root>
		{:else}
			<span class="block truncate">{release.indexerName}</span>
		{/if}
	</Table.Cell>
	<Table.Cell class="whitespace-nowrap">{ageLabel(release)}</Table.Cell>
	<Table.Cell class="w-full min-w-0 max-w-0">
		<ReleaseTitleCell {release} {copiedReleaseId} {onCopy} />
	</Table.Cell>
	<Table.Cell class="whitespace-nowrap">{sizeLabel(release.sizeBytes)}</Table.Cell>
	<Table.Cell class="whitespace-nowrap">
		<Badge variant="secondary" class="bg-muted text-muted-foreground">
			{qualityMatch(release).label}
		</Badge>
	</Table.Cell>
	<Table.Cell class="whitespace-nowrap"><ReleaseScoreCell match={release.match} /></Table.Cell>
	<Table.Cell class="whitespace-nowrap">
		<ReleaseMatchInfo info={release.match} mediaType={item.type} />
	</Table.Cell>
	<Table.Cell class="text-right">
		{#if canManage}
			<ReleaseGrabActions {item} {release} {grabbingKey} {onGrab} />
		{/if}
	</Table.Cell>
</Table.Row>
