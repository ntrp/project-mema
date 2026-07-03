import type { MediaMetadataFact } from '$lib/settings/types';

export interface CrewRolePreview {
	role: string;
	names: string[];
}

const crewRolePriority = ['Creator', 'Director', 'Writer', 'Producer', 'Editor'];
export const crewRoleLabels = [...crewRolePriority];

export function crewRolePreviews(facts: MediaMetadataFact[]): CrewRolePreview[] {
	return crewRoleLabels
		.map((role) => ({
			role,
			names: namesFromFact(facts.find((fact) => fact.label === role)?.value ?? '').slice(0, 3)
		}))
		.filter((preview) => preview.names.length > 0);
}

function namesFromFact(value: string) {
	return value
		.split(',')
		.map((name) => name.trim())
		.filter(Boolean);
}
