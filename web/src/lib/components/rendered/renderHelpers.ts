import { render } from 'svelte/server';
import type { Component as SvelteComponent } from 'svelte';

import RenderWithSidebar from './RenderWithSidebar.svelte';
import RenderWithTooltip from './RenderWithTooltip.svelte';

export function renderWithTooltip<Props extends Record<string, unknown>>(
	component: SvelteComponent<Props>,
	componentProps: Props
) {
	return render(RenderWithTooltip, {
		props: {
			component: component as SvelteComponent<Record<string, unknown>>,
			componentProps
		}
	});
}

export function renderWithSidebar<Props extends Record<string, unknown>>(
	component: SvelteComponent<Props>,
	componentProps: Props
) {
	return render(RenderWithSidebar, {
		props: {
			component: component as SvelteComponent<Record<string, unknown>>,
			componentProps
		}
	});
}
