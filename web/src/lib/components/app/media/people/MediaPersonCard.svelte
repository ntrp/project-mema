<script lang="ts">
	import { imageUrl } from '$lib/components/app/media/detail/mediaDetail';

	interface Props {
		name: string;
		role?: string;
		image?: string;
		href?: string;
	}

	let { name, role, image, href }: Props = $props();

	const profileUrl = $derived(imageUrl(image, 'w185'));
</script>

<article class="group/person relative min-w-0 snap-start">
	{#if href}
		<a class="absolute inset-0 z-[1] rounded-md" {href} aria-label={`Open ${name} details`}></a>
	{/if}
	<div
		class="grid aspect-[2/3] min-w-0 content-center justify-items-center gap-2 rounded-md border border-border bg-card px-3 py-4 text-center transition-[border-color,transform] group-hover/person:-translate-y-1 group-hover/person:border-primary/50 group-focus-within/person:-translate-y-1 group-focus-within/person:border-primary/50"
	>
		<div
			class="grid aspect-square w-[min(74%,126px)] place-items-center overflow-hidden rounded-full bg-muted font-black text-muted-foreground"
		>
			{#if profileUrl}
				<img class="block size-full object-cover" src={profileUrl} alt="" loading="lazy" />
			{:else}
				<span>{name.slice(0, 1)}</span>
			{/if}
		</div>
		<strong class="w-full overflow-hidden text-ellipsis whitespace-nowrap text-sm">{name}</strong>
		{#if role}
			<p class="m-0 line-clamp-2 w-full text-xs whitespace-normal text-muted-foreground">{role}</p>
		{/if}
	</div>
</article>
