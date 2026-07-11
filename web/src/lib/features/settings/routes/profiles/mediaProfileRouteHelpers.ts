import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import { saveMediaProfile } from '$lib/settings/api';
import { listQualitySizeSettings } from '$lib/components/settings/quality/api';
import type { MediaProfileForm } from '$lib/settings/types';

export async function loadMediaProfileQualities() {
	const response = await listQualitySizeSettings();
	return response.qualities;
}

export async function saveMediaProfileForm(form: MediaProfileForm) {
	await saveMediaProfile(form);
}

export async function returnToMediaProfiles() {
	await goto(resolve('/settings/profiles'));
}
