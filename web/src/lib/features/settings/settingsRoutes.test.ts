import { render } from 'svelte/server';
import type { Component as SvelteComponent } from 'svelte';
import { describe, expect, it } from 'vitest';

import { createAppShell } from '$lib/components/rendered/appShellTestData';
import RenderWithAppShell from '$lib/components/rendered/RenderWithAppShell.svelte';
import type { AppShellController } from '$lib/features/app/appShellContext';
import CustomFormatsSettingsRoute from '$lib/features/settings/routes/CustomFormatsSettingsRoute.svelte';
import DownloadClientsSettingsRoute from '$lib/features/settings/routes/DownloadClientsSettingsRoute.svelte';
import GeneralSettingsRoute from '$lib/features/settings/routes/GeneralSettingsRoute.svelte';
import IndexersSettingsRoute from '$lib/features/settings/routes/IndexersSettingsRoute.svelte';
import LanguagesSettingsRoute from '$lib/features/settings/routes/LanguagesSettingsRoute.svelte';
import LibrarySettingsRoute from '$lib/features/settings/routes/LibrarySettingsRoute.svelte';
import MetadataSettingsRoute from '$lib/features/settings/routes/MetadataSettingsRoute.svelte';
import ProfilesSettingsRoute from '$lib/features/settings/routes/ProfilesSettingsRoute.svelte';
import QualitySettingsRoute from '$lib/features/settings/routes/QualitySettingsRoute.svelte';
import SubtitlesSettingsRoute from '$lib/features/settings/routes/SubtitlesSettingsRoute.svelte';
import TagsSettingsRoute from '$lib/features/settings/routes/TagsSettingsRoute.svelte';
import UsersSettingsRoute from '$lib/features/settings/routes/UsersSettingsRoute.svelte';

describe('settings route wrappers (SCN-SETTINGS-001)', () => {
	it('hides settings route content from non-admin users', () => {
		const { body } = renderRoute(UsersSettingsRoute, createAppShell({ isAdmin: false }));

		expect(body).not.toContain('Users');
		expect(body).not.toContain('Add user');
		expect(body).not.toContain('scenario-admin');
	});

	it('renders shell-backed user settings for admins', () => {
		const { body } = renderRoute(UsersSettingsRoute, createAppShell());

		expect(body).toContain('Settings');
		expect(body).toContain('Users');
		expect(body).toContain('scenario-admin');
		expect(body).toContain('Current');
		expect(body).toContain('Add user');
	});
});

describe('settings route wrappers (SCN-SETTINGS-002)', () => {
	it('renders shell-backed tag settings for admins', () => {
		const { body } = renderRoute(TagsSettingsRoute, createAppShell());

		expect(body).toContain('Settings');
		expect(body).toContain('Tags');
		expect(body).toContain('scenario-tag');
		expect(body).toContain('Add tag');
	});
});

describe('settings route wrappers (SCN-SETTINGS-016)', () => {
	it('renders library route content for admins', () => {
		const { body } = renderRoute(LibrarySettingsRoute, createAppShell());

		expect(body).toContain('Library');
		expect(body).toContain('Root Paths');
		expect(body).toContain('Add library folder');
	});
});

describe('settings route wrappers (SCN-SETTINGS-017)', () => {
	it('renders custom format route content for admins', () => {
		const { body } = renderRoute(CustomFormatsSettingsRoute, createAppShell());

		expect(body).toContain('Custom formats');
		expect(body).toContain('Filter by name');
		expect(body).toContain('No custom formats configured');
		expect(body).toContain('Add custom format');
	});
});

describe('settings route wrappers (SCN-SETTINGS-018)', () => {
	it('renders download client route content for admins', () => {
		const { body } = renderRoute(DownloadClientsSettingsRoute, createAppShell());

		expect(body).toContain('Download clients');
		expect(body).toContain('Add download client');
		expect(body).toContain('No download clients configured');
	});
});

describe('settings route wrappers (SCN-SETTINGS-020)', () => {
	it('renders indexer route content for admins', () => {
		const { body } = renderRoute(IndexersSettingsRoute, createAppShell());

		expect(body).toContain('Indexers');
		expect(body).toContain('Add indexer');
		expect(body).toContain('Indexer search settings');
	});
});

describe('settings route wrappers (SCN-SETTINGS-021)', () => {
	it('renders general and quality settings route content for admins', () => {
		const general = renderRoute(GeneralSettingsRoute, createAppShell());
		expect(general.body).toContain('General');
		expect(general.body).toContain('Event retention days');
		expect(general.body).toContain('Write logs to files');

		const quality = renderRoute(QualitySettingsRoute, createAppShell());
		expect(quality.body).toContain('Quality');
		expect(quality.body).toContain('Quality sizes');
		expect(quality.body).toContain('Save sizes');
	});
});

describe('settings route wrappers (SCN-SETTINGS-022)', () => {
	it('renders metadata provider route content for admins', () => {
		const { body } = renderRoute(MetadataSettingsRoute, createAppShell());

		expect(body).toContain('Metadata');
		expect(body).toContain('TMDB');
		expect(body).toContain('TVDB');
	});
});

describe('settings route wrappers (SCN-SETTINGS-024)', () => {
	it('renders subtitle provider route content for admins', () => {
		const { body } = renderRoute(SubtitlesSettingsRoute, createAppShell());

		expect(body).toContain('Subtitles');
		expect(body).toContain('OpenSubtitles');
		expect(body).toContain('Saved API key');
		expect(body).toContain('Saved password');
	});
});

describe('settings route wrappers (SCN-SETTINGS-023)', () => {
	it('renders profile route content for admins', () => {
		const { body } = renderRoute(ProfilesSettingsRoute, createAppShell());

		expect(body).toContain('Profiles');
		expect(body).toContain('Add profile');
		expect(body).toContain('No profiles configured');
	});
});

describe('settings route wrappers (SCN-SETTINGS-003)', () => {
	it('renders language route content for admins', () => {
		const { body } = renderRoute(LanguagesSettingsRoute, createAppShell());

		expect(body).toContain('Languages');
		expect(body).toContain('Add language');
		expect(body).toContain('No languages configured');
	});
});

function renderRoute(component: SvelteComponent, app: AppShellController) {
	return render(RenderWithAppShell, {
		props: {
			app,
			component: component as SvelteComponent<Record<string, unknown>>
		}
	});
}
