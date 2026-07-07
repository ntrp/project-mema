<script lang="ts">
	import ArrowRightIcon from '@lucide/svelte/icons/arrow-right';
	import type { MediaPersonGroup } from '$lib/components/app/media/people/mediaPeople';
	import { mediaPersonHref } from '$lib/components/app/media/people/mediaPeople';

	interface Props {
		groups: MediaPersonGroup[];
		href?: string;
	}

	let { groups, href }: Props = $props();
</script>

{#if groups.length > 0}
	<h3 class="mt-1 mb-0 text-xl text-foreground">
		{#if href}
			<a
				class="inline-flex items-center gap-2 text-inherit no-underline hover:text-primary-hover focus-visible:text-primary-hover focus-visible:outline-none"
				{href}
			>
				<span>Crew</span>
				<ArrowRightIcon aria-hidden="true" />
			</a>
		{:else}
			Crew
		{/if}
	</h3>
	<div class="grid items-start gap-x-7 gap-y-4.5 md:grid-cols-3" aria-label="Crew">
		{#each groups as group (group.title)}
			<div class="grid min-w-0 content-start gap-1">
				<strong class="wrap-anywhere text-foreground">{group.title}</strong>
				<span class="wrap-anywhere text-muted-foreground">
					{#each group.people.slice(0, 3) as person, index (`${group.title}:${person.externalProvider ?? ''}:${person.externalId ?? person.name}:${index}`)}
						{@const personUrl = mediaPersonHref(person)}
						{#if index > 0},&nbsp;
						{/if}{#if personUrl}<a
								class="text-inherit no-underline hover:underline hover:text-primary-hover focus-visible:text-primary-hover focus-visible:outline-none"
								href={personUrl}
							>
								{person.name}</a
							>{:else}{person.name}{/if}
					{/each}
				</span>
			</div>
		{/each}
	</div>
{/if}
