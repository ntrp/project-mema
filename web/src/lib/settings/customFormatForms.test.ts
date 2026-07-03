import { describe, expect, it } from 'vitest';

import {
	customFormatFormFromFormat,
	emptyCustomFormatForm,
	normalizeCustomFormatForm
} from './customFormatForms';
import type { CustomFormat } from './types';

describe('custom format forms (SCN-SETTINGS-012)', () => {
	it('creates an empty editable form', () => {
		expect(emptyCustomFormatForm()).toEqual({
			name: '',
			includeInRenameTemplate: false,
			includeSpecs: [],
			excludeSpecs: []
		});
	});

	it('copies specs so editing the form does not mutate the source format', () => {
		const format = {
			id: 'format-1',
			name: 'WEB Required',
			includeInRenameTemplate: true,
			includeSpecs: [
				{ id: 'source', name: 'Source', type: 'source', value: 'WEB', required: true }
			],
			excludeSpecs: [],
			createdAt: '2026-01-01T00:00:00Z',
			updatedAt: '2026-01-01T00:00:00Z'
		} as CustomFormat;

		const form = customFormatFormFromFormat(format);
		form.includeSpecs[0].value = 'BLURAY';

		expect(format.includeSpecs[0].value).toBe('WEB');
		expect(form.id).toBe('format-1');
	});

	it('trims request fields and drops incomplete specs', () => {
		expect(
			normalizeCustomFormatForm({
				name: '  Anime WEB  ',
				includeInRenameTemplate: true,
				includeSpecs: [
					{ id: ' source ', name: ' Source ', type: 'source', value: ' WEB ', required: true },
					{ id: '', name: 'Empty', type: 'source', value: 'WEB', required: false }
				],
				excludeSpecs: [{ id: 'cam', name: ' CAM ', type: 'source', value: ' ', required: true }]
			})
		).toEqual({
			name: 'Anime WEB',
			includeInRenameTemplate: true,
			includeSpecs: [
				{ id: 'source', name: 'Source', type: 'source', value: 'WEB', required: true }
			],
			excludeSpecs: []
		});
	});
});
