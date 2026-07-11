import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
import { applyMediaRename, previewMediaRename } from '$lib/features/library/filesApi';
import { getFileNamingSettings } from '$lib/components/settings/library/filePoliciesApi';
import { filePolicyKeys } from '$lib/features/settings/resources/filePolicies.svelte';

export const mediaRenameKeys = {
	preview: (id: string) => ['media', id, 'rename-preview'] as const
};

export function createMediaRenameResource(id: () => string) {
	const client = useQueryClient();
	const preview = createQuery(() => ({
		queryKey: mediaRenameKeys.preview(id()),
		queryFn: () => previewMediaRename(id())
	}));
	const naming = createQuery(() => ({
		queryKey: filePolicyKeys.naming(),
		queryFn: getFileNamingSettings
	}));
	const apply = createMutation(() => ({
		mutationFn: (paths: string[]) => applyMediaRename(id(), paths),
		onSuccess: (result) => client.setQueryData(mediaRenameKeys.preview(id()), result)
	}));
	return { preview, naming, apply };
}
