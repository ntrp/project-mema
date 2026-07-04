<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { UserProfile } from '$lib/settings/types';
	import ProfileAvatarPreview from './ProfileAvatarPreview.svelte';

	interface Props {
		profile: UserProfile;
		saving: boolean;
		onSave: (_value: { displayName: string; pictureUrl: string }) => void | Promise<void>;
	}

	let { profile, saving, onSave }: Props = $props();
	// svelte-ignore state_referenced_locally
	let displayName = $state(profile.displayName);
	// svelte-ignore state_referenced_locally
	let pictureUrl = $state(profile.pictureUrl);
	const changed = $derived(
		displayName.trim() !== profile.displayName || pictureUrl.trim() !== profile.pictureUrl
	);

	function submit(event: SubmitEvent) {
		event.preventDefault();
		void onSave({ displayName: displayName.trim(), pictureUrl: pictureUrl.trim() });
	}
</script>

<Card.Root aria-labelledby="profile-form-title">
	<Card.Header>
		<Card.Title id="profile-form-title">Profile</Card.Title>
		<Card.Description>{profile.username} · {profile.role}</Card.Description>
	</Card.Header>
	<Card.Content>
		<form class="grid gap-5" onsubmit={submit}>
			<div class="flex flex-col gap-4 sm:flex-row sm:items-start">
				<ProfileAvatarPreview name={displayName} {pictureUrl} username={profile.username} />
				<div class="grid flex-1 gap-4">
					<div class="grid gap-1.5">
						<Label for="profile-display-name">Name</Label>
						<Input
							id="profile-display-name"
							bind:value={displayName}
							maxlength={200}
							autocomplete="name"
						/>
					</div>
					<div class="grid gap-1.5">
						<Label for="profile-picture-url">Picture URL</Label>
						<Input
							id="profile-picture-url"
							bind:value={pictureUrl}
							maxlength={2000}
							autocomplete="photo"
						/>
					</div>
				</div>
			</div>
			<div class="flex justify-end">
				<Button type="submit" disabled={saving || !changed}>
					{saving ? 'Saving' : 'Save profile'}
				</Button>
			</div>
		</form>
	</Card.Content>
</Card.Root>
