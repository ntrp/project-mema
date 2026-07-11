import { client } from '$lib/api/client';
import { normalizeCustomFormatForm } from '../forms';
import type { CustomFormatForm } from '../types';

export async function saveCustomFormat(form: CustomFormatForm) {
	const body = normalizeCustomFormatForm(form);
	const result = form.id
		? await client.PUT('/settings/custom-formats/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/custom-formats', { body });

	if (result.error) throw new Error(result.error.message);
}

export async function deleteCustomFormat(id: string) {
	const { error } = await client.DELETE('/settings/custom-formats/{id}', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
}
