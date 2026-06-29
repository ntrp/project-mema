<script lang="ts">
	import { resolve } from '$app/paths';

	type MenuIcon = 'discover' | 'movies' | 'series' | 'activity' | 'settings';
	type MenuHref =
		| '/discover'
		| '/requests'
		| '/movies'
		| '/series'
		| '/activity'
		| '/settings/library'
		| '/settings/download-clients'
		| '/settings/indexers'
		| '/settings/quality'
		| '/settings/file-naming'
		| '/settings/profiles'
		| '/settings/metadata'
		| '/settings/tags'
		| '/settings/users'
		| '/settings/system/logs';

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
	let menuOpen = $state(false);

	function selectItem(value: string) {
		onSelect(value);
		menuOpen = false;
	}

	function selectSubmenuItem(value: string) {
		onSubmenuSelect?.(value);
		menuOpen = false;
	}
</script>

<aside class="side-menu" class:menu-open={menuOpen} aria-label={title}>
	<div class="side-menu-header">
		<a class="brand-button" href={resolve('/discover')} onclick={() => selectItem('discover')}>
			<span class="brand-mark" aria-hidden="true">M</span>
			<span class="brand-name">{title}</span>
		</a>
		<button
			type="button"
			class="menu-toggle"
			aria-expanded={menuOpen}
			aria-controls="primary-menu"
			onclick={() => (menuOpen = !menuOpen)}
		>
			{menuOpen ? 'Close' : 'Menu'}
		</button>
	</div>
	<nav id="primary-menu">
		{#each items as item (item.value)}
			<a
				href={resolve(item.href)}
				class:active-menu={active === item.value && !item.children?.length}
				aria-current={active === item.value && !item.children?.length ? 'page' : undefined}
				onclick={() => selectItem(item.value)}
			>
				<span class="menu-icon" aria-hidden="true">
					{#if item.icon === 'discover'}
						<svg viewBox="0 0 24 24">
							<path d="M12 3l1.8 5 5.2 1.9-5.2 1.9L12 17l-1.8-5.2L5 9.9 10.2 8z" />
							<path d="M19 14l.8 2.2L22 17l-2.2.8L19 20l-.8-2.2L16 17l2.2-.8z" />
						</svg>
					{:else if item.icon === 'movies'}
						<svg viewBox="0 0 24 24">
							<rect x="4" y="5" width="16" height="14" rx="2" />
							<path d="M8 5v14M16 5v14M4 9h16M4 15h16" />
						</svg>
					{:else if item.icon === 'series'}
						<svg viewBox="0 0 24 24">
							<rect x="4" y="5" width="16" height="12" rx="2" />
							<path d="M9 21h6M12 17v4" />
						</svg>
					{:else if item.icon === 'activity'}
						<svg viewBox="0 0 24 24">
							<circle cx="12" cy="12" r="8" />
							<path d="M12 8v5l3 2" />
						</svg>
					{:else}
						<svg viewBox="0 0 24 24">
							<circle cx="12" cy="12" r="3" />
							<path
								d="M19 12a7.6 7.6 0 0 0-.1-1l2-1.5-2-3.4-2.4 1a8 8 0 0 0-1.7-1L14.5 3h-5l-.4 3.1a8 8 0 0 0-1.7 1l-2.4-1-2 3.4 2 1.5a7.6 7.6 0 0 0 0 2l-2 1.5 2 3.4 2.4-1a8 8 0 0 0 1.7 1l.4 3.1h5l.4-3.1a8 8 0 0 0 1.7-1l2.4 1 2-3.4-2-1.5c.1-.3.1-.7.1-1z"
							/>
						</svg>
					{/if}
				</span>
				<span>{item.label}</span>
			</a>
			{#if active === item.value && item.children?.length}
				<div class="submenu" aria-label={`${item.label} sections`}>
					{#each item.children as child (child.value)}
						<a
							href={resolve(child.href)}
							class:active-submenu={activeSubmenu === child.value}
							aria-current={activeSubmenu === child.value ? 'page' : undefined}
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
