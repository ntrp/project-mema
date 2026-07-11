import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
import {
	getFileDeleteSettings,
	getFileNamingSettings,
	updateFileDeleteSettings,
	updateFileNamingSettings
} from '$lib/components/settings/library/filePoliciesApi';

export const filePolicyKeys = {
	naming: () => ['settings', 'file-policy', 'naming'] as const,
	deletion: () => ['settings', 'file-policy', 'deletion'] as const
};

export function createFileNamingResource() {
	const client = useQueryClient();
	const query = createQuery(() => ({
		queryKey: filePolicyKeys.naming(),
		queryFn: getFileNamingSettings
	}));
	const save = createMutation(() => ({
		mutationFn: updateFileNamingSettings,
		onSuccess: (settings) => client.setQueryData(filePolicyKeys.naming(), settings)
	}));
	return { query, save };
}

export function createFileDeleteResource() {
	const client = useQueryClient();
	const query = createQuery(() => ({
		queryKey: filePolicyKeys.deletion(),
		queryFn: getFileDeleteSettings
	}));
	const save = createMutation(() => ({
		mutationFn: updateFileDeleteSettings,
		onSuccess: (settings) => client.setQueryData(filePolicyKeys.deletion(), settings)
	}));
	return { query, save };
}
