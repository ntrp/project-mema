import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
export async function returnToMediaProfiles() {
	await goto(resolve('/settings/profiles'));
}
