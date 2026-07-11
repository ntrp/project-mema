import { client } from '$lib/api/client';
import { normalizeLanguageForm, normalizeLanguageUpdateForm } from '../forms';
import type { LanguageForm } from '../types';

export async function saveLanguage(form: LanguageForm) {
	const result = form.originalCode
		? await client.PUT('/settings/languages/{code}', {
				params: { path: { code: form.originalCode } },
				body: normalizeLanguageUpdateForm(form)
			})
		: await client.POST('/settings/languages', { body: normalizeLanguageForm(form) });

	if (result.error) throw new Error(result.error.message);
}

export async function deleteLanguage(code: string) {
	const { error } = await client.DELETE('/settings/languages/{code}', {
		params: { path: { code } }
	});
	if (error) throw new Error(error.message);
}
