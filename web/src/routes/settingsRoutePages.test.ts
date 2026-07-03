import { describe, expect, it } from 'vitest';
import type { Component as SvelteComponent } from 'svelte';

import SettingsPage from './settings/+page.svelte';
import CustomFormatsSettingsPage from './settings/custom-formats/+page.svelte';
import DownloadClientsSettingsPage from './settings/download-clients/+page.svelte';
import GeneralSettingsPage from './settings/general/+page.svelte';
import IndexersSettingsPage from './settings/indexers/+page.svelte';
import LanguagesSettingsPage from './settings/languages/+page.svelte';
import LibrarySettingsPage from './settings/library/+page.svelte';
import MetadataSettingsPage from './settings/metadata/+page.svelte';
import EditProfileSettingsPage from './settings/profiles/[id]/+page.svelte';
import NewProfileSettingsPage from './settings/profiles/new/+page.svelte';
import ProfilesSettingsPage from './settings/profiles/+page.svelte';
import QualitySettingsPage from './settings/quality/+page.svelte';
import TagsSettingsPage from './settings/tags/+page.svelte';
import UsersSettingsPage from './settings/users/+page.svelte';
import { renderPage } from './routeTestHelpers';

type RouteCase = [SvelteComponent<Record<string, unknown>>, string[], Record<string, unknown>?];

describe('settings route pages (SCN-SETTINGS-024)', () => {
	const cases: RouteCase[] = [
		[asRoute(SettingsPage), ['Library', 'Root Paths']],
		[asRoute(LibrarySettingsPage), ['Library', 'Add library folder']],
		[asRoute(CustomFormatsSettingsPage), ['Custom formats', 'Add custom format']],
		[asRoute(DownloadClientsSettingsPage), ['Download clients', 'Add download client']],
		[asRoute(GeneralSettingsPage), ['General', 'Event retention days']],
		[asRoute(IndexersSettingsPage), ['Indexers', 'Indexer search settings']],
		[asRoute(LanguagesSettingsPage), ['Languages', 'Add language']],
		[asRoute(MetadataSettingsPage), ['Metadata', 'TMDB', 'TVDB']],
		[asRoute(ProfilesSettingsPage), ['Profiles', 'Add profile']],
		[asRoute(NewProfileSettingsPage), ['Add profile', 'Create profile']],
		[
			asRoute(EditProfileSettingsPage),
			['Edit profile', 'Update profile'],
			{ params: { id: 'missing' } }
		],
		[asRoute(QualitySettingsPage), ['Quality', 'Quality sizes']],
		[asRoute(TagsSettingsPage), ['Tags', 'scenario-tag']],
		[asRoute(UsersSettingsPage), ['Users', 'scenario-admin']]
	];

	it.each(cases)(
		'renders the route page delegate',
		(component, expectedText, componentProps = {}) => {
			const { body } = renderPage(component, componentProps);

			for (const text of expectedText) {
				expect(body).toContain(text);
			}
		}
	);
});

function asRoute(component: unknown): SvelteComponent<Record<string, unknown>> {
	return component as SvelteComponent<Record<string, unknown>>;
}
