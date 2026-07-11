<script lang="ts">
	import { resolve } from '$app/paths';
	import MediaProfileEditorForm from '$lib/components/settings/profiles/MediaProfileEditorForm.svelte';
	import { createQualitySizeResources } from '$lib/components/settings/quality/resources.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { emptyMediaProfileForm, mediaProfileFormFromProfile } from '$lib/settings/forms';
	import type { MediaProfile, MediaProfileForm } from '$lib/settings/types';
	import { errorMessageFrom } from '$lib/components/app/shell/controller/helpers';
	import { returnToMediaProfiles } from './mediaProfileRouteHelpers';
	import { createProfileEditorResources } from './profileEditorResources.svelte';

	interface Props {
		profileId: string;
	}

	let { profileId }: Props = $props();
	const app = getAppShellContext();
	const resources = createProfileEditorResources();
	const profile = $derived(app.mediaProfiles.find((item: MediaProfile) => item.id === profileId));
	let form = $derived<MediaProfileForm>(
		profile ? mediaProfileFormFromProfile(profile) : emptyMediaProfileForm()
	);
	const qualitySizes = createQualitySizeResources();
	const qualities = $derived(qualitySizes.query.data ?? []);
	const loadingQualities = $derived(qualitySizes.query.isFetching);
	const qualityError = $derived(qualitySizes.query.error?.message ?? '');
	let saving = $state(false);
	let saveError = $state('');

	const notFound = $derived(app.mediaProfiles.length > 0 && !profile);

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
