<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import { selectedFirst } from '$lib/components/shared/multiSelectOrdering';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Label } from '$lib/components/ui/label';
	import type { MediaMetadataEpisode } from '$lib/settings/types';
	import { episodeNumbers, episodeValueFromNumbers } from './releaseOverrideSeriesOptions';
	import { episodeLabel } from './releaseOverrideSeriesOptions';

	interface Props {
		value: string;
		episodes: MediaMetadataEpisode[];
		onChange: (_value: string) => void;
	}

	let { value, episodes, onChange }: Props = $props();

	const selected = $derived(episodeNumbers(value));
	const selectedSet = $derived(new Set(selected));
	const sortedEpisodes = $derived(
		selectedFirst(episodes, selectedSet, (episode) => episode.episodeNumber)
	);
	const selectedBadges = $derived(
		selected.map((number) => {
			const episode = episodes.find((item) => item.episodeNumber === number);
			return episode ? episodeLabel(episode) : `E${String(number).padStart(2, '0')}`;
		})
	);

	function toggle(episodeNumber: number, checked: boolean) {
		const next = checked
			? [...selected, episodeNumber]
			: selected.filter((value) => value !== episodeNumber);
		onChange(episodeValueFromNumbers(next));
	}

	function clear() {
		onChange('');
	}
</script>

<div class="grid gap-1.5">
	<Label for="override-episodes">Episodes</Label>
	<DropdownMenu.Root>
		<DropdownMenu.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					id="override-episodes"
					type="button"
					variant="outline"
					class="min-h-9 w-full justify-between gap-2 py-1.5"
					disabled={episodes.length === 0}
				>
					<span class="flex min-w-0 flex-1 flex-wrap gap-1">
						{#if selectedBadges.length > 0}
							{#each selectedBadges as label (label)}
								<Badge variant="secondary" class="max-w-40 truncate">{label}</Badge>
							{/each}
						{:else}
							<span class="truncate text-muted-foreground">Select episodes</span>
						{/if}
					</span>
					<ChevronDownIcon class="size-4 shrink-0 text-muted-foreground" aria-hidden="true" />
				</Button>
			{/snippet}
		</DropdownMenu.Trigger>
		<DropdownMenu.Content align="start" class="max-h-72 w-96">
			<DropdownMenu.Item onclick={clear}>
				<span class="text-muted-foreground">Clear episodes</span>
			</DropdownMenu.Item>
			<DropdownMenu.Separator />
			{#each sortedEpisodes as episode (episode.episodeNumber)}
				<DropdownMenu.CheckboxItem
					checked={selectedSet.has(episode.episodeNumber)}
					onCheckedChange={(checked) => toggle(episode.episodeNumber, checked === true)}
				>
					<span class="truncate">{episodeLabel(episode)}</span>
				</DropdownMenu.CheckboxItem>
			{/each}
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</div>
