import type { createActivityCache } from '$lib/features/activity/cache';
import type { createLibraryCache } from '$lib/features/library/cache';
import type { createReleaseCache } from '$lib/features/releases/cache';
import type { createDiscoverBlacklistCache } from '$lib/features/discovery/blacklist/cache';
import type { createDiscoverContentCache } from '$lib/features/discovery/content/cache';
import type { createSearchCache } from '$lib/features/search/cache';
import type { MediaAdvancedSearchRequest } from '$lib/features/search/api';
import type { createControllerResourceRuntime } from './controllerResourceRuntime.svelte';

export interface ControllerActionDeps {
	activity: ReturnType<typeof createActivityCache>;
	library: ReturnType<typeof createLibraryCache>;
	releases: ReturnType<typeof createReleaseCache>;
	blacklist: ReturnType<typeof createDiscoverBlacklistCache>;
	discovery: ReturnType<typeof createDiscoverContentCache>;
	search: ReturnType<typeof createSearchCache>;
	resources: ReturnType<typeof createControllerResourceRuntime>;
	setAutocomplete: (_query: string) => void;
	setAdvanced: (_request: MediaAdvancedSearchRequest) => void;
}
