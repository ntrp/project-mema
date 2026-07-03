<script lang="ts">
	import { imageUrl } from '$lib/components/app/media/detail/mediaDetail';
	import type { PersonDetails } from '$lib/settings/types';

	interface Props {
		person: PersonDetails;
	}

	let { person }: Props = $props();

	const profileUrl = $derived(imageUrl(person.profilePath, 'w342'));
	const aliases = $derived(person.alsoKnownAs ?? []);
	const lifespan = $derived(
		[formatDate(person.birthday), formatDate(person.deathday)].filter(Boolean).join(' - ')
	);
	const bornLine = $derived(
		[lifespan ? `Born ${lifespan}` : '', person.placeOfBirth ? `in ${person.placeOfBirth}` : '']
			.filter(Boolean)
			.join(' ')
	);

	function formatDate(value?: string) {
		if (!value) return '';
		const date = new Date(`${value}T00:00:00`);
		if (Number.isNaN(date.getTime())) return value;
		return new Intl.DateTimeFormat(undefined, {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		}).format(date);
	}
</script>

<aside
	class="grid min-w-0 items-start gap-6 min-[781px]:grid-cols-[minmax(220px,290px)_minmax(0,1fr)]"
>
	<div
		class="aspect-square overflow-hidden rounded-full border border-border bg-card shadow-xl max-[780px]:max-w-55 mx-5"
	>
		{#if profileUrl}
			<img class="block size-full object-cover" src={profileUrl} alt="" />
		{:else}
			<div
				class="grid size-full place-items-center bg-muted text-5xl font-black text-muted-foreground"
			>
				{person.name.slice(0, 1)}
			</div>
		{/if}
	</div>

	<div class="grid min-w-0 gap-4 text-sm leading-6 text-muted-foreground">
		<div class="grid gap-2">
			<h1 class="m-0 text-4xl leading-tight font-semibold text-foreground">{person.name}</h1>
			{#if bornLine}
				<p class="m-0 font-medium text-foreground">{bornLine}</p>
			{/if}
			{#if aliases.length > 0}
				<p class="m-0">
					<span class="font-medium text-foreground">Also Known As:</span>
					{aliases.join(', ')}
				</p>
			{/if}
		</div>

		{#if person.biography}
			<p
				class="m-0 max-h-48 max-w-[90ch] overflow-y-auto pr-2 text-base leading-7 text-muted-foreground"
			>
				{person.biography}
			</p>
		{/if}
	</div>
</aside>
