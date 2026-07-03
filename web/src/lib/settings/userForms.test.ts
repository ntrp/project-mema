import { describe, expect, it } from 'vitest';

import {
	emptyUserForm,
	normalizeUserCreateForm,
	normalizeUserUpdateForm,
	userFormFromUser
} from './userForms';
import type { ManagedUser } from './types';

describe('user forms (SCN-SETTINGS-012)', () => {
	it('creates and populates user forms without exposing a password', () => {
		expect(emptyUserForm()).toEqual({ username: '', password: '', role: 'user' });
		expect(
			userFormFromUser({ id: 'user-1', username: 'admin', role: 'admin' } as ManagedUser)
		).toEqual({
			id: 'user-1',
			username: 'admin',
			password: '',
			role: 'admin'
		});
	});

	it('trims usernames and omits blank update passwords', () => {
		expect(
			normalizeUserCreateForm({ username: ' admin ', password: ' secret ', role: 'admin' })
		).toEqual({
			username: 'admin',
			password: ' secret ',
			role: 'admin'
		});
		expect(normalizeUserUpdateForm({ username: ' viewer ', password: ' ', role: 'user' })).toEqual({
			username: 'viewer',
			password: undefined,
			role: 'user'
		});
	});
});
