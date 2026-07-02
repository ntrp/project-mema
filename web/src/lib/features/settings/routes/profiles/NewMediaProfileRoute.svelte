<script lang="ts">
	import { onMount } from 'svelte';
	import MediaProfileEditorForm from '$lib/components/settings/profiles/MediaProfileEditorForm.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { emptyMediaProfileForm } from '$lib/settings/forms';
	import type { MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';
	import { errorMessageFrom } from '$lib/components/app/shell/controller/helpers';
	import {
		loadMediaProfileQualities,
		returnToMediaProfiles,
		saveMediaProfileForm
	} from './mediaProfileRouteHelpers';

	const app = getAppShellContext();
	let form = $state<MediaProfileForm>(emptyMediaProfileForm());
	let qualities = $state<QualitySizeSetting[]>([]);
	let loadingQualities = $state(false);
	let qualityError = $state('');
	let saving = $state(false);
	let saveError = $state('');

	onMount(() => {
		void loadQualities();
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
	<MediaProfileEditorForm
		title="Add profile"
		bind:form
		customFormats={app.customFormats}
		{qualities}
		{loadingQualities}
		{qualityError}
		{saving}
		{saveError}
		submitLabel="Create profile"
		onSubmit={submitProfile}
	/>
{/if}
