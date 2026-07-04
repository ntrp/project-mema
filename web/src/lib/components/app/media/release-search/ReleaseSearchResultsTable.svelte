<script lang="ts">
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import type { MediaItem, ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';
	import ReleaseSearchResultRow from '$lib/components/app/media/release-search/ReleaseSearchResultRow.svelte';
	import ReleaseSearchResultsHeader from '$lib/components/app/media/release-search/table/ReleaseSearchResultsHeader.svelte';
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
		resetKey?: string;
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
		resetKey = '',
		grabbingKey,
		canManage,
		onSort,
		onGrab
	}: Props = $props();

	let copiedReleaseId = $state<string | undefined>();
	let renderLimit = $state(80);
	let scrollContainer = $state<HTMLDivElement | null>(null);
	const releaseSignature = $derived(releases.map((release) => release.id).join('\0'));
	const sortSignature = $derived(`${sort.key ?? ''}:${sort.direction}`);
	const renderedReleases = $derived(releases.slice(0, renderLimit));
	const hiddenReleaseCount = $derived(Math.max(0, releases.length - renderedReleases.length));

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
		sortSignature;
		resetKey;
		renderLimit = 80;
		if (scrollContainer) {
			scrollContainer.scrollTop = 0;
		}
	});
</script>

<Card
	bind:ref={scrollContainer}
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
			<ReleaseSearchResultsHeader {sort} {onSort} />
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
