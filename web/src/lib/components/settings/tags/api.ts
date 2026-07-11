import { client } from '$lib/api/client';
import type { TagForm } from '$lib/settings/types';

export async function saveTag(form: TagForm) {
	const body = { name: form.name.trim() };
	const result = form.id
		? await client.PUT('/settings/tags/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/tags', { body });

	if (result.error) throw new Error(result.error.message);
}

export async function deleteTag(id: string) {
	const { error } = await client.DELETE('/settings/tags/{id}', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
}
