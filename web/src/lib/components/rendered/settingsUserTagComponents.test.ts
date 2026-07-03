import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import TagSettings from '$lib/components/settings/tags/TagSettings.svelte';
import UserForm from '$lib/components/settings/users/UserForm.svelte';
import UsersSettingsSection from '$lib/components/settings/users/UsersSettingsSection.svelte';
import UserTable from '$lib/components/settings/users/UserTable.svelte';
import type { ManagedUser, Tag, UserForm as UserFormValue, UserSummary } from '$lib/settings/types';

describe('rendered user settings components (SCN-SETTINGS-024)', () => {
	it('renders users with the current-user marker and disabled self delete', () => {
		const { body } = render(UserTable, {
			props: {
				users: [managedUser({ id: 'user-1', username: 'admin', role: 'admin' })],
				currentUserId: 'user-1',
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});

		expect(body).toContain('admin');
		expect(body).toContain('Current');
		expect(body).toContain('Delete admin');
		expect(body).toContain('disabled');
	});

	it('renders the users settings section heading and empty table state', () => {
		const { body } = render(UsersSettingsSection, {
			props: {
				users: [],
				currentUser: { id: 'user-1', username: 'admin', role: 'admin' } as UserSummary,
				form: userForm(),
				saving: false,
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});

		expect(body).toContain('Settings');
		expect(body).toContain('Users');
		expect(body).toContain('Add user');
		expect(body).toContain('No users configured');
	});

	it('renders add and edit user forms with password requirements', () => {
		const add = render(UserForm, {
			props: {
				form: userForm(),
				saving: true,
				onSave: vi.fn(),
				onCancel: vi.fn()
			}
		});
		expect(add.body).toContain('Add user');
		expect(add.body).toContain('Password');
		expect(add.body).toContain('required');
		expect(add.body).toContain('Saving');

		const edit = render(UserForm, {
			props: {
				form: userForm({ id: 'user-1', username: 'admin', password: '' }),
				saving: false,
				onSave: vi.fn(),
				onCancel: vi.fn()
			}
		});
		expect(edit.body).toContain('Edit user');
		expect(edit.body).toContain('New password');
		expect(edit.body).toContain('Leave blank to keep current password');
		expect(edit.body).toContain('Cancel');
	});
});

describe('rendered tag settings components (SCN-SETTINGS-002)', () => {
	it('renders tags, deletion state, and empty fallback', () => {
		const populated = render(TagSettings, {
			props: {
				tags: [tag({ id: 'tag-1', name: 'priority' })],
				form: { name: '' },
				saving: false,
				deletingId: 'tag-1',
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});
		expect(populated.body).toContain('Add tag');
		expect(populated.body).toContain('priority');
		expect(populated.body).toContain('Deleting priority');
		expect(populated.body).toContain('disabled');

		const empty = render(TagSettings, {
			props: {
				tags: [],
				form: { name: '' },
				saving: false,
				onSave: vi.fn(),
				onCancel: vi.fn(),
				onEdit: vi.fn(),
				onDelete: vi.fn()
			}
		});
		expect(empty.body).toContain('No tags configured');
	});
});

function managedUser(overrides: Partial<ManagedUser> = {}): ManagedUser {
	return {
		id: 'user-1',
		username: 'scenario-user',
		role: 'user',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}

function userForm(overrides: Partial<UserFormValue> = {}): UserFormValue {
	return {
		username: 'scenario-user',
		password: 'long-password',
		role: 'user',
		...overrides
	};
}

function tag(overrides: Partial<Tag> = {}): Tag {
	return {
		id: 'tag-1',
		name: 'scenario',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:00:00Z',
		...overrides
	};
}
