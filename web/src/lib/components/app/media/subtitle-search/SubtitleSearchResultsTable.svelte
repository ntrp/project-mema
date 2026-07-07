<script lang="ts">
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Card } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { containWheelBoundary } from '$lib/components/app/media/shared/scrollBoundary';
	import SubtitleSearchResultRow from './SubtitleSearchResultRow.svelte';
	import type { SubtitleCandidate } from '$lib/settings/types';

	interface Props {
		candidates: SubtitleCandidate[];
		searching?: boolean;
		grabbingId?: string;
		canManage: boolean;
		onGrab: (_candidate: SubtitleCandidate) => void;
	}

	let { candidates, searching = false, grabbingId, canManage, onGrab }: Props = $props();
</script>

<Card
	class="min-h-56 max-h-[min(460px,calc(100vh-360px))] overflow-auto overscroll-contain p-0"
	onwheel={containWheelBoundary}
>
	{#if searching && candidates.length === 0}
		<div class="grid min-h-56 place-items-center">
			<InlineSpinner label="Searching subtitles" />
		</div>
	{:else}
		<Table.Root class="table-auto">
			<colgroup>
				<col class="w-[1%]" />
				<col class="w-[1%]" />
				<col class="w-full" />
				<col class="w-[1%]" />
				<col class="w-[1%]" />
				<col class="w-[1%]" />
			</colgroup>
			<Table.Header class="sticky top-0 z-10 bg-card">
				<Table.Row>
					<Table.Head>Protocol</Table.Head>
					<Table.Head>Provider</Table.Head>
					<Table.Head>Title / filename</Table.Head>
					<Table.Head>Language</Table.Head>
					<Table.Head>Match</Table.Head>
					<Table.Head><span class="sr-only">Actions</span></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each candidates as candidate (candidate.id)}
					<SubtitleSearchResultRow
						{candidate}
						grabbing={grabbingId === candidate.id}
						{canManage}
						{onGrab}
					/>
				{:else}
					<Table.Row>
						<Table.Cell colspan={6} class="py-8 text-center text-muted-foreground">
							No subtitle results.
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
</Card>
