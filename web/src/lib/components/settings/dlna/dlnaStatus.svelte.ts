import { createQuery } from '@tanstack/svelte-query';
import { getDLNASettings } from './api';
export const dlnaStatusQueryKey = ['settings', 'dlna', 'settings'] as const;

export function createDLNAStatusQuery() {
	return createQuery(() => ({
		queryKey: dlnaStatusQueryKey,
		queryFn: getDLNASettings
	}));
}
