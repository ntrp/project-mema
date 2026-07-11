import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
import type { QualitySizeSettingRequest } from '$lib/settings/types';
import { listQualitySizeSettings, updateQualitySizeSettings } from './api';

export const qualitySizeKey = ['settings', 'quality-sizes'] as const;

export function createQualitySizeResources() {
	const client = useQueryClient();
	return {
		query: createQuery(() => ({
			queryKey: qualitySizeKey,
			queryFn: listQualitySizeSettings,
			select: (response) => response.qualities
		})),
		update: createMutation(() => ({
			mutationFn: (request: QualitySizeSettingRequest[]) => updateQualitySizeSettings(request),
			onSuccess: (response) => client.setQueryData(qualitySizeKey, response)
		}))
	};
}
