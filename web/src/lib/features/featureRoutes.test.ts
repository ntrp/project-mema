import { render } from 'svelte/server';
import type { Component as SvelteComponent } from 'svelte';
import { describe, expect, it } from 'vitest';

import { createAppShell } from '$lib/components/rendered/appShellTestData';
import RenderWithAppShell from '$lib/components/rendered/RenderWithAppShell.svelte';
import type { AppShellController } from '$lib/features/app/appShellContext';
import DiscoverSectionRoute from '$lib/features/discovery/DiscoverSectionRoute.svelte';
import MediaCollectionRoute from '$lib/features/media/MediaCollectionRoute.svelte';
import MetadataDetailRoute from '$lib/features/media/MetadataDetailRoute.svelte';
import AdvancedSearchRoute from '$lib/features/search/AdvancedSearchRoute.svelte';

describe('feature route wrappers (SCN-SEARCH-001)', () => {
	it('renders advanced search with enabled metadata providers from shell state', () => {
		const { body } = renderRoute(asRenderComponent(AdvancedSearchRoute), createAppShell(), {
			initialQuery: 'scenario title'
		});

		expect(body).toContain('Advanced search');
		expect(body).toContain('Metadata providers');
		expect(body).toContain('Scenario Metadata');
		expect(body).toContain('0 results');
	});
});

describe('feature route wrappers (SCN-DISCOVER-001)', () => {
	it('renders discover section fallback when provider results are unavailable', () => {
		const { body } = renderRoute(DiscoverSectionRoute, createAppShell());

		expect(body).toContain('Discovery section not available');
		expect(body).toContain('Could not load this discover section.');
	});
});

describe('feature route wrappers (SCN-METADATA-001)', () => {
	it('renders collection and metadata fallback states from shell data', () => {
		const collection = renderRoute(MediaCollectionRoute, createAppShell());
		expect(collection.body).toContain('Collection not available');
		expect(collection.body).toContain('Could not load this collection.');

		const detail = renderRoute(MetadataDetailRoute, createAppShell());
		expect(detail.body).toContain('Details not available');
		expect(detail.body).toContain('Could not load provider metadata for this item.');
	});
});

function renderRoute(
	component: SvelteComponent<Record<string, unknown>>,
	app: AppShellController,
	componentProps: Record<string, unknown> = {}
) {
	return render(RenderWithAppShell, {
		props: {
			app,
			component: component as SvelteComponent<Record<string, unknown>>,
			componentProps
		}
	});
}

function asRenderComponent(component: unknown): SvelteComponent<Record<string, unknown>> {
	return component as SvelteComponent<Record<string, unknown>>;
}
