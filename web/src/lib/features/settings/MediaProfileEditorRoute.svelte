<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import MediaProfileCustomFormatScores from '$lib/components/settings/MediaProfileCustomFormatScores.svelte';
	import MediaProfileQualitySelector from '$lib/components/settings/MediaProfileQualitySelector.svelte';
	import MediaProfileRules from '$lib/components/settings/MediaProfileRules.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import { listQualitySizeSettings, saveMediaProfile } from '$lib/settings/api';
	import { emptyMediaProfileForm, mediaProfileFormFromProfile } from '$lib/settings/forms';
	import type { MediaProfile, MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';
	import { errorMessageFrom } from '$lib/components/app/shell/controller/helpers';

	interface Props {
		profileId?: string;
	}

	let { profileId }: Props = $props();
	const app = getAppShellContext();
	let form = $state<MediaProfileForm>(emptyMediaProfileForm());
	let activeKey = $state('');
	let qualities = $state<QualitySizeSetting[]>([]);
	let loadingQualities = $state(false);
	let qualityError = $state('');
	let saving = $state(false);
	let saveError = $state('');

	const profile = $derived(
		profileId ? app.mediaProfiles.find((item: MediaProfile) => item.id === profileId) : undefined
	);
	const title = $derived(profileId ? 'Edit profile' : 'Add profile');
	const canSave = $derived(form.name.trim() !== '' && form.qualityIds.length > 0 && !saving);
	const notFound = $derived(profileId !== undefined && app.mediaProfiles.length > 0 && !profile);

	onMount(() => {
		void loadQualities();
	});

	$effect(() => {
		const key = profileId ?? 'new';
		if (activeKey === key) return;
		form = profile ? mediaProfileFormFromProfile(profile) : emptyMediaProfileForm();
		activeKey = key;
	});

	async function loadQualities() {
		loadingQualities = true;
		qualityError = '';
		try {
			const response = await listQualitySizeSettings();
			qualities = response.qualities;
		} catch (error) {
			qualityError = errorMessageFrom(error, 'Could not load qualities');
		} finally {
			loadingQualities = false;
		}
	}

	function updateForm(value: MediaProfileForm) {
		form = value;
	}

	async function submitProfile(event: SubmitEvent) {
		event.preventDefault();
		if (!canSave) return;
		saving = true;
		saveError = '';
		app.clearNotice();
		try {
			await saveMediaProfile(form);
			app.message = 'Profile saved';
			await app.loadSettings();
			await goto(resolve('/settings/profiles'));
		} catch (error) {
			saveError = errorMessageFrom(error, 'Could not save profile');
		} finally {
			saving = false;
		}
	}
</script>

{#if app.isAdmin}
	<PageHeading eyebrow="Settings" {title} titleId="settings-profile-title" />

	{#if notFound}
		<Card class="p-5">
			<p class="m-0 text-sm text-muted-foreground">Profile not found.</p>
			<Button class="mt-4" variant="outline" href={resolve('/settings/profiles')}>Back</Button>
		</Card>
	{:else}
		<Card class="p-5">
			<form class="grid gap-4" onsubmit={submitProfile}>
				{#if saveError}
					<p
						class="m-0 rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2.5 text-sm font-bold text-destructive"
					>
						{saveError}
					</p>
				{/if}
				<div
					class="grid min-w-0 items-start gap-4.5 min-[981px]:grid-cols-[minmax(340px,0.82fr)_minmax(440px,1fr)]"
				>
					<div class="grid min-w-0 gap-3.5">
						<div class="grid gap-1.5">
							<Label>Name</Label>
							<Input bind:value={form.name} type="text" maxlength={200} required />
						</div>

						<MediaProfileRules {form} {qualities} onChange={updateForm} />
						<MediaProfileCustomFormatScores
							{form}
							customFormats={app.customFormats}
							onChange={updateForm}
						/>
					</div>

					<aside class="grid min-w-0 gap-3.5 min-[981px]:sticky min-[981px]:top-0">
						<MediaProfileQualitySelector
							{form}
							{qualities}
							loading={loadingQualities}
							error={qualityError}
							onChange={updateForm}
						/>
					</aside>
				</div>

				<div class="flex items-center justify-end gap-3">
					<Button type="button" variant="outline" href={resolve('/settings/profiles')}
						>Cancel</Button
					>
					<Button type="submit" disabled={!canSave}>
						{saving ? 'Saving' : profileId ? 'Update profile' : 'Create profile'}
					</Button>
				</div>
			</form>
		</Card>
	{/if}
{/if}
