<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import AppDocumentHead from './AppDocumentHead.svelte';
	import AppNav from '../navigation/AppNav.svelte';
	import { routeStateFromPath } from '$lib/components/app/shell/controller/routeState';
	import { createAppShellController } from '$lib/components/app/shell/controller/index.svelte';
	import { imageUrl } from '$lib/components/app/media/detail/mediaDetail';
	import MediaDeleteModal from '$lib/components/app/media/actions/MediaDeleteModal.svelte';
	import MediaActionModal from '$lib/components/app/media/actions/MediaActionModal.svelte';
	import SidebarMenu from '../navigation/SidebarMenu.svelte';
	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import NoticeStack from '$lib/components/settings/shared/NoticeStack.svelte';
	import ScrollTopButton from '$lib/components/shared/ScrollTopButton.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { setAppShellContext } from '$lib/features/app/appShellContext';
	import type { MediaItem, PersonAppearance } from '$lib/settings/types';

	let { children } = $props();
	const route = $derived(routeStateFromPath(page.url.pathname, page.params, page.url.searchParams));
	// svelte-ignore state_referenced_locally
	let app = $state(createAppShellController(route));
	setAppShellContext(app);
	let personBackdropIndex = $state(0);
	const personBackdropPaths = $derived(personAppearanceBackdrops());
	const mainBackdropUrl = $derived(imageUrl(activeMainBackdropPath(), 'original'));

	onMount(() => {
		void app.initialise();
		return app.disconnectEvents;
	});

	$effect(() => {
		void app.applyRoute(route);
	});

	$effect(() => {
		const paths = personBackdropPaths;
		personBackdropIndex = 0;
		if (app.activeView !== 'person-detail' || paths.length < 2) return;
		const interval = window.setInterval(() => {
			personBackdropIndex = (personBackdropIndex + 1) % paths.length;
		}, 6000);
		return () => window.clearInterval(interval);
	});

	function personAppearanceBackdrops(): string[] {
		const appearances: PersonAppearance[] = app.personDetail?.appearances ?? [];
		const paths = appearances
			.map((appearance: PersonAppearance) => appearance.backdropPath)
			.filter((path: string | undefined): path is string => Boolean(path));
		return paths.filter((path: string, index: number) => paths.indexOf(path) === index);
	}

	function activeMainBackdropPath() {
		if (app.activeView === 'person-detail') {
			return personBackdropPaths[personBackdropIndex % personBackdropPaths.length];
		}
		if (app.activeView === 'media-people') {
			return app.mediaPeopleMetadataDetail?.backdropPath;
		}
		if (app.activeView === 'metadata-detail' || app.activeView === 'related-section') {
			return app.metadataDetail?.backdropPath;
		}
		if (!app.selectedMediaItemId || !['movies', 'series'].includes(app.activeHomeSection)) {
			return undefined;
		}
		return app.mediaItems.find(
			(item: MediaItem) =>
				item.id === app.selectedMediaItemId &&
				item.type === (app.activeHomeSection === 'movies' ? 'movie' : 'series')
		)?.backdropPath;
	}
</script>

<AppDocumentHead />

<Tooltip.Provider>
	{#if app.loading}
		<main class="min-h-screen bg-background p-4 min-[641px]:p-8">
			<section class="rounded-md border border-border bg-card p-5 text-card-foreground shadow-xs">
				<p class="text-sm text-muted-foreground">Loading app</p>
			</section>
		</main>
	{:else if !app.authenticated}
		<main
			class="grid min-h-screen content-center justify-items-center gap-5 bg-background p-4 min-[641px]:p-8"
		>
			<div class="flex w-[min(100%,420px)] items-center justify-center gap-3">
				<span
					class="grid size-12 place-items-center rounded-md bg-primary text-2xl font-black text-primary-foreground"
					aria-hidden="true">M</span
				>
				<h1 class="m-0 text-[40px] leading-tight font-semibold text-foreground">mema</h1>
			</div>
			<AuthPanel bind:username={app.username} bind:password={app.password} onLogin={app.login} />
		</main>
	{:else}
		<Sidebar.Provider>
			<SidebarMenu
				title="mema"
				items={app.primaryItems}
				active={app.activePrimarySection}
				activeSubmenu={app.activeSubmenuSection}
				onSelect={app.selectPrimarySection}
				onSubmenuSelect={app.selectSubmenuSection}
			/>
			<Sidebar.Inset>
				<AppNav
					bind:searchQuery={app.searchQuery}
					currentUser={app.currentUser}
					groups={app.autocompleteGroups}
					loading={app.loadingAutocomplete}
					onSearch={app.autocompleteMedia}
					onSelect={app.selectAutocompleteResult}
					onAdvancedSearch={app.openAdvancedSearch}
					onProfile={app.showProfile}
					onLogout={app.logout}
					showNotifications={app.isAdmin}
				/>
				<main
					class="relative isolate min-h-[calc(100vh-76px)] overflow-hidden bg-background px-3.5 py-3.5 min-[641px]:px-4.5 min-[641px]:py-6 min-[641px]:pb-10"
				>
					{#if mainBackdropUrl}
						{#key mainBackdropUrl}
							<div
								class="pointer-events-none absolute inset-x-0 top-0 z-0 h-[min(620px,64vh)] mask-[linear-gradient(to_bottom,black_0%,black_70%,transparent_100%)]"
								transition:fade={{ duration: 650 }}
							>
								<img
									class="absolute inset-x-0 top-0 h-full w-full object-cover opacity-45"
									src={mainBackdropUrl}
									alt=""
								/>
								<div
									class="absolute inset-0 bg-linear-to-b from-background/10 via-background/45 to-background"
								></div>
							</div>
						{/key}
					{/if}
					<div class="relative z-1">
						{@render children?.()}
					</div>
				</main>
			</Sidebar.Inset>
		</Sidebar.Provider>
		{#if app.activeMediaCandidate}
			<MediaActionModal
				candidate={app.activeMediaCandidate}
				isAdmin={app.isAdmin}
				libraryFolders={app.libraryFolders}
				qualityProfiles={app.mediaProfiles}
				tags={app.tags}
				saving={app.savingMediaAction}
				onClose={app.closeMediaAction}
				onConfirm={app.confirmMediaAction}
			/>
		{/if}
		{#if app.mediaDeleteCandidate}
			<MediaDeleteModal
				item={app.mediaDeleteCandidate}
				deleting={app.deletingMediaItemId === app.mediaDeleteCandidate.id}
				onClose={app.closeMediaDelete}
				onDelete={app.confirmMediaDelete}
			/>
		{/if}
		<ScrollTopButton />
	{/if}

	<NoticeStack message={app.message} errorMessage={app.errorMessage} onDismiss={app.clearNotice} />
</Tooltip.Provider>
