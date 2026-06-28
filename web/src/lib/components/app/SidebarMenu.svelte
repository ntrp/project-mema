<script lang="ts">
	import { resolve } from '$app/paths';

	type MenuHref =
		| '/'
		| '/explore'
		| '/movies'
		| '/series'
		| '/activity'
		| '/settings'
		| '/settings/download-clients'
		| '/settings/indexers';

	interface MenuItem<TValue extends string> {
		value: TValue;
		label: string;
		meta?: string;
		href?: MenuHref;
	}

	interface Props<TValue extends string> {
		title: string;
		items: readonly MenuItem<TValue>[];
		active: TValue;
		onSelect: (_value: TValue) => void;
	}

	let { title, items, active, onSelect }: Props<string> = $props();
</script>

<aside class="side-menu" aria-label={title}>
	<h2>{title}</h2>
	<nav>
		{#each items as item (item.value)}
			{#if item.href}
				<a
					href={resolve(item.href)}
					class:active-menu={active === item.value}
					aria-current={active === item.value ? 'page' : undefined}
				>
					<span>{item.label}</span>
					{#if item.meta}
						<small>{item.meta}</small>
					{/if}
				</a>
			{:else}
				<button
					type="button"
					class:active-menu={active === item.value}
					aria-current={active === item.value ? 'page' : undefined}
					onclick={() => onSelect(item.value)}
				>
					<span>{item.label}</span>
					{#if item.meta}
						<small>{item.meta}</small>
					{/if}
				</button>
			{/if}
		{/each}
	</nav>
</aside>
