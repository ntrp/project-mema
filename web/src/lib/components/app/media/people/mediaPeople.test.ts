import { describe, expect, it } from 'vitest';

import { crewPersonGroups } from '$lib/components/app/media/people/mediaPeople';

describe('media people helpers', () => {
	it('deduplicates repeated provider crew entries within a role', () => {
		const groups = crewPersonGroups([
			{
				name: 'Andrew Stanton',
				role: 'Crew',
				externalProvider: 'tvdb',
				externalId: '505481'
			},
			{
				name: 'Other Person',
				role: 'Crew',
				externalProvider: 'tvdb',
				externalId: '42'
			},
			{
				name: 'Andrew Stanton',
				role: 'Crew',
				externalProvider: 'tvdb',
				externalId: '505481'
			}
		]);

		expect(groups.find((group) => group.title === 'Crew')?.people.map((person) => person.name)).toEqual(
			['Andrew Stanton', 'Other Person']
		);
	});
});
