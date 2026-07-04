import { describe, expect, it } from 'vitest';

import { filterCustomFormats } from '$lib/settings/customFormatFilters';
import type { CustomFormat } from '$lib/settings/types';

describe('custom format filters', () => {
	it('returns all custom formats for a blank query', () => {
		const formats = [customFormat('format-1', 'HDR'), customFormat('format-2', 'Anime')];

		expect(filterCustomFormats(formats, '   ')).toEqual(formats);
	});

	it('matches custom formats by name only', () => {
		const formats = [
			customFormat('format-1', 'HDR', { name: 'Required Source', type: 'source', value: 'WEB-DL' }),
			customFormat('format-2', 'Dual Audio', {
				name: 'German Audio',
				type: 'language',
				value: 'German'
			})
		];

		expect(filterCustomFormats(formats, 'hdr').map((format) => format.id)).toEqual(['format-1']);
		expect(filterCustomFormats(formats, 'dual audio').map((format) => format.id)).toEqual([
			'format-2'
		]);
		expect(filterCustomFormats(formats, 'language german')).toEqual([]);
	});
});

function customFormat(
	id: string,
	name: string,
	spec: { name: string; type: 'source' | 'language'; value: string } = {
		name: 'Required Source',
		type: 'source',
		value: 'WEB-DL'
	}
): CustomFormat {
	return {
		id,
		name,
		createdAt: '2026-07-04T00:00:00Z',
		updatedAt: '2026-07-04T00:00:00Z',
		includeInRenameTemplate: false,
		includeSpecs: [{ id: `${id}-spec`, required: true, ...spec }],
		excludeSpecs: []
	};
}
