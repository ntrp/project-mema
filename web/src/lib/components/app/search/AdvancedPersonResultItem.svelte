<script lang="ts">
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import type { PersonSearchResult } from '$lib/settings/types';
	import { imageUrl, personHref } from './advancedSearchResults';

	interface Props {
		person: PersonSearchResult;
	}

	let { person }: Props = $props();
	const href = $derived(personHref(person));
	const knownFor = $derived((person.knownFor ?? []).slice(0, 4).join(', '));
</script>

<article
	class="grid items-center gap-3.5 rounded-md border border-border bg-muted p-2.5 md:grid-cols-[82px_minmax(0,1fr)_auto]"
>
	<div class="aspect-[2/3] overflow-hidden rounded-md bg-card">
		{#if imageUrl(person.profilePath)}
			<img
				class="block h-full w-full object-cover"
				src={imageUrl(person.profilePath)}
				alt=""
				loading="lazy"
			/>
		{:else}
			<div class="flex h-full items-center justify-center text-muted-foreground">
				<UserIcon class="size-8" aria-hidden="true" />
			</div>
		{/if}
	</div>
	<div class="grid min-w-0 gap-2">
		<div>
			<h3 class="m-0 text-base leading-tight">
				<a class="text-foreground no-underline hover:text-primary-hover" {href}>{person.name}</a>
			</h3>
			<p class="m-0 text-sm text-muted-foreground">
				Person{person.popularity ? ` · popularity ${person.popularity.toFixed(1)}` : ''}
			</p>
		</div>
		{#if knownFor}
			<p class="line-clamp-2 m-0 text-sm text-muted-foreground">Known for {knownFor}</p>
		{/if}
	</div>
	<div class="flex items-center justify-end">
		<Button variant="outline" size="sm" {href}>View person</Button>
	</div>
</article>
