<script lang="ts">
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import MediaProfileEditorForm from '$lib/components/settings/profiles/MediaProfileEditorForm.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { emptyMediaProfileForm, mediaProfileFormFromProfile } from '$lib/settings/forms';
	import type { MediaProfile, MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';
	import { errorMessageFrom } from '$lib/components/app/shell/controller/helpers';
	import {
		loadMediaProfileQualities,
		returnToMediaProfiles,
		saveMediaProfileForm
	} from './mediaProfileRouteHelpers';

	interface Props {
		profileId: string;
	}

	let { profileId }: Props = $props();
	const app = getAppShellContext();
	let form = $state<MediaProfileForm>(emptyMediaProfileForm());
	let activeProfileId = $state('');
	let qualities = $state<QualitySizeSetting[]>([]);
	let loadingQualities = $state(false);
	let qualityError = $state('');
	let saving = $state(false);
	let saveError = $state('');

	const profile = $derived(app.mediaProfiles.find((item: MediaProfile) => item.id === profileId));
	const notFound = $derived(app.mediaProfiles.length > 0 && !profile);

	onMount(() => {
		void loadQualities();
	});

	$effect(() => {
		if (!profile || activeProfileId === profile.id) return;
		form = mediaProfileFormFromProfile(profile);
		activeProfileId = profile.id;
	});

	async function loadQualities() {
		loadingQualities = true;
		qualityError = '';
		try {
			qualities = await loadMediaProfileQualities();
		} catch (error) {
			qualityError = errorMessageFrom(error, 'Could not load qualities');
		} finally {
			loadingQualities = false;
		}
	}

	async function submitProfile(event: SubmitEvent) {
		event.preventDefault();
		saving = true;
		saveError = '';
		app.clearNotice();
		try {
			await saveMediaProfileForm(form);
			app.message = 'Profile saved';
			await app.loadSettings();
			await returnToMediaProfiles();
		} catch (error) {
			saveError = errorMessageFrom(error, 'Could not save profile');
		} finally {
			saving = false;
		}
	}
</script>

{#if app.isAdmin}
	{#if notFound}
		<PageHeading eyebrow="Settings" title="Edit profile" titleId="settings-profile-title" />
		<Card class="p-5">
			<p class="m-0 text-sm text-muted-foreground">Profile not found.</p>
			<Button class="mt-4" variant="outline" href={resolve('/settings/profiles')}>Back</Button>
		</Card>
	{:else}
		<MediaProfileEditorForm
			title="Edit profile"
			bind:form
			customFormats={app.customFormats}
			languages={app.languages}
			{qualities}
			{loadingQualities}
			{qualityError}
			{saving}
			{saveError}
			submitLabel="Update profile"
			onSubmit={submitProfile}
		/>
	{/if}
{/if}
