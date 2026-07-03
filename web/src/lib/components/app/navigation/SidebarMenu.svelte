<script lang="ts">
	import ActivityIcon from '@lucide/svelte/icons/clock-3';
	import CompassIcon from '@lucide/svelte/icons/compass';
	import ComputerIcon from '@lucide/svelte/icons/monitor';
	import EyeOffIcon from '@lucide/svelte/icons/eye-off';
	import MoviesIcon from '@lucide/svelte/icons/clapperboard';
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import SeriesIcon from '@lucide/svelte/icons/tv';
	import { resolve } from '$app/paths';
	import * as Sidebar from '$lib/components/ui/sidebar';

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
		| '/activity/history'
		| '/activity/blocklist'
		| '/settings/general'
		| '/settings/library'
		| '/settings/download-clients'
		| '/settings/indexers'
		| '/settings/quality'
		| '/settings/profiles'
		| '/settings/custom-formats'
		| '/settings/metadata'
		| '/settings/languages'
		| '/settings/tags'
		| '/settings/users'
		| '/system/status'
		| '/system/indexing'
		| '/system/metadata'
		| '/system/jobs'
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

	function resolveHref(href: MenuHref) {
		return resolve(href as '/discover');
	}

	const navButtonClass =
		'gap-4 text-[20px] font-semibold data-[active=true]:bg-primary data-[active=true]:text-primary-foreground [&_svg]:size-6';
	const navSubButtonClass = 'text-[15px] font-semibold';
</script>

<Sidebar.Root collapsible="icon" aria-label={title}>
	<Sidebar.Header>
		<a
			class="flex min-w-0 items-center gap-3 rounded-md px-2 py-2 text-sidebar-foreground no-underline transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
			href={resolve('/discover')}
			onclick={() => selectItem('discover')}
		>
			<span
				class="grid size-8 shrink-0 place-items-center rounded-md bg-primary text-base font-black text-primary-foreground min-[641px]:size-9"
				aria-hidden="true">M</span
			>
			<span class="truncate text-lg font-black group-data-[collapsible=icon]:hidden">{title}</span>
		</a>
	</Sidebar.Header>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupContent>
				<Sidebar.Menu id="primary-menu">
					{#each items as item (item.value)}
						{@const itemActive = active === item.value}
						{@const pageActive = itemActive && !item.children?.length}
						{@const Icon = iconComponent(item.icon)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton
								size="lg"
								isActive={pageActive}
								tooltipContent={item.label}
								class={navButtonClass}
							>
								{#snippet child({ props })}
									<a
										{...props}
										href={resolveHref(item.href)}
										aria-current={pageActive ? 'page' : undefined}
										onclick={() => selectItem(item.value)}
									>
										<Icon aria-hidden="true" />
										<span>{item.label}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
						{#if itemActive && item.children?.length}
							<Sidebar.MenuSub aria-label={`${item.label} sections`}>
								{#each item.children as child (child.value)}
									{@const childActive = activeSubmenu === child.value}
									<Sidebar.MenuSubItem>
										<Sidebar.MenuSubButton
											href={resolveHref(child.href)}
											isActive={childActive}
											aria-current={childActive ? 'page' : undefined}
											class={navSubButtonClass}
											onclick={() => selectSubmenuItem(child.value)}
										>
											<span>{child.label}</span>
										</Sidebar.MenuSubButton>
									</Sidebar.MenuSubItem>
								{/each}
							</Sidebar.MenuSub>
						{/if}
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Rail />
</Sidebar.Root>
