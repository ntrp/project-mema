import { getContext, setContext } from 'svelte';
import type { createAppShellController } from '$lib/components/app/shell/controller/index.svelte';

const APP_SHELL_CONTEXT = Symbol('app-shell');

export type AppShellController = ReturnType<typeof createAppShellController>;

export function setAppShellContext(app: AppShellController) {
	setContext(APP_SHELL_CONTEXT, app);
}

export function getAppShellContext() {
	return getContext<AppShellController>(APP_SHELL_CONTEXT);
}
