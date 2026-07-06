import { client } from '$lib/api/client';
import type { components } from '$lib/api/generated/schema';

type MediaComponentAssemblyRequest = components['schemas']['MediaComponentAssemblyRequest'];
type MediaComponentAssemblyEnqueueResponse =
	components['schemas']['MediaComponentAssemblyEnqueueResponse'];

export async function enqueueMediaComponentAssembly(
	id: string,
	request: MediaComponentAssemblyRequest
): Promise<MediaComponentAssemblyEnqueueResponse> {
	const { data, error } = await client.POST('/media/items/{id}/assemblies', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Component assembly job was not returned');
	}
	return data;
}
