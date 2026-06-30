<script lang="ts">
	import { onMount } from 'svelte';
	import AppDocumentHead from '$lib/components/app/AppDocumentHead.svelte';
	import AppMainContent from '$lib/components/app/AppMainContent.svelte';
	import AppNav from '$lib/components/app/AppNav.svelte';
	import {
		createAppShellController,
		type AppShellOptions
	} from '$lib/components/app/appShellController.svelte';
	import MediaDeleteModal from '$lib/components/app/MediaDeleteModal.svelte';
	import MediaActionModal from '$lib/components/app/MediaActionModal.svelte';
	import SidebarMenu from '$lib/components/app/SidebarMenu.svelte';
	import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
	import NoticeStack from '$lib/components/settings/NoticeStack.svelte';
	import '$lib/settings/styles.css';

	let props: AppShellOptions = $props();
	// svelte-ignore state_referenced_locally
	let app = $state(createAppShellController(props));

	onMount(() => {
		void app.initialise();
		return app.disconnectEvents;
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
				<AppMainContent bind:app options={props} />
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
