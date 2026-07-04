import { render } from 'svelte/server';
import type { Component as SvelteComponent } from 'svelte';
import { describe, expect, it, vi } from 'vitest';

import AuthPanel from '$lib/components/settings/AuthPanel.svelte';
import { createAppShell } from '$lib/components/rendered/appShellTestData';
import { mediaProfile } from '$lib/components/rendered/appShellTestValues';
import RenderWithAppShell from '$lib/components/rendered/RenderWithAppShell.svelte';
import ProfileRoute from '$lib/features/profile/ProfileRoute.svelte';
import EditMediaProfileRoute from '$lib/features/settings/routes/profiles/EditMediaProfileRoute.svelte';
import NewMediaProfileRoute from '$lib/features/settings/routes/profiles/NewMediaProfileRoute.svelte';
import type { MediaProfile } from '$lib/settings/types';

describe('rendered auth and profile routes (SCN-AUTH-001, SCN-SETTINGS-023)', () => {
	it('renders the admin login form with required credentials fields', () => {
		const { body } = render(AuthPanel, {
			props: {
				username: 'admin',
				password: '',
				onLogin: vi.fn()
			}
		});

		expect(body).toContain('Admin login');
		expect(body).toContain('Username');
		expect(body).toContain('Password');
		expect(body).toContain('Log in');
		expect(body).toContain('autocomplete="username"');
		expect(body).toContain('autocomplete="current-password"');
	});

	it('renders profile creation and edit states from shell data', () => {
		const create = renderProfileRoute(NewMediaProfileRoute);
		expect(create.body).toContain('Add profile');
		expect(create.body).toContain('Create profile');
		expect(create.body).toContain('Qualities');
		expect(create.body).toContain('Target languages');

		const edit = renderProfileRoute(EditMediaProfileRoute, { profileId: 'profile-1' }, [
			mediaProfile()
		]);
		expect(edit.body).toContain('Edit profile');
		expect(edit.body).toContain('Update profile');

		const missing = renderProfileRoute(EditMediaProfileRoute, { profileId: 'missing' }, [
			mediaProfile()
		]);
		expect(missing.body).toContain('Profile not found.');
		expect(missing.body).toContain('Back');
	});

	it('renders the current user profile page', () => {
		const { body } = renderProfileRoute(ProfileRoute);

		expect(body).toContain('Profile');
		expect(body).toContain('Scenario Admin');
		expect(body).toContain('Picture URL');
		expect(body).toContain('Save profile');
	});
});

function renderProfileRoute(
	component: unknown,
	componentProps: Record<string, unknown> = {},
	mediaProfiles: MediaProfile[] = []
) {
	return render(RenderWithAppShell, {
		props: {
			app: createAppShell({ mediaProfiles }),
			component: component as SvelteComponent<Record<string, unknown>>,
			componentProps
		}
	});
}
