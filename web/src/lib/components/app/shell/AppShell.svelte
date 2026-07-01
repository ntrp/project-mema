<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import AppDocumentHead from './AppDocumentHead.svelte';
	import AppNav from '../navigation/AppNav.svelte';
	import { routeStateFromPath } from '$lib/components/app/shell/controller/routeState';
	import { createAppShellController } from '$lib/components/app/shell/controller/index.svelte';
	import MediaDeleteModal from '../media/MediaDeleteModal.svelte';
	import MediaActionModal from '../media/MediaActionModal.svelte';
	import SidebarMenu from '../navigation/SidebarMenu.svelte';
	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import NoticeStack from '$lib/components/settings/shared/NoticeStack.svelte';
	import ScrollTopButton from '$lib/components/shared/ScrollTopButton.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { setAppShellContext } from '$lib/features/app/appShellContext';

	let { children } = $props();
	const route = $derived(routeStateFromPath(page.url.pathname, page.params, page.url.searchParams));
	// svelte-ignore state_referenced_locally
	let app = $state(createAppShellController(route));
	setAppShellContext(app);

	onMount(() => {
		void app.initialise();
		return app.disconnectEvents;
	});

	$effect(() => {
		void app.applyRoute(route);
	});
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
		<div
			class="grid min-h-screen w-full max-w-[100vw] grid-cols-1 overflow-x-clip bg-background min-[981px]:grid-cols-[220px_minmax(0,1fr)]"
		>
			<SidebarMenu
				title="mema"
				items={app.primaryItems}
				active={app.activePrimarySection}
				activeSubmenu={app.activeSubmenuSection}
				onSelect={app.selectPrimarySection}
				onSubmenuSelect={app.selectSubmenuSection}
			/>
			<div class="min-w-0 bg-background">
				<AppNav
					bind:searchQuery={app.searchQuery}
					groups={app.autocompleteGroups}
					loading={app.loadingAutocomplete}
					onSearch={app.autocompleteMedia}
					onSelect={app.selectAutocompleteResult}
					onAdvancedSearch={app.openAdvancedSearch}
					onProfile={app.showProfile}
					onLogout={app.logout}
					showNotifications={app.isAdmin}
				/>
				<main class="px-3.5 py-3.5 min-[641px]:px-[18px] min-[641px]:py-6 min-[641px]:pb-10">
					{@render children?.()}
				</main>
			</div>
		</div>
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
