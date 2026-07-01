<script lang="ts">
	import { cn, type WithElementRef } from "$lib/utils.js";
	import type { HTMLLiAttributes } from "svelte/elements";

	let {
		ref = $bindable(null),
		children,
		class: className,
		...restProps
	}: WithElementRef<HTMLLiAttributes> = $props();
</script>

<li
	bind:this={ref}
	data-slot="sidebar-menu-sub-item"
	data-sidebar="menu-sub-item"
	class={cn("group/menu-sub-item relative", className)}
	{...restProps}
>
	{@render children?.()}
</li>
