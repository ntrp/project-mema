import { createGetSystemStatus } from '$lib/api/generated/tanstack';

export function createSystemStatusQuery() {
	return createGetSystemStatus();
}
