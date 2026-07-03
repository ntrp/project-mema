<script lang="ts">
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import MediaMetadataHero from '$lib/components/app/media/metadata/MediaMetadataHero.svelte';
	import MediaMetadataShell from '$lib/components/app/media/metadata/MediaMetadataShell.svelte';
	import MediaPersonCard from '$lib/components/app/media/people/MediaPersonCard.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import {
		castPeople,
		crewPersonGroups,
		mediaPersonHref,
		type MediaPersonCardData,
		type MediaPersonGroup
	} from '$lib/components/app/media/people/mediaPeople';
	import type { MediaMetadataDetails } from '$lib/settings/types';

	interface Props {
		detail?: MediaMetadataDetails;
		kind?: 'cast' | 'crew';
		loading: boolean;
	}

	let { detail, kind = 'cast', loading }: Props = $props();

	const groups = $derived(detail ? peopleGroups(detail, kind) : []);
	const pageTitle = $derived(kind === 'crew' ? 'Crew' : 'Cast');

	function peopleGroups(
		details: MediaMetadataDetails,
		sectionKind: 'cast' | 'crew'
	): MediaPersonGroup[] {
		const cast = castPeople(details.cast ?? []);
		const crew = crewPersonGroups(details.crew ?? [], details.facts ?? []);
		return (sectionKind === 'cast' ? [{ title: 'Cast', people: cast }] : crew).filter(
			(group) => group.people.length > 0
		);
	}

	function personHref(person: MediaPersonCardData) {
		return mediaPersonHref(person);
	}
</script>

{#if loading}
	<section class="min-h-[220px] rounded-md border border-border bg-card p-5">
		<p class="m-0 text-sm leading-6 text-muted-foreground">Loading {pageTitle.toLowerCase()}</p>
	</section>
{:else if !detail}
	<EmptyState
		title={`${pageTitle} not available`}
		description="Could not load people for this item."
	/>
{:else}
	<MediaMetadataShell labelledby="media-people-title">
		<MediaMetadataHero
			{detail}
			titleId="media-people-title"
			showMonitorBookmark={false}
			showTrailerButton={false}
		/>

		<div class="grid items-start gap-7">
			<main class="grid min-w-0 gap-6 [&>section]:grid [&>section]:min-w-0 [&>section]:gap-2.5">
				{#each groups as group (group.title)}
					<section aria-labelledby={`people-${group.title.toLowerCase().replaceAll(' ', '-')}`}>
						<SectionHeading
							title={group.title}
							titleId={`people-${group.title.toLowerCase().replaceAll(' ', '-')}`}
						>
							{#snippet actions()}
								<span>{group.people.length}</span>
							{/snippet}
						</SectionHeading>
						<div class="grid grid-cols-[repeat(auto-fill,minmax(231px,1fr))] gap-4">
							{#each group.people as person (`${group.title}:${person.name}:${person.role ?? ''}`)}
								<MediaPersonCard
									name={person.name}
									role={person.role}
									image={person.image}
									href={personHref(person)}
								/>
							{/each}
						</div>
					</section>
				{/each}
			</main>
		</div>
	</MediaMetadataShell>
{/if}
