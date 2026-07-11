import { describe, expect, it, vi } from 'vitest';

import { loadAppRouteData, type RouteDataDeps } from './routeData';
import { defaultRouteState } from './routeState';

describe('route data loading', () => {
	it('loads only the active settings section', async () => {
		const dependencies = deps();
		await loadAppRouteData(
			{ ...defaultRouteState(), view: 'settings', settingsSection: 'languages' },
			true,
			dependencies
		);

		expect(dependencies.loadSettingsSection).toHaveBeenCalledWith('languages');
		expect(dependencies.loadSystemSettings).not.toHaveBeenCalled();
	});

	it('does not load settings resources for non-admin users', async () => {
		const dependencies = deps();
		await loadAppRouteData(
			{ ...defaultRouteState(), view: 'settings', settingsSection: 'indexers' },
			false,
			dependencies
		);
		expect(dependencies.loadSettingsSection).not.toHaveBeenCalled();
	});

	it('leaves profile loading to its route query', async () => {
		const dependencies = deps();
		await loadAppRouteData({ ...defaultRouteState(), view: 'profile' }, false, dependencies);
		expect(dependencies.loadProfile).not.toHaveBeenCalled();
	});
});

function deps(): RouteDataDeps {
	return {
		loadSettingsSection: vi.fn(),
		loadSystemSettings: vi.fn(),
		loadMediaActionSettings: vi.fn(),
		loadProfile: vi.fn()
	};
}
