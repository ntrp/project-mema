import { client } from '$lib/api/client';
import { normalizeUserCreateForm, normalizeUserUpdateForm } from '../forms';
import type { UserForm } from '../types';

export async function saveUser(form: UserForm) {
	const result = form.id
		? await client.PUT('/settings/users/{id}', {
				params: { path: { id: form.id } },
				body: normalizeUserUpdateForm(form)
			})
		: await client.POST('/settings/users', { body: normalizeUserCreateForm(form) });

	if (result.error) throw new Error(result.error.message);
}

export async function deleteUser(id: string) {
	const { error } = await client.DELETE('/settings/users/{id}', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
}
