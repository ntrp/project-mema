<script lang="ts">
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { UserProfile } from '$lib/settings/types';
	import ProfileForm from './ProfileForm.svelte';

	interface Props {
		profile?: UserProfile;
		loading: boolean;
		saving: boolean;
		errorMessage: string;
		onRetry: () => void | Promise<void>;
		onSave: (_value: { displayName: string; pictureUrl: string }) => void | Promise<void>;
	}

	let { profile, loading, saving, errorMessage, onRetry, onSave }: Props = $props();
</script>

<div class="mx-auto grid w-full max-w-3xl gap-4">
	<PageHeading eyebrow="Account" title="Profile" titleId="profile-title" />

	{#if errorMessage}
		<div class="rounded-md border border-destructive/40 bg-destructive/10 p-4 text-destructive">
			<p class="m-0 text-sm font-bold">{errorMessage}</p>
			<Button type="button" variant="outline" class="mt-3" onclick={() => void onRetry()}>
				Retry
			</Button>
		</div>
	{/if}

	{#if loading}
		<div class="rounded-md border border-border bg-card p-5 text-sm text-muted-foreground">
			Loading profile
		</div>
	{:else if profile}
		{#key `${profile.id}:${profile.updatedAt}`}
			<ProfileForm {profile} {saving} {onSave} />
		{/key}
	{/if}
</div>
