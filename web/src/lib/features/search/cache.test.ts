import type { QueryClient } from '@tanstack/svelte-query';
import { describe, expect, it, vi } from 'vitest';
import { createSearchCache } from './cache';
import { searchKeys } from './queries.svelte';

describe('search cache', () => {
	it('clears only search-owned query keys', () => {
		const removeQueries = vi.fn();
		createSearchCache({ removeQueries } as unknown as QueryClient).clear();
		expect(removeQueries).toHaveBeenCalledWith({ queryKey: searchKeys.all });
	});
});
