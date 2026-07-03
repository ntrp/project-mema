import { resolve } from '$app/paths';
import type { MediaMetadataDetails, MediaMetadataFact } from '$lib/settings/types';

type MediaMetadataPerson = NonNullable<MediaMetadataDetails['cast']>[number];

export interface CrewRolePreview {
	role: string;
	names: string[];
}

export interface MediaPersonCardData {
	name: string;
	role?: string;
	image?: string;
	externalProvider?: string;
	externalId?: string;
}

export interface MediaPersonGroup {
	title: string;
	people: MediaPersonCardData[];
}

const crewRolePriority = ['Creator', 'Director', 'Writer', 'Producer', 'Editor'];
export const crewRoleLabels = [...crewRolePriority];

export function castPeople(cast: MediaMetadataPerson[] = []): MediaPersonCardData[] {
	return cast.map(personData);
}

export function crewPersonGroups(
	crew: MediaMetadataPerson[] = [],
	facts: MediaMetadataFact[] = []
): MediaPersonGroup[] {
	const groups = crew.length > 0 ? structuredCrewGroups(crew) : factCrewGroups(facts);
	return groups.filter((group) => group.people.length > 0);
}

export function crewRolePreviews(
	facts: MediaMetadataFact[],
	crew: MediaMetadataPerson[] = []
): CrewRolePreview[] {
	return crewPersonGroups(crew, facts)
		.map((group) => ({
			role: group.title,
			names: group.people.map((person) => person.name).slice(0, 3)
		}))
		.filter((preview) => preview.names.length > 0);
}

export function mediaPersonHref(person: { externalProvider?: string; externalId?: string }) {
	if (!person.externalProvider || !person.externalId) {
		return undefined;
	}
	return resolve('/people/[provider]/[personId]', {
		provider: person.externalProvider,
		personId: person.externalId
	});
}

function structuredCrewGroups(crew: MediaMetadataPerson[]): MediaPersonGroup[] {
	const mapped = new Map<string, MediaPersonCardData[]>();
	for (const person of crew) {
		const card = personData(person);
		const role = card.role ?? 'Crew';
		mapped.set(role, [...(mapped.get(role) ?? []), card]);
	}
	const orderedRoles = [
		...crewRoleLabels,
		...[...mapped.keys()].filter((role) => !crewRoleLabels.includes(role)).sort()
	];
	return orderedRoles.map((title) => ({ title, people: mapped.get(title) ?? [] }));
}

function factCrewGroups(facts: MediaMetadataFact[]): MediaPersonGroup[] {
	return crewRoleLabels.map((title) => ({
		title,
		people: namesFromFact(facts.find((fact) => fact.label === title)?.value ?? '').map((name) => ({
			name
		}))
	}));
}

function personData(person: MediaMetadataPerson): MediaPersonCardData {
	return {
		name: person.name,
		role: person.role,
		image: person.profilePath,
		externalProvider: person.externalProvider,
		externalId: person.externalId
	};
}

function namesFromFact(value: string) {
	return value
		.split(',')
		.map((name) => name.trim())
		.filter(Boolean);
}
