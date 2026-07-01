<script lang="ts">
	import ActivityIcon from '@lucide/svelte/icons/clock-3';
	import CompassIcon from '@lucide/svelte/icons/compass';
	import ComputerIcon from '@lucide/svelte/icons/monitor';
	import EyeOffIcon from '@lucide/svelte/icons/eye-off';
	import MoviesIcon from '@lucide/svelte/icons/clapperboard';
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import SeriesIcon from '@lucide/svelte/icons/tv';
	import { resolve } from '$app/paths';
	import { cn } from '$lib/utils';

	type MenuIcon =
		| 'discover'
		| 'movies'
		| 'series'
		| 'activity'
		| 'settings'
		| 'computer'
		| 'visibility_off';
	type MenuHref =
		| '/discover'
		| `/discover/${string}`
		| '/blacklist'
		| '/requests'
		| '/movies'
		| '/series'
		| '/wanted'
		| '/activity'
		| '/settings/general'
		| '/settings/library'
		| '/settings/download-clients'
		| '/settings/indexers'
		| '/settings/quality'
		| '/settings/file-naming'
		| '/settings/profiles'
		| '/settings/custom-formats'
		| '/settings/metadata'
		| '/settings/tags'
		| '/settings/users'
		| '/system/status'
		| '/system/logs'
		| '/system/events';

	interface SubmenuItem<TValue extends string> {
		value: TValue;
		label: string;
		href: MenuHref;
	}

	interface MenuItem<TValue extends string> {
		value: TValue;
		label: string;
		icon: MenuIcon;
		href: MenuHref;
		children?: readonly SubmenuItem<string>[];
	}

	interface Props<TValue extends string> {
		title: string;
		items: readonly MenuItem<TValue>[];
		active: TValue;
		activeSubmenu?: string;
		onSelect: (_value: TValue) => void;
		onSubmenuSelect?: (_value: string) => void;
	}

	let { title, items, active, activeSubmenu, onSelect, onSubmenuSelect }: Props<string> = $props();

	function selectItem(value: string) {
		onSelect(value);
	}

	function selectSubmenuItem(value: string) {
		onSubmenuSelect?.(value);
	}

	type MenuIconComponent = typeof CompassIcon;

	const menuIcons: Record<MenuIcon, MenuIconComponent> = {
		discover: CompassIcon,
		movies: MoviesIcon,
		series: SeriesIcon,
		activity: ActivityIcon,
		settings: SettingsIcon,
		computer: ComputerIcon,
		visibility_off: EyeOffIcon
	};

	function iconComponent(icon: MenuIcon) {
		return menuIcons[icon];
	}

	function navLinkClass(isActive: boolean) {
		return cn(
			'relative flex min-h-10 w-full items-center gap-3 rounded-md border border-transparent px-2.5 py-2 text-left font-bold text-muted-foreground no-underline transition-colors hover:bg-muted hover:text-foreground',
			'max-[640px]:h-[50px] max-[640px]:min-h-[50px] max-[640px]:w-16 max-[640px]:min-w-16 max-[640px]:shrink-0 max-[640px]:flex-col max-[640px]:justify-center max-[640px]:gap-1 max-[640px]:px-1 max-[640px]:py-1 max-[640px]:text-center max-[640px]:text-[10px] max-[640px]:leading-none',
			isActive &&
				'border-primary/60 bg-primary text-primary-foreground shadow-md hover:bg-primary hover:text-primary-foreground before:absolute before:left-0 before:top-1/2 before:h-5 before:w-0.5 before:-translate-y-1/2 before:rounded-md before:bg-primary-foreground max-[640px]:before:inset-x-2 max-[640px]:before:bottom-0.5 max-[640px]:before:top-auto max-[640px]:before:h-0.5 max-[640px]:before:w-auto max-[640px]:before:translate-y-0'
		);
	}

	function submenuLinkClass(isActive: boolean) {
		return cn(
			'relative flex min-h-8 items-center rounded-md border border-transparent px-2.5 py-1.5 text-sm font-bold text-muted-foreground no-underline transition-colors hover:bg-primary/10 hover:text-foreground max-[640px]:min-h-[34px] max-[640px]:shrink-0 max-[640px]:whitespace-nowrap',
			isActive &&
				'border-primary/60 bg-primary/15 text-foreground ring-1 ring-primary/20 before:absolute before:left-0 before:top-1/2 before:h-5 before:w-0.5 before:-translate-y-1/2 before:rounded-md before:bg-foreground max-[640px]:before:inset-x-2 max-[640px]:before:bottom-0.5 max-[640px]:before:top-auto max-[640px]:before:h-0.5 max-[640px]:before:w-auto max-[640px]:before:translate-y-0'
		);
	}
</script>

<aside
	class="sticky top-0 z-20 grid content-start gap-2.5 border-b border-border bg-card px-3 py-2 shadow-lg min-[641px]:gap-4 min-[641px]:px-2.5 min-[641px]:py-6 min-[981px]:min-h-screen min-[981px]:gap-8 min-[981px]:border-r min-[981px]:border-b-0 min-[981px]:shadow-none"
	aria-label={title}
>
	<div class="flex min-w-0 items-center justify-center gap-2.5 min-[981px]:justify-between">
		<a
			class="flex min-w-0 items-center gap-3 text-foreground no-underline"
			href={resolve('/discover')}
			onclick={() => selectItem('discover')}
		>
			<span
				class="grid size-8 shrink-0 place-items-center rounded-md bg-primary text-base font-black text-primary-foreground min-[641px]:size-9"
				aria-hidden="true">M</span
			>
			<span class="truncate text-lg font-black min-[981px]:text-2xl">{title}</span>
		</a>
	</div>
	<nav
		id="primary-menu"
		class="grid gap-1.5 max-[640px]:flex max-[640px]:min-w-0 max-[640px]:overflow-x-auto max-[640px]:overflow-y-hidden max-[640px]:pb-0.5 max-[640px]:[scrollbar-width:none] min-[641px]:max-[980px]:grid-cols-[repeat(auto-fit,minmax(110px,1fr))]"
	>
		{#each items as item (item.value)}
			{@const itemActive = active === item.value && !item.children?.length}
			{@const Icon = iconComponent(item.icon)}
			<a
				href={resolve(item.href)}
				class={navLinkClass(itemActive)}
				aria-current={itemActive ? 'page' : undefined}
				onclick={() => selectItem(item.value)}
			>
				<Icon aria-hidden="true" class="size-5 max-[640px]:size-[18px]" />
				<span class="max-[640px]:block max-[640px]:w-full max-[640px]:truncate">{item.label}</span>
			</a>
			{#if active === item.value && item.children?.length}
				<div
					class="grid gap-1 border-l border-border pl-2.5 min-[641px]:max-[980px]:col-span-full min-[641px]:max-[980px]:grid-cols-[repeat(auto-fit,minmax(130px,1fr))] min-[641px]:max-[980px]:border-t min-[641px]:max-[980px]:border-l-0 min-[641px]:max-[980px]:pt-2 min-[641px]:max-[980px]:pl-0 max-[640px]:flex max-[640px]:shrink-0 max-[640px]:border-0 max-[640px]:pl-1"
					aria-label={`${item.label} sections`}
				>
					{#each item.children as child (child.value)}
						{@const childActive = activeSubmenu === child.value}
						<a
							href={resolve(child.href)}
							class={submenuLinkClass(childActive)}
							aria-current={childActive ? 'page' : undefined}
							onclick={() => selectSubmenuItem(child.value)}
						>
							{child.label}
						</a>
					{/each}
				</div>
			{/if}
		{/each}
	</nav>
</aside>
