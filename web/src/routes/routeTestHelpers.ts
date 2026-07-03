import { render } from 'svelte/server';
import type { Component as SvelteComponent } from 'svelte';

import { createAppShell } from '$lib/components/rendered/appShellTestData';
import RenderWithAppShell from '$lib/components/rendered/RenderWithAppShell.svelte';
import type { AppShellController } from '$lib/features/app/appShellContext';

export function renderPage(
	component: SvelteComponent,
	componentProps: Record<string, unknown> = {},
	app: AppShellController = createAppShell()
) {
	return render(RenderWithAppShell, {
		props: {
			app,
			component: component as SvelteComponent<Record<string, unknown>>,
			componentProps
		}
	});
}
