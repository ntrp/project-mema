import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { settingsCatalogKeys } from '$lib/features/settings/catalog/queries.svelte';
import { saveMediaProfile } from '$lib/settings/api';
import type { MediaProfileForm } from '$lib/settings/types';

export function createProfileEditorResources() {
	const client = useQueryClient();
	return {
		save: createMutation(() => ({
			mutationFn: (form: MediaProfileForm) => saveMediaProfile(form),
			onSuccess: () => client.invalidateQueries({ queryKey: settingsCatalogKeys.mediaProfiles() })
		}))
	};
}
