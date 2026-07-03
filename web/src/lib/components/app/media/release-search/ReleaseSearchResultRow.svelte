<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import type { MediaItem, ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';
	import ReleaseGrabActions from '$lib/components/app/media/release-display/ReleaseGrabActions.svelte';
	import ReleaseMatchInfo from '$lib/components/app/media/release-display/ReleaseMatchInfo.svelte';
	import ReleaseScoreCell from '$lib/components/app/media/release-display/ReleaseScoreCell.svelte';
	import ReleaseTitleCell from '$lib/components/app/media/release-display/ReleaseTitleCell.svelte';
	import {
		ageLabel,
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
					class="absolute -right-2 -bottom-2 rounded-[3px] border border-background bg-background px-1 text-[9px] leading-3 font-black text-foreground shadow-sm"
				>
					{peers}
				</span>
			{/if}
		</Badge>
	</Table.Cell>
	<Table.Cell class="max-w-[160px] truncate whitespace-nowrap">{release.indexerName}</Table.Cell>
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
