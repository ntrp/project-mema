import { client } from '$lib/api/client';
import {
	normalizeMediaProfileForm,
	normalizeMetadataProviderForm,
	normalizeSubtitleProviderForm
} from '../forms';
import type { MediaProfileForm, MetadataProviderForm, SubtitleProviderForm } from '../types';

export async function saveMetadataProvider(form: MetadataProviderForm) {
	const body = normalizeMetadataProviderForm(form);
	const result = form.id
		? await client.PUT('/settings/metadata-providers/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/metadata-providers', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function saveSubtitleProvider(form: SubtitleProviderForm) {
	const body = normalizeSubtitleProviderForm(form);
	const result = form.id
		? await client.PUT('/settings/subtitle-providers/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/subtitle-providers', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function saveMediaProfile(form: MediaProfileForm) {
	const body = normalizeMediaProfileForm(form);
	const result = form.id
		? await client.PUT('/settings/profiles/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/profiles', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function testMetadataProvider(id: string) {
	const { data, error } = await client.POST('/settings/metadata-providers/{id}/test', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Metadata provider test did not return a result');
	}
	return data;
}

export async function testSubtitleProvider(id: string) {
	const { data, error } = await client.POST('/settings/subtitle-providers/{id}/test', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Subtitle provider test did not return a result');
	}
	return data;
}
