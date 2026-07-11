import { client } from '$lib/api/client';
import type { CustomFormat, MediaProfile } from '../types';

export async function listMediaProfiles(): Promise<MediaProfile[]> {
	const { data, error } = await client.GET('/settings/profiles');

	if (error) {
		throw new Error(error.message);
	}
	return data?.profiles ?? [];
}

export async function listCustomFormats(): Promise<CustomFormat[]> {
	const { data, error } = await client.GET('/settings/custom-formats');

	if (error) {
		throw new Error(error.message);
	}
	return data?.formats ?? [];
}

export async function testCustomFormatParsing(fileName: string) {
	const { data, error } = await client.POST('/settings/custom-formats/test-parsing', {
		body: { fileName }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Parsing result was not returned');
	}
	return data;
}
