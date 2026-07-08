import { client } from '$lib/api/client';
import type { components } from '$lib/api/generated/schema';

export type ManualFulfillmentAction = components['schemas']['ManualFulfillmentAction'];

export async function listManualFulfillmentActions() {
	const { data, error } = await client.GET('/media/manual-fulfillment-actions');
	if (error) {
		throw new Error(error.message);
	}
	return data?.actions ?? [];
}
