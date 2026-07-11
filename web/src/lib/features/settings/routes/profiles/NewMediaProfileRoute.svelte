<script lang="ts">
	import MediaProfileEditorForm from '$lib/components/settings/profiles/MediaProfileEditorForm.svelte';
	import { createQualitySizeResources } from '$lib/components/settings/quality/resources.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { emptyMediaProfileForm } from '$lib/settings/forms';
	import type { MediaProfileForm } from '$lib/settings/types';
	import { errorMessageFrom } from '$lib/components/app/shell/controller/helpers';
	import { returnToMediaProfiles } from './mediaProfileRouteHelpers';
	import { createProfileEditorResources } from './profileEditorResources.svelte';

	const app = getAppShellContext();
	const resources = createProfileEditorResources();
	let form = $state<MediaProfileForm>(emptyMediaProfileForm());
	const qualitySizes = createQualitySizeResources();
	const qualities = $derived(qualitySizes.query.data ?? []);
	const loadingQualities = $derived(qualitySizes.query.isFetching);
	const qualityError = $derived(qualitySizes.query.error?.message ?? '');
	let saving = $state(false);
	let saveError = $state('');

	async function submitProfile(event: SubmitEvent) {
		event.preventDefault();
		saving = true;
		saveError = '';
		app.clearNotice();
		try {
			await resources.save.mutateAsync(form);
			app.message = 'Profile saved';
			await app.loadSettingsSection('profiles');
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
		languages={app.languages}
		{qualities}
		{loadingQualities}
		{qualityError}
		{saving}
		{saveError}
		submitLabel="Create profile"
		onSubmit={submitProfile}
	/>
{/if}
