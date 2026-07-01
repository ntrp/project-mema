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
	import { setAppShellContext } from '$lib/features/app/appShellContext';
	import '$lib/settings/styles.css';

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

{#if app.loading}
	<main class="shell">
		<section class="panel">
			<p class="muted">Loading app</p>
		</section>
	</main>
{:else if !app.authenticated}
	<main class="shell login-shell">
		<div class="login-brand">
			<span class="brand-mark large" aria-hidden="true">M</span>
			<h1>mema</h1>
		</div>
		<AuthPanel bind:username={app.username} bind:password={app.password} onLogin={app.login} />
	</main>
{:else}
	<div class="app-frame">
		<SidebarMenu
			title="mema"
			items={app.primaryItems}
			active={app.activePrimarySection}
			activeSubmenu={app.activeSubmenuSection}
			onSelect={app.selectPrimarySection}
			onSubmenuSelect={app.selectSubmenuSection}
		/>
		<div class="app-main">
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
			<main class="app-content">
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
{/if}

<NoticeStack message={app.message} errorMessage={app.errorMessage} onDismiss={app.clearNotice} />
