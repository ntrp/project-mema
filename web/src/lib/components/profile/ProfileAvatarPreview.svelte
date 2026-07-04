<script lang="ts">
	interface Props {
		name: string;
		pictureUrl: string;
		username: string;
	}

	let { name, pictureUrl, username }: Props = $props();
	const fallback = $derived((name || username || '?').slice(0, 1).toUpperCase());
</script>

<div
	class="grid size-24 shrink-0 place-items-center overflow-hidden rounded-md border border-border bg-muted text-3xl font-black text-muted-foreground"
	aria-label="Profile picture preview"
>
	{#if pictureUrl}
		<img class="h-full w-full object-cover" src={pictureUrl} alt="" />
	{:else}
		<span aria-hidden="true">{fallback}</span>
	{/if}
</div>
