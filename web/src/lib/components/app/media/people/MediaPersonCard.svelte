<script lang="ts">
	import { imageUrl } from '$lib/components/app/media/detail/mediaDetail';

	interface Props {
		name: string;
		role?: string;
		image?: string;
	}

	let { name, role, image }: Props = $props();

	const profileUrl = $derived(imageUrl(image, 'w185'));
</script>

<article
	class="grid aspect-[2/3] min-w-0 snap-start content-center justify-items-center gap-2 rounded-md border border-border bg-card px-3 py-4 text-center"
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
</article>
