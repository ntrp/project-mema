import { describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({ createQuery: vi.fn((options: () => unknown) => options()) }));
vi.mock('@tanstack/svelte-query', () => ({ createQuery: mocks.createQuery }));
vi.mock('./api', () => ({ listLanguages: vi.fn(), listTags: vi.fn(), listUsers: vi.fn() }));

import {
	createLanguagesQuery,
	createTagsQuery,
	createUsersQuery,
	settingsCatalogKeys
} from './queries.svelte';

describe('settings catalog queries', () => {
	it('uses stable keys and authentication guards', () => {
		expect(createLanguagesQuery(() => true)).toMatchObject({
			queryKey: settingsCatalogKeys.languages(),
			enabled: true
		});
		expect(createTagsQuery(() => false)).toMatchObject({ enabled: false });
		expect(createUsersQuery(() => true)).toMatchObject({ queryKey: settingsCatalogKeys.users() });
	});
});
