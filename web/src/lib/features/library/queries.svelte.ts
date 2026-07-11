import { createQuery } from '@tanstack/svelte-query';
import { listMediaItems, listMediaRequests } from './api';

export const libraryKeys = {
	all: ['library'] as const,
	items: () => [...libraryKeys.all, 'items'] as const,
	requests: () => [...libraryKeys.all, 'requests'] as const
};

export function createMediaItemsQuery(enabled: () => boolean = () => true) {
	return createQuery(() => ({
		queryKey: libraryKeys.items(),
		queryFn: ({ signal }) => listMediaItems({ signal }),
		enabled: enabled(),
		select: (response) => response.items
	}));
}

export function createMediaRequestsQuery() {
	return createQuery(() => ({
		queryKey: libraryKeys.requests(),
		queryFn: ({ signal }) => listMediaRequests({ signal }),
		select: (response) => response.requests
	}));
}
