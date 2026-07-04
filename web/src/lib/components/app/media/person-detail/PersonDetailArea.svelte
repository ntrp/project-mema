<script lang="ts">
	import CalendarClockIcon from '@lucide/svelte/icons/calendar-clock';
	import LayoutGridIcon from '@lucide/svelte/icons/layout-grid';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { MediaItem, MediaSearchResult, PersonDetails } from '$lib/settings/types';
	import PersonAppearancesGrid from './PersonAppearancesGrid.svelte';
	import PersonAppearancesTimeline from './PersonAppearancesTimeline.svelte';
	import PersonDetailsPanel from './PersonDetailsPanel.svelte';
	import { filteredAppearances, type PersonAppearanceFilter } from './personDetail';

	interface Props {
		person?: PersonDetails;
		loading: boolean;
		mediaItems?: MediaItem[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { person, loading, mediaItems = [], addingKey, actionLabel, onAdd }: Props = $props();

	let filter = $state<PersonAppearanceFilter>('all');
	let viewMode = $state<'grid' | 'timeline'>('grid');
	const appearances = $derived(person?.appearances ?? []);
	const filtered = $derived(filteredAppearances(appearances, filter));
	const counts = $derived({
		all: appearances.length,
		movie: appearances.filter((appearance) => appearance.type === 'movie').length,
		series: appearances.filter((appearance) => appearance.type === 'serie').length
	});

	const filterOptions: { value: PersonAppearanceFilter; label: string }[] = [
		{ value: 'all', label: 'All' },
		{ value: 'movie', label: 'Movies' },
		{ value: 'series', label: 'Series' }
	];
	const viewOptions = [
		{ value: 'grid', label: 'Grid', icon: LayoutGridIcon },
		{ value: 'timeline', label: 'Timeline', icon: CalendarClockIcon }
	] as const;
</script>

{#if loading}
	<section class="min-h-[260px] rounded-md border border-border bg-card p-5">
		<p class="m-0 text-sm leading-6 text-muted-foreground">Loading person details</p>
	</section>
{:else if !person}
	<EmptyState title="Person not available" description="Could not load details for this person." />
{:else}
	<div class="grid gap-8">
		<section class="grid items-start gap-7">
			<PersonDetailsPanel {person} />

			<div class="grid min-w-0 gap-7">
				<section aria-labelledby="person-appearances-title" class="grid min-w-0 gap-4">
					<div class="flex flex-wrap items-center justify-between gap-3">
						<h2 id="person-appearances-title" class="m-0 text-3xl font-semibold text-foreground">
							Appearances
						</h2>
						<div class="flex flex-wrap items-center justify-end gap-2">
							<div
								class="inline-grid grid-cols-3 overflow-hidden rounded-md border border-border bg-card"
								aria-label="Appearance type filter"
							>
								{#each filterOptions as option (option.value)}
									<Button
										type="button"
										variant={filter === option.value ? 'default' : 'ghost'}
										class="min-h-10 rounded-none px-4 text-sm"
										onclick={() => (filter = option.value)}
									>
										{option.label}
										<span class="ml-1 text-xs opacity-75">{counts[option.value]}</span>
									</Button>
								{/each}
							</div>
							<div
								class="inline-grid grid-cols-2 overflow-hidden rounded-md border border-border bg-card"
								aria-label="Appearance view mode"
							>
								{#each viewOptions as option (option.value)}
									{@const Icon = option.icon}
									<Button
										type="button"
										variant={viewMode === option.value ? 'default' : 'ghost'}
										class="min-h-10 min-w-10 rounded-none px-3 text-sm"
										aria-label={`${option.label} view`}
										onclick={() => (viewMode = option.value)}
									>
										<Icon aria-hidden="true" />
									</Button>
								{/each}
							</div>
						</div>
					</div>

					{#if filtered.length > 0}
						{#if viewMode === 'timeline'}
							<PersonAppearancesTimeline
								appearances={filtered}
								{mediaItems}
								{addingKey}
								{actionLabel}
								{onAdd}
							/>
						{:else}
							<PersonAppearancesGrid
								appearances={filtered}
								{mediaItems}
								{addingKey}
								{actionLabel}
								{onAdd}
							/>
						{/if}
					{:else}
						<EmptyState
							title="No appearances"
							description="There are no appearances matching this filter."
						/>
					{/if}
				</section>
			</div>
		</section>
	</div>
{/if}
