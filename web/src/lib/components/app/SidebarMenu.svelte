<script lang="ts">
	interface MenuItem<TValue extends string> {
		value: TValue;
		label: string;
		meta?: string;
	}

	interface Props<TValue extends string> {
		title: string;
		items: MenuItem<TValue>[];
		active: TValue;
		onSelect: (_value: TValue) => void;
	}

	let { title, items, active, onSelect }: Props<string> = $props();
</script>

<aside class="side-menu" aria-label={title}>
	<h2>{title}</h2>
	<nav>
		{#each items as item (item.value)}
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
		{/each}
	</nav>
</aside>
