<script lang="ts">
	import MediaMetadataHero from './MediaMetadataHero.svelte';
	import { imageUrl } from './mediaDetail';
	import type { MediaMetadataDetails, MediaMetadataFact } from '$lib/settings/types';

	interface Props {
		detail?: MediaMetadataDetails;
		loading: boolean;
	}

	type PersonCard = {
		name: string;
		role?: string;
		image?: string;
	};

	let { detail, loading }: Props = $props();

	const crewLabels = ['Creator', 'Director', 'Writer', 'Editor', 'Producer'];
	const groups = $derived(detail ? peopleGroups(detail) : []);

	function peopleGroups(details: MediaMetadataDetails): { title: string; people: PersonCard[] }[] {
		const cast: PersonCard[] = (details.cast ?? []).map((person) => ({
			name: person.name,
			role: person.role,
			image: person.profilePath
		}));
		const crew = crewLabels
			.map((label) => ({
				title: label,
				people: peopleFromFact((details.facts ?? []).find((fact) => fact.label === label))
			}))
			.filter((group) => group.people.length > 0);
		return [{ title: 'Cast', people: cast }, ...crew].filter((group) => group.people.length > 0);
	}

	function peopleFromFact(fact: MediaMetadataFact | undefined): PersonCard[] {
		return (fact?.value ?? '')
			.split(',')
			.map((name) => name.trim())
			.filter(Boolean)
			.map((name) => ({ name }));
	}
</script>

{#if loading}
	<section class="metadata-detail-loading panel">
		<p class="muted">Loading cast</p>
	</section>
{:else if !detail}
	<section class="empty-state">
		<h2>Cast not available</h2>
		<p>Could not load people for this item.</p>
	</section>
{:else}
	<section
		class="metadata-detail"
		aria-labelledby="media-people-title"
		style:--backdrop-url={imageUrl(detail.backdropPath, 'original')
			? `url("${imageUrl(detail.backdropPath, 'original')}")`
			: undefined}
	>
		<MediaMetadataHero {detail} titleId="media-people-title" />

		<div class="metadata-body">
			<main class="metadata-main">
				{#each groups as group (group.title)}
					<section aria-labelledby={`people-${group.title.toLowerCase().replaceAll(' ', '-')}`}>
						<div class="section-heading">
							<h2 id={`people-${group.title.toLowerCase().replaceAll(' ', '-')}`}>
								{group.title}
							</h2>
							<span>{group.people.length}</span>
						</div>
						<div class="metadata-people-grid">
							{#each group.people as person (`${group.title}:${person.name}:${person.role ?? ''}`)}
								<article class="metadata-cast-card">
									<div>
										{#if person.image && imageUrl(person.image, 'w185')}
											<img src={imageUrl(person.image, 'w185')} alt="" loading="lazy" />
										{:else}
											<span>{person.name.slice(0, 1)}</span>
										{/if}
									</div>
									<strong>{person.name}</strong>
									{#if person.role}
										<p>{person.role}</p>
									{/if}
								</article>
							{/each}
						</div>
					</section>
				{/each}
			</main>
		</div>
	</section>
{/if}
